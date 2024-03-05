package psql_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"

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

func (s *TaskPsqlStorage) Create(ctx context.Context, createTask *models.CreateTask) (int64, error) {

	return 999, nil
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
