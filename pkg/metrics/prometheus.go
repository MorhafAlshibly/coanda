package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	service  string
	metrics  map[string]*prometheusCommandMetrics
	registry *prometheus.Registry
}

func NewPrometheusMetrics(reg *prometheus.Registry, service string, port uint16) (*PrometheusMetrics, error) {
	m := &PrometheusMetrics{
		service:  service,
		metrics:  make(map[string]*prometheusCommandMetrics),
		registry: reg,
	}
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return m, nil
}

func (m *PrometheusMetrics) Record(ctx context.Context, command string, latency time.Duration, err error) error {
	if _, ok := m.metrics[command]; !ok {
		m.metrics[command] = m.newPrometheusCommandMetrics(m.registry, command)
	}
	err = m.metrics[command].Record(ctx, latency, err)
	return err
}

type prometheusCommandMetrics struct {
	command       string
	totalRequests prometheus.Counter
	totalErrors   prometheus.Counter
	latency       prometheus.Histogram
}

func (m *PrometheusMetrics) newPrometheusCommandMetrics(reg *prometheus.Registry, command string) *prometheusCommandMetrics {
	metrics := &prometheusCommandMetrics{
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
	reg.MustRegister(metrics.totalRequests)
	reg.MustRegister(metrics.totalErrors)
	reg.MustRegister(metrics.latency)
	return metrics
}

func (m *prometheusCommandMetrics) Record(ctx context.Context, latency time.Duration, err error) error {
	m.totalRequests.Inc()
	if err != nil {
		m.totalErrors.Inc()
	}
	m.latency.Observe(float64(latency.Milliseconds()))
	return nil
}
