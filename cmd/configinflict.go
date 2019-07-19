package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cInflictCmd)
}

var cInflictCmd = &cobra.Command{
	Use:   "inflict",
	Short: "push tracked values to the workspace",
	Long: longUsage(`Copy the values of tracked keys from commute's config
  into the workspace's git config.`),
	RunE: cInflictFn,
}


func cInflictFn(cmd *cobra.Command, args []string) error {
  return nil
}
