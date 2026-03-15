---
phase: 02-codebase-mapping
plan: "03"
subsystem: api
tags: [go, prompts, ast, codegen, gsd-docs]

# Dependency graph
requires:
  - phase: 02-codebase-mapping
    provides: CodebaseSnapshot struct and WalkProject from 02-01; TDD RED tests from 02-02
provides:
  - 7 exported prompt builder functions in internal/mapper/prompts.go
  - StackPrompt, ArchitecturePrompt, StructurePrompt, ConventionsPrompt, TestingPrompt, IntegrationsPrompt, ConcernsPrompt
affects: [02-04-cmd-map, phase-03-mcp-server]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Each prompt builder slices only its relevant snapshot fields — no full snapshot dumps"
    - "Hardcoded GSD heading structure embedded in prompt template — AI fills content, never chooses headings"
    - "fmt.Sprintf with backtick multi-line template strings for prompt construction"

key-files:
  created: []
  modified:
    - internal/mapper/prompts.go

key-decisions:
  - "Each builder receives *CodebaseSnapshot and slices only relevant fields to keep prompts targeted and token-efficient"
  - "ConcernsPrompt instructs Claude to include Severity: High/Medium/Low per concern entry — matches CONTEXT.md locked decision"
  - "TestingPrompt filters _test.go files by suffix — avoids passing non-test AST data to the testing analysis prompt"

patterns-established:
  - "Prompt template pattern: hardcoded headings block + Data: section with sliced snapshot fields + fill instruction"
  - "IntegrationsPrompt uses Imports directly (WalkProject already stores only external/non-stdlib imports)"

requirements-completed: [MAP-04]

# Metrics
duration: 10min
completed: 2026-03-15
---

# Phase 2 Plan 03: Prompt Builders Summary

**7 GSD document prompt builders using fmt.Sprintf templates with hardcoded heading structure and per-function snapshot field slicing**

## Performance

- **Duration:** ~10 min
- **Started:** 2026-03-15T16:49:00Z
- **Completed:** 2026-03-15T16:49:55Z
- **Tasks:** 1
- **Files modified:** 1

## Accomplishments
- All 7 prompt builders implemented with hardcoded GSD heading structure — AI fills content, never chooses headings
- TestStackDocHeadings passes — StackPrompt returns string with all 7 required STACK.md headings
- TestConcernsDocSeverity passes — ConcernsPrompt returns string containing "Severity:"
- Each builder slices only its relevant snapshot fields (no full snapshot dumps)
- Full mapper test suite passes with race detector (6/6 tests)

## Task Commits

1. **Task 1: Implement 7 document prompt builders** - `df0d895` (feat)

**Plan metadata:** (docs commit — see below)

## Files Created/Modified
- `internal/mapper/prompts.go` - All 7 prompt builder functions implemented (was stub returning "")

## Decisions Made
- Each prompt builder slices only its relevant `*CodebaseSnapshot` fields — keeps prompts targeted and avoids token waste on irrelevant AST data
- ConcernsPrompt template includes `**Severity:** High/Medium/Low` instruction per concern — matches locked CONTEXT.md decision
- TestingPrompt filters using `strings.HasSuffix(f, "_test.go")` — clean separation of test vs non-test data
- IntegrationsPrompt uses `snap.Imports` directly — WalkProject already stores only external imports per 02-01 implementation

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- All 7 prompt builders are ready for dispatch in cmd/map (plan 02-04)
- AIClient.Generate caller owns full prompt construction — prompt builders provide the complete prompt string
- go build ./... succeeds, all mapper tests GREEN

---
*Phase: 02-codebase-mapping*
*Completed: 2026-03-15*
