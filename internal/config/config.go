package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Atlassian AtlassianConfig `mapstructure:"atlassian"`
	VCS       VCSConfig       `mapstructure:"vcs"`
	AI        AIConfig        `mapstructure:"ai"`
	SpecRepo  SpecRepoConfig  `mapstructure:"spec_repo"`
}

type AtlassianConfig struct {
	Domain     string `mapstructure:"domain"`
	Email      string `mapstructure:"email"`
	APIToken   string `mapstructure:"api_token"`
	ProjectKey string `mapstructure:"project_key"`
}

type VCSConfig struct {
	Provider  string          `mapstructure:"provider"` // github, gitlab, bitbucket
	GitHub    GitHubConfig    `mapstructure:"github"`
	GitLab    GitLabConfig    `mapstructure:"gitlab"`
	Bitbucket BitbucketConfig `mapstructure:"bitbucket"`
}

type GitHubConfig struct {
	Token string `mapstructure:"token"`
	Owner string `mapstructure:"owner"`
	Repo  string `mapstructure:"repo"`
}

type GitLabConfig struct {
	Domain    string `mapstructure:"domain"`
	Token     string `mapstructure:"token"`
	Group     string `mapstructure:"group"`
	ProjectID string `mapstructure:"project_id"`
}

type BitbucketConfig struct {
	Domain      string `mapstructure:"domain"`
	Username    string `mapstructure:"username"`
	AppPassword string `mapstructure:"app_password"`
	Workspace   string `mapstructure:"workspace"`
}

type AIConfig struct {
	Provider  string `mapstructure:"provider"` // claude, opencode
	Model     string `mapstructure:"model"`
	MaxTokens int    `mapstructure:"max_tokens"`
}

type SpecRepoConfig struct {
	URL    string `mapstructure:"url"`
	Branch string `mapstructure:"branch"`
}

func Load() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not find home directory: %v\n", err)
	}

	viper.AddConfigPath(home + "/.specforge")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("atlassian.domain", "")
	viper.SetDefault("vcs.provider", "github")
	viper.SetDefault("vcs.bitbucket.domain", "bitbucket.org")
	viper.SetDefault("ai.provider", "claude")
	viper.SetDefault("ai.model", "sonnet-4-20250514")
	viper.SetDefault("ai.max_tokens", 4000)
	viper.SetDefault("spec_repo.branch", "main")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
	}

	return &cfg
}
