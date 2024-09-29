package worker

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Worker struct {
	tickDuration time.Duration
}

func NewWorker(tickDuration time.Duration) *Worker {
	return &Worker{tickDuration: tickDuration}
}

func (w *Worker) Run(executionFunc func(ctx context.Context) error) {
	log.Println("RUN WORKER...")

	ticker := time.NewTicker(w.tickDuration)
	defer ticker.Stop()

	for range ticker.C {
		err := executionFunc(context.TODO())
		if err != nil {
			err = fmt.Errorf("run worker err: %w", err)
			log.Println(err)
		}
	}
}
