package controller

import (
	"context"

	"github.com/balchua/bopbag/pkg/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskById(ctx context.Context, id int64) (*domain.Task, error)
	GetAllTasks(ctx context.Context) (*[]domain.Task, error)
}

type ClusterService interface {
	GetClusterInfo() ([]domain.ClusterInfo, error)
}
