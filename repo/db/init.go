package db

import (
	"gorm.io/gorm"

	"github.com/ars0915/gogolook-exercise/repo"
)

// New func implements the storage interface for app
func New(db *gorm.DB) *AppRepo {
	return &AppRepo{
		db: db,
	}
}

// AppRepo is interface structure
type AppRepo struct {
	db *gorm.DB
}

// Begin begin a transaction
func (s *AppRepo) Begin() repo.App {
	return &AppRepo{
		db: s.db.Begin(),
	}
}

// Commit commit a transaction
func (s *AppRepo) Commit() error {
	return s.db.Commit().Error
}

// Rollback rollback a transaction
func (s *AppRepo) Rollback() error {
	return s.db.Rollback().Error
}

// Debug Debug log
func (s *AppRepo) Debug() {
	s.db = s.db.Debug()
}

// type jsonLogger struct{}

// func (*jsonLogger) Print(values ...interface{}) {
// 	if len(values) > 1 {
// 		logrus.WithFields(log.Fields{"type": values[0], "full": values[0:]}).Info("GORM Log")
// 	}
// }
