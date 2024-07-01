package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	ID          uint           `json:"id" gorm:"primaryKey" `
	Source      string         `json:"source" gorm:"type:varchar(64);not null;default:''"`
	Destination string         `json:"destination" gorm:"type:varchar(64);not null;default:''"`
	DepartureAt time.Time      `json:"departure_at"`
	Status      string         `json:"status" gorm:"type:varchar(32);not null;default:'available'"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`

	Classes []Class `json:"classes"`
}

func (t Flight) MarshalJSON() ([]byte, error) {
	type Alias Flight
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

type ListFlightParam struct {
	Source        *string
	Destination   *string
	DepartureDate *time.Time
	SortBy        *string
	Offset        *int
	Limit         *int

	WithClass   bool
	WithBooking bool
}
