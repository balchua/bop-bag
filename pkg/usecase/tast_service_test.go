package usecase

import (
	"context"
	"testing"

	"github.com/balchua/uncapsizable/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
  Test objects
*/

type MockedTaskRepository struct {
	mock.Mock
}

func (m *MockedTaskRepository) Add(task *domain.Task) (*domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockedTaskRepository) FindById(id int64) (*domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockedTaskRepository) FindAll() (*[]domain.Task, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Task), args.Error(1)
}
func TestSuccessfulTaskCreation(t *testing.T) {
	assert := assert.New(t)

	// create an instance of our test object
	mockTaskRepo := new(MockedTaskRepository)
	task := &domain.Task{
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}
	newTask := task
	newTask.Id = 1

	// setup expectations
	mockTaskRepo.On("Add", task).Return(newTask, nil)

	service := NewTaskService(mockTaskRepo)

	response, err := service.CreateTask(context.Background(), task)
	assert.NotNil(response)
	assert.Nil(err)
}

func TestShoudReturnTaskWhenIdIsPassed(t *testing.T) {
	assert := assert.New(t)

	// create an instance of our test object
	mockTaskRepo := new(MockedTaskRepository)
	newTask := &domain.Task{
		Id:          999,
		Title:       "test",
		Details:     "test",
		CreatedDate: "20210926",
	}

	// setup expectations
	mockTaskRepo.On("FindById", int64(999)).Return(newTask, nil)

	service := NewTaskService(mockTaskRepo)

	response, err := service.GetTaskById(context.Background(), 999)
	assert.NotNil(response)
	assert.Nil(err)
}

func TestShoudReturnAllTasks(t *testing.T) {
	assert := assert.New(t)
	// create an instance of our test object
	mockTaskRepo := new(MockedTaskRepository)

	var tasks *[]domain.Task
	_tasks := make([]domain.Task, 2)

	for i := 0; i < 2; i++ {
		newTask := domain.Task{
			Id:          int64(i),
			Title:       "test",
			Details:     "test",
			CreatedDate: "20210926",
		}
		_tasks[i] = newTask
	}
	tasks = &_tasks
	// setup expectations
	mockTaskRepo.On("FindAll").Return(tasks, nil)

	service := NewTaskService(mockTaskRepo)

	response, err := service.GetAllTasks(context.Background())
	assert.Equal(len(*response), 2)
	assert.Nil(err)
}
