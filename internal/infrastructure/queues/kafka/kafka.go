package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	taskschemasregistry "github.com/blr-coder/task_system_schemas_registry/task_svc"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Sender struct {
	kafkaProducer sarama.SyncProducer
}

func NewKafkaSender(config config.KafkaConfig) (*Sender, error) {
	newConfig := sarama.NewConfig()
	newConfig.Producer.RequiredAcks = sarama.WaitForAll
	newConfig.Producer.Partitioner = sarama.NewHashPartitioner
	newConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer([]string{config.Address}, newConfig)
	if err != nil {
		return nil, err
	}
	return &Sender{
		kafkaProducer: syncProducer,
	}, nil
}

func (s *Sender) Send(ctx context.Context, event *events.Event) error {
	log.Println("Send")
	_, _, err := s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     string(event.Topic),
		Key:       sarama.StringEncoder(event.Name),
		Value:     sarama.ByteEncoder(event.Data),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Println("SENT")
	return nil
}

func (s *Sender) SendTaskCreated(ctx context.Context, task *models.Task) error {
	log.Println("TaskCreatedSend")

	event := &taskschemasregistry.BaseEvent{
		Id:        uuid.New().String(),
		CreatedAt: timestamppb.New(time.Now().UTC()),
		TopicType: taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING,
		EventType: taskschemasregistry.EventType_EVENT_TYPE_CREATED,
		Data:      domainTaskToPB(task),
	}

	// TODO: Think about proto.Marshal(event)
	/*kafkaMessageValue, err := proto.Marshal(event)
	if err != nil {
		return err
	}*/

	kafkaMessageValue, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     taskschemasregistry.TopicType_name[int32(taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING)],
		Key:       sarama.StringEncoder(taskschemasregistry.EventType_EVENT_TYPE_CREATED),
		Value:     sarama.ByteEncoder(kafkaMessageValue),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Println("TaskCreated SENT")
	return nil
}

func (s *Sender) SendTaskUpdated(ctx context.Context, task *models.Task) error {
	log.Println("TaskUpdated")

	event := &taskschemasregistry.BaseEvent{
		Id:        uuid.New().String(),
		CreatedAt: timestamppb.New(time.Now().UTC()),
		TopicType: taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING,
		EventType: taskschemasregistry.EventType_EVENT_TYPE_UPDATED,
		Data:      domainTaskToPB(task),
	}

	kafkaMessageValue, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     taskschemasregistry.TopicType_name[int32(taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING)],
		Key:       sarama.StringEncoder(taskschemasregistry.EventType_EVENT_TYPE_UPDATED),
		Value:     sarama.ByteEncoder(kafkaMessageValue),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Println("TaskUpdated SENT")
	return nil
}

func (s *Sender) SendTaskDeleted(ctx context.Context, task *models.Task) error {
	log.Println("TaskDeleted")

	event := &taskschemasregistry.BaseEvent{
		Id:        uuid.New().String(),
		CreatedAt: timestamppb.New(time.Now().UTC()),
		TopicType: taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING,
		EventType: taskschemasregistry.EventType_EVENT_TYPE_DELETED,
		Data:      domainTaskToPB(task),
	}

	kafkaMessageValue, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     taskschemasregistry.TopicType_name[int32(taskschemasregistry.TopicType_TOPIC_TYPE_CUD_STREAMING)],
		Key:       sarama.StringEncoder(taskschemasregistry.EventType_EVENT_TYPE_DELETED),
		Value:     sarama.ByteEncoder(kafkaMessageValue),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Println("TaskDeleted SENT")
	return nil
}

func domainTaskStatusToPB(domainStatus models.TaskStatus) (pbStatus taskschemasregistry.TaskStatus) {
	switch domainStatus {
	case models.PendingStatus:
		pbStatus = taskschemasregistry.TaskStatus_TASK_STATUS_PENDING
	case models.ProcessingStatus:
		pbStatus = taskschemasregistry.TaskStatus_TASK_STATUS_PROCESSING
	case models.DoneStatus:
		pbStatus = taskschemasregistry.TaskStatus_TASK_STATUS_DONE
	}
	return pbStatus
}

func domainTaskToPB(task *models.Task) *taskschemasregistry.TaskV1 {
	return &taskschemasregistry.TaskV1{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CustomerId:  task.CustomerID.String(),
		ExecutorId:  task.ExecutorID.String(),
		Status:      domainTaskStatusToPB(task.Status),
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}
}
