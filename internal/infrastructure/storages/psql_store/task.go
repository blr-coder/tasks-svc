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
	return &models.Task{
		ID:          20000,
		Title:       "Fix2 errors handling",
		Description: "Fix2 errors handling description",
		CustomerID:  uuid.UUID([]byte(`13bb16c2-9d81-4697-bf43-430142f38ab5`)),
		ExecutorID:  uuid.UUID([]byte(`13bb16c2-9d81-4697-bf43-430142f38ab5`)),
		Status:      models.PendingStatus,
	}, nil
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

func (s *TaskPsqlStorage) Delete(ctx context.Context, taskID int64) error {
	log.Println("Delete in Storage", taskID)
	return nil
}
