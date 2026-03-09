# Roadmap: SpecForge

## Overview

SpecForge is extended from a partially broken Go CLI into a full GSD workflow integration layer. The journey runs in strict dependency order: fix two critical bugs that make AI features broken and wave execution insecure, then build codebase mapping as the data-generation foundation, then add memory scaffolding and GSD document management as CLI commands, and finally wrap everything into an MCP stdio server so Claude Code can invoke all capabilities as tools. Each phase is a gate — nothing in phase N+1 is started until phase N is verified working.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Debt and Security Clearance** - Fix broken Anthropic API client and shell injection before any new feature work
- [ ] **Phase 2: Codebase Mapping** - `specforge map` command producing all 7 GSD codebase documents from Go AST analysis
- [ ] **Phase 3: Memory and GSD Docs** - Memory scaffolding and GSD planning document generation as CLI commands
- [ ] **Phase 4: MCP Server and Claude Code Integration** - Wrap all CLI capabilities as MCP tools; register with Claude Code

## Phase Details

### Phase 1: Debt and Security Clearance
**Goal**: The existing Anthropic API client works correctly and wave execution is safe from shell injection
**Depends on**: Nothing (first phase)
**Requirements**: FOUND-01, FOUND-02, FOUND-03
**Success Criteria** (what must be TRUE):
  1. `specforge dev discuss` (or any AI command) returns non-empty content from the Anthropic API — no silent empty responses
  2. Wave execution passes prompts via stdin pipe to the subprocess — no user-supplied content interpolated into shell strings
  3. All HTTP clients (`ai`, `jira`, `vcs`) enforce a timeout; a request to an unreachable host fails within 30s rather than hanging indefinitely
**Plans**: 3 plans

Plans:
- [ ] 01-01-PLAN.md — Wave 0 test scaffolds: failing tests for FOUND-01, FOUND-02, FOUND-03
- [ ] 01-02-PLAN.md — Fix ClaudeResponse struct (FOUND-01) and HTTP timeouts for all five clients (FOUND-03)
- [ ] 01-03-PLAN.md — Replace bash shell injection with exec.Command + StdinPipe (FOUND-02)

### Phase 2: Codebase Mapping
**Goal**: A developer can run `specforge map` on any Go project and get all 7 GSD codebase documents written to `.planning/codebase/`
**Depends on**: Phase 1
**Requirements**: MAP-01, MAP-02, MAP-03, MAP-04
**Success Criteria** (what must be TRUE):
  1. Running `specforge map` on the SpecForge repo creates all 7 files (STACK.md, ARCHITECTURE.md, STRUCTURE.md, CONVENTIONS.md, TESTING.md, INTEGRATIONS.md, CONCERNS.md) under `.planning/codebase/`
  2. Progress messages appear on stderr while the command runs; stdout contains no stray output
  3. Running `specforge map --update CONCERNS.md` regenerates only that file without modifying the other six
  4. Each generated document matches the GSD codebase map template structure (correct headings, key-value tables, file path references)
**Plans**: TBD

### Phase 3: Memory and GSD Docs
**Goal**: A developer can scaffold persistent AI memory and generate the full GSD planning document set from the terminal
**Depends on**: Phase 2
**Requirements**: MEM-01, MEM-02, MEM-03, MEM-04, DOCS-01, DOCS-02, DOCS-03, DOCS-04
**Success Criteria** (what must be TRUE):
  1. `specforge memory init` creates `.memory/` with PROJECT.md, DECISIONS.md, PATTERNS.md, and PROGRESS.md — each containing H2-headed sections and a `Last Updated` timestamp
  2. `specforge memory update DECISIONS.md "New decision text"` appends or replaces the correct section without touching other sections
  3. `specforge memory status` prints each memory file's name and last-modified timestamp
  4. `specforge docs init` generates all four planning documents (PROJECT.md, REQUIREMENTS.md, ROADMAP.md, STATE.md) under `.planning/` — each matching the GSD template heading structure
**Plans**: TBD

### Phase 4: MCP Server and Claude Code Integration
**Goal**: Claude Code can invoke codebase mapping, memory reads/writes, and doc generation as MCP tools via `specforge serve`
**Depends on**: Phase 3
**Requirements**: MCP-01, MCP-02, MCP-03, MCP-04
**Success Criteria** (what must be TRUE):
  1. Piping a raw JSON-RPC `initialize` request to `specforge serve` returns a valid JSON-RPC response on stdout with no other bytes mixed in
  2. Claude Code can call the `specforge_map`, `specforge_docs`, and `specforge_memory` tools and get correct results — verified by MCP tool call and response in Claude Code's tool use panel
  3. Running `specforge install` adds an `mcpServers` entry to `.claude/settings.local.json` and creates `.claude/commands/specforge/` with slash command definitions
  4. No specforge CLI command (outside of `specforge serve`) writes anything to stdout that would corrupt a JSON-RPC channel — all human-readable output routes to stderr or is gated on non-MCP mode
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Debt and Security Clearance | 0/3 | Not started | - |
| 2. Codebase Mapping | 0/TBD | Not started | - |
| 3. Memory and GSD Docs | 0/TBD | Not started | - |
| 4. MCP Server and Claude Code Integration | 0/TBD | Not started | - |
