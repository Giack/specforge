package analyzer

// AnalysisResult is the top-level JSON output of specforge analyze.
// It is a language-agnostic superset of CodebaseSnapshot.
type AnalysisResult struct {
	Language     string       `json:"language"`
	Module       string       `json:"module"`
	Version      string       `json:"version"`
	Packages     []Package    `json:"packages"`
	Exports      []Export     `json:"exports"`
	Types        []Type       `json:"types"`
	Imports      []Import     `json:"imports"`
	Files        []string     `json:"files"`
	EntryPoints  []string     `json:"entry_points"`
	Dependencies Dependencies `json:"dependencies"`
	Errors       []string     `json:"errors"`
}

// Package represents a named group of source files (Go package, JS module, Python module).
type Package struct {
	Name  string   `json:"name"`
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

// Export represents an exported symbol.
type Export struct {
	Package string `json:"package"`
	Name    string `json:"name"`
	Kind    string `json:"kind"` // "func", "type", "var", "const"
	File    string `json:"file"`
}

// Type represents a type declaration.
type Type struct {
	Package  string `json:"package"`
	Name     string `json:"name"`
	Kind     string `json:"kind"` // "struct", "interface", "alias"
	File     string `json:"file"`
	Exported bool   `json:"exported"`
}

// Import represents an external import used in the project.
type Import struct {
	Path     string `json:"path"`
	Internal bool   `json:"internal"` // true if it's the project's own module
}

// Dependencies holds dependency manifest information.
type Dependencies struct {
	ManifestFile string   `json:"manifest_file"`
	DirectDeps   []string `json:"direct_deps"`
}
