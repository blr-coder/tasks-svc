package psql_store

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTaskPsqlStorage_Get(t *testing.T) {
	testDB := NewTaskPsqlStorage(NewTestPsqlDB())

	testCases := []struct {
		name    string
		taskID  int64
		want    models.Task
		wantErr error
	}{
		{
			name:    "ok",
			taskID:  0,
			want:    models.Task{},
			wantErr: nil,
		},
		{
			name:    "not found err",
			taskID:  0,
			want:    models.Task{},
			wantErr: errs.DomainNotFoundError{},
		},
		{
			name:    "some db error",
			taskID:  0,
			want:    models.Task{},
			wantErr: nil, // ??????
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// TODO: cleanupTasks()
			// TODO: addTestsTasks()

			got, err := testDB.Get(context.Background(), testCase.taskID)
			require.NoError(t, err)

			require.Equal(t, got.ID, testCase.taskID)
			// TODO: requireTasksEqual()
		})
	}
}
