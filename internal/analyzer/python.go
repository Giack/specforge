package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AnalyzePython performs a regex-based analysis of a Python codebase.
func AnalyzePython(root string, excludeDirs []string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Language: "python",
		Errors:   []string{},
	}

	// Try pyproject.toml or setup.py for module name
	if data, err := os.ReadFile(filepath.Join(root, "pyproject.toml")); err == nil {
		content := string(data)
		if m := regexp.MustCompile(`name\s*=\s*"([^"]+)"`).FindStringSubmatch(content); len(m) > 1 {
			result.Module = m[1]
		}
		if m := regexp.MustCompile(`version\s*=\s*"([^"]+)"`).FindStringSubmatch(content); len(m) > 1 {
			result.Version = m[1]
		}
	}

	// Read requirements.txt
	if f, err := os.Open(filepath.Join(root, "requirements.txt")); err == nil {
		defer f.Close()
		deps := Dependencies{ManifestFile: "requirements.txt", DirectDeps: []string{}}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				pkg := regexp.MustCompile(`^([A-Za-z0-9_.-]+)`).FindString(line)
				if pkg != "" {
					deps.DirectDeps = append(deps.DirectDeps, pkg)
				}
			}
		}
		result.Dependencies = deps
	}

	excludeSet := buildExcludeSet(excludeDirs)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "__pycache__" || name == ".venv" || name == "venv" || name == "dist" || excludeSet[name] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		result.Files = append(result.Files, rel)

		if f, err := os.Open(path); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			importRe := regexp.MustCompile(`^(?:import|from)\s+([a-zA-Z0-9_.]+)`)
			classDef := regexp.MustCompile(`^class\s+([A-Z]\w*)`)
			funcDef := regexp.MustCompile(`^def\s+([a-z]\w*)`)
			pkg := filepath.Dir(rel)
			for scanner.Scan() {
				line := scanner.Text()
				if m := importRe.FindStringSubmatch(line); len(m) > 1 {
					base := strings.Split(m[1], ".")[0]
					if base != "" {
						result.Imports = append(result.Imports, Import{Path: base, Internal: false})
					}
				}
				if m := classDef.FindStringSubmatch(line); len(m) > 1 {
					result.Types = append(result.Types, Type{Package: pkg, Name: m[1], Kind: "class", File: rel, Exported: true})
				}
				if m := funcDef.FindStringSubmatch(line); len(m) > 1 {
					result.Exports = append(result.Exports, Export{Package: pkg, Name: m[1], Kind: "func", File: rel})
				}
			}
		}
		return nil
	})
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
	}

	result.Imports = deduplicateImports(result.Imports)
	return result, nil
}
