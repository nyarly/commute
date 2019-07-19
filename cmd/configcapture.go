package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cCaptureCmd)
}

var cCaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "capture configs from the workspace",
	Long: longUsage(`copy the values of all tracked keys
  from the current workspace to the commute configuration.`),
	RunE: cCaptureFn,
}


func cCaptureFn(cmd *cobra.Command, args []string) error {
  return nil
}
