package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
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
	eventSender events.IEventSender
}

func NewTaskService(taskStorage storages.ITaskStorage, eventSender events.IEventSender) *TaskService {
	return &TaskService{
		taskStorage: taskStorage,
		eventSender: eventSender,
	}
}

func (ts *TaskService) sendEvent(ctx context.Context, task *models.Task, topic events.Topic, name events.Name) {
	jTask, err := task.ToJson()
	if err != nil {
		// TODO: Handle errors
		log.Println("sendEvent ERR", err)
		return
	}

	// TODO: Think about worker pool in TaskService for sending events
	err = ts.eventSender.Send(ctx, &events.Event{
		Topic: topic,
		Name:  name,
		Data:  jTask,
	})
	if err != nil {
		// TODO: Log something
		log.Println("sendEvent ERR", err)
	}

	return
}

func (ts *TaskService) Create(ctx context.Context, input *models.CreateTask) (int64, error) {
	log.Println("create in TaskService")
	task, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	go ts.sendEvent(ctx, task, events.CUDTopic, events.TaskCreated)

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
	log.Println("delete in TaskService")
	err := ts.taskStorage.Delete(ctx, taskID)
	if err != nil {
		// TODO: Handle errors
		return err
	}

	go ts.sendEvent(ctx, &models.Task{ID: taskID}, events.CUDTopic, events.TaskDeleted)

	return nil
}
