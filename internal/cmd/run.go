package cmd

import "github.com/spf13/cobra"

func NewRunCmd() *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run collector daemon",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// initialize all collectors in the path
			// schedule all collectors at their desired intervals
		},
	}

	return runCmd
}
