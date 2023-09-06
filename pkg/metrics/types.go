package metrics

import (
	"context"
	"time"
)

type Metrics interface {
	Record(ctx context.Context, command string, latency time.Duration, err error) error
}
