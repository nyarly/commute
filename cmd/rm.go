package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [remote]",
	Short: "remote repo from tracking",
	Long:  longUsage(`Remove either the current project (or the provided remote) from commmute tracking.`),
	Args:  cobra.RangeArgs(0, 1),
	RunE:  rmFn,
}

func rmFn(cmd *cobra.Command, args []string) error {
	var remotes []gitRemote
	var found bool
	if len(args) == 0 {
    var err error
    remotes, err = getRepoRemotes()
		if err != nil {
			return err
		}
	} else {
		remotes = []gitRemote{{url: args[0]}}
	}

	rem, found := chooseRemote(cfg, remotes)

	if !found {
		return fmt.Errorf("no remote %q recorded in config", rem)
	}

	for i, crem := range cfg.Remotes {
		if sameRemote(crem, rem) {
			cfg.Remotes = append(cfg.Remotes[:i], cfg.Remotes[i+1:]...)
			break
		}
	}

	lp, err := rem.linkPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(lp); err != nil {
		fmt.Printf("remote not yet accounted for")
	}
	return os.Remove(lp)
}
