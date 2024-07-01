package db

import (
	"github.com/ars0915/tonfura-exercise/entity"
)

func (s *AppRepo) CreateBooking(booking entity.Booking) (entity.Booking, error) {
	if err := s.db.Create(&booking).Error; err != nil {
		return booking, err
	}
	return booking, nil
}
