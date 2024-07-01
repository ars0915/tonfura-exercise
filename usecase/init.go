package usecase

import "github.com/ars0915/tonfura-exercise/repo"

func InitHandler(db repo.App, redis repo.Redis) Handler {
	flight := NewFlightHandler(db)
	booking := NewBookingHandler(db, redis)

	h := newHandler(
		WithFlight(flight),
		WithBooking(booking),
	)

	return h
}
