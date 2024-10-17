package handlers

import (
	"fmt"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func PbCreateTaskToDomain(createRequest *taskpbv1.CreateTaskRequest) (*models.CreateTask, error) {
	customerID, err := uuid.Parse(createRequest.GetCustomerId())
	if err != nil {
		return nil, err
	}

	createTask := &models.CreateTask{
		Title:       createRequest.GetTitle(),
		Description: createRequest.GetDescription(),
		CustomerID:  customerID,
	}

	if createRequest.ExecutorId != nil {
		eID, err := uuid.Parse(createRequest.GetExecutorId())
		if err != nil {
			return nil, err
		}

		createTask.ExecutorID = &eID
	}

	if createRequest.GetPrice() != nil {
		currency, err := PBCurrencyToDomainCurrency(createRequest.GetPrice().GetCurrency())
		if err != nil {
			return nil, err
		}

		createTask.Amount = utils.Pointer(createRequest.GetPrice().GetAmount())
		createTask.Currency = utils.Pointer(currency)
	}

	return createTask, nil
}

func PBCurrencyToDomainCurrency(pbCurrency string) (currency models.Currency, err error) {
	// TODO: Add ENUM to PROTO or another solution
	switch pbCurrency {
	case "EUR":
		currency = models.CurrencyEUR
	case "USD":
		currency = models.CurrencyUSD
	case "PLN":
		currency = models.CurrencyPLN
	default:
		return "", fmt.Errorf("unknown currency, %s", pbCurrency)
	}

	return currency, nil
}

func DomainTaskStatusToPB(status models.TaskStatus) (pbStatus taskpbv1.TaskStatus) {
	switch status {
	case models.PendingStatus:
		pbStatus = taskpbv1.TaskStatus_TASK_STATUS_PENDING
	case models.ProcessingStatus:
		pbStatus = taskpbv1.TaskStatus_TASK_STATUS_PROCESSING
	case models.DoneStatus:
		pbStatus = taskpbv1.TaskStatus_TASK_STATUS_DONE
	default:
		pbStatus = taskpbv1.TaskStatus_TASK_STATUS_UNSPECIFIED
	}

	return pbStatus
}

func PbTaskStatusToDomain(pbStatus taskpbv1.TaskStatus) (dStatus *models.TaskStatus) {
	switch pbStatus {
	case taskpbv1.TaskStatus_TASK_STATUS_PENDING:
		dStatus = utils.Pointer(models.PendingStatus)
	case taskpbv1.TaskStatus_TASK_STATUS_PROCESSING:
		dStatus = utils.Pointer(models.ProcessingStatus)
	case taskpbv1.TaskStatus_TASK_STATUS_DONE:
		dStatus = utils.Pointer(models.DoneStatus)
	default:
		// TODO: UNSPECIFIED == nil it's OK???
		dStatus = nil
	}

	return dStatus
}

func domainTasksToPBTasks(domainTasks []*models.Task) []*taskpbv1.Task {
	ls := make([]*taskpbv1.Task, len(domainTasks))
	for i, task := range domainTasks {
		ls[i] = domainTaskToPBTask(task)
	}

	return ls
}

func domainTaskToPBTask(task *models.Task) *taskpbv1.Task {
	t := &taskpbv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CustomerId:  task.CustomerID.String(),
		Status:      DomainTaskStatusToPB(task.Status),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
		IsActive:    task.IsActive,
		Price:       &taskpbv1.Price{},
	}

	if task.ExecutorID != nil {
		t.ExecutorId = task.ExecutorID.String()
	}

	if task.Currency != nil {
		t.Price.Currency = task.Currency.String()
	}

	if task.Amount != nil {
		t.Price.Amount = utils.Value(task.Amount)
	}

	return t
}

func PbListLimitingToDB(pbLimiting *taskpbv1.Limiting) *models.Limiting {
	listLimiting := &models.Limiting{}

	if nil == pbLimiting {
		return listLimiting
	}

	if pbLimiting.GetLimit() != 0 {
		listLimiting.Limit = pbLimiting.GetLimit()
	}

	if pbLimiting.GetOffset() != 0 {
		listLimiting.Offset = pbLimiting.GetOffset()
	}

	return listLimiting
}

func PbTasksSortingToDomain(pbSorting *taskpbv1.TaskSorting) *models.Sorting {
	domainSorting := &models.Sorting{}

	for _, sortBy := range pbSorting.GetSortBy() {
		domainSorting.SortBy = append(domainSorting.SortBy, strings.ToLower(sortBy.String()))
	}

	if pbSorting.GetSortOrder().String() != "" {
		domainSorting.SortOrder = pbSorting.GetSortOrder().String()
	}

	return domainSorting
}

func PbListTasksFilteringToDomain(pbFiltering *taskpbv1.TaskFiltering) *models.TaskFiltering {
	filter := &models.TaskFiltering{}

	if pbFiltering != nil {
		domainStatuses := make([]models.TaskStatus, 0, len(pbFiltering.GetStatuses()))

		for _, pbStatus := range pbFiltering.GetStatuses() {
			domainStatus := PbTaskStatusToDomain(pbStatus)
			if domainStatus != nil {
				domainStatuses = append(domainStatuses, *domainStatus)
			}
		}

		filter.Statuses = domainStatuses

		if pbFiltering.Currency != nil {
			switch pbFiltering.GetCurrency() {
			case "EUR":
				filter.Currency = utils.Pointer(models.CurrencyEUR)
			case "USD":
				filter.Currency = utils.Pointer(models.CurrencyUSD)
			case "PLN":
				filter.Currency = utils.Pointer(models.CurrencyPLN)
			default:
				filter.Currency = nil
			}
		}

		filter.IsActive = pbFiltering.IsActive
	}

	return filter
}
