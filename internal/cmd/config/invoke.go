package config

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbolytics/collector/internal/collector"
	"go.uber.org/zap"
)

func NewInvokeCmd() *cobra.Command {
	var configPath string

	var invokeCmd = &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a signal collection",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			ctx := context.Background()

			c, err := collector.NewFromFile(
				configPath,
				collector.WithJustValidation(true),
				collector.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			if err = c.Invoke(ctx); err != nil {
				panic(err)
			}
		},
	}

	invokeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")
	invokeCmd.MarkFlagRequired("config")

	return invokeCmd
}
