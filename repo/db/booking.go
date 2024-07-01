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

func (s *AppRepo) GetBooking(bookingID uint) (booking entity.Booking, err error) {
	err = s.db.Where("id = ?", bookingID).Preload("Flight").Preload("Class").First(&booking).Error
	return
}

func (s *AppRepo) UpdateBooking(bookingID uint, booking entity.Booking) error {
	err := s.db.Model(&booking).Where("id = ?", bookingID).Updates(booking).Error
	return err
}
