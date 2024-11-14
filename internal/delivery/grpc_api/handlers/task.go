package handlers

import (
	"context"
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

type TaskGRPCHandler struct {
	taskpbv1.UnimplementedTaskServiceServer
	Validator   *protovalidate.Validator
	taskService services.ITaskService
}

func NewTaskGRPCHandler(validator *protovalidate.Validator, taskService services.ITaskService) *TaskGRPCHandler {
	return &TaskGRPCHandler{
		Validator:   validator,
		taskService: taskService,
	}
}

func (s *TaskGRPCHandler) CreateTask(ctx context.Context, createRequest *taskpbv1.CreateTaskRequest) (*taskpbv1.CreateTaskResponse, error) {
	log.Println("create in TaskServiceServer")

	task, err := s.validateCreateTaskAndConvertToDomain(createRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newTaskId, err := s.taskService.CreateWithTransaction(ctx, task)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpbv1.CreateTaskResponse{
		NewTaskId: newTaskId,
	}, nil
}

func (s *TaskGRPCHandler) GetTask(ctx context.Context, getRequest *taskpbv1.GetTaskRequest) (*taskpbv1.GetTaskResponse, error) {
	log.Println("get")

	task, err := s.taskService.Get(ctx, getRequest.GetTaskId())
	if err != nil {
		// TODO: Check domain errors. If err is NotFoundError - return nil, status.Error(codes.NotFound, err.Error())
		return nil, err
	}

	log.Println(task)

	return &taskpbv1.GetTaskResponse{Task: domainTaskToPBTask(task)}, nil
}

func (s *TaskGRPCHandler) ListTasks(ctx context.Context, listRequest *taskpbv1.ListTasksRequest) (*taskpbv1.ListTasksResponse, error) {
	domainTasks, err := s.taskService.List(ctx, &models.TasksFilter{
		Filtering: PbListTasksFilteringToDomain(listRequest.GetFiltering()),
		Sorting:   PbTasksSortingToDomain(listRequest.GetSorting()),
		Limiting:  PbListLimitingToDB(listRequest.GetLimiting()),
	})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.ListTasksResponse{Tasks: domainTasksToPBTasks(domainTasks)}, nil
}

func (s *TaskGRPCHandler) TotalTasks(ctx context.Context, totalRequest *taskpbv1.TotalTasksRequest) (*taskpbv1.TotalTasksResponse, error) {
	total, err := s.taskService.Count(ctx, &models.TasksFilter{
		Filtering: PbListTasksFilteringToDomain(totalRequest.GetFiltering()),
		Sorting:   PbTasksSortingToDomain(totalRequest.GetSorting()),
		Limiting:  PbListLimitingToDB(totalRequest.GetLimiting()),
	})
	if err != nil {
		return nil, err
	}

	return &taskpbv1.TotalTasksResponse{Total: total}, nil
}

func (s *TaskGRPCHandler) UpdateTask(ctx context.Context, updateRequest *taskpbv1.UpdateTaskRequest) (*taskpbv1.UpdateTaskResponse, error) {
	log.Println("update")

	var (
		customerID, executorID uuid.UUID
		err                    error
		updateTask             = &models.UpdateTask{
			ID:          updateRequest.GetTaskId(),
			Title:       updateRequest.Title,
			Description: updateRequest.Description,
			Status:      PbTaskStatusToDomain(updateRequest.GetStatus()),
			Currency:    nil,
			Amount:      nil,
			IsActive:    updateRequest.IsActive,
		}
	)

	if err = s.Validator.Validate(updateRequest); err != nil {
		return nil, err
	}

	if updateRequest.CustomerId != nil {
		customerID, err = uuid.Parse(updateRequest.GetCustomerId())
		if err != nil {
			return nil, err
		}

		updateTask.CustomerID = &customerID
	}

	if updateRequest.ExecutorId != nil {
		executorID, err = uuid.Parse(updateRequest.GetExecutorId())
		if err != nil {
			return nil, err
		}

		updateTask.ExecutorID = utils.Pointer(executorID)
	}

	if updateRequest.Price != nil {
		currency, err := PBCurrencyToDomainCurrency(updateRequest.GetPrice().GetCurrency())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		updateTask.Currency = utils.Pointer(currency)
		updateTask.Amount = utils.Pointer(updateRequest.GetPrice().GetAmount())
	}

	task, err := s.taskService.Update(ctx, updateTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpbv1.UpdateTaskResponse{
		Task: domainTaskToPBTask(task),
	}, nil
}

func (s *TaskGRPCHandler) DeleteTask(ctx context.Context, deleteRequest *taskpbv1.DeleteTaskRequest) (*taskpbv1.DeleteTaskResponse, error) {
	log.Println("delete")

	if err := s.taskService.Delete(ctx, deleteRequest.GetTaskId()); err != nil {
		return nil, err
	}

	return &taskpbv1.DeleteTaskResponse{}, nil
}

func (s *TaskGRPCHandler) AssignExecutor(ctx context.Context, request *taskpbv1.AssignExecutorRequest) (*taskpbv1.AssignExecutorResponse, error) {

	return &taskpbv1.AssignExecutorResponse{}, nil
}

func (s *TaskGRPCHandler) SetStatus(ctx context.Context, request *taskpbv1.SetStatusRequest) (*taskpbv1.SetStatusResponse, error) {
	log.Println("SetStatus")

	if err := s.Validator.Validate(request); err != nil {
		return nil, err
	}

	_, err := s.taskService.Update(ctx, &models.UpdateTask{
		ID:     request.GetTaskId(),
		Status: PbTaskStatusToDomain(request.GetStatus()),
	})
	if err != nil {
		// TODO: Add check for domainError - IsUpdatePossible
		return nil, err
	}

	return &taskpbv1.SetStatusResponse{}, nil
}

func (s *TaskGRPCHandler) AssignRandomExecutor(ctx context.Context, request *taskpbv1.AssignRandomExecutorRequest) (*taskpbv1.AssignRandomExecutorResponse, error) {
	log.Println("AssignRandomExecutor")

	if err := s.Validator.Validate(request); err != nil {
		return nil, err
	}

	if err := s.taskService.AssignRandomExecutor(ctx, request.GetTaskId()); err != nil {
		return nil, err
	}

	return &taskpbv1.AssignRandomExecutorResponse{}, nil
}

func (s *TaskGRPCHandler) validateCreateTaskAndConvertToDomain(createRequest *taskpbv1.CreateTaskRequest) (*models.CreateTask, error) {
	if err := s.Validator.Validate(createRequest); err != nil {
		return nil, err
	}

	return PbCreateTaskToDomain(createRequest)
}
