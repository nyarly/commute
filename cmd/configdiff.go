package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cDiffCmd)
	queryCommand(cDiffCmd)
}

var cDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "display the diff of configs",
	Long: longUsage(`Output the difference between commute's version of tracked
  git configuration values and the values current in the workspace.`),
	RunE: cdiffFn,
}

func cdiffFn(cmd *cobra.Command, args []string) error {
	tracked, err := repoConfigs()
	if err != nil {
		return err
	}

	for name, tv := range tracked {
		// find value in config
		gv, err := workspaceValues(name)
		if err != nil {
			return err
		}

		if !valuesEqual(gv, tv) {
			normal("%s:\n  Workspace: %s\n  Tracked:%s\n", name, strings.Join(gv, ", "), strings.Join(tv, ", "))
		} else {
			verbose("%s: same\n")
		}
	}

	verbose("%d entries", len(tracked))
	return nil
}
