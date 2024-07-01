package db

import (
	"github.com/ars0915/tonfura-exercise/entity"
)

func (s *AppRepo) UpdateClass(id uint, class entity.Class) (err error) {
	if err = s.db.Where("id = ?", id).Updates(class).Error; err != nil {
		return err
	}
	return nil
}
