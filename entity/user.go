package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(64);not null;default:''"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Bookings []Booking `gorm:"foreignKey:UserID" json:"bookings"`
}

func (t User) MarshalJSON() ([]byte, error) {
	type Alias User
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
