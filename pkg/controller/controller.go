package controller

import (
	"context"
	"strconv"

	"github.com/balchua/bopbag/pkg/domain"
	fiber "github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskService TaskService
}

func NewTaskController(taskService TaskService) *TaskController {

	return &TaskController{
		taskService: taskService,
	}
}

func (q *TaskController) NewTask(c *fiber.Ctx) error {

	ctx := context.TODO()
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
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
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
