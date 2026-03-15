package mapcommand

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"specforge/internal/ai"
	"specforge/internal/config"
	"specforge/internal/mapper"

	"github.com/spf13/cobra"
)

// validDocs is the ordered list of the 7 GSD codebase documents.
var validDocs = []string{
	"STACK.md",
	"ARCHITECTURE.md",
	"STRUCTURE.md",
	"CONVENTIONS.md",
	"TESTING.md",
	"INTEGRATIONS.md",
	"CONCERNS.md",
}

// outputDir is the directory where documents are written. Overridden in tests.
var outputDir = filepath.Join(".", ".planning", "codebase")

// newAIClient is a factory for creating an AIClient. Tests can replace this
// to inject a mock that doesn't require ANTHROPIC_API_KEY.
var newAIClient = func(cfg config.AIConfig) generator {
	return ai.NewAIClient(cfg)
}

// generator is the interface satisfied by *ai.AIClient for document generation.
type generator interface {
	Generate(prompt string) (string, error)
}

// docEntry maps a document name to its prompt builder function.
type docEntry struct {
	name     string
	promptFn func(*mapper.CodebaseSnapshot) string
}

// dispatchTable maps each valid doc name to its corresponding prompt builder.
var dispatchTable = []docEntry{
	{"STACK.md", mapper.StackPrompt},
	{"ARCHITECTURE.md", mapper.ArchitecturePrompt},
	{"STRUCTURE.md", mapper.StructurePrompt},
	{"CONVENTIONS.md", mapper.ConventionsPrompt},
	{"TESTING.md", mapper.TestingPrompt},
	{"INTEGRATIONS.md", mapper.IntegrationsPrompt},
	{"CONCERNS.md", mapper.ConcernsPrompt},
}

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

// runMap implements the map command logic.
func runMap(cfg *config.Config, updateDoc string) error {
	// Determine which documents to generate.
	var docsToGenerate []docEntry
	if updateDoc == "" {
		docsToGenerate = dispatchTable
	} else {
		// Validate updateDoc against the list of valid document names.
		valid := false
		for _, entry := range dispatchTable {
			if entry.name == updateDoc {
				docsToGenerate = []docEntry{entry}
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid --update value %q: valid documents are: %s",
				updateDoc, strings.Join(validDocs, ", "))
		}
	}

	// Walk the project to build the snapshot.
	snap, err := mapper.WalkProject(".")
	if err != nil {
		return fmt.Errorf("walk project: %w", err)
	}

	// Create the AI client.
	aiClient := newAIClient(cfg.AI)

	// Resolve the output directory.
	outDir := outputDir

	// Launch a goroutine per document; collect errors with a mutex-protected slice.
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []string
	)

	for _, doc := range docsToGenerate {
		wg.Add(1)
		go func(d docEntry) {
			defer wg.Done()

			fmt.Fprintf(os.Stderr, "Analyzing %s...\n", d.name)

			prompt := d.promptFn(snap)

			content, genErr := aiClient.Generate(prompt)
			if genErr != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: %v", d.name, genErr))
				mu.Unlock()
				fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", d.name, genErr)
				return
			}

			writeErr := mapper.WriteDocument(outDir, d.name, content)
			if writeErr != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: %v", d.name, writeErr))
				mu.Unlock()
				fmt.Fprintf(os.Stderr, "Warning: failed to write %s: %v\n", d.name, writeErr)
				return
			}

			fmt.Fprintf(os.Stderr, "Analyzing %s... done\n", d.name)
		}(doc)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("%d document(s) failed: %s", len(errs), strings.Join(errs, "; "))
	}

	return nil
}
