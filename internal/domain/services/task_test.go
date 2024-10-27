package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/transaction"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TaskServiceTestsSuite struct {
	suite.Suite

	mockController            *gomock.Controller
	mockITaskStorage          *psql_store.MockITaskStorage
	mockICurrencyStorage      *psql_store.MockICurrencyStorage
	mockIDBTransactionManager *transaction.MockIDBTransactionManager
	mockIEventSender          *events.MockIEventSender
	taskService               *TaskService
}

func (tss *TaskServiceTestsSuite) SetupTest() {
	tss.mockController = gomock.NewController(tss.T())
	tss.mockITaskStorage = psql_store.NewMockITaskStorage(tss.mockController)
	tss.mockICurrencyStorage = psql_store.NewMockICurrencyStorage(tss.mockController)
	tss.mockIDBTransactionManager = transaction.NewMockIDBTransactionManager(tss.mockController)
	tss.mockIEventSender = events.NewMockIEventSender(tss.mockController)
	tss.taskService = NewTaskService(tss.mockITaskStorage, tss.mockICurrencyStorage, tss.mockIEventSender, tss.mockIDBTransactionManager)
}

func TestTaskServiceTests(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TaskServiceTestsSuite))
}

func (tss *TaskServiceTestsSuite) TearDownTest() {
	tss.mockController.Finish()
}

func (tss *TaskServiceTestsSuite) TestGetTaskOK() {
	ctx := context.Background()

	mockTask := models.NewMockTask(tss.T())

	tss.Run("ok", func() {
		gomock.InOrder(
			tss.mockITaskStorage.EXPECT().Get(ctx, mockTask.ID).
				Return(mockTask, nil),
		)

		res, err := tss.taskService.Get(ctx, mockTask.ID)
		tss.Require().NoError(err)
		tss.Require().Equal(mockTask, res)
	})
}

func (tss *TaskServiceTestsSuite) TestGetTaskERROR() {
	ctx := context.Background()

	mockTask := models.NewMockTask(tss.T())

	domainErr := errs.NewDomainNotFoundError().WithParam("task_id", fmt.Sprint(mockTask.ID))

	tss.Run("err", func() {
		gomock.InOrder(
			tss.mockITaskStorage.EXPECT().Get(ctx, mockTask.ID).
				Return(nil, domainErr),
		)

		_, err := tss.taskService.Get(ctx, mockTask.ID)
		tss.Require().EqualError(err, domainErr.Error())
	})
}
