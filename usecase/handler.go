package usecase

import (
	"github.com/ars0915/gogolook-exercise/repo"
)

type AppHandler struct {
	Task
}

type NewHandlerOption func(*AppHandler)

func newHandler(optFn ...NewHandlerOption) *AppHandler {
	h := &AppHandler{}

	for _, o := range optFn {
		o(h)
	}

	return h
}

type TaskHandler struct {
	db repo.App
}

func NewTaskHandler(db repo.App) *TaskHandler {
	return &TaskHandler{
		db: db,
	}
}

func WithTask(i Task) func(h *AppHandler) {
	return func(h *AppHandler) { h.Task = i }
}
