package usecase

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/ars0915/gogolook-exercise/entity"
	mocks "github.com/ars0915/gogolook-exercise/mocks/repo"
	"github.com/ars0915/gogolook-exercise/util/ctest"
)

type taskTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	h  *TaskHandler
	db *mocks.MockApp
}

func Test_taskTestSuite(t *testing.T) {
	suite.Run(t, &taskTestSuite{})
}

func (s *taskTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.db = mocks.NewMockApp(s.ctrl)
	s.h = NewTaskHandler(s.db)
}

func (s *taskTestSuite) TearDownTest(t *testing.T) {
	defer s.ctrl.Finish()
}

func (s *taskTestSuite) Test_Create() {
	task := entity.Task{
		Name:   pointer.ToString("task1"),
		Status: pointer.ToUint8(0),
	}

	s.db.EXPECT().CreateTask(ctest.DiffWrapper(task)).Return(task, nil)

	actualTask, err := s.h.CreateTask(context.Background(), task)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), task, actualTask)
}

func (s *taskTestSuite) Test_List() {
	tasks := []entity.Task{
		{
			Name:   pointer.ToString("task1"),
			Status: pointer.ToUint8(0),
		},
		{
			Name:   pointer.ToString("task2"),
			Status: pointer.ToUint8(1),
		},
	}

	s.db.EXPECT().ListTasks(gomock.Any()).Return(tasks, nil)
	s.db.EXPECT().GetTasksCount().Return(int64(len(tasks)), nil)

	actualTasks, actualCount, err := s.h.ListTasks(context.Background(), entity.ListTaskParam{})
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), tasks, actualTasks)
	assert.Equal(s.T(), int64(len(tasks)), actualCount)
}

func (s *taskTestSuite) Test_Get_Success() {
	task := entity.Task{
		ID:     uint(1),
		Name:   pointer.ToString("task1"),
		Status: pointer.ToUint8(0),
	}

	s.db.EXPECT().GetTask(uint(1)).Return(task, nil)

	actualTask, err := s.h.GetTask(context.Background(), uint(1))
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), task, actualTask)
}

func (s *taskTestSuite) Test_Get_NotFound() {
	s.db.EXPECT().GetTask(uint(1)).Return(entity.Task{}, gorm.ErrRecordNotFound)

	_, err := s.h.GetTask(context.Background(), uint(1))
	assert.Equal(s.T(), ErrorTaskNotFound, err)
}

func (s *taskTestSuite) Test_Update() {
	task := entity.Task{
		ID:     uint(1),
		Name:   pointer.ToString("task_rename"),
		Status: pointer.ToUint8(1),
	}

	s.db.EXPECT().UpdateTask(uint(1), ctest.DiffWrapper(task)).Return(nil)
	s.db.EXPECT().GetTask(uint(1)).Return(task, nil)

	actualTask, err := s.h.UpdateTask(context.Background(), uint(1), task)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), task, actualTask)
}

func (s *taskTestSuite) Test_Delete() {
	s.db.EXPECT().DeleteTask(uint(1)).Return(nil)

	err := s.h.DeleteTask(context.Background(), uint(1))
	assert.Nil(s.T(), err)
}
