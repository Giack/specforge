package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AnalyzeKotlin performs a regex-based analysis of a Kotlin codebase.
func AnalyzeKotlin(root string, excludeDirs []string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Language: "kotlin",
		Errors:   []string{},
	}

	// Try build.gradle.kts or build.gradle for version
	for _, manifest := range []string{"build.gradle.kts", "build.gradle"} {
		if data, err := os.ReadFile(filepath.Join(root, manifest)); err == nil {
			content := string(data)
			if m := regexp.MustCompile(`version\s*=\s*"([^"]+)"`).FindStringSubmatch(content); len(m) > 1 {
				result.Version = m[1]
			}
			result.Dependencies = Dependencies{ManifestFile: manifest, DirectDeps: []string{}}
			depRe := regexp.MustCompile(`implementation\s*[\("']([^"')]+)`)
			for _, dep := range depRe.FindAllStringSubmatch(content, -1) {
				result.Dependencies.DirectDeps = append(result.Dependencies.DirectDeps, dep[1])
			}
			break
		}
	}

	excludeSet := buildExcludeSet(excludeDirs)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "build" || name == ".gradle" || excludeSet[name] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".kt") && !strings.HasSuffix(path, ".kts") {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		result.Files = append(result.Files, rel)

		if f, err := os.Open(path); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			importRe := regexp.MustCompile(`^import\s+([a-zA-Z0-9_.]+)`)
			classDef := regexp.MustCompile(`^(?:class|data class|object|interface|enum class)\s+(\w+)`)
			funcDef := regexp.MustCompile(`^fun\s+(\w+)`)
			pkg := filepath.Dir(rel)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if m := importRe.FindStringSubmatch(line); len(m) > 1 {
					result.Imports = append(result.Imports, Import{Path: m[1], Internal: false})
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
