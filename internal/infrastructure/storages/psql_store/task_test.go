package psql_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testUUID = `13bb16c2-9d81-4697-bf43-430142f38ab5`

func TestTaskPsqlStorage_Get(t *testing.T) {
	testDB := NewTaskPsqlStorage(dbConnTest)

	testCases := []struct {
		name    string
		taskID  int64
		want    *models.Task
		wantErr error
	}{
		{
			name:   "ok",
			taskID: 1,
			want: &models.Task{
				ID:          1,
				Title:       "First test task title",
				Description: "First test task description",
				CustomerID:  uuid.UUID([]byte(testUUID)),
				ExecutorID:  nil,
				Status:      models.PendingStatus,
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			wantErr: nil,
		},
		{
			name:    "not found err",
			taskID:  999,
			want:    nil,
			wantErr: errs.NewDomainNotFoundError().WithParam("task_id", fmt.Sprint(999)),
		},
		/*{
			name:    "some db error ???",
			taskID:  0,
			want:    models.Task{},
			wantErr: nil, // ??????
		},*/
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cleanupTasks(t)
			addTestsTasks(t, testDB)

			got, err := testDB.Get(context.Background(), testCase.taskID)
			if !errors.Is(err, testCase.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			requireTasksEqual(t, got, testCase.want)
		})
	}
}

func requireTasksEqual(t *testing.T, got, want *models.Task) {
	if got != nil && want != nil {
		require.Equal(t, want.ID, got.ID)
		require.Equal(t, want.Title, got.Title)
		require.Equal(t, want.Description, got.Description)
		require.Equal(t, want.CustomerID, got.CustomerID)
		require.Equal(t, want.ExecutorID, got.ExecutorID)
		require.Equal(t, want.Status, got.Status)
		require.NotZero(t, got.CreatedAt)
		require.NotZero(t, got.UpdatedAt)
	}
}

func addTestsTasks(t *testing.T, s *TaskPsqlStorage) {
	_, err := s.Create(context.Background(), &models.CreateTask{
		Title:       "First test task title",
		Description: "First test task description",
		CustomerID:  uuid.UUID([]byte(testUUID)),
		ExecutorID:  nil,
	})
	require.NoError(t, err)

	_, err = s.Create(context.Background(), &models.CreateTask{
		Title:       "Second test task title",
		Description: "Second test task description",
		CustomerID:  uuid.UUID([]byte(testUUID)),
		ExecutorID:  nil,
	})
	require.NoError(t, err)
}

func cleanupTasks(t *testing.T) {
	_, err := dbConnTest.Exec("TRUNCATE TABLE tasks CASCADE")
	require.NoError(t, err)
}

func TestTaskPsqlStorage_Create(t *testing.T) {
	testDB := NewTaskPsqlStorage(dbConnTest)

	testCases := []struct {
		name       string
		createTask *models.CreateTask
		want       *models.Task
		wantErr    error
	}{
		{
			name: "ok",
			createTask: &models.CreateTask{
				Title:       "title 3",
				Description: "Description 3333 test !!!",
				CustomerID:  uuid.UUID([]byte(testUUID)),
				ExecutorID:  nil,
			},
			want: &models.Task{
				ID:          3,
				Title:       "title 3",
				Description: "Description 3333 test !!!",
				CustomerID:  uuid.UUID([]byte(testUUID)),
				ExecutorID:  nil,
				Status:      models.PendingStatus,
			},
			wantErr: nil,
		},
		{
			name: "already exist",
			createTask: &models.CreateTask{
				Title:       "First test task title",
				Description: "It's not important",
				CustomerID:  uuid.UUID([]byte(testUUID)),
				ExecutorID:  nil,
			},
			want: nil,
			wantErr: errs.NewDomainDuplicateError().WithParam("create_task", fmt.Sprint(&models.CreateTask{
				Title:       "First test task title",
				Description: "It's not important",
				CustomerID:  uuid.UUID([]byte(testUUID)),
				ExecutorID:  nil,
			})),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cleanupTasks(t)
			addTestsTasks(t, testDB)

			got, err := testDB.Create(context.Background(), testCase.createTask)
			if !errors.Is(err, testCase.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			requireTasksEqual(t, got, testCase.want)
		})
	}
}
