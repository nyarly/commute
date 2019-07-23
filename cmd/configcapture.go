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
  Args: cobra.ExactArgs(0),
}


func cCaptureFn(cmd *cobra.Command, args []string) error {
	tracked, err := repoConfigs()
	if err != nil {
		return err
	}

	for name := range tracked {
		// find value in config
		gv, err := workspaceValues(".", name)
		if err != nil {
			return err
		}

    tracked[name] = gv
	}

	return nil
}
