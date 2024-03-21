package services

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages"
	"github.com/google/uuid"
	"log"
)

type ITaskService interface {
	Create(ctx context.Context, input *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)
	List(ctx context.Context, filter *models.ListTasksFilter) ([]*models.Task, error)
	Count(ctx context.Context, filter *models.ListTasksFilter) (uint64, error)
	Update(ctx context.Context, input *models.UpdateTask) error
	Delete(ctx context.Context, taskID int64) error

	AssignExecutor(ctx context.Context, taskID int64, executorID uuid.UUID) error
	AssignRandomExecutor(ctx context.Context, taskId int64) error
	SetStatus(ctx context.Context, taskID int64, status models.Status) error
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

func (ts *TaskService) Create(ctx context.Context, input *models.CreateTask) (int64, error) {
	log.Println("create in TaskService")
	task, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		// TODO: Handle errors
		return 0, err
	}

	go ts.eventSender.TaskCreated(ctx, task)

	return task.ID, nil
}

func (ts *TaskService) Get(ctx context.Context, taskID int64) (*models.Task, error) {
	task, err := ts.taskStorage.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) List(ctx context.Context, filter *models.ListTasksFilter) ([]*models.Task, error) {

	return nil, nil
}

func (ts *TaskService) Count(ctx context.Context, filter *models.ListTasksFilter) (uint64, error) {

	return 0, nil
}

func (ts *TaskService) Update(ctx context.Context, input *models.UpdateTask) error {

	return nil
}

func (ts *TaskService) Delete(ctx context.Context, taskID int64) error {
	log.Println("delete in TaskService")
	err := ts.taskStorage.Delete(ctx, taskID)
	if err != nil {
		// TODO: Handle errors
		return err
	}

	go ts.eventSender.TaskDeleted(ctx, &models.Task{ID: taskID})

	return nil
}

func (ts *TaskService) AssignExecutor(ctx context.Context, taskID int64, executorID uuid.UUID) error {

	return nil
}

func (ts *TaskService) AssignRandomExecutor(ctx context.Context, taskId int64) error {

	return nil
}

func (ts *TaskService) SetStatus(ctx context.Context, taskID int64, status models.Status) error {

	return nil
}
