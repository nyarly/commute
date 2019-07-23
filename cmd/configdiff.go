package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
  cDiffCmd.Flags().StringP("remote", "r", "", "check diff for a tracked remote")
  cDiffCmd.Flags().StringP("workspace", "w", "", "check diff for a tracked workspace directory")
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
  var tracked gitconfig
  var err error
  workspace := "."

  if rem, _ := cmd.Flags().GetString("remote"); rem == "" {
    tracked, err = repoConfigs()
    if ws, _ := cmd.Flags().GetString("workspace"); ws != "" {
      workspace = ws
    }
  } else {
    tracked, err = trackedConfigs(remote(rem))

    r := remote(rem)
    path, err := r.localPath()
    if err != nil {
      return err
    }
    workspace = path
  }
	if err != nil {
		return err
	}

	for name, tv := range tracked {
		// find value in config
		gv, err := workspaceValues(workspace, name)
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
