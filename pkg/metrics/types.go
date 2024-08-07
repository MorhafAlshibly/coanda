package metrics

import (
	"context"
	"time"
)

type Metrics interface {
	Record(ctx context.Context, command string, latency time.Duration, err error) error
}

type MockMetrics struct {
	RecordFunc func(ctx context.Context, command string, latency time.Duration, err error) error
}

func (m *MockMetrics) Record(ctx context.Context, command string, latency time.Duration, err error) error {
	return m.RecordFunc(ctx, command, latency, err)
}
