package grpc

import (
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

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

func PbTaskStatusToDomain(pbStatus taskpbv1.TaskStatus) (dStatus models.TaskStatus) {
	switch pbStatus {
	case taskpbv1.TaskStatus_TASK_STATUS_PENDING:
		dStatus = models.PendingStatus
	case taskpbv1.TaskStatus_TASK_STATUS_PROCESSING:
		dStatus = models.ProcessingStatus
	case taskpbv1.TaskStatus_TASK_STATUS_DONE:
		dStatus = models.DoneStatus
	default:
		// TODO: UNSPECIFIED???
		dStatus = models.PendingStatus
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

func PbListTasksRequestToDomain(listRequest *taskpbv1.ListTasksRequest) *models.TasksFilter {
	return &models.TasksFilter{
		Filtering: PbListTasksFilteringToDomain(listRequest.GetFiltering()),
		Sorting:   PbTasksSortingToDomain(listRequest.GetSorting()),
		Limiting:  PbListLimitingToDB(listRequest.GetLimiting()),
	}
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
			domainStatuses = append(domainStatuses, PbTaskStatusToDomain(pbStatus))
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