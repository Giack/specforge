package main

import (
	"fmt"
	"os"

	analyzecmd "specforge/cmd/analyze"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:     "specforge",
		Short:   "Codebase structure analyzer for Claude Code plugin",
		Long:    `specforge analyze — outputs structured JSON of a codebase's packages, types, funcs, and imports.`,
		Version: version,
	}

	rootCmd.AddCommand(analyzecmd.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
