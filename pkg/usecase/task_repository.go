package usecase

import "github.com/balchua/bopbag/pkg/domain"

type TaskRepository interface {
	Add(task *domain.Task) (*domain.Task, error)

	FindById(id int64) (*domain.Task, error)

	FindAll() (*[]domain.Task, error)
}
