package handlers

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues/kafka"
)

type TaskActionHandler struct {
	receiver          *kafka.Receiver
	taskActionService services.ITaskActionService
}

func NewTaskActionHandler(actionReceiver *kafka.Receiver, actionService services.ITaskActionService) *TaskActionHandler {
	return &TaskActionHandler{
		receiver:          actionReceiver,
		taskActionService: actionService,
	}
}

func (h *TaskActionHandler) Handle(ctx context.Context) error {
	actionFunc := func(ctx context.Context, event any) error {
		taskAction, ok := event.(models.TaskAction)
		if !ok {
			return fmt.Errorf("expected TaskAction, got %T", event)
		}

		if err := h.taskActionService.Save(ctx, &taskAction); err != nil {
			return fmt.Errorf("failed to save task action: %w", err)
		}

		return nil
	}

	return h.receiver.ReceiveWithAction(ctx, actionFunc)
}
