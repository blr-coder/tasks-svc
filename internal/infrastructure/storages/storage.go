package storages

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
)

type IStorage interface {
	ITaskStorage
}

type ITaskStorage interface {
	Create(ctx context.Context, createTask *models.CreateTask) (*models.Task, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)
	List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error)
	Count(ctx context.Context, filter *models.TasksFilter) (uint64, error)
	Update(ctx context.Context, input *models.Task) (*models.Task, error)
	Delete(ctx context.Context, taskID int64) error
}
