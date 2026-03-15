package mapper

// CodebaseSnapshot holds all structural information extracted from a Go project
// by WalkProject. Document generators slice this struct to build targeted prompts.
type CodebaseSnapshot struct {
	ModuleName string
	GoVersion  string
	Packages   []PackageInfo
	Imports    []ImportInfo // deduplicated external imports
	Types      []TypeInfo   // structs, interfaces
	Funcs      []FuncInfo   // exported and unexported funcs
	Files      []string     // all .go file paths relative to project root
}

// PackageInfo describes a Go package within the project.
type PackageInfo struct {
	Name  string
	Path  string   // e.g., "internal/ai"
	Files []string // .go file paths in this package
}

// FuncInfo describes a function or method declaration.
type FuncInfo struct {
	Package      string
	Name         string
	File         string
	IsExported   bool
	HasReceiver  bool
	ReceiverType string
}

// TypeInfo describes a type declaration (struct or interface).
type TypeInfo struct {
	Package    string
	Name       string
	Kind       string // "struct", "interface"
	File       string
	IsExported bool
}

// ImportInfo describes an external import used in the project.
type ImportInfo struct {
	Path string // e.g., "github.com/spf13/cobra"
}
