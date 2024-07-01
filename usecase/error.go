package usecase

import (
	"net/http"

	"github.com/ars0915/tonfura-exercise/util/cGin"
)

var (
	ErrorBookingNotFound = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusNotFound,
		Message:  "Booking not found",
	}

	ErrorFlightSoldOut = cGin.CustomError{
		Code:     1002,
		HTTPCode: http.StatusBadRequest,
		Message:  "Flight sold out",
	}

	ErrorClassNotFound = cGin.CustomError{
		Code:     1003,
		HTTPCode: http.StatusNotFound,
		Message:  "Class not found",
	}

	ErrorNoAvailableSeat = cGin.CustomError{
		Code:     1004,
		HTTPCode: http.StatusBadRequest,
		Message:  "No available seat",
	}

	ErrorClassNotBelongToFlight = cGin.CustomError{
		Code:     1005,
		HTTPCode: http.StatusBadRequest,
		Message:  "Class not belong to flight",
	}
)
