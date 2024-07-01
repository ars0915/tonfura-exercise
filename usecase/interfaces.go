// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ars0915/gogolook-exercise/entity"
)

type (
	Handler interface {
		Task
	}
)

type (
	Task interface {
		ListTasks(ctx context.Context, param entity.ListTaskParam) (tasks []entity.Task, count int64, err error)
		GetTask(ctx context.Context, id uint) (task entity.Task, err error)
		CreateTask(ctx context.Context, t entity.Task) (entity.Task, error)
		UpdateTask(ctx context.Context, id uint, t entity.Task) (entity.Task, error)
		DeleteTask(ctx context.Context, id uint) error
	}
)
