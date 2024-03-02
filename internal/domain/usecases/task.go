package usecases

import (
	"context"
	"tasks-svc/internal/domain/models"
)

type ITaskUseCase interface {
	Create(ctx context.Context, create *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)

	Delete(ctx context.Context, taskID int64) (int64, error)
}
