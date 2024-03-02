package grpc

import (
	"context"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/google/uuid"
	"log"
	"tasks-svc/internal/domain/models"
	"tasks-svc/internal/domain/usecases"
)

type TaskServiceServer struct {
	taskpbv1.UnimplementedTaskServiceServer
	taskUseCase usecases.ITaskUseCase
}

func NewTaskServiceServer(useCase usecases.ITaskUseCase) *TaskServiceServer {
	return &TaskServiceServer{
		taskUseCase: useCase,
	}
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, createRequest *taskpbv1.CreateTaskRequest) (*taskpbv1.CreateTaskResponse, error) {
	log.Println("create")

	customerID, err := uuid.Parse(createRequest.GetCustomerId())
	if err != nil {
		return nil, err
	}

	executorID, err := uuid.Parse(createRequest.GetExecutorId())
	if err != nil {
		return nil, err
	}

	newId, err := s.taskUseCase.Create(ctx, &models.CreateTask{
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

	task, err := s.taskUseCase.Get(ctx, getRequest.GetTaskId())
	if err != nil {
		return nil, err
	}

	log.Println(task)

	return &taskpbv1.GetTaskResponse{}, nil
}
