package db

import (
	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ars0915/gogolook-exercise/config"
)

func NewDB(config config.ConfENV) (*gorm.DB, error) {
	gDB, err := gorm.Open(sqlite.Open(config.SQLite.Database), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "init db")
	}

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, errors.Wrap(err, "connect db")
	}
	sqlDB.SetMaxOpenConns(config.SQLite.MaxConn)

	return gDB, nil
}
