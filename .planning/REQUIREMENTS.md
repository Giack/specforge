# Requirements: SpecForge GSD Integration

**Defined:** 2026-03-08
**Core Value:** A developer can run `specforge map` on any codebase and get a GSD-style structured project plan ready for AI-assisted execution — without leaving their terminal.

## v1 Requirements

### Foundation (Bug Fixes & Security)

- [x] **FOUND-01**: `ClaudeResponse` struct correctly decodes Anthropic API shape (`content[].text`) instead of OpenAI shape (`choices[].message.content`)
- [x] **FOUND-02**: Wave execution replaces bash heredoc string interpolation with `exec.Command` + stdin pipe to eliminate shell injection
- [x] **FOUND-03**: All HTTP clients (`ai`, `jira`, `vcs`) initialized with a timeout (≥ 30s)

### Codebase Mapping

- [x] **MAP-01**: `specforge map` command analyzes a Go project and produces all 7 GSD codebase documents in `.planning/codebase/` (STACK.md, ARCHITECTURE.md, STRUCTURE.md, CONVENTIONS.md, TESTING.md, INTEGRATIONS.md, CONCERNS.md)
- [x] **MAP-02**: `specforge map` runs analysis in parallel (goroutines per document category) and reports progress to stderr
- [x] **MAP-03**: `specforge map --update <doc>` regenerates a single document without overwriting others
- [x] **MAP-04**: Each generated document follows the GSD codebase map template structure (headings, key-value tables, file path references)

### GSD Document Management

- [ ] **DOCS-01**: `specforge docs init` generates `.planning/PROJECT.md` from interactive prompts or `--from <file>` flag
- [ ] **DOCS-02**: `specforge docs init` generates `.planning/REQUIREMENTS.md` with REQ-ID format and traceability table
- [ ] **DOCS-03**: `specforge docs init` generates `.planning/ROADMAP.md` with phase structure derived from requirements
- [ ] **DOCS-04**: `specforge docs init` generates `.planning/STATE.md` initialized to Phase 1 pending

### Memory Scaffolding

- [ ] **MEM-01**: `specforge memory init` creates `.memory/` directory with 4 standard MD files: PROJECT.md (re-hydration summary), DECISIONS.md (settled choices), PATTERNS.md (conventions), PROGRESS.md (current state)
- [ ] **MEM-02**: `.memory/` files use plain Markdown with H2 headings as named sections — no proprietary format
- [ ] **MEM-03**: `specforge memory update <section> <content>` appends or replaces a named H2 section in the appropriate memory file
- [ ] **MEM-04**: `specforge memory status` prints a summary of memory files with last-modified timestamps

### Claude Code Plugin Integration

- [ ] **MCP-01**: `specforge serve` starts an MCP stdio server — reads JSON-RPC from stdin, writes to stdout, logs to stderr only
- [ ] **MCP-02**: MCP server exposes 3 tools: `specforge_map` (runs map analysis), `specforge_docs` (generates planning docs), `specforge_memory` (read/write memory sections)
- [ ] **MCP-03**: `specforge install` registers the MCP server in `.claude/settings.local.json` under `mcpServers` and creates the `.claude/commands/specforge/` directory with updated slash command definitions
- [ ] **MCP-04**: No stdout output from any CLI initialization path — all human-readable output uses stderr or is gated behind `os.Stdout` only when MCP mode is off

## v2 Requirements

### Document Updates

- **DOCS-05**: `specforge docs update --section <name>` — in-place section update without clobbering manual edits (requires section-aware merge parser)
- **DOCS-06**: `specforge docs sync` — syncs `.memory/PROGRESS.md` from current STATE.md automatically

### Developer Experience

- **DX-01**: `specforge doctor` — checks environment (Claude API key, VCS credentials, MCP registration, timeout config) and reports issues
- **DX-02**: `specforge map --watch` — incremental re-analysis on file changes
- **MCP-05**: `specforge_phase_context` MCP tool — aggregates current phase context (STATE.md + relevant PLAN.md) in one call

## Out of Scope

| Feature | Reason |
|---------|--------|
| Web UI or dashboard | Terminal-first tool; SaaS complexity out of scope |
| Real-time collaboration | Single-user CLI, not a platform |
| Non-Anthropic AI providers (beyond existing OpenCode) | Scope to what's already integrated |
| Confluence write-back | Read-only sync; write path is high risk |
| Section-aware merge in v1 | High complexity; defer until manual-edit preservation pain is reported |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| FOUND-01 | Phase 1 | Complete |
| FOUND-02 | Phase 1 | Complete |
| FOUND-03 | Phase 1 | Complete |
| MAP-01 | Phase 2 | Complete |
| MAP-02 | Phase 2 | Complete |
| MAP-03 | Phase 2 | Complete |
| MAP-04 | Phase 2 | Complete |
| DOCS-01 | Phase 3 | Pending |
| DOCS-02 | Phase 3 | Pending |
| DOCS-03 | Phase 3 | Pending |
| DOCS-04 | Phase 3 | Pending |
| MEM-01 | Phase 3 | Pending |
| MEM-02 | Phase 3 | Pending |
| MEM-03 | Phase 3 | Pending |
| MEM-04 | Phase 3 | Pending |
| MCP-01 | Phase 4 | Pending |
| MCP-02 | Phase 4 | Pending |
| MCP-03 | Phase 4 | Pending |
| MCP-04 | Phase 4 | Pending |

**Coverage:**
- v1 requirements: 19 total
- Mapped to phases: 19
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-08*
*Last updated: 2026-03-08 after initial definition*
