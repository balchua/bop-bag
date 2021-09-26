package controller

import (
	"context"
	"strconv"

	"github.com/balchua/uncapsizable/pkg/domain"
	"github.com/balchua/uncapsizable/pkg/repository"
	"github.com/balchua/uncapsizable/pkg/usecase"
	fiber "github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskService *usecase.TaskService
}

func NewTaskController(taskRepo *repository.TaskRepository) *TaskController {

	ts := usecase.NewTaskService(taskRepo)
	return &TaskController{
		taskService: ts,
	}
}

func (q *TaskController) NewTask(c *fiber.Ctx) error {

	ctx := context.Background()
	task := new(domain.Task)
	if err := c.BodyParser(task); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "marshalling error!")
	}
	newTask, err := q.taskService.CreateTask(ctx, task)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(newTask)
}

func (q *TaskController) FindById(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}
	task, queryError := q.taskService.GetTaskById(ctx, id)
	if queryError != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, queryError.Error())
	}

	return c.JSON(task)

}

func (q *TaskController) FindAll(c *fiber.Ctx) error {
	ctx := context.Background()
	tasks, queryError := q.taskService.GetAllTasks(ctx)
	if queryError != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, queryError.Error())
	}

	return c.JSON(tasks)

}
