---
phase: 01-debt-and-security-clearance
plan: "03"
subsystem: security
tags: [go, exec, shell-injection, subprocess, stdin-pipe]

requires:
  - phase: 01-01
    provides: TDD test scaffolding (TestExecuteWithClaudeNoShell, TestExecuteWithOpenCodeNoShell) in RED state

provides:
  - Shell-injection-free executeWithClaude using exec.Command("claude") + StdinPipe
  - Shell-injection-free executeWithOpenCode using exec.Command("opencode") + StdinPipe
  - No temp file written to /tmp/specforge_exec.sh

affects:
  - Phase 4 MCP integration — wave execution is now safe for autonomous invocation

tech-stack:
  added: []
  patterns:
    - "Direct subprocess invocation via exec.Command + StdinPipe instead of bash wrapper"
    - "Package-level execCommand factory for testability of subprocess calls"

key-files:
  created: []
  modified:
    - cmd/dev/wave.go

key-decisions:
  - "Keep execCommand factory var for testability — StdinPipe() works on the *exec.Cmd returned by the mock"
  - "Add 'io' import for io.WriteString; retain 'strings' and 'os' (used elsewhere in file)"

patterns-established:
  - "Subprocess pattern: execCommand(binary, args...) + StdinPipe + Start + WriteString + Close + Wait"

requirements-completed:
  - FOUND-02

duration: 10min
completed: 2026-03-09
---

# Phase 1 Plan 03: Shell Injection Fix Summary

**Eliminated bash heredoc shell injection in wave execution — both executeWithClaude and executeWithOpenCode now invoke claude/opencode directly via exec.Command + StdinPipe with no temp file written**

## Performance

- **Duration:** 10 min
- **Started:** 2026-03-09T20:10:00Z
- **Completed:** 2026-03-09T20:20:00Z
- **Tasks:** 1
- **Files modified:** 1

## Accomplishments

- Removed `fmt.Sprintf` bash heredoc that interpolated AI prompts into shell scripts (shell injection vector)
- Both execute functions now call the target binary directly with prompt delivered via stdin pipe
- No temp file `/tmp/specforge_exec.sh` created
- All tests pass: TestExecuteWithClaudeNoShell PASS, TestExecuteWithOpenCodeNoShell PASS

## Task Commits

1. **Task 1: Replace bash wrapper with exec.Command + StdinPipe** - `a558ff9` (feat)

**Plan metadata:** (included in task commit)

## Files Created/Modified

- `cmd/dev/wave.go` - Replaced bash wrapper with direct subprocess invocation + StdinPipe in executeWithClaude and executeWithOpenCode

## Decisions Made

- Retained `execCommand` factory variable so existing test mocks (`execCommand = func(...)`) continue to intercept calls; `StdinPipe()` works on the `*exec.Cmd` the mock returns since tests mock to `exec.Command("true")` which is a real process.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None — `go build` succeeded on first attempt; both tests passed immediately after implementation.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- FOUND-02 fixed; shell injection surface eliminated from wave execution
- Phase 4 MCP integration can proceed without risk of prompt content executing as shell code
- All 3 FOUND-0x findings (FOUND-01, FOUND-02, FOUND-03) are now resolved

---
*Phase: 01-debt-and-security-clearance*
*Completed: 2026-03-09*
