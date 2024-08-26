package invoker

import (
	"context"
	"testing"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/metric"
)

func Test_MetricInvoker_Invoked_TimeRecorded(t *testing.T) {
	i := NewMetricInvoker(&metric.MockMetric{
		RecordFunc: func(ctx context.Context, command string, latency time.Duration, err error) error {
			return nil
		},
	})
	c := &MockCommand{
		ExecuteFunc: func(ctx context.Context) error {
			return nil
		},
	}
	if err := i.Invoke(context.Background(), c); err != nil {
		t.Error("Expected nil")
	}
}
