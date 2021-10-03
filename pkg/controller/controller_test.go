package controller

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/balchua/bopbag/pkg/domain"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Add(task *domain.Task) (*domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindById(id int64) (*domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll() (*[]domain.Task, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Task), args.Error(1)
}

func setupApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestInsertNewTask(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	app := setupApp()

	// create an instance of our test object
	mockTaskRepo := new(MockTaskRepository)

	var tasks []domain.Task

	task := domain.Task{
		Id: 1234,
	}
	tasks = append(tasks, task)

	mockTaskRepo.On("FindAll").Return(&tasks, nil)

	controller := NewTaskController(mockTaskRepo, 1, logger)

	app.Get("/api/v1/tasks", controller.FindAll)
	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Get All")

}

func TestMustNotAllowInsertNewTask(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	app := setupApp()
	// create an instance of our test object
	mockTaskRepo := new(MockTaskRepository)

	var tasks []domain.Task

	tasks = make([]domain.Task, 0) //append(tasks, task)

	mockTaskRepo.On("FindAll").Return(&tasks, fmt.Errorf("db error"))

	controller := NewTaskController(mockTaskRepo, 1, logger)

	app.Get("/api/v1/tasks", controller.FindAll)
	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	resp, _ := app.Test(req, 1)
	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get All")
}

func TestMustBeAbleToFindById(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	task := &domain.Task{
		Id: 1234,
	}
	// create an instance of our test object
	mockTaskRepo := new(MockTaskRepository)

	mockTaskRepo.On("FindById", int64(1234)).Return(task, nil)

	controller := NewTaskController(mockTaskRepo, 1, logger)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/1234", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "Get By Id")

}

func TestFindInvalidId(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	task := &domain.Task{
		Id: 1234,
	}
	// create an instance of our test object
	mockTaskRepo := new(MockTaskRepository)

	mockTaskRepo.On("FindById", 1234).Return(task, nil)

	controller := NewTaskController(mockTaskRepo, 1, logger)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/abc", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get By Id")

}

func TestUnableToFindId(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	task := &domain.Task{}
	// create an instance of our test object
	mockTaskRepo := new(MockTaskRepository)

	mockTaskRepo.On("FindById", int64(1234)).Return(task, fmt.Errorf("database error"))

	controller := NewTaskController(mockTaskRepo, 1, logger)
	app := setupApp()
	app.Get("/api/v1/task/:id", controller.FindById)

	req := httptest.NewRequest("GET", "/api/v1/task/1234", nil)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "Get By Id")

}

func TestMustBeAbleToInsertATask(t *testing.T) {
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskRepo := new(MockTaskRepository)
	mockTaskRepo.On("Add", task).Return(task, nil)
	controller := NewTaskController(mockTaskRepo, 1, logger)

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
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskRepo := new(MockTaskRepository)
	mockTaskRepo.On("Add", task).Return(task, fmt.Errorf("json parse error"))
	controller := NewTaskController(mockTaskRepo, 1, logger)

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
	logger, logerr := zap.NewProduction()
	if logerr != nil {
		t.Fatal(logerr)
	}
	// prepare the mock
	task := &domain.Task{
		Id:          0,
		Title:       "My First Task",
		Details:     "Here you go, this is what i should do",
		CreatedDate: "2021-10-25",
	}
	mockTaskRepo := new(MockTaskRepository)
	mockTaskRepo.On("Add", task).Return(task, fmt.Errorf("database error"))
	controller := NewTaskController(mockTaskRepo, 1, logger)

	//set up fiber
	app := setupApp()
	app.Post("/api/v1/task/", controller.NewTask)

	var jsonData = `{ "title": "My First Task", "details": "Here you go, this is what i should do", "createdDate": "2021-10-25"}`
	req := httptest.NewRequest("POST", "/api/v1/task/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	// Verify, if the status code is as expected
	assert.Equalf(t, 503, resp.StatusCode, "New task")
	logger.Info("status mesage", zap.String("msg", resp.Status))

}
