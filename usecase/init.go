package usecase

import "github.com/ars0915/gogolook-exercise/repo"

func InitHandler(db repo.App) Handler {
	task := NewTaskHandler(db)

	h := newHandler(
		WithTask(task),
	)

	return h
}
