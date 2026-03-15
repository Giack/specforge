package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"specforge/internal/mapper"
)

// AnalyzeGo runs the Go AST walker on root and converts the result to AnalysisResult.
// root is resolved to an absolute path so that WalkProject's hidden-dir guard
// does not skip the root itself when root == ".".
func AnalyzeGo(root string, excludeDirs []string) (*AnalysisResult, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	snap, err := mapper.WalkProject(absRoot)
	if err != nil {
		return nil, err
	}
	return goSnapshotToResult(snap, absRoot), nil
}

// goSnapshotToResult converts a CodebaseSnapshot to the language-agnostic AnalysisResult.
func goSnapshotToResult(snap *mapper.CodebaseSnapshot, root string) *AnalysisResult {
	result := &AnalysisResult{
		Language: "go",
		Module:   snap.ModuleName,
		Version:  snap.GoVersion,
		Files:    snap.Files,
		Errors:   []string{},
	}

	// Packages
	for _, p := range snap.Packages {
		result.Packages = append(result.Packages, Package{
			Name:  p.Name,
			Path:  p.Path,
			Files: p.Files,
		})
	}

	// Exports (exported funcs and types)
	for _, f := range snap.Funcs {
		if f.IsExported {
			result.Exports = append(result.Exports, Export{
				Package: f.Package,
				Name:    f.Name,
				Kind:    "func",
				File:    f.File,
			})
		}
	}

	// Types
	for _, t := range snap.Types {
		result.Types = append(result.Types, Type{
			Package:  t.Package,
			Name:     t.Name,
			Kind:     t.Kind,
			File:     t.File,
			Exported: t.IsExported,
		})
	}

	// Imports
	for _, imp := range snap.Imports {
		internal := snap.ModuleName != "" && strings.HasPrefix(imp.Path, snap.ModuleName)
		result.Imports = append(result.Imports, Import{
			Path:     imp.Path,
			Internal: internal,
		})
	}

	// Entry points: any main package
	for _, p := range snap.Packages {
		if p.Name == "main" {
			for _, f := range p.Files {
				if strings.HasSuffix(f, "main.go") {
					result.EntryPoints = append(result.EntryPoints, f)
				}
			}
		}
	}

	// Dependencies from go.mod
	result.Dependencies = parseGoMod(root)

	return result
}

// parseGoMod reads go.mod and extracts direct dependencies.
func parseGoMod(root string) Dependencies {
	deps := Dependencies{
		ManifestFile: "go.mod",
		DirectDeps:   []string{},
	}
	modPath := filepath.Join(root, "go.mod")
	f, err := os.Open(modPath)
	if err != nil {
		return deps
	}
	defer f.Close()

	inRequire := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if inRequire && line != "" && !strings.HasPrefix(line, "//") {
			// e.g., "github.com/spf13/cobra v1.10.2"
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				deps.DirectDeps = append(deps.DirectDeps, parts[0])
			}
		}
		// single-line require: "require github.com/foo/bar v1.0.0"
		if strings.HasPrefix(line, "require ") && !strings.HasSuffix(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps.DirectDeps = append(deps.DirectDeps, parts[1])
			}
		}
	}
	return deps
}
