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
)
