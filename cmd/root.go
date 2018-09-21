package cmd

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/samsalisbury/yaml"
	"github.com/spf13/cobra"
)

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "commute",
	Short: "commute maintains the mapping of remote repositories to local workspaces",
	Long: longUsage(`A utility for keeping track of where repos are on you work machine,
		suitable for use by humans and scripts.`),
	PersistentPreRunE: setupStuff,
	SilenceUsage:      true,
}

func setupStuff(cmd *cobra.Command, args []string) error {
	err := setupPaths()
	if err != nil {
		return err
	}

	return loadConfig()
}

func setupPaths() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	configDir = filepath.Join(u.HomeDir, relConfigDir)
	configFile = filepath.Join(configDir, relConfigFile)
	return nil
}

func loadConfig() error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}

	return nil
}
