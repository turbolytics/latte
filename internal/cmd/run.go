package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbolytics/latte/internal/collector/initializer"
	"github.com/turbolytics/latte/internal/invoker"
	"github.com/turbolytics/latte/internal/obs"
	"github.com/turbolytics/latte/internal/service"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"log"
)

func NewRunCmd() *cobra.Command {
	var configsGlob string
	var otelExporter string

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run collector daemon",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// Create resource.
			res, err := obs.NewResource()
			if err != nil {
				panic(err)
			}

			// Create a meter provider.
			// You can pass this instance directly to your instrumented code if it
			// accepts a MeterProvider instance.
			meterProvider, err := obs.NewMeterProvider(obs.Exporter(otelExporter), res)
			if err != nil {
				panic(err)
			}

			// Handle shutdown properly so nothing leaks.
			defer func() {
				if err := meterProvider.Shutdown(context.Background()); err != nil {
					log.Println(err)
				}
			}()

			otel.SetMeterProvider(meterProvider)

			go otelruntime.Start(otelruntime.WithMeterProvider(meterProvider))

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			if obs.Exporter(otelExporter) == obs.ExporterPrometheus {
				go obs.ServeMetrics(logger, ":12223")
			}

			logger.Info(
				"loading configs",
				zap.String("path", configsGlob),
			)

			// initialize all collectors in the path
			collectors, err := initializer.NewCollectorsFromGlob(
				configsGlob,
				initializer.RootWithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			var invokers []*invoker.Invoker
			for _, coll := range collectors {
				i, err := invoker.New(coll,
					invoker.WithLogger(logger),
				)
				if err != nil {
					panic(err)
				}
				invokers = append(invokers, i)
			}
			logger.Info(
				"initialized invokers",
				zap.Int("num_invokers", len(invokers)),
			)

			// invocation all collectors at their desired intervals
			s, err := service.NewService(
				invokers,
				service.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			defer s.Shutdown()

			ctx := context.Background()
			if err := s.Run(ctx); err != nil {
				panic(err)
			}
		},
	}

	runCmd.Flags().StringVarP(&configsGlob, "configs", "c", "", "Path to config directory")
	runCmd.Flags().StringVarP(&otelExporter, "otel-exporter", "", "prometheus", "Opentelemetry exporter: 'console', prometheus")
	runCmd.MarkFlagRequired("config")

	return runCmd
}
