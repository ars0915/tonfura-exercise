package repo

import (
	"context"

	"github.com/ars0915/gogolook-exercise/entity"
)

//go:generate mockgen -destination=../mocks/repo/app_repo.go -package=mocks github.com/ars0915/gogolook-exercise/repo App

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
	}
)

type txKey struct{}

// injectTx injects transaction to context
func InjectTx(ctx context.Context, tx App) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func ExtractTx(ctx context.Context) App {
	if tx, ok := ctx.Value(txKey{}).(App); ok {
		return tx
	}
	return nil
}
