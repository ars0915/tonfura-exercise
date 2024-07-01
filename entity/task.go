package entity

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Name      *string        `json:"name" gorm:"type:varchar(255);not null;default:''"`
	Status    *uint8         `json:"status" gorm:"not null;default:0"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task
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

type ListTaskParam struct {
	Offset *int
	Limit  *int
}
