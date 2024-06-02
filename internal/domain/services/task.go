package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages"
	"github.com/google/uuid"
	"log"
)

type ITaskService interface {
	Create(ctx context.Context, input *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)
	List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error)
	Count(ctx context.Context, filter *models.TasksFilter) (uint64, error)
	Update(ctx context.Context, input *models.Task) error
	Delete(ctx context.Context, taskID int64) error

	AssignExecutor(ctx context.Context, taskID int64, executorID uuid.UUID) error
	AssignRandomExecutor(ctx context.Context, taskId int64) error
	SetStatus(ctx context.Context, taskID int64, status models.TaskStatus) error
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
		return 0, fmt.Errorf("create new task err, %w", err)
	}

	err = ts.eventSender.SendTaskCreated(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("sending event err, %w", err)
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

func (ts *TaskService) List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error) {

	return ts.taskStorage.List(ctx, filter)
}

func (ts *TaskService) Count(ctx context.Context, filter *models.TasksFilter) (uint64, error) {

	return ts.taskStorage.Count(ctx, filter)
}

func (ts *TaskService) Update(ctx context.Context, input *models.Task) error {
	log.Println("update in TaskService")
	task, err := ts.taskStorage.Update(ctx, input)
	if err != nil {
		return err
	}

	err = ts.eventSender.SendTaskUpdated(ctx, task)
	if err != nil {
		return fmt.Errorf("sending event err, %w", err)
	}

	return nil
}

func (ts *TaskService) Delete(ctx context.Context, taskID int64) error {
	log.Println("delete in TaskService")
	err := ts.taskStorage.Delete(ctx, taskID)
	if err != nil {
		// TODO: Handle errors
		return err
	}

	err = ts.eventSender.SendTaskDeleted(ctx, taskID)
	if err != nil {
		return fmt.Errorf("sending event err, %w", err)
	}

	return nil
}

func (ts *TaskService) AssignExecutor(ctx context.Context, taskID int64, executorID uuid.UUID) error {

	return nil
}

func (ts *TaskService) AssignRandomExecutor(ctx context.Context, taskId int64) error {
	// TODO: Go to user/auth svc for getting random executor

	return nil
}

func (ts *TaskService) SetStatus(ctx context.Context, taskID int64, status models.TaskStatus) error {

	return nil
}
