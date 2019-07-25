package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func workspaceValues(workspace string, cfgName string) (gitvalue, error) {
	gitconfigcmd := exec.Command("git", "config", "--local", "--get-all", cfgName)
  gitconfigcmd.Dir = workspace
	value, err := gitconfigcmd.Output()
	if err != nil {
		return gitvalue{}, err
	}
	return strings.Split(string(value), "\n"), nil
}

func setValues(cfgName string, values []string) error {
  first, rest := values[0], values[1:len(values)]

	err := exec.Command("git", "config", "--local", "--replace-all", cfgName, first).Run()
  if err != nil {
    return err
  }
  for _, v := range rest {
    err := exec.Command("git", "config", "--local", "--add", cfgName, v).Run()
    if err != nil {
      return err
    }
  }

  return nil
}

func repoConfigs() (gitconfig, error) {
	// get values from commute config
	reponame, has := chooseRepoRemote(cfg)
	if !has {
		return nil, fmt.Errorf("couldn't find remote name (have you `commute add`?) (you may need to sometimes `commute cleanup`")
	}

  return trackedConfigs(reponame)
}

func trackedConfigs(reponame remote) (gitconfig, error){
	tracked, has := cfg.GitConfigs[reponame]
	if !has {
		tracked := gitconfig{}
		cfg.GitConfigs[reponame] = tracked
	}
	return tracked, nil
}

func valuesEqual(left, right gitvalue) bool {
	if len(left) != len(right) {
		return false
	}

	for n, l := range left {
		r := right[n]
		if l != r {
			return false
		}
	}
	return true
}
