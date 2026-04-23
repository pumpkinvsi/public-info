package outbox

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, value []byte) error
}

type WorkerPool struct {
	workers      int
	pollInterval time.Duration
	outbox       Outbox
	publisher    Publisher
}

func NewWorkerPool(workers int, pollInterval time.Duration, outbox Outbox, publisher Publisher) *WorkerPool {
	return &WorkerPool{
		workers:      workers,
		pollInterval: pollInterval,
		outbox:       outbox,
		publisher:    publisher,
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := range wp.workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			wp.runWorker(ctx, id)
		}(i)
	}

	wg.Wait()
	slog.Info("outbox worker pool stopped")
}

func (wp *WorkerPool) runWorker(ctx context.Context, id int) {
	slog.Info("outbox worker started", "worker_id", id)

	for {
		processed, err := wp.outbox.ProcessNext(ctx, func(msg *Message) error {
			return wp.publisher.Publish(ctx, msg.Payload)
		})

		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("outbox: process message failed", "worker_id", id, "error", err)
		}

		if processed {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(wp.pollInterval):
		}
	}
}
