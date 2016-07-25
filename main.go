package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/samsalisbury/yaml"
)

const (
	configDir  = "~/.config/commute"
	configFile = "~/.config/commute/config.yaml"
	docs       = `commute: transit git projects back and forth
Usage:
  commute list   Ensure that the config maps to projects
	commute add    Add the current git project
`
)

type (
	config struct {
		Remotes []remote
	}

	remote string
)

var remoteNameRE = regexp.MustCompile(`[^/]+/[^/]+\.git$`)

func (r *remote) name() string {
	return remoteNameRE.FindString(string(*r))
}

func (r *remote) linkPath() string {
	return filepath.Join(configDir, r.name())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(docs)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	default:
		fmt.Println(docs)
		os.Exit(1)
	case `list`:
		err = doList()
	case `add`:
		err = doAdd()
	}

	if err != nil {
		fmt.Print(err)
	}

}

func doList() error {
	var cfg config
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

	for _, remote := range cfg.Remotes {
		st, err := os.Stat(remote.linkPath())
		fmt.Printf("%# v %T %v %T", st, st, err, err)
	}
	return nil
}

func doAdd() error {
	return nil
}
