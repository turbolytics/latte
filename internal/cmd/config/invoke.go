package config

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/turbolytics/latte/internal/collector/initializer"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
)

func NewInvokeCmd() *cobra.Command {
	var configPath string

	var invokeCmd = &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a signal collection",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			lConfig := zap.NewProductionConfig()
			lConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			logger := zap.Must(lConfig.Build())
			defer logger.Sync() // flushes buffer, if any

			ctx := context.Background()

			c, err := initializer.NewCollectorFromFile(
				configPath,
				initializer.WithJustValidation(false),
				initializer.RootWithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			i, err := invoker.New(c,
				invoker.WithLogger(logger),
			)
			if err != nil {
				panic(err)
			}

			if err = i.Invoke(ctx); err != nil {
				panic(err)
			}
		},
	}

	invokeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")
	invokeCmd.MarkFlagRequired("config")

	return invokeCmd
}
