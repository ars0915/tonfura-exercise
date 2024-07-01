package repo

import (
	"github.com/ars0915/tonfura-exercise/entity"
)

//go:generate mockgen -destination=../mocks/repo/app_repo.go -package=mocks github.com/ars0915/tonfura-exercise/repo App

type (
	App interface {
		Migrate()
		Debug()

		// transaction
		Begin() App
		Commit() error
		Rollback() error

		ListTasks(param entity.ListTaskParam) (t []entity.Task, err error)
		GetTasksCount() (count int64, err error)
		GetTask(id uint) (task entity.Task, err error)
		CreateTask(t entity.Task) (entity.Task, error)
		UpdateTask(id uint, t entity.Task) error
		DeleteTask(id uint) (err error)

		ListFlights(param entity.ListFlightParam) (f []entity.Flight, err error)
		GetFlightsCount(param entity.ListFlightParam) (count int64, err error)
		GetFlight(id uint) (flight entity.Flight, err error)
	}
)
