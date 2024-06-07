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

	var executorID *uuid.UUID

	if createRequest.GetExecutorId() != nil {
		eID, err := uuid.Parse(createRequest.GetExecutorId().GetValue())
		if err != nil {
			return nil, err
		}

		executorID = &eID
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
	log.Println("get")

	task, err := s.taskService.Get(ctx, getRequest.GetTaskId())
	if err != nil {
		return nil, err
	}

	log.Println(task)

	return &taskpbv1.GetTaskResponse{Task: domainTaskToPBTask(task)}, nil
}

func (s *TaskServiceServer) ListTasks(ctx context.Context, listRequest *taskpbv1.ListTasksRequest) (*taskpbv1.ListTasksResponse, error) {
	domainTasks, err := s.taskService.List(ctx, &models.TasksFilter{})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.ListTasksResponse{Tasks: domainTasksToPBTasks(domainTasks)}, nil
}

func (s *TaskServiceServer) TotalTasks(ctx context.Context, totalRequest *taskpbv1.TotalTasksRequest) (*taskpbv1.TotalTasksResponse, error) {
	total, err := s.taskService.Count(ctx, &models.TasksFilter{})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.TotalTasksResponse{Total: total}, nil
}

func (s *TaskServiceServer) UpdateTask(ctx context.Context, updateRequest *taskpbv1.UpdateTaskRequest) (*taskpbv1.UpdateTaskResponse, error) {
	log.Println("update")

	if err := s.Validator.Validate(updateRequest); err != nil {
		return nil, err
	}

	customerID, err := uuid.Parse(updateRequest.GetCustomerId())
	if err != nil {
		return nil, err
	}

	var executorID *uuid.UUID

	if updateRequest.GetExecutorId() != nil {
		eID, err := uuid.Parse(updateRequest.GetExecutorId().GetValue())
		if err != nil {
			return nil, err
		}

		executorID = &eID
	}

	task, err := s.taskService.Update(ctx, &models.UpdateTask{
		ID:          updateRequest.GetTaskId(),
		Title:       updateRequest.GetTitle(),
		Description: updateRequest.GetDescription(),
		CustomerID:  customerID,
		ExecutorID:  executorID,
		Status:      PbTaskStatusToDomain(updateRequest.GetStatus()),
	})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.UpdateTaskResponse{
		Task: domainTaskToPBTask(task),
	}, nil
}

func (s *TaskServiceServer) DeleteTask(ctx context.Context, deleteRequest *taskpbv1.DeleteTaskRequest) (*taskpbv1.DeleteTaskResponse, error) {
	log.Println("delete")

	if err := s.taskService.Delete(ctx, deleteRequest.GetTaskId()); err != nil {
		return nil, err
	}

	return &taskpbv1.DeleteTaskResponse{}, nil
}

func (s *TaskServiceServer) AssignExecutor(ctx context.Context, request *taskpbv1.AssignExecutorRequest) (*taskpbv1.AssignExecutorResponse, error) {

	return &taskpbv1.AssignExecutorResponse{}, nil
}

func (s *TaskServiceServer) SetStatus(ctx context.Context, request *taskpbv1.SetStatusRequest) (*taskpbv1.SetStatusResponse, error) {

	return &taskpbv1.SetStatusResponse{}, nil
}
