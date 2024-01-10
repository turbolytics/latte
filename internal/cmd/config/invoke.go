package config

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbolytics/collector/internal/collector"
	"github.com/turbolytics/collector/internal/config"
	"github.com/turbolytics/collector/internal/config/sc"
	"go.uber.org/zap"
)

func NewInvokeCmd() *cobra.Command {
	var configPath string
	var scConfigPath string

	var invokeCmd = &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a signal collection",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			ctx := context.Background()
			config, err := config.NewFromFile(
				configPath,
				config.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			scConfig, err := sc.NewFromFile(
				scConfigPath,
				sc.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			c, err := collector.New(
				config,
				collector.WithLogger(logger),
				collector.WithStateStorer(scConfig.StateStore.Storer),
			)
			if err != nil {
				panic(err)
			}

			_, err = c.Invoke(ctx)
			if err != nil {
				panic(err)
			}
		},
	}

	invokeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")
	invokeCmd.Flags().StringVarP(&scConfigPath, "sc-config", "", "", "Path to signals collector config file")
	invokeCmd.MarkFlagRequired("config")
	invokeCmd.MarkFlagRequired("sc-config")

	return invokeCmd
}
