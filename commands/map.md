---
description: Map a codebase to Skills 2.0 domain artifacts (SKILL.md per domain, enriched CLAUDE.md, architecture diagrams)
argument-hint: [path to codebase, default "."]
allowed-tools:
  - Read
  - Bash
  - Glob
  - Grep
  - Write
  - Agent
---

# /specforge:map

Map the codebase at `$ARGUMENTS` (default: `.`) to Skills 2.0 domain artifacts.

## Step 1: Gather Context

Ask the user these questions (all at once, formatted as a numbered list):

1. **Codebase path**: "What is the path to the codebase? (default: `.`)"
2. **Primary language**: "Primary language? `go` / `ts` / `python` / `kotlin` / `auto-detect`"
3. **Domains**: "What are the main architectural domains? (e.g., `auth, api, database, worker`) — or write `auto-detect`"
4. **Existing CLAUDE.md**: "Is there an existing CLAUDE.md? (`yes` to enrich it / `no` to create new)"
5. **Exclude dirs**: "Directories to exclude? (e.g., `vendor, dist, __pycache__`) — press enter to skip"

Wait for the user's answers before proceeding.

## Step 2: Run specforge analyze

Run the analysis:

```bash
specforge analyze [path] --lang [lang] --format json > /tmp/specforge-analysis.json 2>/tmp/specforge-analyze-stderr.txt
```

**If `specforge` binary is not found:**

Tell the user:
```
specforge binary not found. Install options:

1. Download release: curl -sSL https://raw.githubusercontent.com/specforge/specforge/main/install.sh | bash
2. Build from source: cd [path-to-specforge-repo] && GOROOT=/usr/local/go go build -o ~/.local/bin/specforge ./cmd/specforge
```

Do not proceed until the binary is available.

**If analyze exits with "not yet implemented" (stub):**
Continue with an empty analysis structure — the agents will explore the codebase directly.

## Step 3: Detect or Confirm Domains

**If user said `auto-detect`:**
- Read `/tmp/specforge-analysis.json`
- Extract the `packages` array
- Group by top-level directory under module root:
  - `internal/auth/...` → domain `auth`
  - `internal/api/...` → domain `api`
  - `cmd/worker/...` → domain `worker`
- Present to user: "Detected domains: [list]. Confirm? (yes / edit)"

**If user provided domains:**
Proceed with those domains.

## Step 4: Prepare Output Directories

For each domain, create the output directory:
```bash
mkdir -p .claude/skills/[domain]
```

## Step 5: Spawn Parallel Mapper Agents

For EACH domain, spawn these 4 agents in parallel using the Agent tool:

**specforge-tech-mapper:**
```
domain: [domain]
codebase_path: [path]
analysis_json_path: /tmp/specforge-analysis.json
output_dir: .claude/skills/[domain]
```

**specforge-arch-mapper:**
```
domain: [domain]
codebase_path: [path]
analysis_json_path: /tmp/specforge-analysis.json
output_dir: .claude/skills/[domain]
```

**specforge-quality-mapper:**
```
domain: [domain]
codebase_path: [path]
analysis_json_path: /tmp/specforge-analysis.json
output_dir: .claude/skills/[domain]
```

**specforge-concerns-mapper:**
```
domain: [domain]
codebase_path: [path]
analysis_json_path: /tmp/specforge-analysis.json
output_dir: .claude/skills/[domain]
```

**Critical:** Dispatch ALL N×4 agents in a single parallel batch. Do not spawn sequentially.

Wait for all agents to complete.

## Step 6: Spawn Synthesizer

After all mapper agents complete, spawn `specforge-synthesizer`:

```
domains: [comma-separated domain list]
codebase_path: [path]
has_existing_claude_md: [yes/no]
skills_dir: .claude/skills
```

Wait for the synthesizer to complete.

## Step 7: Report to User

```
✅ Codebase mapping complete!

Domains mapped: [list]

Files created:
  .claude/skills/[domain]/{SKILL,STACK,ARCHITECTURE,STRUCTURE,CONVENTIONS,TESTING,CONCERNS,INTEGRATIONS}.md
  ...
  .claude/skills/plugin.json
  .claude/skills/README.md
  CLAUDE.md (enriched with ## Architecture block)

Skills will activate automatically in future Claude Code sessions for domain-specific guidance.
```
