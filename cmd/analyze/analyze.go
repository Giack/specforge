package analyze

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCommand returns the analyze cobra command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze [path]",
		Short: "Analyze codebase structure and output JSON",
		Long:  `Walks a codebase and outputs structured JSON of packages, types, funcs, and imports.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAnalyze,
	}
	cmd.Flags().String("lang", "auto", "Language: go|ts|kotlin|python|auto")
	cmd.Flags().String("format", "json", "Output format: json|markdown")
	cmd.Flags().String("exclude", "", "Comma-separated directories to exclude")
	return cmd
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(os.Stderr, "analyze: not yet implemented — coming in next phase")
	return nil
}
