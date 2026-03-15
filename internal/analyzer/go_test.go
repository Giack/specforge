package analyzer

import (
	"path/filepath"
	"runtime"
	"testing"
)

// testdataSimple returns the path to mapper's testdata/simple fixture.
func testdataSimple() string {
	_, file, _, _ := runtime.Caller(0)
	// Navigate from internal/analyzer/ up two levels to internal/, then to mapper/testdata/simple
	return filepath.Join(filepath.Dir(file), "..", "mapper", "testdata", "simple")
}

func TestAnalyzeGo_BasicResult(t *testing.T) {
	result, err := AnalyzeGo(testdataSimple(), nil)
	if err != nil {
		t.Fatalf("AnalyzeGo returned error: %v", err)
	}
	if result.Language != "go" {
		t.Errorf("expected language=go, got %q", result.Language)
	}
	if len(result.Packages) == 0 {
		t.Error("expected at least one package")
	}
	if result.Errors == nil {
		t.Error("errors field should be non-nil (empty slice)")
	}
}

func TestAnalyzeGo_HasTypes(t *testing.T) {
	result, err := AnalyzeGo(testdataSimple(), nil)
	if err != nil {
		t.Fatalf("AnalyzeGo returned error: %v", err)
	}
	if len(result.Types) == 0 {
		t.Error("expected at least one type (Widget struct in testdata/simple)")
	}
}

func TestAnalyzeGo_HasImports(t *testing.T) {
	result, err := AnalyzeGo(testdataSimple(), nil)
	if err != nil {
		t.Fatalf("AnalyzeGo returned error: %v", err)
	}
	if len(result.Imports) == 0 {
		t.Error("expected at least one import")
	}
}
