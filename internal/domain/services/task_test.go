package services

import (
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
