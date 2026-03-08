package dev

import (
	"fmt"

	"specforge/internal/config"
	"specforge/internal/vcs"

	"github.com/spf13/cobra"
)

func prCmd(cfg *config.Config) *cobra.Command {
	var repoSlug string
	var sourceBranch string
	var targetBranch string
	var title string

	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Create a Pull Request (GitHub, GitLab, or Bitbucket)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repoSlug == "" {
				return fmt.Errorf("repo-slug is required (use --repo flag)")
			}
			if sourceBranch == "" {
				return fmt.Errorf("source-branch is required (use --source flag)")
			}
			if targetBranch == "" {
				targetBranch = "main"
			}
			if title == "" {
				title = "Feature: Updated via SpecForge"
			}

			vcsClient, err := vcs.NewVCSClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create VCS client: %w", err)
			}

			if vcsClient == nil {
				return fmt.Errorf("no VCS provider configured. Please set 'vcs.provider' in config.yaml")
			}

			fmt.Printf("🔗 Creating Pull Request on %s...\n", vcsClient.GetProviderName())
			fmt.Printf("   Repo: %s\n", repoSlug)
			fmt.Printf("   Source: %s -> Target: %s\n", sourceBranch, targetBranch)
			fmt.Printf("   Title: %s\n", title)

			url, err := vcsClient.CreatePullRequest(repoSlug, title, sourceBranch, targetBranch)
			if err != nil {
				return fmt.Errorf("failed to create PR: %w", err)
			}

			fmt.Printf("\n✅ Pull Request created!\n")
			fmt.Printf("   🔗 %s\n", url)

			return nil
		},
	}

	cmd.Flags().StringVar(&repoSlug, "repo", "", "Repository slug (e.g., my-project)")
	cmd.Flags().StringVar(&sourceBranch, "source", "", "Source branch name")
	cmd.Flags().StringVar(&targetBranch, "target", "main", "Target branch (default: main)")
	cmd.Flags().StringVar(&title, "title", "", "PR title")

	return cmd
}
