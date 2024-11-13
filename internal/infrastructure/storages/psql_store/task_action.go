package psql_store

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type ITaskActionStorage interface {
	Create(ctx context.Context, action *models.TaskAction) error
}

type TaskActionPsqlStorage struct {
	db *sqlx.DB
}

func NewTaskActionPsqlStorage(database *sqlx.DB) *TaskActionPsqlStorage {
	return &TaskActionPsqlStorage{
		db: database,
	}
}

func (as *TaskActionPsqlStorage) Create(ctx context.Context, action *models.TaskAction) error {
	log.Println("Create in TaskActionPsqlStorage")

	query := `
		INSERT INTO task_actions (external_id, task_id, type, url) 
		VALUES (:external_id, :task_id, :type, :url)
`

	_, err := as.db.NamedExecContext(ctx, query, action)
	if err != nil {
		return fmt.Errorf("create task action in storage err: %w", err)
	}

	return nil
}