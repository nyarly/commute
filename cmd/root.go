package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "print nothing about what commute is doing")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print out a bunch of junk about what commute is doing")
	rootCmd.PersistentFlags().BoolP("dryrun", "d", false, "take no action")
}

var (
	rootCmd = &cobra.Command{
		Use:   "commute",
		Short: "commute maintains the mapping of remote repositories to local workspaces",
		Long: longUsage(`A utility for keeping track of where repos are on you work machine,
		suitable for use by humans and scripts.`),
		PersistentPreRunE:  setupStuff,
		PersistentPostRunE: saveConfig,
		SilenceUsage:       true,
	}
	normal   = func(fmt string, as ...interface{}) {}
	verbose  = func(fmt string, as ...interface{}) {}
	actually = func(f func() error) error { return f() }
)

func needsConfig(cmd *cobra.Command) bool {
  if (cmd.Use[0:4] == "help") {
    return false
  }
  if _, stop := cmd.Annotations[dontLoadConfig]; stop{
    return false
  }
  return true
}

func changesConfig(cmd *cobra.Command) bool {
  if (cmd.Use[0:4] == "help") {
    return false
  }
  if _, stop := cmd.Annotations[dontWriteConfig]; stop{
    return false
  }
  return true
}

func setupStuff(cmd *cobra.Command, args []string) error {
	err := setupBasics(cmd)
	if err != nil {
		return err
	}

	err = setupPaths()
	if err != nil {
		return err
	}

  if !needsConfig(cmd) {
    return nil
  }

	return loadConfig()
}

func setupBasics(cmd *cobra.Command) error {
	verbFlag, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	quietFlag, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		return err
	}
	dryFlag, err:= cmd.Flags().GetBool("dryrun")
	if err != nil {
		return err
	}

	if verbFlag && quietFlag {
		return fmt.Errorf("can't specify both quiet and verbose")
	}

	if !quietFlag {
		normal = func(tmpl string, as ...interface{}) {
			log.Printf(tmpl, as...)
		}
	}

	if verbFlag {
		log.SetFlags(log.Ltime | log.Lshortfile)
		verbose = func(tmpl string, as ...interface{}) {
			log.Printf(tmpl, as...)
		}
	} else {
		log.SetFlags(0)
	}

	if dryFlag {
		actually = func(func() error) error { return nil }
	}
	return nil
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

func saveConfig(cmd *cobra.Command, args []string) error {
  if !changesConfig(cmd) {
    return nil
  }

	return cfgEnvelope.save()
}
