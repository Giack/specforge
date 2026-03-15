package mapcommand

import (
	"os"
	"path/filepath"
	"strings"
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

// mockGenerator is a test double for the generator interface.
// It returns a fixed content string for every Generate call.
type mockGenerator struct {
	content string
	err     error
}

func (m *mockGenerator) Generate(_ string) (string, error) {
	return m.content, m.err
}

// TestMapCommandWritesAllDocs verifies that runMap writes all 7 GSD codebase
// documents to the output directory when given a valid project root and a mock AI.
func TestMapCommandWritesAllDocs(t *testing.T) {
	outDir := t.TempDir()
	cfg := &config.Config{}

	// Inject mock AI client so no real HTTP calls are made.
	origNewAIClient := newAIClient
	newAIClient = func(_ config.AIConfig) generator {
		return &mockGenerator{content: "# Mock content"}
	}
	defer func() { newAIClient = origNewAIClient }()

	// Inject temp output dir so files land where the test can check them.
	origOutputDir := outputDir
	outputDir = outDir
	defer func() { outputDir = origOutputDir }()

	// runMap with empty updateDoc should generate all 7 documents.
	err := runMap(cfg, "")
	if err != nil {
		t.Fatalf("runMap returned unexpected error: %v", err)
	}

	// Verify all 7 documents were written to the output directory.
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

	// The error message should mention "valid documents".
	if err != nil && !strings.Contains(err.Error(), "valid documents") {
		t.Errorf("error message %q does not contain 'valid documents'", err.Error())
	}
}
