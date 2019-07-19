package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cTrackCmd)
}

var cTrackCmd = &cobra.Command{
	Use:   "track",
	Short: "start tracking a git config value",
	Long: longUsage(`add a git config value to the list of values
  tracked by commute for this repo.`),
	RunE: trackFn,
	Args:  cobra.ExactArgs(1),
}


func trackFn(cmd *cobra.Command, args []string) error {
  // find value in config
  // add value to commute config
  return nil
}
