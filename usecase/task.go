package usecase

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ars0915/gogolook-exercise/entity"
)

func (h TaskHandler) ListTasks(ctx context.Context, param entity.ListTaskParam) (tasks []entity.Task, count int64, err error) {
	if tasks, err = h.db.ListTasks(param); err != nil {
		return tasks, 0, errors.Wrap(err, "list task")
	}

	count, err = h.db.GetTasksCount()
	if err != nil {
		return nil, 0, errors.Wrap(err, "count tasks")
	}

	return
}

func (h TaskHandler) GetTask(ctx context.Context, id uint) (task entity.Task, err error) {
	task, err = h.db.GetTask(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task, ErrorTaskNotFound
		}
		return task, errors.Wrap(err, "get task")
	}
	return
}

func (h TaskHandler) CreateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	return h.db.CreateTask(t)
}

func (h TaskHandler) UpdateTask(ctx context.Context, id uint, t entity.Task) (entity.Task, error) {
	if err := h.db.UpdateTask(id, t); err != nil {
		return entity.Task{}, errors.Wrap(err, "update task")
	}

	task, err := h.db.GetTask(id)
	if err != nil {
		return entity.Task{}, errors.Wrap(err, "get task")
	}

	return task, nil
}

func (h TaskHandler) DeleteTask(ctx context.Context, id uint) error {
	return h.db.DeleteTask(id)
}
