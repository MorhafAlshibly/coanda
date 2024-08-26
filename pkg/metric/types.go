package metric

import (
	"context"
	"time"
)

type Metric interface {
	Record(ctx context.Context, command string, latency time.Duration, err error) error
}

type MockMetric struct {
	RecordFunc func(ctx context.Context, command string, latency time.Duration, err error) error
}

func (m *MockMetric) Record(ctx context.Context, command string, latency time.Duration, err error) error {
	return m.RecordFunc(ctx, command, latency, err)
}
