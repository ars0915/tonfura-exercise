package usecase

import "github.com/ars0915/tonfura-exercise/repo"

func InitHandler(db repo.App) Handler {
	task := NewTaskHandler(db)

	h := newHandler(
		WithTask(task),
	)

	return h
}
