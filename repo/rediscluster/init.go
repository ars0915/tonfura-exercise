package redis

import "gorm.io/gorm"

// AppRepo is interface structure
type AppRepo struct {
	db *gorm.DB
}

// New func implements the storage interface for app
func New(db *gorm.DB) *AppRepo {
	return &AppRepo{
		db: db,
	}
}
