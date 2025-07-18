package queue

import "context"

type Queue interface {
	Enqueue(EmailJob) error
	StartWorkers(ctx context.Context)
	Close()
	Wait()
}
