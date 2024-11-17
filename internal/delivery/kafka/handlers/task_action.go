package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/events"
)

type TaskActionHandler struct {
	receiver events.IEventReceiver
	//taskActionService          services.ITaskActionService
	taskActionProcessingRunner events.IProcessRunner
}

func NewTaskActionHandler(actionReceiver events.IEventReceiver, actionService services.ITaskActionService) *TaskActionHandler {
	return &TaskActionHandler{
		receiver: actionReceiver,
		//taskActionService: actionService,
		taskActionProcessingRunner: &TaskActionProcessingRunner{
			TaskActionService: actionService,
		},
	}
}

type TaskActionProcessingRunner struct {
	TaskActionService services.ITaskActionService
}

func (r *TaskActionProcessingRunner) Run(ctx context.Context, event any) error {
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

	if err := r.TaskActionService.Save(ctx, taskActionDTO.ToDomainAction()); err != nil {
		return fmt.Errorf("failed to save task action: %w", err)
	}

	return nil
}

func (h *TaskActionHandler) Handle(ctx context.Context) error {
	/*actionFunc := func(ctx context.Context, event any) error {
		// event is map[string]interface{}

		// TODO: Add TaskAction to schema registry
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
	}*/

	//return h.receiver.ReceiveWithAction(ctx, h.taskActionProcessingRunner.Run)
	return h.receiver.ReceiveWithRunner(ctx, h.taskActionProcessingRunner)
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
