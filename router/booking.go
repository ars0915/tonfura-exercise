package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ars0915/tonfura-exercise/constant"
	"github.com/ars0915/tonfura-exercise/usecase"
	"github.com/ars0915/tonfura-exercise/util/cGin"
)

type createBookingBody struct {
	FlightID uint `json:"flight_id" binding:"required"`
	UserID   uint `json:"user_id" binding:"required"`
	ClassID  uint `json:"class_id" binding:"required"`
	Price    uint `json:"price" binding:"required"`
	Amount   uint `json:"amount" binding:"required"`
}

func (rH *HttpHandler) createBookingHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)
	body := createBookingBody{}
	if err := c.BindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "parse error")
		return
	}

	param := usecase.CreateBookingParam{
		FlightID: body.FlightID,
		UserID:   body.UserID,
		ClassID:  body.ClassID,
		Price:    body.Price,
		Amount:   body.Amount,
	}

	data, err := rH.h.CreateBooking(ctx, param)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Create Booking Failed")
		return
	}

	ctx.WithData(data).Response(http.StatusOK, "")
}

type checkInBody struct {
	BookingID uint `json:"booking_id" binding:"required"`
}

type checkInResponse struct {
	Status     string             `json:"status"`
	Suggestion *checkInSuggestion `json:"suggestion"`
}

type checkInSuggestion struct {
	FlightID uint `json:"flight_id"`
	ClassID  uint `json:"class_id"`
}

func (rH *HttpHandler) checkInHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)
	body := checkInBody{}
	if err := c.BindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "parse error")
		return
	}

	result, err := rH.h.CheckInBooking(ctx, body.BookingID)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Check-in Failed")
		return
	}

	response := checkInResponse{
		Status: result.Status,
	}
	if result.Status == constant.BookingCheckInFail {
		response.Suggestion = &checkInSuggestion{
			FlightID: result.SuggestFlightID,
			ClassID:  result.SuggestClassID,
		}
	}

	ctx.WithData(response).Response(http.StatusOK, "")
}

type giveUpBody struct {
	BookingID uint `json:"booking_id" binding:"required"`
}

func (rH *HttpHandler) giveUpHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)
	body := giveUpBody{}
	if err := c.BindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "parse error")
		return
	}

	result, err := rH.h.GiveUpBooking(ctx, body.BookingID)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Give up Failed")
		return
	}

	ctx.WithData(result).Response(http.StatusOK, "")
}
