package mapcommand

import (
	"specforge/internal/config"

	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for "specforge map".
func NewCommand(cfg *config.Config) *cobra.Command {
	var updateDoc string

	cmd := &cobra.Command{
		Use:   "map",
		Short: "Analyze a Go project and produce GSD codebase documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMap(cfg, updateDoc)
		},
	}

	cmd.Flags().StringVar(&updateDoc, "update", "", "Regenerate only this document (e.g., CONCERNS.md)")
	return cmd
}

// runMap is the stub implementation — not yet implemented.
func runMap(cfg *config.Config, updateDoc string) error {
	return nil
}
