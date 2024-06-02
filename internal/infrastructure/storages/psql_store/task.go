package psql_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/google/uuid"
	"log"

	// DB driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
		    (title, description, customer_id, executor_id, status)
		VALUES 
		    ($1, $2, $3, $4, $5)
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
	)
	if err != nil {
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
			e := errs.NewDomainNotFoundError()
			return nil, fmt.Errorf("data not found: %w", e.WithParam("task_id", fmt.Sprint(taskID)))
		}
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
	eID := uuid.UUID([]byte(`13bb16c2-9d81-4697-bf43-430142f38ab5`))
	return &models.Task{
		ID:          20000,
		Title:       "Fix2 errors handling UPDATED",
		Description: "Fix2 errors handling description UPDATED",
		CustomerID:  uuid.UUID([]byte(`13bb16c2-9d81-4697-bf43-430142f38ab5`)),
		ExecutorID:  &eID,
		Status:      models.PendingStatus,
	}, nil
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
