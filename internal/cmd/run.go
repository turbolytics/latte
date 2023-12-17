package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbolytics/collector/internal"
	"github.com/turbolytics/collector/internal/collector"
	"github.com/turbolytics/collector/internal/collector/service"
	"go.uber.org/zap"
)

func NewRunCmd() *cobra.Command {
	var configDir string

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run collector daemon",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			logger.Info(
				"loading configs",
				zap.String("path", configDir),
			)

			// initialize all collectors in the path
			confs, err := internal.NewConfigsFromDir(configDir)
			if err != nil {
				panic(err)
			}

			cs, err := collector.NewFromConfigs(
				confs,
				collector.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			logger.Info(
				"initialized collectors",
				zap.Int("num_collectors", len(cs)),
			)

			// schedule all collectors at their desired intervals
			s, err := service.NewService(
				cs,
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

	runCmd.Flags().StringVarP(&configDir, "config-dir", "c", "", "Path to config directory")
	runCmd.MarkFlagRequired("config")

	return runCmd
}
