package mapper

import (
	"fmt"
	"strings"
)

// StackPrompt builds the prompt for generating STACK.md from the given snapshot.
// It uses only snap.ModuleName, snap.GoVersion, and snap.Imports (not full snapshot).
func StackPrompt(snap *CodebaseSnapshot) string {
	var imports []string
	for _, imp := range snap.Imports {
		imports = append(imports, imp.Path)
	}
	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# Technology Stack
## Languages
## Runtime
## Frameworks
## Key Dependencies
## Configuration
## Platform Requirements

Data:
Module: %s
Go version: %s
External imports: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		snap.ModuleName,
		snap.GoVersion,
		strings.Join(imports, "\n"),
	)
}

// ArchitecturePrompt builds the prompt for generating ARCHITECTURE.md.
// It uses snap.Packages, snap.Funcs (exported, with receivers), and snap.Types.
func ArchitecturePrompt(snap *CodebaseSnapshot) string {
	var packages []string
	for _, pkg := range snap.Packages {
		packages = append(packages, fmt.Sprintf("%s (%s)", pkg.Name, pkg.Path))
	}

	var exportedFuncs []string
	for _, fn := range snap.Funcs {
		if fn.IsExported {
			if fn.HasReceiver {
				exportedFuncs = append(exportedFuncs, fmt.Sprintf("(%s).%s in %s", fn.ReceiverType, fn.Name, fn.Package))
			} else {
				exportedFuncs = append(exportedFuncs, fmt.Sprintf("%s.%s", fn.Package, fn.Name))
			}
		}
	}

	var types []string
	for _, t := range snap.Types {
		types = append(types, fmt.Sprintf("%s.%s (%s)", t.Package, t.Name, t.Kind))
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# ARCHITECTURE.md
## Pattern
## Entry Points
## Command Groups (Layers)
## Data Flow
## Abstractions
## Configuration
## Key Files

Data:
Packages: %s
Exported functions and methods: %s
Types: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		strings.Join(packages, "\n"),
		strings.Join(exportedFuncs, "\n"),
		strings.Join(types, "\n"),
	)
}

// StructurePrompt builds the prompt for generating STRUCTURE.md.
// It uses snap.Files, snap.Packages, and snap.ModuleName.
func StructurePrompt(snap *CodebaseSnapshot) string {
	var packages []string
	for _, pkg := range snap.Packages {
		packages = append(packages, fmt.Sprintf("%s (%s)", pkg.Name, pkg.Path))
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# STRUCTURE.md
## Directory Layout
## Key Locations
## Naming Conventions
## Go Module

Data:
Module: %s
Files: %s
Packages: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		snap.ModuleName,
		strings.Join(snap.Files, "\n"),
		strings.Join(packages, "\n"),
	)
}

// ConventionsPrompt builds the prompt for generating CONVENTIONS.md.
// It uses snap.Funcs (all), snap.Types (all), and snap.Files.
func ConventionsPrompt(snap *CodebaseSnapshot) string {
	var funcs []string
	for _, fn := range snap.Funcs {
		if fn.HasReceiver {
			funcs = append(funcs, fmt.Sprintf("(%s).%s [exported=%v] in %s", fn.ReceiverType, fn.Name, fn.IsExported, fn.Package))
		} else {
			funcs = append(funcs, fmt.Sprintf("%s.%s [exported=%v]", fn.Package, fn.Name, fn.IsExported))
		}
	}

	var types []string
	for _, t := range snap.Types {
		types = append(types, fmt.Sprintf("%s.%s (%s) [exported=%v]", t.Package, t.Name, t.Kind, t.IsExported))
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# Coding Conventions
## Naming Patterns
## Code Style
## Import Organization
## Error Handling
## Logging
## Comments
## Function Design
## Module Design
## Struct Initialization

Data:
Functions: %s
Types: %s
Files: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		strings.Join(funcs, "\n"),
		strings.Join(types, "\n"),
		strings.Join(snap.Files, "\n"),
	)
}

// TestingPrompt builds the prompt for generating TESTING.md.
// It uses only test files (names ending in _test.go) and FuncInfo entries from test files.
func TestingPrompt(snap *CodebaseSnapshot) string {
	var testFiles []string
	for _, f := range snap.Files {
		if strings.HasSuffix(f, "_test.go") {
			testFiles = append(testFiles, f)
		}
	}

	var testFuncs []string
	for _, fn := range snap.Funcs {
		if strings.HasSuffix(fn.File, "_test.go") {
			testFuncs = append(testFuncs, fmt.Sprintf("%s.%s", fn.Package, fn.Name))
		}
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# TESTING.md
## Status
## Test Infrastructure
## Only Testing Artifact / Test Patterns
## Testable Seams
## What Testing Would Look Like
## Recommended Test Framework
## Coverage

Data:
Test files: %s
Test functions: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		strings.Join(testFiles, "\n"),
		strings.Join(testFuncs, "\n"),
	)
}

// IntegrationsPrompt builds the prompt for generating INTEGRATIONS.md.
// It uses only snap.Imports where IsStdlib=false (non-stdlib external imports).
func IntegrationsPrompt(snap *CodebaseSnapshot) string {
	var externalImports []string
	for _, imp := range snap.Imports {
		// All ImportInfo entries in the snapshot are already deduplicated external imports.
		// The snapshot only stores non-stdlib imports per WalkProject implementation.
		externalImports = append(externalImports, imp.Path)
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# External Integrations
## APIs & External Services
## Data Storage
## Authentication & Identity
## Monitoring & Observability
## CI/CD & Deployment
## Environment Configuration
## Webhooks & Callbacks

Data:
External imports: %s

Fill each section based on the data above. Use concise key-value or bullet format.`,
		strings.Join(externalImports, "\n"),
	)
}

// ConcernsPrompt builds the prompt for generating CONCERNS.md.
// It uses snap.Packages, snap.Types, and snap.Funcs.
// The template instructs Claude to include Severity: High/Medium/Low per concern entry.
func ConcernsPrompt(snap *CodebaseSnapshot) string {
	var packages []string
	for _, pkg := range snap.Packages {
		packages = append(packages, fmt.Sprintf("%s (%s)", pkg.Name, pkg.Path))
	}

	var types []string
	for _, t := range snap.Types {
		types = append(types, fmt.Sprintf("%s.%s (%s) [exported=%v]", t.Package, t.Name, t.Kind, t.IsExported))
	}

	var funcs []string
	for _, fn := range snap.Funcs {
		if fn.HasReceiver {
			funcs = append(funcs, fmt.Sprintf("(%s).%s [exported=%v] in %s", fn.ReceiverType, fn.Name, fn.IsExported, fn.Package))
		} else {
			funcs = append(funcs, fmt.Sprintf("%s.%s [exported=%v]", fn.Package, fn.Name, fn.IsExported))
		}
	}

	return fmt.Sprintf(`You are analyzing a Go project to produce a GSD codebase document.

Output ONLY the following document with these exact headings filled in.
Do NOT add extra headings. Do NOT include explanations outside the document.

# CONCERNS.md
## Critical Bugs
## Security Issues
## Performance Issues
## Technical Debt
## Fragile Areas

For each concern, include: **Severity:** High/Medium/Low

Data:
Packages: %s
Types: %s
Functions: %s

Fill each section based on the data above. For each issue identified, rate it with Severity: High, Medium, or Low.`,
		strings.Join(packages, "\n"),
		strings.Join(types, "\n"),
		strings.Join(funcs, "\n"),
	)
}
