package models

import (
	"github.com/blr-coder/tasks-svc/pkg/utils"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/jaswdr/faker"
)

func NewMockTask(t *testing.T) *Task {
	t.Helper()
	f := faker.New()
	return &Task{
		ID:          int64(f.RandomDigit()),
		Title:       f.RandomLetter(),
		Description: f.RandomLetter(),
		CustomerID:  uuid.UUID([]byte(f.UUID().V4())),
		ExecutorID:  utils.Pointer(uuid.UUID([]byte(f.UUID().V4()))),
		Status: TaskStatus(f.RandomStringElement([]string{
			PendingStatus.String(),
			ProcessingStatus.String(),
			DoneStatus.String(),
		})),
		IsActive: f.Boolean().Bool(),
		Currency: utils.Pointer(
			Currency(f.RandomStringElement([]string{
				CurrencyEUR.String(),
				CurrencyUSD.String(),
				CurrencyPLN.String(),
			})),
		),
		Amount:    utils.Pointer(float64(f.RandomDigit())),
		CreatedAt: f.Time().Time(time.Now()).UTC(),
		UpdatedAt: f.Time().Time(time.Now()).UTC(),
	}
}
