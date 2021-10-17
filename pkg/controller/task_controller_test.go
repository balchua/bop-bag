package controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/balchua/bopbag/pkg/domain"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockTaskService) GetTaskById(ctx context.Context, id int64) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskService) GetAllTasks(ctx context.Context) (*[]domain.Task, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Task), args.Error(1)
}

func setupApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestMustReturnAllTasks(t *testing.T) {
	app := setupApp()

	// create an instance of our test object
	mockTaskService := new(MockTaskService)

	var tasks []domain.Task

	task := domain.Task{
		Id: 1234,
	}
	tasks = append(tasks, task)

	mockTaskService.On("GetAllTasks",
		mock.Anything).Return(&tasks, nil)

	controller := NewTaskController(mockTaskService)

	app.Get("/api/v1/tasks", controller.FindAll)
	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Get All")

}

func TestMustReturnErrorWhenGetAllTaskFail(t *testing.T) {

	app := setupApp()
	// create an instance of our test object
	mockTaskService := new(MockTaskService)

	var tasks []domain.Task

	tasks = make([]domain.Task, 0) //append(tasks, task)

	mockTaskService.On("GetAllTasks", mock.Anything).Return(&tasks, fmt.Errorf("db error"))

	controller := NewTaskController(mockTaskService)

	app.Get("/api/v1/tasks", controller.FindAll)
	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get All")
}

func TestMustBeAbleToFindById(t *testing.T) {

	task := &domain.Task{
		Id: 1234,
	}
	// create an instance of our test object
	mockTaskService := new(MockTaskService)

	mockTaskService.On("GetTaskById", mock.Anything, int64(1234)).Return(task, nil)

	controller := NewTaskController(mockTaskService)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/1234", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Get By Id")

}

func TestFindInvalidId(t *testing.T) {

	task := &domain.Task{
		Id: 1234,
	}
	// create an instance of our test object
	mockTaskService := new(MockTaskService)

	mockTaskService.On("GetTaskById", mock.Anything, 1234).Return(task, nil)

	controller := NewTaskController(mockTaskService)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/abc", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get By Id")

}

func TestUnableToFindId(t *testing.T) {

	task := &domain.Task{}
	// create an instance of our test object
	mockTaskService := new(MockTaskService)

	mockTaskService.On("GetTaskById", mock.Anything, int64(1234)).Return(task, fmt.Errorf("database error"))

	controller := NewTaskController(mockTaskService)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/1234", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get By Id")

}

func TestMustBeAbleToInsertATask(t *testing.T) {

	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskService := new(MockTaskService)
	mockTaskService.On("CreateTask", mock.Anything, task).Return(task, nil)
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Post("/api/v1/task/", controller.NewTask)

	var jsonData = `{ "title": "My First Task", "details": "Here you go, this is what i should do", "createdDate": "2021-10-25"}`
	req := httptest.NewRequest("POST", "/api/v1/task/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "New task")

}

func TestMustFailWhenTaskJsonIsInvalid(t *testing.T) {
	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskService := new(MockTaskService)
	mockTaskService.On("CreateTask", mock.Anything, task).Return(task, fmt.Errorf("json parse error"))
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Post("/api/v1/task/", controller.NewTask)

	var jsonData = `{ "title": "My First Task" "details": "Here you go, this is what i should do", "createdDate": "2021-10-25"}`
	req := httptest.NewRequest("POST", "/api/v1/task/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "New task")

}

func TestMustFailWhenUnableToInsertATask(t *testing.T) {
	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskService := new(MockTaskService)
	mockTaskService.On("CreateTask", mock.Anything, task).Return(task, fmt.Errorf("database error"))
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Post("/api/v1/task/", controller.NewTask)

	var jsonData = `{ "title": "My First Task", "details": "Here you go, this is what i should do", "createdDate": "2021-10-25"}`
	req := httptest.NewRequest("POST", "/api/v1/task/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "New task")

}

func TestMustDeleteTask(t *testing.T) {
	id := int64(1)
	// prepare the mock
	mockTaskService := new(MockTaskService)
	mockTaskService.On("DeleteTask", id).Return(nil)
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Delete("/api/v1/task/:id", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/api/v1/task/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Delete task")

}

func TestFailDeleteTask(t *testing.T) {
	id := int64(1)
	// prepare the mock
	mockTaskService := new(MockTaskService)
	mockTaskService.On("DeleteTask", id).Return(fmt.Errorf("service unable to delete task"))
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Delete("/api/v1/task/:id", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/api/v1/task/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Delete task")

}

func TestMustBeAbleToUpdateTask(t *testing.T) {

	// prepare the mock
	task := &domain.Task{
		Id:      1,
		Title:   "Ok task",
		Details: "Updating my task",
	}
	mockTaskService := new(MockTaskService)
	mockTaskService.On("UpdateTask", mock.Anything, task).Return(task, nil)
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Put("/api/v1/task/:id", controller.UpdateTask)

	var jsonData = `{ "title": "Ok task", "details": "Updating my task"}`
	req := httptest.NewRequest("PUT", "/api/v1/task/1", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Update task")

}

func TestFailUpdateTask(t *testing.T) {

	task := &domain.Task{
		Id:      1,
		Title:   "Ok task",
		Details: "Updating my task",
	}
	// prepare the mock
	mockTaskService := new(MockTaskService)
	mockTaskService.On("UpdateTask", mock.Anything, task).Return(fmt.Errorf("service unable to delete task"))
	controller := NewTaskController(mockTaskService)

	//set up fiber
	app := setupApp()
	app.Put("/api/v1/task/:id", controller.UpdateTask)

	req := httptest.NewRequest("PUT", "/api/v1/task/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Fail to update task")

}
