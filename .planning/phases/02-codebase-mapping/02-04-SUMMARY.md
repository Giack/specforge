---
phase: 02-codebase-mapping
plan: "04"
subsystem: cli
tags: [cobra, sync, goroutines, mapper, ai-client, parallel]

# Dependency graph
requires:
  - phase: 02-codebase-mapping/02-01
    provides: cmd/map/map.go stub and map_test.go (TDD RED scaffold)
  - phase: 02-codebase-mapping/02-02
    provides: internal/ai.AIClient.Generate method
  - phase: 02-codebase-mapping/02-03
    provides: internal/mapper prompt builders (StackPrompt...ConcernsPrompt)
provides:
  - specforge map Cobra command, fully wired with parallel goroutine dispatch
  - --update flag with validation against 7 valid document names
  - cmd/map package with generator interface for testability
  - main.go registration of map command alongside pm, em, dev
affects:
  - 03-slash-commands (uses specforge map as underlying engine)
  - 04-mcp-server (map command stdout cleanliness required for MCP use)

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Package-level var injection (newAIClient, outputDir) for test isolation without interfaces overhead
    - generator interface (single Generate method) for mock injection in white-box tests
    - sync.WaitGroup + mutex-protected error slice for parallel document generation

key-files:
  created: []
  modified:
    - cmd/map/map.go
    - cmd/map/map_test.go
    - cmd/specforge/main.go

key-decisions:
  - "Package-level var newAIClient factory allows test injection of mockGenerator without changing runMap signature"
  - "Package-level var outputDir allows test to redirect writes to t.TempDir() without passing outDir through function signature"
  - "generator interface (single method) kept internal to cmd/map — not exported, no need for external use"

patterns-established:
  - "All cmd/map output via fmt.Fprintf(os.Stderr, ...) — stdout reserved for MCP Phase 4"
  - "Parallel goroutines: wg.Add(1) + go func(d docEntry){defer wg.Done(); ...}(doc) — errors collected with mutex"

requirements-completed: [MAP-01, MAP-02, MAP-03]

# Metrics
duration: 10min
completed: 2026-03-15
---

# Phase 02 Plan 04: Map Command Integration Summary

**`specforge map` Cobra command wired end-to-end: parallel goroutine dispatch of 7 prompt builders, --update flag validation, and stdout-clean output for future MCP use**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-03-15T16:53:00Z
- **Completed:** 2026-03-15T16:54:18Z
- **Tasks:** 2
- **Files modified:** 3

## Accomplishments

- Implemented `runMap` with `sync.WaitGroup` parallel dispatch for all 7 GSD document goroutines
- Added `newAIClient` and `outputDir` package-level vars enabling white-box test injection without API keys
- Defined `generator` interface so `mockGenerator` can replace `*ai.AIClient` in tests
- Registered `mapcmd.NewCommand(cfg)` in `main.go`; `specforge map --help` shows --update flag
- Full test suite passes with race detector; `specforge map --update TYPO.md` exits 1 with valid-doc list

## Task Commits

Each task was committed atomically:

1. **Task 1: Implement cmd/map/map.go with parallel document generation** - `e45b606` (feat)
2. **Task 2: Register map command in main.go** - `45ed5d1` (feat)

**Plan metadata:** _(docs commit follows)_

## Files Created/Modified

- `cmd/map/map.go` - Full runMap implementation: validDocs, dispatchTable, newAIClient var, outputDir var, generator interface, parallel WaitGroup dispatch
- `cmd/map/map_test.go` - Updated with mockGenerator, newAIClient injection, outputDir injection; TestMapCommandWritesAllDocs and TestUpdateFlagInvalidDoc both GREEN
- `cmd/specforge/main.go` - Added mapcmd import alias and rootCmd.AddCommand(mapcmd.NewCommand(cfg))

## Decisions Made

- **Package-level var injection for AI client:** `newAIClient func(config.AIConfig) generator` allows tests to provide `mockGenerator{content: "# Mock content"}` without real HTTP calls. Consistent with existing `execCommand` and `newHTTPClient` patterns.
- **Package-level var for outputDir:** Tests set `outputDir = t.TempDir()` to capture written files without passing an extra parameter through `runMap`. Clean interface preserved.
- **generator interface stays internal:** The one-method interface is an implementation detail of cmd/map. No need to export it.
- **go 1.24 WaitGroup pattern:** `sync.WaitGroup.Go()` is not available (requires 1.25+). Used `wg.Add(1)` + `go func(d docEntry){ defer wg.Done(); ... }(doc)` to avoid loop-variable capture.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Test `TestMapCommandWritesAllDocs` had wiring gap**
- **Found during:** Task 1
- **Issue:** The existing test created `outDir := t.TempDir()` but never injected it into `runMap`. Files would be written to `.planning/codebase` relative to test working dir, not to `outDir`. Test would always fail the stat checks.
- **Fix:** Added `outputDir` package-level var (default `".planning/codebase"`). Updated test to set `outputDir = outDir` before calling `runMap` and restore it via defer.
- **Files modified:** `cmd/map/map_test.go`
- **Verification:** `TestMapCommandWritesAllDocs PASS` with all 7 files present in outDir
- **Committed in:** `e45b606` (Task 1 commit)

**2. [Rule 1 - Bug] sync.WaitGroup.Go() unavailable in go 1.24**
- **Found during:** Task 1 implementation
- **Issue:** Plan spec referenced `sync.WaitGroup.Go()` (added in Go 1.25). Module is `go 1.24.0`.
- **Fix:** Used standard `wg.Add(1)` + `go func(d docEntry){ defer wg.Done(); ... }(doc)` pattern with explicit loop variable capture.
- **Files modified:** `cmd/map/map.go`
- **Verification:** `go vet ./cmd/map/... && go build ./cmd/map/...` both exit 0
- **Committed in:** `e45b606` (Task 1 commit)

---

**Total deviations:** 2 auto-fixed (both Rule 1 - Bug)
**Impact on plan:** Both fixes necessary for correctness. No scope creep.

## Issues Encountered

None beyond the auto-fixed deviations above.

## User Setup Required

None — no external service configuration required for this plan.

## Next Phase Readiness

- `specforge map` command is fully functional; runs on any Go project root
- Requires `ANTHROPIC_API_KEY` set at runtime for actual document generation (not a test-time requirement)
- Phase 3 (slash commands) can reference `specforge map` as the underlying engine
- Phase 4 (MCP server) can rely on stdout cleanliness — all progress output goes to stderr

---
*Phase: 02-codebase-mapping*
*Completed: 2026-03-15*

## Self-Check: PASSED

- FOUND: cmd/map/map.go
- FOUND: cmd/map/map_test.go
- FOUND: cmd/specforge/main.go
- FOUND: e45b606 (Task 1 commit)
- FOUND: 45ed5d1 (Task 2 commit)
