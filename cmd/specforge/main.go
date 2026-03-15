package main

import (
	"fmt"
	"os"

	"specforge/cmd/dev"
	"specforge/cmd/em"
	mapcmd "specforge/cmd/map"
	"specforge/cmd/pm"
	"specforge/internal/config"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	cfg := config.Load()

	rootCmd := &cobra.Command{
		Use:   "specforge",
		Short: "Enterprise Spec-Driven Development Tool",
		Long: `SpecForge bridges business requirements (Confluence/Jira) 
with AI execution engines (Claude Code/OpenCode).

Workflow:
  pm    - Product Manager: Sync requirements from Confluence/Jira
  em    - Engineering Manager: Architect roadmaps & verify
  dev   - Developer: Execute specs & open PRs`,
		Version: version,
	}

	rootCmd.AddCommand(pm.NewCommand(cfg))
	rootCmd.AddCommand(em.NewCommand(cfg))
	rootCmd.AddCommand(dev.NewCommand(cfg))
	rootCmd.AddCommand(mapcmd.NewCommand(cfg))
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
