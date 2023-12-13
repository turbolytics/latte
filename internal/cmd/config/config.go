package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "config ",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config")
		},
	}
	return configCmd
}
