package router

import (
	"github.com/ars0915/gogolook-exercise/config"
	"github.com/ars0915/gogolook-exercise/usecase"
)

type HttpHandler struct {
	conf config.ConfENV
	h    usecase.Handler
}

func newHttpHandler(conf config.ConfENV, h usecase.Handler) *HttpHandler {
	return &HttpHandler{
		conf: conf,
		h:    h,
	}
}

func (h *HttpHandler) Usecase() usecase.Handler {
	return h.h
}

type Handler struct {
	http *HttpHandler
}

func NewHandler(conf config.ConfENV, h usecase.Handler) Handler {
	return Handler{
		http: newHttpHandler(conf, h),
	}
}
