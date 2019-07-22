package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(cInflictCmd)
  queryCommand(cInflictCmd)
}

var cInflictCmd = &cobra.Command{
	Use:   "inflict",
	Short: "push tracked values to the workspace",
	Long: longUsage(`Copy the values of tracked keys from commute's config
  into the workspace's git config.`),
	RunE: cInflictFn,
}


func cInflictFn(cmd *cobra.Command, args []string) error {
	tracked, err := repoConfigs()
	if err != nil {
		return err
	}

	for name, tv := range tracked {
    err := setValues(name, tv)
    if err != nil {
      return err
    }
	}

  return nil
}
