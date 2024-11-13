package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"log"
)

type ITaskActionService interface {
	Save(ctx context.Context, action *models.TaskAction) error
}

type TaskActionService struct {
	actionStorage psql_store.ITaskActionStorage
}

func NewTaskActionService(actionStorage psql_store.ITaskActionStorage) *TaskActionService {
	return &TaskActionService{
		actionStorage: actionStorage,
	}
}

func (as *TaskActionService) Save(ctx context.Context, action *models.TaskAction) error {
	log.Println("Save in TaskActionService")

	err := as.actionStorage.Create(ctx, action)
	if err != nil {
		return err
	}

	return nil
}
