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
	Args: cobra.ExactArgs(1),
}


func untrackFn(cmd *cobra.Command, args []string) error {
	tracked, err := repoConfigs()
	if err != nil {
		return err
	}

  delete(tracked, args[0])

  return nil
}
