package controller

import (
	"strconv"

	"github.com/balchua/uncapsizable/pkg/repository"
	fiber "github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskRepo *repository.TaskRepository
}

func NewQueryController(taskRepo *repository.TaskRepository) *TaskController {
	return &TaskController{
		taskRepo: taskRepo,
	}
}

func (q *TaskController) NewTask(c *fiber.Ctx) error {

	task := new(repository.Task)
	if err := c.BodyParser(task); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "marshalling error!")
	}
	newTask, err := q.taskRepo.Add(task)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(newTask)
}

func (q *TaskController) FindById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}
	task, queryError := q.taskRepo.FindById(id)
	if queryError != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, queryError.Error())
	}

	return c.JSON(task)

}

func (q *TaskController) FindAll(c *fiber.Ctx) error {
	tasks, queryError := q.taskRepo.FindAll()
	if queryError != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, queryError.Error())
	}

	return c.JSON(tasks)

}
