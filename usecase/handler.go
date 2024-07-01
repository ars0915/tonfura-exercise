package usecase

import (
	"github.com/ars0915/tonfura-exercise/repo"
)

type AppHandler struct {
	Flight
	Booking
}

type NewHandlerOption func(*AppHandler)

func newHandler(optFn ...NewHandlerOption) *AppHandler {
	h := &AppHandler{}

	for _, o := range optFn {
		o(h)
	}

	return h
}

type FlightHandler struct {
	db repo.App
}

func NewFlightHandler(db repo.App) *FlightHandler {
	return &FlightHandler{
		db: db,
	}
}

func WithFlight(i *FlightHandler) func(h *AppHandler) {
	return func(h *AppHandler) { h.Flight = i }
}

type BookingHandler struct {
	db    repo.App
	redis repo.Redis
}

func NewBookingHandler(db repo.App, redis repo.Redis) *BookingHandler {
	return &BookingHandler{
		db:    db,
		redis: redis,
	}
}

func WithBooking(i *BookingHandler) func(h *AppHandler) {
	return func(h *AppHandler) { h.Booking = i }
}
