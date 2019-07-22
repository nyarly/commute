package cmd

import (
	"fmt"

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
	Args: cobra.ExactArgs(1),
}

func trackFn(cmd *cobra.Command, args []string) error {
	// find value in config
  cfgName := args[0]
	values, err := workspaceValues(cfgName)
  if err != nil {
    return err
  }

	// add value to commute config
	reponame, has := chooseRepoRemote(cfg)
	if !has {
		return fmt.Errorf("couldn't find remote name (maybe needs commute add?)")
	}

	tracked, has := cfg.GitConfigs[reponame]
	if !has {
		tracked := gitconfig{}
		cfg.GitConfigs[reponame] = tracked
	}
	if _, exists := tracked[cfgName]; exists {
		return fmt.Errorf("%q is already tracked for %q", cfgName, reponame)
	}
	tracked[cfgName] = gitvalue(values)

	return nil
}
