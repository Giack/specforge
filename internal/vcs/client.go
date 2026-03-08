package vcs

import "specforge/internal/config"

type PRProvider interface {
	CreatePullRequest(repoSlug, title, sourceBranch, targetBranch string) (string, error)
	GetProviderName() string
}

func NewVCSClient(cfg *config.Config) (PRProvider, error) {
	switch cfg.VCS.Provider {
	case "github":
		return NewGitHubClient(cfg.VCS.GitHub), nil
	case "gitlab":
		return NewGitLabClient(cfg.VCS.GitLab), nil
	case "bitbucket":
		return NewBitbucketClient(cfg.VCS.Bitbucket), nil
	default:
		return nil, nil
	}
}
