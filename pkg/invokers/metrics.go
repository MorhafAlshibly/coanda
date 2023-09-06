package invokers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/metrics"
)

type MetricsInvoker struct {
	invoker Invoker
	metrics metrics.Metrics
}

func NewMetricsInvoker(metrics metrics.Metrics) *MetricsInvoker {
	return &MetricsInvoker{
		invoker: &BasicInvoker{},
		metrics: metrics,
	}
}

func (i *MetricsInvoker) SetInvoker(invoker Invoker) *MetricsInvoker {
	i.invoker = invoker
	return i
}

func (i *MetricsInvoker) Invoke(ctx context.Context, command Command) error {
	start := time.Now()
	err := i.invoker.Invoke(ctx, command)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	commandNameSplit := strings.Split(strings.Replace(fmt.Sprintf("%T", command), "*", "", -1), ".")
	go i.metrics.Record(context.Background(), commandNameSplit[len(commandNameSplit)-1], elapsed, err)
	return nil
}
