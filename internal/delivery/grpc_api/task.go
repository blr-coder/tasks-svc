package grpc

import (
	"context"
	"fmt"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// TODO: Add s.validate(r *taskpbv1.CreateTaskRequest) func, for validation createRequest
	if err := s.Validator.Validate(createRequest); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customerID, err := uuid.Parse(createRequest.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		executorID     *uuid.UUID
		domainCurrency *models.Currency
		amount         *float64
	)

	if createRequest.GetExecutorId() != nil {
		eID, err := uuid.Parse(createRequest.GetExecutorId().GetValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		executorID = &eID
	}

	if createRequest.GetPrice() != nil {
		currency, err := ProtoCurrencyToDomainCurrency(createRequest.GetPrice().GetCurrency())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		amount = utils.Pointer(createRequest.GetPrice().GetAmount())
		domainCurrency = utils.Pointer(currency)
	}

	newId, err := s.taskService.Create(ctx, &models.CreateTask{
		Title:       createRequest.GetTitle(),
		Description: createRequest.GetDescription(),
		CustomerID:  customerID,
		ExecutorID:  executorID,
		Amount:      amount,
		Currency:    domainCurrency,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpbv1.CreateTaskResponse{
		NewTaskId: newId,
	}, nil
}

func ProtoCurrencyToDomainCurrency(pbCurrency string) (currency models.Currency, err error) {
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

func (s *TaskServiceServer) GetTask(ctx context.Context, getRequest *taskpbv1.GetTaskRequest) (*taskpbv1.GetTaskResponse, error) {
	log.Println("get")

	task, err := s.taskService.Get(ctx, getRequest.GetTaskId())
	if err != nil {
		// TODO: Check domain errors. If err is NotFoundError - return nil, status.Error(codes.NotFound, err.Error())
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

	var (
		customerID, executorID uuid.UUID
		err                    error
		updateTask             = &models.UpdateTask{
			ID: updateRequest.GetTaskId(),
		}
	)

	if err = s.Validator.Validate(updateRequest); err != nil {
		return nil, err
	}

	if updateRequest.GetCustomerId() != nil {
		customerID, err = uuid.Parse(updateRequest.GetCustomerId().GetValue())
		if err != nil {
			return nil, err
		}

		updateTask.CustomerID = &customerID
	}

	if updateRequest.GetExecutorId() != nil {
		executorID, err = uuid.Parse(updateRequest.GetExecutorId().GetValue())
		if err != nil {
			return nil, err
		}

		updateTask.ExecutorID = utils.Pointer(executorID)
	}

	if updateRequest.GetTitle() != nil {
		updateTask.Title = utils.Pointer(updateRequest.GetTitle().GetValue())
	}

	if updateRequest.GetDescription() != nil {
		updateTask.Description = utils.Pointer(updateRequest.GetDescription().GetValue())
	}

	if updateRequest.GetStatus() != taskpbv1.TaskStatus_TASK_STATUS_UNSPECIFIED {
		updateTask.Status = utils.Pointer(PbTaskStatusToDomain(updateRequest.GetStatus()))
	}

	if updateRequest.Price != nil {
		currency, err := ProtoCurrencyToDomainCurrency(updateRequest.GetPrice().GetCurrency())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		updateTask.Currency = utils.Pointer(currency)
		updateTask.Amount = utils.Pointer(updateRequest.GetPrice().GetAmount())
	}

	if updateRequest.IsActive != nil {
		updateTask.IsActive = utils.Pointer(updateRequest.IsActive.Value)
	}

	task, err := s.taskService.Update(ctx, updateTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
	log.Println("SetStatus")

	if err := s.Validator.Validate(request); err != nil {
		return nil, err
	}

	_, err := s.taskService.Update(ctx, &models.UpdateTask{
		ID:     request.GetTaskId(),
		Status: utils.Pointer(PbTaskStatusToDomain(request.GetStatus())),
	})
	if err != nil {
		// TODO: Add check for domainError - IsUpdatePossible
		return nil, err
	}

	return &taskpbv1.SetStatusResponse{}, nil
}

func (s *TaskServiceServer) AssignRandomExecutor(ctx context.Context, request *taskpbv1.AssignRandomExecutorRequest) (*taskpbv1.AssignRandomExecutorResponse, error) {
	log.Println("AssignRandomExecutor")

	if err := s.Validator.Validate(request); err != nil {
		return nil, err
	}

	if err := s.taskService.AssignRandomExecutor(ctx, request.GetTaskId()); err != nil {
		return nil, err
	}

	return &taskpbv1.AssignRandomExecutorResponse{}, nil
}
