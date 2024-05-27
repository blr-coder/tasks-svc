package grpc

import (
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func domainTasksToPBTasks(domainTasks []*models.Task) []*taskpbv1.Task {
	ls := make([]*taskpbv1.Task, len(domainTasks))
	for i, task := range domainTasks {
		ls[i] = domainTaskToPBTask(task)
	}

	return ls
}

func domainTaskToPBTask(task *models.Task) *taskpbv1.Task {
	return &taskpbv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CustomerId:  task.CustomerID.String(),
		ExecutorId:  task.ExecutorID.String(),
		Status:      DomainTaskStatusToPB(task.Status),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
		IsActive:    true,
	}
}
