package router

import (
	"net/http"
	"strconv"

	"github.com/ars0915/tonfura-exercise/util/cGin"
)

func resourceCheck(rH HttpHandler) cGin.HandlerFunc {
	return func(ctx *cGin.Context) {
		bookingIDStr := ctx.Param("bookingID")

		var bookingID int
		var err error
		if len(bookingIDStr) > 0 {
			bookingID, err = strconv.Atoi(bookingIDStr)
			if err != nil {
				ctx.WithError(err).Response(http.StatusBadRequest, "Invalid ID")
				return
			}

			if _, err = rH.Usecase().GetBooking(ctx, uint(bookingID)); err != nil {
				ctx.WithError(err).Response(http.StatusInternalServerError, "Check Task failed")
				return
			}
		}

		ctx.Next()
	}
}
