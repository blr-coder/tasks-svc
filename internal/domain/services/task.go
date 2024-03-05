package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages"
	"log"
)

type ITaskService interface {
	Create(ctx context.Context, input *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)

	Delete(ctx context.Context, taskID int64) error
}

type TaskService struct {
	taskStorage storages.ITaskStorage
	eventSender queues.IQueueEventSender
}

func NewTaskService(taskStorage storages.ITaskStorage, eventSender queues.IQueueEventSender) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
		eventSender: eventSender,
	}
}

func (ts *TaskService) Create(ctx context.Context, input *models.CreateTask) (int64, error) {
	log.Println("create in TaskService")
	id, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	// TODO: Send events like "new task created", topic - "?", partition - "?"
	err = ts.eventSender.Send(ctx, []byte("testK"), []byte("testV"))
	if err != nil {
		// TODO: Log something
	}

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
