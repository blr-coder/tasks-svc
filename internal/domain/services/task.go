package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/transaction"
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"log"
)

type ITaskService interface {
	Create(ctx context.Context, input *models.CreateTask) (int64, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)
	List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error)
	Count(ctx context.Context, filter *models.TasksFilter) (uint64, error)
	Update(ctx context.Context, input *models.UpdateTask) (*models.Task, error)
	Delete(ctx context.Context, taskID int64) error

	AssignRandomExecutor(ctx context.Context, taskId int64) error

	// CreateWithTransaction only for testing transactionManager
	CreateWithTransaction(ctx context.Context, input *models.CreateTask) (int64, error)
}

type TaskService struct {
	taskStorage        psql_store.ITaskStorage
	eventSender        events.IEventSender
	transactionManager transaction.DBTransactionManager
	currencyStorage    psql_store.ICurrencyStorage
}

func NewTaskService(taskStorage psql_store.ITaskStorage, currencyStorage psql_store.ICurrencyStorage, eventSender events.IEventSender, transactionManager transaction.DBTransactionManager) *TaskService {
	return &TaskService{
		taskStorage:        taskStorage,
		eventSender:        eventSender,
		transactionManager: transactionManager,
		currencyStorage:    currencyStorage,
	}
}

func (ts *TaskService) recalculateCurrencyToEUR(ctx context.Context, currency *models.Currency, amount *float64) (*float64, error) {
	var res float64

	rate, err := ts.currencyStorage.GetRate(ctx, currency.String())
	if err != nil {
		// TODO: Normal domain err
		return nil, err
	}

	// TODO: Add round
	res = *amount / rate

	return &res, nil
}

func (ts *TaskService) CreateWithTransaction(ctx context.Context, input *models.CreateTask) (int64, error) {
	log.Println("CreateWithTransaction START")
	tx, err := ts.transactionManager.StartTx(ctx)
	if err != nil {
		// TODO: Normal err handling
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	// TODO: It's OK?
	ts.taskStorage = ts.taskStorage.WithTransaction(tx.GetTx())

	task, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return 0, fmt.Errorf("failed to rollback transaction: %v, error: %w", rollbackErr, err)
		}

		return 0, err
	}

	// Recalculate amount before sending event(for kafka use only EUR)
	if input.Currency != nil && input.Currency != &models.DefaultCurrency && input.Amount != nil {
		task.Amount, err = ts.recalculateCurrencyToEUR(ctx, task.Currency, task.Amount)
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return 0, fmt.Errorf("failed to rollback transaction: %v, error: %w", rollbackErr, err)
			}

			return 0, err
		}
	}

	task.Currency = &models.DefaultCurrency

	err = ts.eventSender.SendTaskCreated(ctx, task)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return 0, fmt.Errorf("failed to rollback transaction: %v, error: %w", rollbackErr, err)
		}

		return 0, err
	}

	if err = tx.Finish(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return task.ID, nil
}

func (ts *TaskService) Create(ctx context.Context, input *models.CreateTask) (int64, error) {
	log.Println("create in TaskService")
	task, err := ts.taskStorage.Create(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("create new task err, %w", err)
	}

	// TODO: CheckCurrency() if currency != EUR

	err = ts.eventSender.SendTaskCreated(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("sending event err, %w", err)
	}

	return task.ID, nil
}

func (ts *TaskService) Get(ctx context.Context, taskID int64) (*models.Task, error) {
	task, err := ts.taskStorage.Get(ctx, taskID)
	if err != nil {
		// TODO: Check DomainNotFoundError
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

func (ts *TaskService) Update(ctx context.Context, input *models.UpdateTask) (*models.Task, error) {
	log.Println("update in TaskService")

	task, err := ts.taskStorage.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.CustomerID != nil {
		task.CustomerID = *input.CustomerID
	}

	if input.ExecutorID != nil {
		task.ExecutorID = input.ExecutorID
	}

	// TODO: Check status. Cannot update a task with the status DONE - task.IsUpdatePossible()
	if input.Status != nil {
		task.Status = *input.Status
	}

	if input.Amount != nil {
		task.Amount = input.Amount
	}

	if input.Currency != nil {
		task.Currency = input.Currency
	}

	if input.IsActive != nil {
		task.IsActive = utils.Value(input.IsActive)
	}

	updatedTask, err := ts.taskStorage.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	err = ts.eventSender.SendTaskUpdated(ctx, updatedTask)
	if err != nil {
		return nil, fmt.Errorf("sending event err, %w", err)
	}

	return updatedTask, nil
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

func (ts *TaskService) AssignRandomExecutor(ctx context.Context, taskId int64) error {
	// TODO: Go to user/auth svc for getting random executor

	return nil
}
