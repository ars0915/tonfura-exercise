package db

import (
	"github.com/ars0915/gogolook-exercise/entity"
)

func (s *AppRepo) ListTasks(param entity.ListTaskParam) (t []entity.Task, err error) {
	db := s.db.Model(entity.Task{})
	if param.Offset != nil && param.Limit != nil {
		db = db.Limit(*param.Limit)
		db = db.Offset(*param.Offset)
	}

	if err = db.Find(&t).Error; err != nil {
		return
	}
	return t, nil
}

func (s *AppRepo) GetTasksCount() (count int64, err error) {
	db := s.db.Model(entity.Task{})
	db = db.Limit(-1)
	db = db.Offset(-1)
	err = db.Count(&count).Error
	return
}

func (s *AppRepo) CreateTask(t entity.Task) (entity.Task, error) {
	if err := s.db.Create(&t).Error; err != nil {
		return t, err
	}

	return t, nil
}

func (s *AppRepo) UpdateTask(id uint, t entity.Task) error {
	return s.db.Model(entity.Task{}).Where(`"id" = ?`, id).Updates(t).Error
}

func (s *AppRepo) GetTask(id uint) (task entity.Task, err error) {
	if err = s.db.Where(`"id" = ?`, id).First(&task).Error; err != nil {
		return task, err
	}
	return
}

func (s *AppRepo) DeleteTask(id uint) (err error) {
	if err = s.db.Where(`"id" = ?`, id).Delete(&entity.Task{}).Error; err != nil {
		return err
	}
	return
}
