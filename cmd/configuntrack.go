package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cUntrackCmd)
}

var cUntrackCmd = &cobra.Command{
	Use:   "untrack",
	Short: "stop tracking a git config value",
	Long: longUsage(`remove a git config value from the list of values
  tracked by commute for this repo.`),
	RunE: untrackFn,
}


func untrackFn(cmd *cobra.Command, args []string) error {
  return nil
}
