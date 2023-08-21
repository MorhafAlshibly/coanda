package queue

import "context"

// Queuer is the interface that implements the queue
type Queuer interface {
	Enqueue(ctx context.Context, data []byte) error
}

// MockQueue is a mock implementation of the queue
