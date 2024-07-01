package usecase

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ars0915/tonfura-exercise/constant"
	"github.com/ars0915/tonfura-exercise/entity"
	"github.com/ars0915/tonfura-exercise/repo"
	"github.com/ars0915/tonfura-exercise/util/cTypes"
)

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

		if err = h.redis.Lock(ctx, genLockKey(param.FlightID)); err != nil {
			return err
		}
		defer h.redis.UnLock(ctx, genLockKey(param.FlightID))

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
			FlightID: param.FlightID,
			UserID:   param.UserID,
			ClassID:  param.ClassID,
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

func genLockKey(flightID uint) string {
	return fmt.Sprintf("lock:%d", flightID)
}
