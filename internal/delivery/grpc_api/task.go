package grpc

import (
	"context"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"log"
)

type TaskServiceServer struct {
	taskpbv1.UnimplementedTaskServiceServer
	Validator   *protovalidate.Validator
	taskService services.ITaskService
}

func NewTaskServiceServer(validator *protovalidate.Validator, taskService services.ITaskService) *TaskServiceServer {
	return &TaskServiceServer{
		Validator:   validator,
		taskService: taskService,
	}
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, createRequest *taskpbv1.CreateTaskRequest) (*taskpbv1.CreateTaskResponse, error) {
	log.Println("create in TaskServiceServer")

	if err := s.Validator.Validate(createRequest); err != nil {
		return nil, err
	}

	customerID, err := uuid.Parse(createRequest.GetCustomerId())
	if err != nil {
		return nil, err
	}

	executorID, err := uuid.Parse(createRequest.GetExecutorId())
	if err != nil {
		return nil, err
	}

	newId, err := s.taskService.Create(ctx, &models.CreateTask{
		Title:       createRequest.GetTitle(),
		Description: createRequest.GetDescription(),
		CustomerID:  customerID,
		ExecutorID:  executorID,
	})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.CreateTaskResponse{
		NewTaskId: newId,
	}, nil
}

func (s *TaskServiceServer) GetTask(ctx context.Context, getRequest *taskpbv1.GetTaskRequest) (*taskpbv1.GetTaskResponse, error) {
	log.Println("gat")

	task, err := s.taskService.Get(ctx, getRequest.GetTaskId())
	if err != nil {
		return nil, err
	}

	log.Println(task)

	return &taskpbv1.GetTaskResponse{}, nil
}
