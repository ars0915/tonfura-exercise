package usecase

import "github.com/ars0915/tonfura-exercise/repo"

func InitHandler(db repo.App, redis repo.Redis) Handler {
	task := NewTaskHandler(db)
	flight := NewFlightHandler(db)

	h := newHandler(
		WithTask(task),
		WithFlight(flight),
	)

	return h
}
