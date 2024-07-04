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
}

type TaskPsqlStorage struct {
	db *sqlx.DB
}

func NewTaskPsqlStorage(database *sqlx.DB) *TaskPsqlStorage {
	return &TaskPsqlStorage{db: database}
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

	err := s.db.GetContext(
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

	err := s.db.GetContext(ctx, task, query, taskID)
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

	err = s.db.SelectContext(ctx, &tasks, query, args...)
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

	err = s.db.GetContext(ctx, &count, query, args...)
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

	if err := s.db.QueryRowxContext(
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
		return nil, fmt.Errorf("updating task in DB error: %w", err)
	}

	return updatedTask, nil
}

func (s *TaskPsqlStorage) Delete(ctx context.Context, taskID int64) error {
	log.Println("Delete in Storage", taskID)

	q := `DELETE FROM tasks WHERE id = $1`

	res, err := s.db.ExecContext(ctx, s.db.Rebind(q), taskID)
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
		query = `SELECT id, title, description, customer_id, executor_id, status, created_at, updated_at FROM tasks`
	} else {
		query = `SELECT count(*) FROM tasks`
	}

	query = s.db.Rebind(query)

	return query, args, err
}
