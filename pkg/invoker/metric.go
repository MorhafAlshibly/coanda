package invoker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MorhafAlshibly/coanda/pkg/metric"
)

type MetricInvoker struct {
	invoker Invoker
	metric  metric.Metric
}

func NewMetricInvoker(metric metric.Metric) *MetricInvoker {
	return &MetricInvoker{
		invoker: &BasicInvoker{},
		metric:  metric,
	}
}

func (i *MetricInvoker) SetInvoker(invoker Invoker) *MetricInvoker {
	i.invoker = invoker
	return i
}

func (i *MetricInvoker) Invoke(ctx context.Context, command Command) error {
	start := time.Now()
	err := i.invoker.Invoke(ctx, command)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	commandNameSplit := strings.Split(strings.Replace(fmt.Sprintf("%T", command), "*", "", -1), ".")
	go i.metric.Record(context.Background(), commandNameSplit[len(commandNameSplit)-1], elapsed, err)
	return nil
}
