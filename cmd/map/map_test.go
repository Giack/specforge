package mapcommand

import (
	"os"
	"path/filepath"
	"testing"

	"specforge/internal/config"
)

// expectedDocs is the list of 7 GSD codebase documents the map command must produce.
var expectedDocs = []string{
	"STACK.md",
	"ARCHITECTURE.md",
	"STRUCTURE.md",
	"CONVENTIONS.md",
	"TESTING.md",
	"INTEGRATIONS.md",
	"CONCERNS.md",
}

// TestMapCommandWritesAllDocs verifies that runMap writes all 7 GSD codebase
// documents to the output directory when given a valid project root and a mock AI.
func TestMapCommandWritesAllDocs(t *testing.T) {
	outDir := t.TempDir()
	cfg := &config.Config{}

	// runMap with empty updateDoc should generate all 7 documents.
	err := runMap(cfg, "")
	if err != nil {
		t.Fatalf("runMap returned unexpected error: %v", err)
	}

	// Verify all 7 documents were written to the output directory.
	// Note: outDir is not yet wired into runMap — this test will fail (RED) until
	// runMap accepts an output directory and actually writes the files.
	for _, doc := range expectedDocs {
		path := filepath.Join(outDir, doc)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			t.Errorf("expected document %q was not written to %q", doc, outDir)
		}
	}
}

// TestUpdateFlagInvalidDoc verifies that runMap returns a non-nil error when
// the updateDoc parameter is not one of the 7 valid document names.
func TestUpdateFlagInvalidDoc(t *testing.T) {
	cfg := &config.Config{}

	err := runMap(cfg, "TYPO.md")
	if err == nil {
		t.Error("expected non-nil error for invalid --update value 'TYPO.md', got nil")
	}
}
