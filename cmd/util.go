package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/samsalisbury/yaml"
)

const (
	relConfigDir  = ".config/commute"
	relConfigFile = "config.yaml"
	docs          = `commute: transit git projects back and forth
Usage:
  commute list         Ensure that the config maps to projects
	commute add          Add the current git project
	commute rm [remote]  Remove the current project (or named remote)
	                       from the commute config.
`
)

type (
	config struct {
		Remotes remotes
	}
)

var (
	configDir    string
	configFile   string
	cfg          config
	remoteNameRE = regexp.MustCompile(`([^/:]+/[^/.]+)(?:\.git)?$`)
	fieldsRE     = regexp.MustCompile(`\s+`)
)

func lookup(start, tgt string) (string, error) {
	for from, _ := filepath.Abs(start); !(from == "" || from == "/"); from = filepath.Dir(from) {
		cb := filepath.Join(from, tgt)
		_, err := os.Lstat(cb)
		if err == nil {
			return from, nil
		}
	}
	return "", fmt.Errorf("No %s found above %s", tgt, start)
}

func (c *config) save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(b)
	return nil
}

var longFix = regexp.MustCompile(`(?m)^[ \t]*`)

func longUsage(s string) string {
	return longFix.ReplaceAllString(s, "")
}
