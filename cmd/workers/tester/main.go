package main

import (
	"context"
	"github.com/blr-coder/tasks-svc/pkg/worker"
	"log"
	"time"
)

func main() {
	err := runTestWorker()
	if err != nil {
		log.Fatal(err)
	}
}

func runTestWorker() error {
	w := worker.NewWorker(3 * time.Second)

	w.Run(func(ctx context.Context) error {

		log.Println("TEST WORKER IS RUNNING >>>", time.Now().Format(time.RFC3339))

		return nil
	})

	return nil
}
