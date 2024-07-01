package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ars0915/tonfura-exercise/constant"
	"github.com/ars0915/tonfura-exercise/entity"
	"github.com/ars0915/tonfura-exercise/repo"
	"github.com/ars0915/tonfura-exercise/util/cTypes"
)

func (h BookingHandler) GetBooking(ctx context.Context, id uint) (booking entity.Booking, err error) {
	booking, err = h.db.GetBooking(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return booking, ErrorBookingNotFound
		}
		return booking, errors.Wrap(err, "get booking")
	}
	return
}

type CreateBookingParam struct {
	FlightID uint
	UserID   uint
	ClassID  uint
	Price    uint
	Amount   uint
}

func (h BookingHandler) CreateBooking(ctx context.Context, param CreateBookingParam) (booking entity.Booking, err error) {
	err = repo.WithinTransaction(ctx, h.db, func(txCtx context.Context) error {
		tx := repo.ExtractTx(txCtx)

		flight, err := tx.GetFlight(param.FlightID)
		if err != nil {
			return errors.Wrap(err, "get flight")
		}

		if err = h.redis.Lock(ctx, genFlightLockKey(param.FlightID)); err != nil {
			return err
		}
		defer h.redis.UnLock(ctx, genFlightLockKey(param.FlightID))

		var updateClass entity.Class
		var foundClass bool
		for _, class := range flight.Classes {
			if class.ID != param.ClassID {
				continue
			}

			foundClass = true
			updateClass.Sold = cTypes.Uint(*class.Sold + param.Amount)
			if *updateClass.Sold > *class.SeatAmount+*class.OversellAmount {
				return ErrorFlightSoldOut
			}

			if *updateClass.Sold == *class.SeatAmount+*class.OversellAmount {
				updateClass.Status = cTypes.String(constant.StatusSoldOut)
			}
			break
		}

		if !foundClass {
			return ErrorClassNotFound
		}

		if booking, err = tx.CreateBooking(entity.Booking{
			FlightID: cTypes.Uint(param.FlightID),
			UserID:   cTypes.Uint(param.UserID),
			ClassID:  cTypes.Uint(param.ClassID),
			Price:    cTypes.Uint(param.Price),
			Amount:   cTypes.Uint(param.Amount),
		}); err != nil {
			return errors.Wrap(err, "create booking")
		}

		if err = tx.UpdateClass(param.ClassID, updateClass); err != nil {
			return errors.Wrap(err, "update class")
		}

		return nil
	})

	return
}

func genFlightLockKey(flightID uint) string {
	return fmt.Sprintf("flight:%d", flightID)
}

func genClassLockKey(classID uint) string {
	return fmt.Sprintf("class:%d", classID)
}

type CheckInResult struct {
	Status          string
	SuggestFlightID uint
	SuggestClassID  uint
}

func (h BookingHandler) CheckInBooking(ctx context.Context, bookingID uint) (result CheckInResult, err error) {
	result.Status = constant.BookingCheckInSuccess
	err = repo.WithinTransaction(ctx, h.db, func(txCtx context.Context) error {
		tx := repo.ExtractTx(txCtx)

		booking, err := tx.GetBooking(bookingID)
		if err != nil {
			return errors.Wrap(err, "get booking")
		}

		flight, err := tx.GetFlight(*booking.FlightID)
		if err != nil {
			return errors.Wrap(err, "get flight")
		}

		class := booking.Class

		if err = h.redis.Lock(ctx, genClassLockKey(class.ID)); err != nil {
			return err
		}
		defer h.redis.UnLock(ctx, genClassLockKey(class.ID))

		if *class.CheckInAmount+*booking.Amount <= *class.SeatAmount {
			if err = tx.UpdateBooking(bookingID, entity.Booking{Status: cTypes.String(constant.BookingStatusCheckedIn)}); err != nil {
				return errors.Wrap(err, "update booking")
			}
			if err = tx.UpdateClass(class.ID, entity.Class{CheckInAmount: cTypes.Uint(*class.CheckInAmount + *booking.Amount)}); err != nil {
				return errors.Wrap(err, "update class")
			}
			return nil
		}

		if result, err = handleOversell(tx, flight, booking); err != nil {
			return errors.Wrap(err, "handle oversell")
		}
		result.Status = constant.BookingCheckInFail

		return nil
	})

	return
}

