package analyzer

import (
	"os"
	"path/filepath"
)

// DetectLanguage detects the primary language of the project at root
// by checking for language-specific manifest files.
func DetectLanguage(root string) string {
	checks := []struct {
		file string
		lang string
	}{
		{"go.mod", "go"},
		{"package.json", "ts"},
		{"pyproject.toml", "python"},
		{"setup.py", "python"},
		{"requirements.txt", "python"},
		{"build.gradle", "kotlin"},
		{"build.gradle.kts", "kotlin"},
	}
	for _, c := range checks {
		if _, err := os.Stat(filepath.Join(root, c.file)); err == nil {
			return c.lang
		}
	}
	return "go" // fallback
}

// Analyze runs the appropriate analyzer for the given language.
func Analyze(root, lang string, excludeDirs []string) (*AnalysisResult, error) {
	switch lang {
	case "go":
		return AnalyzeGo(root, excludeDirs)
	case "ts", "js", "typescript", "javascript":
		return AnalyzeTS(root, excludeDirs)
	case "python", "py":
		return AnalyzePython(root, excludeDirs)
	case "kotlin", "kt":
		return AnalyzeKotlin(root, excludeDirs)
	default:
		return AnalyzeGo(root, excludeDirs)
	}
}
