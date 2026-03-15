package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AnalyzeTS performs a regex-based analysis of a TypeScript/JavaScript codebase.
func AnalyzeTS(root string, excludeDirs []string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Language: "ts",
		Errors:   []string{},
	}

	// Read package.json for module name and version
	if data, err := os.ReadFile(filepath.Join(root, "package.json")); err == nil {
		pkgJSON := string(data)
		if m := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`).FindStringSubmatch(pkgJSON); len(m) > 1 {
			result.Module = m[1]
		}
		if m := regexp.MustCompile(`"version"\s*:\s*"([^"]+)"`).FindStringSubmatch(pkgJSON); len(m) > 1 {
			result.Version = m[1]
		}
		result.Dependencies = parsePkgJSON(pkgJSON)
	}

	excludeSet := buildExcludeSet(excludeDirs)

	// Walk and find .ts/.js files
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == ".git" || name == "dist" || name == "build" || excludeSet[name] {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".ts" && ext != ".tsx" && ext != ".js" && ext != ".jsx" && ext != ".mjs" {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		result.Files = append(result.Files, rel)

		// Scan for exports and imports
		if f, err := os.Open(path); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			importRe := regexp.MustCompile(`^import\s+.*from\s+['"]([^'"]+)['"]`)
			exportFuncRe := regexp.MustCompile(`^export\s+(async\s+)?function\s+(\w+)`)
			exportClassRe := regexp.MustCompile(`^export\s+(default\s+)?class\s+(\w+)`)
			exportInterfaceRe := regexp.MustCompile(`^export\s+interface\s+(\w+)`)
			pkg := filepath.Dir(rel)
			for scanner.Scan() {
				line := scanner.Text()
				if m := importRe.FindStringSubmatch(line); len(m) > 1 {
					if !strings.HasPrefix(m[1], ".") {
						result.Imports = append(result.Imports, Import{Path: m[1], Internal: false})
					}
				}
				if m := exportFuncRe.FindStringSubmatch(line); len(m) > 2 {
					result.Exports = append(result.Exports, Export{Package: pkg, Name: m[2], Kind: "func", File: rel})
				}
				if m := exportClassRe.FindStringSubmatch(line); len(m) > 2 {
					name := m[2]
					if m[1] == "default " {
						name = m[2]
					}
					result.Exports = append(result.Exports, Export{Package: pkg, Name: name, Kind: "type", File: rel})
					result.Types = append(result.Types, Type{Package: pkg, Name: name, Kind: "class", File: rel, Exported: true})
				}
				if m := exportInterfaceRe.FindStringSubmatch(line); len(m) > 1 {
					result.Types = append(result.Types, Type{Package: pkg, Name: m[1], Kind: "interface", File: rel, Exported: true})
				}
			}
		}
		return nil
	})
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
	}

	// Dedup imports
	result.Imports = deduplicateImports(result.Imports)

	return result, nil
}

func parsePkgJSON(content string) Dependencies {
	deps := Dependencies{ManifestFile: "package.json", DirectDeps: []string{}}
	inDeps := false
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"dependencies"`) || strings.Contains(line, `"devDependencies"`) {
			inDeps = true
			continue
		}
		if inDeps && line == "}" {
			inDeps = false
			continue
		}
		if inDeps {
			re := regexp.MustCompile(`"(@?[^"]+)"\s*:`)
			if m := re.FindStringSubmatch(line); len(m) > 1 {
				deps.DirectDeps = append(deps.DirectDeps, m[1])
			}
		}
	}
	return deps
}

func buildExcludeSet(dirs []string) map[string]bool {
	s := make(map[string]bool)
	for _, d := range dirs {
		s[d] = true
	}
	return s
}

func deduplicateImports(imports []Import) []Import {
	seen := make(map[string]bool)
	var out []Import
	for _, imp := range imports {
		if !seen[imp.Path] {
			seen[imp.Path] = true
			out = append(out, imp)
		}
	}
	return out
}
