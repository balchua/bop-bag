package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
	"github.com/balchua/bopbag/pkg/applog"
	"github.com/balchua/bopbag/pkg/domain"
	"go.uber.org/zap"
)

type TaskService struct {
	taskRepo domain.TaskRepository
	lg       *applog.Logger
	retries  uint
}

func NewTaskService(repo domain.TaskRepository, retries uint, lg *applog.Logger) *TaskService {
	return &TaskService{
		taskRepo: repo,
		lg:       lg,
		retries:  retries,
	}

}

func (t *TaskService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Millisecond)
	defer cancel()

	var newTask *domain.Task
	//validate the fields as part of the business requirement
	if t.isValidTask(task) {
		seed := time.Now().UnixNano()
		random := rand.New(rand.NewSource(seed))
		action := func(attempt uint) error {
			var addErr error
			newTask, addErr = t.taskRepo.Add(task)
			t.lg.Log.Info("Create task Attempt", zap.Uint("attempt", attempt))
			if addErr != nil {
				t.lg.Log.Info("Unable to add the task", zap.Error(addErr))
			}
			return addErr
		}

		err := retry.Retry(
			action,
			strategy.Limit(t.retries),
			strategy.BackoffWithJitter(
				backoff.BinaryExponential(10*time.Millisecond),
				jitter.Deviation(random, 0.5),
			),
		)
		if err != nil {
			return nil, err
		}
		return newTask, nil
	}
	return nil, fmt.Errorf("invalid task passed")
}

func (t *TaskService) GetTaskById(ctx context.Context, id int64) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Millisecond)
	defer cancel()
	task, err := t.taskRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TaskService) GetAllTasks(ctx context.Context) (*[]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Millisecond)
	defer cancel()
	tasks, err := t.taskRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) isValidTask(task *domain.Task) bool {
	if task.Title == "" || task.Details == "" {
		return false
	}
	return true
}

func (t *TaskService) DeleteTask(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Millisecond)
	defer cancel()

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))
	delete := func(attempt uint) error {
		var deleteErr error
		deleteErr = t.taskRepo.Delete(id)
		t.lg.Log.Info("Deleting task Attempt", zap.Uint("attempt", attempt))
		if deleteErr != nil {
			t.lg.Log.Info("Unable to delete the task", zap.Error(deleteErr))
		}
		return deleteErr
	}

	err := retry.Retry(
		delete,
		strategy.Limit(t.retries),
		strategy.BackoffWithJitter(
			backoff.BinaryExponential(10*time.Millisecond),
			jitter.Deviation(random, 0.5),
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Millisecond)
	defer cancel()

	var updatedTask *domain.Task

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))
	action := func(attempt uint) error {
		var updateErr error
		updatedTask, updateErr = t.taskRepo.Update(task)
		t.lg.Log.Info("Update task Attempt", zap.Uint("attempt", attempt))
		if updateErr != nil {
			t.lg.Log.Info("Unable to update the task", zap.Error(updateErr))
		}
		return updateErr
	}

	err := retry.Retry(
		action,
		strategy.Limit(t.retries),
		strategy.BackoffWithJitter(
			backoff.BinaryExponential(10*time.Millisecond),
			jitter.Deviation(random, 0.5),
		),
	)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil

}
