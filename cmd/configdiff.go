package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cDiffCmd)
}

var cDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "display the diff of configs",
	Long: longUsage(`Output the difference between commute's version of tracked
  git configuration values and the values current in the workspace.`),
	RunE: cdiffFn,
}


func cdiffFn(cmd *cobra.Command, args []string) error {
  return nil
}
