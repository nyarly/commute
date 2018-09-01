package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
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
		Remotes []remote
	}

	remote string
)

var (
	configDir    string
	configFile   string
	cfg          config
	remoteNameRE = regexp.MustCompile(`([^/:]+/[^/.]+)(?:\.git)?$`)
	fieldsRE     = regexp.MustCompile(`\s+`)
)

func (r *remote) name() string {
	m := remoteNameRE.FindStringSubmatch(string(*r))
	if m == nil || len(m) < 2 {
		panic("Badly formatted git remote: " + string(*r))
	}
	return m[1]
}

func (r *remote) linkPath() string {
	return filepath.Join(configDir, r.name())
}

func (r *remote) localPath() (string, error) {
	return os.Readlink(r.linkPath())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(docs)
		os.Exit(1)
	}

	err := setupPaths()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = loadConfig()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	switch os.Args[1] {
	default:
		fmt.Println(docs)
		os.Exit(1)
	case `list`:
		err = doList()
	case `add`:
		err = doAdd()
	case `rm`:
		err = doRemove(os.Args[2:])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

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
