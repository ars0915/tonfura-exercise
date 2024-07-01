package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	FlightID       uint           `json:"flight_id"`
	Type           *string        `json:"type" gorm:"type:varchar(32);not null;default:'economy'"`
	SeatAmount     *uint          `json:"seat_amount" gorm:"not null;default:0"`
	OversellAmount *uint          `json:"oversell_amount" gorm:"not null;default:0"`
	Price          *uint          `json:"price" gorm:"not null;default:0"`
	Status         *string        `json:"status" gorm:"type:varchar(32);not null;default:'available'"`
	Sold           *uint          `json:"sold" gorm:"not null;default:0"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`

	Bookings []Booking `gorm:"foreignKey:ClassID" json:"bookings"`
}

func (t Class) MarshalJSON() ([]byte, error) {
	type Alias Class
	s := struct {
		*Alias
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		DeletedAt *int64 `json:"deleted_at,omitempty"`
	}{
		Alias:     (*Alias)(&t),
		CreatedAt: t.CreatedAt.Unix(),
		UpdatedAt: t.UpdatedAt.Unix(),
	}

	if t.DeletedAt.Valid {
		deletedAt := t.DeletedAt.Time.Unix()
		s.DeletedAt = &deletedAt
	}

	return json.Marshal(&s)
}
