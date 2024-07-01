// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ars0915/tonfura-exercise/entity"
)

type (
	Handler interface {
		Flight
		Booking
	}
)

type (
	Flight interface {
		ListFlights(ctx context.Context, param entity.ListFlightParam) (flights []entity.Flight, count int64, err error)
	}

	Booking interface {
		GetBooking(ctx context.Context, id uint) (booking entity.Booking, err error)
		CreateBooking(ctx context.Context, param CreateBookingParam) (booking entity.Booking, err error)
		CheckInBooking(ctx context.Context, bookingID uint) (result CheckInResult, err error)
		GiveUpBooking(ctx context.Context, bookingID uint) (booking entity.Booking, err error)
		UpdateBooking(ctx context.Context, bookingID uint, data entity.Booking) (booking entity.Booking, err error)
	}
)
