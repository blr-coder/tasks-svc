package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/transaction"
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
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

func (tss *TaskServiceTestsSuite) TestCreateWithTransactionOK() {
	ctx := context.Background()

	// Создаём mockTask для проверки возвращаемого значения
	mockTask := models.NewMockTask(tss.T())

	tss.Run("ok with default currency EUR", func() {
		// Currency always EUR
		mockTask.Currency = &models.DefaultCurrency

		// Входные данные для создания задачи
		input := &models.CreateTask{
			Title:       mockTask.Title,
			Description: mockTask.Description,
			CustomerID:  mockTask.CustomerID,
			ExecutorID:  mockTask.ExecutorID,
			Currency:    mockTask.Currency,
			Amount:      mockTask.Amount,
		}

		// Создаём mock для транзакции
		tx := transaction.NewMockITransaction(tss.mockController)
		sqlTx := &sqlx.Tx{} // Создаем пустой экземпляр *sqlx.Tx для теста

		gomock.InOrder(
			// Ожидаем старт транзакции
			tss.mockIDBTransactionManager.EXPECT().StartTx(ctx).Return(tx, nil),

			// Ожидаем вызов GetTx и передаём результат в WithTransaction
			tx.EXPECT().GetTx().Return(sqlTx),

			// Связываем транзакцию с mock хранилищем задач
			tss.mockITaskStorage.EXPECT().WithTransaction(sqlTx).Return(tss.mockITaskStorage),

			// Создаём задачу в хранилище
			tss.mockITaskStorage.EXPECT().Create(ctx, input).Return(mockTask, nil),

			// Отправляем событие
			tss.mockIEventSender.EXPECT().SendTaskCreated(ctx, mockTask).Return(nil),

			// Завершаем транзакцию
			tx.EXPECT().Finish(ctx).Return(nil),
		)

		// Запускаем тестируемый метод и проверяем результаты
		res, err := tss.taskService.CreateWithTransaction(ctx, input)
		tss.Require().NoError(err)
		tss.Require().Equal(mockTask.ID, res)
	})

	tss.Run("ok with different currency USD", func() {
		// Currency always USD
		mockTask.Currency = utils.Pointer(models.CurrencyUSD)

		// Входные данные для создания задачи
		input := &models.CreateTask{
			Title:       mockTask.Title,
			Description: mockTask.Description,
			CustomerID:  mockTask.CustomerID,
			ExecutorID:  mockTask.ExecutorID,
			Currency:    mockTask.Currency,
			Amount:      mockTask.Amount,
		}

		// Создаём mock для транзакции
		tx := transaction.NewMockITransaction(tss.mockController)
		sqlTx := &sqlx.Tx{} // Создаем пустой экземпляр *sqlx.Tx для теста

		gomock.InOrder(
			tss.mockIDBTransactionManager.EXPECT().StartTx(ctx).Return(tx, nil),
			tx.EXPECT().GetTx().Return(sqlTx),
			tss.mockITaskStorage.EXPECT().WithTransaction(sqlTx).Return(tss.mockITaskStorage),
			tss.mockITaskStorage.EXPECT().Create(ctx, input).Return(mockTask, nil),
			// Ожидаем пересчет валюты
			tss.mockICurrencyStorage.EXPECT().GetRateByEUR(ctx, models.CurrencyUSD).Return(&models.CurrencyRate{}, nil),
			// Проверка отправки события
			tss.mockIEventSender.EXPECT().SendTaskCreated(ctx, mockTask).Return(nil),
			tx.EXPECT().Finish(ctx).Return(nil),
		)

		// Запуск и проверка
		res, err := tss.taskService.CreateWithTransaction(ctx, input)
		tss.Require().NoError(err)
		tss.Require().Equal(mockTask.ID, res)
	})
}
