---
phase: 02-codebase-mapping
plan: "01"
subsystem: testing
tags: [go, ast, mapper, tdd, scaffolding, codebase-mapping]

# Dependency graph
requires: []
provides:
  - "internal/mapper package with stub types (CodebaseSnapshot, PackageInfo, FuncInfo, TypeInfo, ImportInfo)"
  - "WalkProject stub in ast.go"
  - "7 prompt builder stubs in prompts.go"
  - "WriteDocument stub in writer.go"
  - "8 RED failing tests covering MAP-01 through MAP-04"
  - "testdata/simple Go fixture (Widget struct, NewWidget func, fmt import)"
  - "cmd/map/map.go stub Cobra command"
affects:
  - 02-02-ast-walker
  - 02-03-prompt-builders
  - 02-04-map-command

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "TDD RED scaffold: stubs compile, tests fail at runtime (not compile-fail)"
    - "testdata fixture pattern: minimal Go module at internal/mapper/testdata/simple/"
    - "White-box test pattern: map_test.go is package mapcommand (same package as cmd/map/map.go)"

key-files:
  created:
    - internal/mapper/snapshot.go
    - internal/mapper/ast.go
    - internal/mapper/prompts.go
    - internal/mapper/writer.go
    - internal/mapper/mapper_test.go
    - internal/mapper/testdata/simple/main.go
    - internal/mapper/testdata/simple/go.mod
    - cmd/map/map.go
    - cmd/map/map_test.go
  modified: []

key-decisions:
  - "testdata/simple uses go.mod module 'testfixture' with go 1.24 to match actual project Go version (not 1.26 as in research)"
  - "cmd/map package named 'mapcommand' to avoid collision with 'map' builtin"
  - "map_test.go is white-box (package mapcommand) to access unexported runMap — consistent with Phase 1 test pattern"
  - "WalkProject stub returns (nil, error) to make all 4 WalkProject tests fail at runtime"
  - "TestMapCommandWritesAllDocs uses outDir from t.TempDir() but runMap does not yet accept it — documents fail because files were never written"

patterns-established:
  - "TDD RED: stubs return zero values or errors; tests fail at test level, not compile level"
  - "Fixture module: testdata/simple is an isolated Go module (own go.mod) so WalkProject can treat it as a standalone project root"

requirements-completed:
  - MAP-01
  - MAP-02
  - MAP-03
  - MAP-04

# Metrics
duration: 2min
completed: 2026-03-15
---

# Phase 2 Plan 01: Wave 0 Test Scaffold Summary

**8 RED failing tests and stub internal/mapper package establishing CodebaseSnapshot contracts for MAP-01 through MAP-04**

## Performance

- **Duration:** 2 min
- **Started:** 2026-03-15T16:34:40Z
- **Completed:** 2026-03-15T16:36:57Z
- **Tasks:** 1 (single TDD RED task)
- **Files modified:** 9 created

## Accomplishments

- Created `internal/mapper` package with all 5 contract types (CodebaseSnapshot, PackageInfo, FuncInfo, TypeInfo, ImportInfo) matching the research spec exactly
- Wrote 8 failing tests across `mapper_test.go` (6 tests) and `cmd/map/map_test.go` (2 tests) — all fail at runtime, none fail at compile time
- Created `testdata/simple` Go fixture module with Widget struct, NewWidget func, and fmt import for deterministic AST tests
- `go build ./...` succeeds; `go test ./internal/mapper/... ./cmd/map/...` exits non-zero with FAIL

## Task Commits

Each task was committed atomically:

1. **Task 1: Wave 0 TDD RED scaffold** - `4fdabd2` (test)

## Files Created/Modified

- `internal/mapper/snapshot.go` - CodebaseSnapshot, PackageInfo, FuncInfo, TypeInfo, ImportInfo struct definitions (full contracts, not stubs)
- `internal/mapper/ast.go` - WalkProject stub returning (nil, fmt.Errorf("not implemented"))
- `internal/mapper/prompts.go` - 7 prompt builder stubs each returning ""
- `internal/mapper/writer.go` - WriteDocument stub returning nil
- `internal/mapper/mapper_test.go` - 6 RED unit tests (4 WalkProject, TestStackDocHeadings, TestConcernsDocSeverity)
- `internal/mapper/testdata/simple/main.go` - Minimal Go fixture: Widget struct, NewWidget func, fmt import
- `internal/mapper/testdata/simple/go.mod` - Module "testfixture" go 1.24
- `cmd/map/map.go` - Cobra stub command (NewCommand + runMap stub)
- `cmd/map/map_test.go` - 2 RED integration tests (TestMapCommandWritesAllDocs, TestUpdateFlagInvalidDoc)

## Decisions Made

- Used `package mapcommand` for `cmd/map/` to avoid collision with Go's built-in `map` keyword in package names
- `map_test.go` is white-box (same package) to access unexported `runMap` — consistent with Phase 1 pattern for testing unexported functions
- testdata/simple uses `go 1.24` (matching the actual project's go.mod) not `go 1.26` as stated in research notes

## Deviations from Plan

None — plan executed exactly as written. One minor note: the plan specified `go 1.26` for the testdata fixture, but the project's actual `go.mod` specifies `go 1.24.0`, so `go 1.24` was used in the fixture to avoid any future toolchain compatibility surprises.

## Issues Encountered

- The `go` binary at `/usr/local/go/bin/go` (version 1.24.0) had GOROOT pointing to the mise Go 1.26.0 installation, causing "no such tool compile" errors. Used `/Users/gsortino/.local/share/mise/installs/go/1.26.0/bin/go` explicitly for build and test verification.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- All 8 tests are in RED state and will guide the implementations in plans 02-02 through 02-04
- Plan 02-02 (AST walker) implements `WalkProject` against the 4 WalkProject tests
- Plan 02-03 (prompt builders) implements the 7 prompt functions against TestStackDocHeadings and TestConcernsDocSeverity
- Plan 02-04 (map command) implements `runMap` fully against TestMapCommandWritesAllDocs and TestUpdateFlagInvalidDoc

---
*Phase: 02-codebase-mapping*
*Completed: 2026-03-15*

## Self-Check: PASSED

All 10 created files verified present. Commit 4fdabd2 verified in git log.
