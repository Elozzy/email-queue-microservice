package queue

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type EmailQueue struct {
	queue       chan EmailJob
	wg          sync.WaitGroup
	workerCount int
}

func NewEmailQueue(cfg Config) *EmailQueue {
	return &EmailQueue{
		queue:       make(chan EmailJob, cfg.QueueSize),
		workerCount: cfg.WorkerCount,
	}
}

func (q *EmailQueue) Enqueue(job EmailJob) error {
	select {
	case q.queue <- job:
		return nil
	default:
		return errors.New("Queue is full")
	}
}

func (q *EmailQueue) StartWorkers(ctx context.Context) {
	for i := 0; i < q.workerCount; i++ {
		q.wg.Add(1)
		go q.worker(ctx, i+1)
	}
}

func (q *EmailQueue) worker(ctx context.Context, id int) {
	defer q.wg.Done()
	log.Printf("Worker %d started", id)
	for {
		select {
		case job, ok := <-q.queue:
			if !ok {
				log.Printf("Worker %d: queue closed", id)
				return
			}
			log.Printf("Worker %d: Sending email to %s", id, job.To)
			time.Sleep(1 * time.Second)
			log.Printf("Worker %d: Email sent to %s", id, job.To)
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		}
	}
}

func (q *EmailQueue) Close() {
	close(q.queue)
}

func (q *EmailQueue) Wait() {
	q.wg.Wait()
}
