package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FlightID  uint           `json:"flight_id"`
	ClassID   uint           `json:"class_id"`
	UserID    uint           `json:"user_id"`
	Price     *uint          `json:"price" gorm:not null;default:0""`
	Amount    *uint          `json:"amount" gorm:not null;default:0""`
	Status    *string        `json:"status" gorm:"type:varchar(64);not null;default:'to_be_confirmed'"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Flight Flight `gorm:"foreignKey:FlightID" json:"-"`
	Class  Class  `gorm:"foreignKey:ClassID" json:"-"`
}

func (t Booking) MarshalJSON() ([]byte, error) {
	type Alias Booking
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
