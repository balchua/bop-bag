package controller

import (
	"context"
	"fmt"
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

func (q *TaskController) UpdateTask(c *fiber.Ctx) error {

	ctx := context.TODO()
	task := new(domain.Task)
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err := c.BodyParser(task); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "marshalling error!")
	}
	task.Id = id
	newTask, err := q.taskService.UpdateTask(ctx, task)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(newTask)
}

func (q *TaskController) DeleteTask(c *fiber.Ctx) error {
	ctx := context.TODO()
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}
	if err := q.taskService.DeleteTask(ctx, id); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(fmt.Sprintf("task %d is deleted", id))
}
