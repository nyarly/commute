package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/google/go-github/github"
	"github.com/samsalisbury/yaml"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type (
	config struct {
		Remotes    remotes
		GitConfigs map[remote]gitconfig
	}

	gitconfig map[string]gitvalue
	gitvalue  []string

	envelope struct {
		OauthToken string
		GistID     string
		Version    string
		Config     config
	}
)

var (
	configDir   string
	configFile  string
	cfgEnvelope *envelope
	cfg         *config
	gh          *github.Client
	clientOnce  = &sync.Once{}
)

func queryCommand(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[dontWriteConfig] = "query command: doesn't update config"
}

func configOblivious(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[dontWriteConfig] = "non-config command: doesn't update config"
	cmd.Annotations[dontLoadConfig] = "non-config command: doesn't read config"
}

func getClient(token string) *github.Client {
	if token == "" {
		return nil
	}
	clientOnce.Do(func() {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		gh = github.NewClient(tc)
	})

	return gh
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

	env := &envelope{}

	err = yaml.Unmarshal(b, env)
	if err != nil {
		return err
	}

	env.fetchConfigGist()

	cfg = &env.Config
	cfgEnvelope = env

	return nil
}

const configYamlFilename = "config.yaml"

func (env *envelope) fetchConfigGist() {
	if env.GistID == "" {
		return
	}
	client := getClient(env.OauthToken)
	if client == nil {
		fmt.Printf("Oauth token not set, working with local config.")
		return
	}
	ctx := context.Background()
	gist, _, err := client.Gists.Get(ctx, env.GistID)
	if err != nil {
		fmt.Printf("Using local cache becasue there was an error fetching gist: %v", err)
		return
	}

	configYaml, present := gist.Files[configYamlFilename]
	if !present {
		fmt.Printf("There is a gist for %q but it doesn't have a file named %q.", env.GistID, configYamlFilename)
		return
	}

	err = yaml.Unmarshal([]byte(configYaml.GetContent()), &env.Config)
	if err != nil {
		fmt.Printf("There was a problem parsing the YAML in Gist %q, file %q.", env.GistID, configYamlFilename)
		return
	}

	err = env.write()
	if err != nil {
		fmt.Printf("Trouble saving the config file: %q", err)
	}
}

func (env *envelope) save() error {
	env.write()
	b, err := yaml.Marshal(env.Config)
	if err != nil {
		return err
	}
	gist := &github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			configYamlFilename: github.GistFile{
				Content: github.String(string(b)),
			},
		},
	}
	client := getClient(env.OauthToken)
	if client == nil {
		fmt.Printf("Oauth token not set, working with local config.")
		return nil
	}
	ctx := context.Background()
	if env.GistID == "" {
		gist.Public = github.Bool(false)
		gist, _, err = client.Gists.Create(ctx, gist)
	} else {
		gist, _, err = client.Gists.Edit(ctx, env.GistID, gist)
	}
	if err != nil {
		return err
	}
	if gid := gist.GetID(); gid != "" {
		env.GistID = gid
	}
	env.write()

	return nil
}

func (env *envelope) write() error {
	b, err := yaml.Marshal(env)
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
