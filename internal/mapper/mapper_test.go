package mapper

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// testdataSimple returns the absolute path to the testdata/simple fixture directory.
func testdataSimple() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", "simple")
}

// TestWalkProject_ExtractsPackages verifies that WalkProject returns a snapshot
// containing at least one package named "main" from testdata/simple.
func TestWalkProject_ExtractsPackages(t *testing.T) {
	snap, err := WalkProject(testdataSimple())
	if err != nil {
		t.Fatalf("WalkProject returned error: %v", err)
	}
	for _, pkg := range snap.Packages {
		if pkg.Name == "main" {
			return
		}
	}
	t.Errorf("expected at least one package named 'main', got: %v", snap.Packages)
}

// TestWalkProject_ExtractsImports verifies that the snapshot contains deduplicated
// external imports from testdata/simple (at minimum "fmt" from main.go).
func TestWalkProject_ExtractsImports(t *testing.T) {
	snap, err := WalkProject(testdataSimple())
	if err != nil {
		t.Fatalf("WalkProject returned error: %v", err)
	}
	if len(snap.Imports) == 0 {
		t.Error("expected at least one import, got none")
	}
}

// TestWalkProject_ExtractsFuncs verifies that the snapshot contains at least one
// FuncInfo extracted from testdata/simple.
func TestWalkProject_ExtractsFuncs(t *testing.T) {
	snap, err := WalkProject(testdataSimple())
	if err != nil {
		t.Fatalf("WalkProject returned error: %v", err)
	}
	if len(snap.Funcs) == 0 {
		t.Error("expected at least one FuncInfo, got none")
	}
}

// TestWalkProject_ExtractsTypes verifies that the snapshot contains at least one
// TypeInfo (Widget struct) from testdata/simple.
func TestWalkProject_ExtractsTypes(t *testing.T) {
	snap, err := WalkProject(testdataSimple())
	if err != nil {
		t.Fatalf("WalkProject returned error: %v", err)
	}
	if len(snap.Types) == 0 {
		t.Error("expected at least one TypeInfo, got none")
	}
}

// TestStackDocHeadings verifies that StackPrompt returns a string containing all
// 7 required STACK.md headings.
func TestStackDocHeadings(t *testing.T) {
	snap := &CodebaseSnapshot{
		ModuleName: "testmodule",
		GoVersion:  "1.24",
	}
	result := StackPrompt(snap)

	requiredHeadings := []string{
		"# Technology Stack",
		"## Languages",
		"## Runtime",
		"## Frameworks",
		"## Key Dependencies",
		"## Configuration",
		"## Platform Requirements",
	}
	for _, heading := range requiredHeadings {
		if !strings.Contains(result, heading) {
			t.Errorf("StackPrompt output missing required heading: %q", heading)
		}
	}
}

// TestConcernsDocSeverity verifies that ConcernsPrompt returns a string containing
// the "Severity:" label required by CONTEXT.md.
func TestConcernsDocSeverity(t *testing.T) {
	snap := &CodebaseSnapshot{}
	result := ConcernsPrompt(snap)
	if !strings.Contains(result, "Severity:") {
		t.Errorf("ConcernsPrompt output missing required 'Severity:' label, got: %q", result)
	}
}
