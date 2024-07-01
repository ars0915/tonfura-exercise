package usecase

import (
	"net/http"

	"github.com/ars0915/tonfura-exercise/util/cGin"
)

var (
	ErrorTaskNotFound = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusNotFound,
		Message:  "Task not found",
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
)
