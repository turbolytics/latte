package obs

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Exporter string

const (
	ExporterConsole    Exporter = "console"
	ExporterPrometheus Exporter = "prometheus"
)

func ErrToStatus(err error) string {
	if err != nil {
		return "ERROR"
	}
	return "OK"
}

func NewResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("signals-collector"),
		))
}

func NewMeterProvider(exporter Exporter, res *resource.Resource) (*metric.MeterProvider, error) {
	var meterProvider *metric.MeterProvider
	switch exporter {
	case ExporterConsole:
		metricExporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}

		meterProvider = metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(metricExporter,
				// Default is 1m. Set to 3s for demonstrative purposes.
				metric.WithInterval(60*time.Second))),
		)
	case ExporterPrometheus:
		exporter, err := prometheus.New()
		if err != nil {
			return nil, err
		}
		meterProvider = metric.NewMeterProvider(metric.WithReader(exporter))
	default:
		return nil, fmt.Errorf("unsupported exporter: %q", exporter)
	}

	return meterProvider, nil
}

func ServeMetrics(l *zap.Logger, addr string) {
	l.Info(
		"serving metrics at /metrics",
		zap.String("addr", addr),
	)

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
