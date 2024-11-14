package handlers

import (
	"context"
	"encoding/json"
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
		// map[string]interface{}

		var taskActionDTO TaskActionDTO
		// TODO: Think about json.Marshal and json.Unmarshal. It's difficult. Maybe using "https://github.com/mitchellh/mapstructure" will be better, or something else.
		data, err := json.Marshal(event) // Кодируем в JSON
		if err != nil {
			return fmt.Errorf("failed to marshal event to JSON: %w", err)
		}

		// Декодируем JSON в нужную структуру
		if err := json.Unmarshal(data, &taskActionDTO); err != nil {
			return fmt.Errorf("failed to unmarshal JSON to TaskAction: %w", err)
		}

		if err := h.taskActionService.Save(ctx, taskActionDTO.ToDomainAction()); err != nil {
			return fmt.Errorf("failed to save task action: %w", err)
		}

		return nil
	}

	return h.receiver.ReceiveWithAction(ctx, actionFunc)
}

type TaskActionDTO struct {
	ID     int64  `json:"id"`
	TaskID int64  `json:"task_id"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

func (d *TaskActionDTO) ToDomainAction() *models.TaskAction {
	return &models.TaskAction{
		ExternalID: d.ID,
		TaskID:     d.TaskID,
		Type:       models.ActionType(d.Type),
		URL:        d.URL,
	}
}
