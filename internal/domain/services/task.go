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
	task, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	jTask, err := task.ToJson()
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	// TODO: Send events like "new task created", topic - "?", partition - "?"
	err = ts.eventSender.Send(ctx, &queues.Event{ // Move to ts.SendEvent(task)
		Topic: "async_arc_topic",
		Name:  "task.created",
		Data:  jTask,
	})
	if err != nil {
		// TODO: Log something
	}

	return task.ID, nil
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
