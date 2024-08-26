package metric

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetric struct {
	service  string
	metric   map[string]*prometheusCommandMetric
	registry *prometheus.Registry
}

func NewPrometheusMetric(reg *prometheus.Registry, service string, port uint16) (*PrometheusMetric, error) {
	m := &PrometheusMetric{
		service:  service,
		metric:   make(map[string]*prometheusCommandMetric),
		registry: reg,
	}
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metric", promHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return m, nil
}

func (m *PrometheusMetric) Record(ctx context.Context, command string, latency time.Duration, err error) error {
	if _, ok := m.metric[command]; !ok {
		m.metric[command] = m.newPrometheusCommandMetric(m.registry, command)
	}
	err = m.metric[command].Record(ctx, latency, err)
	return err
}

type prometheusCommandMetric struct {
	command       string
	totalRequests prometheus.Counter
	totalErrors   prometheus.Counter
	latency       prometheus.Histogram
}

func (m *PrometheusMetric) newPrometheusCommandMetric(reg *prometheus.Registry, command string) *prometheusCommandMetric {
	metric := &prometheusCommandMetric{
		command: command,
		totalRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: m.service,
			Subsystem: command,
			Name:      "total_requests",
			Help:      "Total number of requests",
		}),
		totalErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: m.service,
			Subsystem: command,
			Name:      "total_errors",
			Help:      "Total number of errors",
		}),
		latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: m.service,
			Subsystem: command,
			Name:      "latency",
			Help:      "Latency in milliseconds",
		}),
	}
	reg.MustRegister(metric.totalRequests)
	reg.MustRegister(metric.totalErrors)
	reg.MustRegister(metric.latency)
	return metric
}

func (m *prometheusCommandMetric) Record(ctx context.Context, latency time.Duration, err error) error {
	m.totalRequests.Inc()
	if err != nil {
		m.totalErrors.Inc()
	}
	m.latency.Observe(float64(latency.Milliseconds()))
	return nil
}
