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
		{http.MethodGet, "/tasks", h.ListTasksHandler},
		{http.MethodPost, "/tasks", h.CreateTaskHandler},
		{http.MethodPut, "/tasks/:taskID", h.UpdateTaskHandler},
		{http.MethodDelete, "/tasks/:taskID", h.DeleteTaskHandler},
	}
}
