package metric

import (
	"context"
	"testing"
	"time"
)

func Test_MockMetric_Recorded_FunctionCalled(t *testing.T) {
	m := &MockMetric{
		RecordFunc: func(ctx context.Context, command string, latency time.Duration, err error) error {
			return nil
		},
	}
	if err := m.Record(context.Background(), "", 0, nil); err != nil {
		t.Error("Expected nil")
	}
}
