package db

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ars0915/gogolook-exercise/entity"
)

func (s *AppRepo) Migrate() {
	s.db = s.db.Debug()
	logrus.Info("Start Migration")

	m := gormigrate.New(s.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202403031000",
			Migrate: func(tx *gorm.DB) error {
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	m.InitSchema(func(tx *gorm.DB) error {
		logrus.Info("Create Tables...")
		if err := s.db.AutoMigrate(&entity.Task{}); err != nil {
			return err
		}
		return nil
	})

	if err := m.Migrate(); err != nil {
		logrus.WithError(err).Error("Database Migration Failed")
	}

	logrus.Info("Migrate Finished")
	s.db.Logger = s.db.Logger.LogMode(logger.Silent)
}
