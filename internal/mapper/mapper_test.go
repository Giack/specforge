package mapper

import (
	"path/filepath"
	"runtime"
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
