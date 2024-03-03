package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/storage"
)

type ITaskService interface {
	Create(ctx context.Context, input *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)

	Delete(ctx context.Context, taskID int64) error
}

type TaskService struct {
	taskStorage storage.ITaskStorage
}

func NewTaskService(taskStorage storage.ITaskStorage) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
	}
}

func (ts *TaskService) Create(ctx context.Context, input *models.CreateTask) (int64, error) {
	id, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	// TODO: Send event like "new task created", topic - "?", partition - "?"

	return id, nil
}

func (ts *TaskService) Get(ctx context.Context, taskID int64) (*models.Task, error) {
	task, err := ts.taskStorage.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) Delete(ctx context.Context, taskID int64) error {

	return nil
}
