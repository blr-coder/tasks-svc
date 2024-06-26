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
			name:    "some db error",
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
			require.Equal(t, got.ID, testCase.taskID)
			// TODO: requireTasksEqual()
		})
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
