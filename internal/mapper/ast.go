package mapper

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// WalkProject analyzes the Go project rooted at root and returns a CodebaseSnapshot
// containing all extracted structural information.
func WalkProject(root string) (*CodebaseSnapshot, error) {
	snap := &CodebaseSnapshot{}

	// Step 1: Extract ModuleName and GoVersion from go.mod via text scan.
	modPath := filepath.Join(root, "go.mod")
	if f, err := os.Open(modPath); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "module ") {
				snap.ModuleName = strings.TrimPrefix(line, "module ")
			}
			if strings.HasPrefix(line, "go ") {
				snap.GoVersion = strings.TrimPrefix(line, "go ")
			}
		}
		f.Close()
	}

	// pkgMap deduplicates packages by their directory path.
	pkgMap := make(map[string]*PackageInfo)
	// importSeen deduplicates imports by import path.
	importSeen := make(map[string]struct{})

	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor, .git, and hidden directories.
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files; skip _test.go files.
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Get relative path from root.
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("rel path: %w", err)
		}
		snap.Files = append(snap.Files, relPath)

		// Parse the file.
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			// Non-fatal: skip unparseable files.
			return nil
		}

		// Step 4: Extract package info.
		pkgName := f.Name.Name
		pkgDir := filepath.Dir(relPath)
		if pkgDir == "." {
			pkgDir = ""
		}
		if _, ok := pkgMap[pkgDir]; !ok {
			pkgMap[pkgDir] = &PackageInfo{
				Name: pkgName,
				Path: pkgDir,
			}
		}
		pkgMap[pkgDir].Files = append(pkgMap[pkgDir].Files, relPath)

		// Step 5: Extract imports.
		for _, imp := range f.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if _, seen := importSeen[importPath]; !seen {
				importSeen[importPath] = struct{}{}
				snap.Imports = append(snap.Imports, ImportInfo{
					Path: importPath,
				})
			}
		}

		// Step 6: Extract funcs.
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			fi := FuncInfo{
				Package:    pkgName,
				Name:       fn.Name.Name,
				File:       relPath,
				IsExported: fn.Name.IsExported(),
				HasReceiver: fn.Recv != nil,
			}
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				fi.ReceiverType = fmt.Sprintf("%v", fn.Recv.List[0].Type)
			}
			snap.Funcs = append(snap.Funcs, fi)
			return true
		})

		// Step 7: Extract types.
		ast.Inspect(f, func(n ast.Node) bool {
			gd, ok := n.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				return true
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				kind := ""
				switch ts.Type.(type) {
				case *ast.StructType:
					kind = "struct"
				case *ast.InterfaceType:
					kind = "interface"
				}
				snap.Types = append(snap.Types, TypeInfo{
					Package:    pkgName,
					Name:       ts.Name.Name,
					Kind:       kind,
					File:       relPath,
					IsExported: ts.Name.IsExported(),
				})
			}
			return true
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk project: %w", err)
	}

	// Convert pkgMap to slice.
	for _, pkg := range pkgMap {
		snap.Packages = append(snap.Packages, *pkg)
	}

	return snap, nil
}
