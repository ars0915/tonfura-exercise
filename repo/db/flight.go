package db

import (
	"gorm.io/gorm"

	"github.com/ars0915/tonfura-exercise/entity"
)

func (s *AppRepo) ListFlights(param entity.ListFlightParam) (f []entity.Flight, err error) {
	db := s.getFlightDB(param)
	if param.Offset != nil && param.Limit != nil {
		db = db.Limit(*param.Limit)
		db = db.Offset(*param.Offset)
	}

	if param.WithClass {
		db.Preload("Classes")
	}

	if param.WithBooking {
		db.Preload("Bookings")
	}

	if err = db.Find(&f).Error; err != nil {
		return
	}
	return f, nil
}

func (s *AppRepo) GetFlightsCount(param entity.ListFlightParam) (count int64, err error) {
	db := s.getFlightDB(param)
	db = db.Limit(-1)
	db = db.Offset(-1)
	err = db.Count(&count).Error
	return
}

func (s *AppRepo) getFlightDB(param entity.ListFlightParam) *gorm.DB {
	db := s.db.Model(entity.Flight{})
	if param.Source != nil {
		db.Where("source = ?", *param.Source)
	}

	if param.Destination != nil {
		db.Where("destination = ?", *param.Destination)
	}

	if param.DepartureDate != nil {
		db.Where("departure_at >= ? AND departure_at < ?", param.DepartureDate, param.DepartureDate.AddDate(0, 0, 1))
	}

	if param.DepartureAfter != nil {
		db.Where("departure_at > ?", param.DepartureAfter)
	}

	if param.SortBy != nil {
		db.Order(*param.SortBy)
	} else {
		db.Order("id ASC")
	}

	return db
}

func (s *AppRepo) GetFlight(id uint) (flight entity.Flight, err error) {
	if err = s.db.Where(`"id" = ?`, id).Preload("Classes").First(&flight).Error; err != nil {
		return flight, err
	}
	return
}

func (s *AppRepo) UpdateFlight(id uint, flight entity.Flight) (err error) {
	if err = s.db.Where("id = ?", id).Updates(flight).Error; err != nil {
		return err
	}
	return nil
}
