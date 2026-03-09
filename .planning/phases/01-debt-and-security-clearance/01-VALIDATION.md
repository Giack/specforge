---
phase: 1
slug: debt-and-security-clearance
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-09
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — Wave 0 creates test files |
| **Quick run command** | `go test ./internal/ai/... ./internal/jira/... ./internal/vcs/... ./cmd/dev/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/ai/... ./internal/jira/... ./internal/vcs/... ./cmd/dev/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 1-01-01 | 01 | 0 | FOUND-03 | unit | `go test ./internal/ai/...` | ❌ W0 | ⬜ pending |
| 1-01-02 | 01 | 0 | FOUND-03 | unit | `go test ./internal/jira/...` | ❌ W0 | ⬜ pending |
| 1-01-03 | 01 | 0 | FOUND-03 | unit | `go test ./internal/vcs/...` | ❌ W0 | ⬜ pending |
| 1-02-01 | 02 | 0 | FOUND-01 | unit | `go test ./internal/ai/...` | ❌ W0 | ⬜ pending |
| 1-03-01 | 03 | 0 | FOUND-02 | unit | `go test ./cmd/dev/...` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/ai/client_test.go` — stubs for FOUND-01 and FOUND-03 (ai client timeout)
- [ ] `internal/jira/client_test.go` — stubs for FOUND-03 (jira/confluence timeout)
- [ ] `internal/vcs/github_test.go` — stubs for FOUND-03 (vcs timeout)
- [ ] `cmd/dev/wave_test.go` — stubs for FOUND-02 (shell injection fix)

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| `specforge dev discuss` returns real content | FOUND-01 | Requires live Anthropic API key | Run `specforge dev discuss` with a real task and confirm non-empty AI response |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