func handleOversell(tx repo.App, flight entity.Flight, booking entity.Booking) (result CheckInResult, err error) {
	class := booking.Class
	if isCloseToDeparture(*flight.DepartureAt) {
		upgradeClassID, upgrade := checkHigherClassCanCheckIn(flight, class, *booking.Amount)
		if upgrade {
			result = CheckInResult{
				SuggestFlightID: flight.ID,
				SuggestClassID:  upgradeClassID,
			}
			return
		}

		downgradeClassID, downgrade := checkLowerClassCanCheckIn(flight, class, *booking.Amount)
		if downgrade {
			result = CheckInResult{
				SuggestFlightID: flight.ID,
				SuggestClassID:  downgradeClassID,
			}
			return
		}

		result, err = handleNextFlight(tx, flight, booking)
		if err != nil {
			return result, err
		}
	} else {
		var sameFlightChangeOK bool
		result, sameFlightChangeOK = checkFlightAvailable(flight, booking)
		if sameFlightChangeOK {
			return
		}

		result, err = handleNextFlight(tx, flight, booking)
		if err != nil {
			return result, err
		}
	}
	return
}

func isCloseToDeparture(departureTime time.Time) bool {
	return time.Until(departureTime) < 1*time.Hour
}

func checkFlightAvailable(flight entity.Flight, booking entity.Booking) (CheckInResult, bool) {
	upgradeClassID, upgrade := checkHigherClassAvailable(flight, booking.Class, *booking.Amount)
	if upgrade {
		return CheckInResult{
			Status:          constant.BookingCheckInFail,
			SuggestFlightID: flight.ID,
			SuggestClassID:  upgradeClassID,
		}, true
	}

	downgradeClassID, downgrade := checkLowerClassAvailable(flight, booking.Class, *booking.Amount)
	if downgrade {
		return CheckInResult{
			Status:          constant.BookingCheckInFail,
			SuggestFlightID: flight.ID,
			SuggestClassID:  downgradeClassID,
		}, true
	}

	return CheckInResult{}, false
}

func checkHigherClassAvailable(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] <= constant.ClassTypes[*currentClass.Type] {
			continue
		}
		if *class.Sold+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func checkHigherClassCanCheckIn(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] <= constant.ClassTypes[*currentClass.Type] {
			continue
		}

		if *class.CheckInAmount+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func checkLowerClassAvailable(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] >= constant.ClassTypes[*currentClass.Type] {
			continue
		}
		if *class.Sold+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func checkLowerClassCanCheckIn(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] >= constant.ClassTypes[*currentClass.Type] {
			continue
		}
		if *class.CheckInAmount+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func checkSameClassAvailable(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] != constant.ClassTypes[*currentClass.Type] {
			continue
		}
		if *class.Sold+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func checkSameClassCanCheckIn(flight entity.Flight, currentClass entity.Class, amount uint) (classID uint, ok bool) {
	for _, class := range flight.Classes {
		if constant.ClassTypes[*class.Type] != constant.ClassTypes[*currentClass.Type] {
			continue
		}
		if *class.CheckInAmount+amount <= *class.SeatAmount {
			return class.ID, true
		}
	}
	return 0, false
}

func handleNextFlight(tx repo.App, currentFlight entity.Flight, booking entity.Booking) (CheckInResult, error) {
	flights, err := tx.ListFlights(entity.ListFlightParam{
		Source:         currentFlight.Source,
		Destination:    currentFlight.Destination,
		DepartureAfter: currentFlight.DepartureAt,
		SortBy:         cTypes.String("departure_at ASC"),
		WithClass:      true,
	})
	if err != nil {
		return CheckInResult{}, errors.Wrap(err, "failed to list flights")
	}

	for _, flight := range flights {
		changeClassID, change := checkSameClassAvailable(flight, booking.Class, *booking.Amount)
		if change {
			return CheckInResult{
				SuggestFlightID: flight.ID,
				SuggestClassID:  changeClassID,
			}, nil
		}

		flightSuggest, flightAvailable := checkFlightAvailable(flight, booking)
		if flightAvailable {
			return flightSuggest, nil
		}

		changeClassID, change = checkSameClassCanCheckIn(flight, booking.Class, *booking.Amount)
		if change {
			return CheckInResult{
				SuggestFlightID: flight.ID,
				SuggestClassID:  changeClassID,
			}, nil
		}
	}
	return CheckInResult{}, ErrorNoAvailableSeat
}

