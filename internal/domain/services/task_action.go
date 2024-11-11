package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
)

type ITaskActionService interface {
	Save(ctx context.Context, action *models.TaskAction) error
}

type TaskActionService struct {
}

func NewTaskActionService() *TaskActionService {

	return &TaskActionService{}
}

func (as *TaskActionService) Save(ctx context.Context, action *models.TaskAction) error {

	// TODO: Implement me:)

	return nil
}
