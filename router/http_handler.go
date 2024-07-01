package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type appRouter struct {
	method   string
	endpoint string
	worker   gin.HandlerFunc
}

func (h HttpHandler) getRouter() (routes []appRouter) {
	return []appRouter{
		{http.MethodGet, "/flight/", h.listFlightsHandler},

		{http.MethodPost, "/booking/", h.createBookingHandler},
		{http.MethodPatch, "/booking/:bookingID/", h.updateBookingHandler},
		{http.MethodPost, "/check-in/", h.checkInHandler},
		{http.MethodPost, "/give-up/", h.giveUpHandler},
	}
}
