package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/turbolytics/collector/internal/cmd/config"
	"os"
)

func init() {
	configCmd := config.NewConfigCmd()
	configCmd.AddCommand(
		config.NewValidateCmd(),
		config.NewInvokeCmd(),
	)
	rootCmd.AddCommand(
		configCmd,
		NewRunCmd(),
	)
}

var rootCmd = &cobra.Command{
	Use:   "collector",
	Short: "Signals Collector queries source systems to capture and emit signals",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
