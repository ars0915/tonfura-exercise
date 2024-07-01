package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ars0915/tonfura-exercise/entity"
	"github.com/ars0915/tonfura-exercise/util/cGin"
)

type listFlightsQueryString struct {
	Source        string  `form:"source" binding:"required"`
	Destination   string  `form:"destination"  binding:"required"`
	DepartureDate string  `form:"departure_date" binding:"required"`
	SortBy        *string `form:"sort_by"`
}

func (rH *HttpHandler) listFlightsHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)
	queryString := listFlightsQueryString{}
	if err := c.BindQuery(&queryString); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "parse error")
		return
	}

	departureDate, err := time.Parse("2006-01-02", queryString.DepartureDate)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "parse error")
		return
	}

	page := ctx.GetPaginator()
	param := entity.ListFlightParam{
		Source:        &queryString.Source,
		Destination:   &queryString.Destination,
		DepartureDate: &departureDate,
		SortBy:        queryString.SortBy,
		Offset:        &page.Offset,
		Limit:         &page.Limit,
	}

	param.WithClass = true

	data, count, err := rH.h.ListFlights(ctx, param)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "List Flights Failed")
		return
	}
	page.SetTotalCount(int(count))

	ctx.WithPaginator(page).WithData(data).Response(http.StatusOK, "")
}
