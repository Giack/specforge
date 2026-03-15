package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"specforge/internal/analyzer"

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
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	lang, _ := cmd.Flags().GetString("lang")
	format, _ := cmd.Flags().GetString("format")
	excludeRaw, _ := cmd.Flags().GetString("exclude")

	var excludeDirs []string
	if excludeRaw != "" {
		for _, d := range strings.Split(excludeRaw, ",") {
			d = strings.TrimSpace(d)
			if d != "" {
				excludeDirs = append(excludeDirs, d)
			}
		}
	}

	if lang == "auto" {
		lang = analyzer.DetectLanguage(root)
		fmt.Fprintf(os.Stderr, "detected language: %s\n", lang)
	}

	result, err := analyzer.Analyze(root, lang, excludeDirs)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	case "markdown":
		fmt.Fprintf(os.Stdout, "# Codebase Analysis\n\n")
		fmt.Fprintf(os.Stdout, "- Language: %s\n", result.Language)
		fmt.Fprintf(os.Stdout, "- Module: %s\n", result.Module)
		fmt.Fprintf(os.Stdout, "- Packages: %d\n", len(result.Packages))
		fmt.Fprintf(os.Stdout, "- Types: %d\n", len(result.Types))
		fmt.Fprintf(os.Stdout, "- Exports: %d\n", len(result.Exports))
		return nil
	default:
		return fmt.Errorf("unknown format %q: use json or markdown", format)
	}
}