func (h BookingHandler) GiveUpBooking(ctx context.Context, bookingID uint) (booking entity.Booking, err error) {
	err = repo.WithinTransaction(ctx, h.db, func(txCtx context.Context) error {
		tx := repo.ExtractTx(txCtx)

		booking, err = tx.GetBooking(bookingID)
		if err != nil {
			return errors.Wrap(err, "get booking")
		}

		flight, err := tx.GetFlight(*booking.FlightID)
		if err != nil {
			return errors.Wrap(err, "get flight")
		}

		suggestion, err := handleNextFlight(tx, flight, booking)
		if err != nil {
			return err
		}

		if err = tx.UpdateBooking(bookingID, entity.Booking{
			FlightID: cTypes.Uint(suggestion.SuggestFlightID),
			ClassID:  cTypes.Uint(suggestion.SuggestClassID),
			Status:   cTypes.String(constant.BookingStatusCheckedIn),
		}); err != nil {
			return errors.Wrap(err, "update booking")
		}

		class, err := tx.GetClass(suggestion.SuggestClassID)
		if err != nil {
			return errors.Wrap(err, "get class")
		}

		updateClass := entity.Class{
			Sold:          cTypes.Uint(*class.Sold + *booking.Amount),
			CheckInAmount: cTypes.Uint(*class.CheckInAmount + *booking.Amount),
		}
		if *class.Sold+*booking.Amount >= *class.SeatAmount {
			updateClass.Status = cTypes.String(constant.StatusSoldOut)
		}

		if err = tx.UpdateClass(suggestion.SuggestClassID, updateClass); err != nil {
			return errors.Wrap(err, "update class")
		}

		booking, err = tx.GetBooking(bookingID)
		if err != nil {
			return errors.Wrap(err, "get booking")
		}

		return nil
	})

	return
}

func (h BookingHandler) UpdateBooking(ctx context.Context, bookingID uint, data entity.Booking) (booking entity.Booking, err error) {
	err = repo.WithinTransaction(ctx, h.db, func(txCtx context.Context) error {
		tx := repo.ExtractTx(txCtx)

		booking, err = tx.GetBooking(bookingID)
		if err != nil {
			return errors.Wrap(err, "get booking")
		}

		var class entity.Class
		if data.ClassID != nil {
			class, err = tx.GetClass(*data.ClassID)
			if err != nil {
				return errors.Wrap(err, "get class")
			}
		}

		if data.FlightID != nil {
			if data.ClassID != nil && *data.FlightID != class.FlightID {
				return ErrorClassNotBelongToFlight
			}

			if data.ClassID == nil && *data.FlightID != booking.Class.FlightID {
				return ErrorClassNotBelongToFlight
			}
		} else if data.ClassID != nil && *booking.FlightID != class.FlightID {
			return ErrorClassNotBelongToFlight
		}

		if err = h.redis.Lock(ctx, genClassLockKey(class.ID)); err != nil {
			return err
		}
		defer h.redis.UnLock(ctx, genClassLockKey(class.ID))

		if *class.CheckInAmount+*booking.Amount <= *class.SeatAmount {
			if err = tx.UpdateBooking(bookingID, entity.Booking{
				FlightID: data.FlightID,
				ClassID:  data.ClassID,
				Status:   data.Status,
			}); err != nil {
				return errors.Wrap(err, "update booking")
			}

			updateClass := entity.Class{
				Sold:          cTypes.Uint(*class.Sold + *booking.Amount),
				CheckInAmount: cTypes.Uint(*class.CheckInAmount + *booking.Amount),
			}
			if *class.Sold+*booking.Amount >= *class.SeatAmount {
				updateClass.Status = cTypes.String(constant.StatusSoldOut)
			}

			if err = tx.UpdateClass(class.ID, updateClass); err != nil {
				return errors.Wrap(err, "update class")
			}
			return nil
		} else {
			return ErrorNoAvailableSeat
		}
	})
	if err != nil {
		return entity.Booking{}, err
	}

	booking, err = h.db.GetBooking(bookingID)
	if err != nil {
		return entity.Booking{}, errors.Wrap(err, "get booking")
	}

	return
}
