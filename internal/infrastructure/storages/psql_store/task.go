package psql_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/lib/pq"
	"log"

	// DB driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ITaskStorage interface {
	Create(ctx context.Context, createTask *models.CreateTask) (*models.Task, error)
	Get(ctx context.Context, taskID int64) (*models.Task, error)
	List(ctx context.Context, filter *models.TasksFilter) ([]*models.Task, error)
	Count(ctx context.Context, filter *models.TasksFilter) (uint64, error)
	Update(ctx context.Context, input *models.Task) (*models.Task, error)
	Delete(ctx context.Context, taskID int64) error

	WithTransaction(tx *sqlx.Tx) *TaskPsqlStorage
}

type IStorageExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type TaskPsqlStorage struct {
	db       *sqlx.DB
	executor IStorageExecutor
}

func NewTaskPsqlStorage(database *sqlx.DB) *TaskPsqlStorage {
	return &TaskPsqlStorage{
		db:       database,
		executor: database,
	}
}

func (s *TaskPsqlStorage) WithTransaction(tx *sqlx.Tx) *TaskPsqlStorage {
	return &TaskPsqlStorage{
		db:       s.db,
		executor: tx,
	}
}

func (s *TaskPsqlStorage) Create(ctx context.Context, createTask *models.CreateTask) (*models.Task, error) {
	log.Println("Create in Storage")

	query := `
		INSERT INTO tasks 
		    (title, description, customer_id, executor_id, status, currency, amount)
		VALUES 
		    ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`

	task := &models.Task{}

	err := s.executor.GetContext(
		ctx,
		task,
		query,
		createTask.Title,
		createTask.Description,
		createTask.CustomerID,
		createTask.ExecutorID,
		models.PendingStatus,
		createTask.Currency,
		createTask.Amount,
	)
	if err != nil {
		// TODO: handlePostgresError
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				// Это ошибка дублирования ключа
				return nil, errs.NewDomainDuplicateError().WithParam("create_task", fmt.Sprint(createTask))
			}
		}

		return nil, fmt.Errorf("create task in storage err: %w", err)
	}

	return task, nil
}

func (s *TaskPsqlStorage) Get(ctx context.Context, taskID int64) (*models.Task, error) {
	query := `
		SELECT
			*
		FROM tasks
		WHERE
			id = $1
	`

	task := &models.Task{}

	err := s.executor.GetContext(ctx, task, query, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewDomainNotFoundError().WithParam("task_id", fmt.Sprint(taskID))
		}

		return nil, err
	}

	return task, nil
}

func (s *TaskPsqlStorage) List(ctx context.Context, filter *models.TasksFilter) (tasks []*models.Task, err error) {
	query, args, err := s.buildQueryFromTasksFilter(filter, true)
	if err != nil {
		return tasks, fmt.Errorf("fetching tasks from DB error: %w", err)
	}

	err = s.executor.SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetching tasks from DB error: %w", err)
	}

	return tasks, nil
}

func (s *TaskPsqlStorage) Count(ctx context.Context, filter *models.TasksFilter) (uint64, error) {
	var count uint64
	query, args, err := s.buildQueryFromTasksFilter(filter, false)
	if err != nil {
		return 0, fmt.Errorf("fetching tasks from DB error: %w", err)
	}

	err = s.executor.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("fetching tasks from DB error: %w", err)
	}

	return count, nil
}

func (s *TaskPsqlStorage) Update(ctx context.Context, input *models.Task) (*models.Task, error) {
	log.Println("Update in Storage")

	var (
		query = `UPDATE tasks SET 
                 	title=$2, 
                 	description=$3,
                 	customer_id=$4,
                 	executor_id=$5,
                 	status=$6,
                 	updated_at=CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
                 	is_active=$7,
                 	amount=$8,
                 	currency=$9
             WHERE id=$1 RETURNING *`

		updatedTask = &models.Task{}
	)

	if err := s.executor.QueryRowxContext(
		ctx,
		query,
		input.ID,
		input.Title,
		input.Description,
		input.CustomerID,
		input.ExecutorID,
		input.Status,
		input.IsActive,
		input.Amount,
		input.Currency,
	).StructScan(updatedTask); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewDomainNotFoundError().WithParam("task_id", fmt.Sprint(input.ID))
		}

		return nil, fmt.Errorf("updating task in DB error: %w", err)
	}

	return updatedTask, nil
}

func (s *TaskPsqlStorage) Delete(ctx context.Context, taskID int64) error {
	log.Println("Delete in Storage", taskID)

	q := `DELETE FROM tasks WHERE id = $1`

	res, err := s.executor.ExecContext(ctx, s.db.Rebind(q), taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		//return errs.NewDomainDeleteError()
		return fmt.Errorf("delete task from DB error: %w", err)
	}

	return nil
}

func (s *TaskPsqlStorage) buildQueryFromTasksFilter(filter *models.TasksFilter, isList bool) (string, []any, error) {
	var (
		query string
		args  []any
		err   error
	)

	if isList {
		query = `SELECT id, title, description, customer_id, executor_id, status, is_active, currency, amount, created_at, updated_at FROM tasks`
	} else {
		query = `SELECT count(*) FROM tasks`
	}

	query = fmt.Sprintf("%s WHERE TRUE", query)

	if len(filter.Statuses) > 0 {
		query, args, err = sqlx.In(fmt.Sprintf("%s AND status IN (?)", query), filter.Statuses)
		if err != nil {
			return "", nil, err
		}
	}

	if filter.Currency != nil {
		query = fmt.Sprintf("%s AND currency = ?", query)

		args = append(args, filter.Currency.String())
	}

	if isList {
		// TODO: Add norm sorting and paging
		var (
			sortBy    = "id"
			sortOrder = "ASC"
		)

		query = fmt.Sprintf("%s ORDER BY %s %s", query, sortBy, sortOrder)
	}

	query = s.db.Rebind(query)

	return query, args, err
}
