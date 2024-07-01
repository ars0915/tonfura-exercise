package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ars0915/tonfura-exercise/config"
)

func NewDB(config config.ConfENV) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.DB.Host,
		config.DB.Username,
		config.DB.Password,
		config.DB.Database,
		config.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, errors.Wrap(err, "init db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "connect db")
	}

	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	err = sqlDB.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping db")
	}

	return db, nil
}
