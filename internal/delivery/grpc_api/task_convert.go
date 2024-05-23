package grpc

import (
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
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
