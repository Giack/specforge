---
name: specforge-mapper-workflow
description: Use this skill when mapping a codebase domain as a specforge mapper agent (tech, arch, quality, or concerns). Provides the workflow protocol for exploring, analyzing, and writing domain artifacts.
version: 0.1.0
---

# SpecForge Mapper Workflow

This skill is loaded by specforge mapper agents to guide their codebase exploration and document writing workflow.

## Core Protocol

### 1. Parse Your Inputs
Extract from your prompt:
- `domain`: the domain name (e.g., `auth`, `api`, `database`)
- `codebase_path`: root path of the codebase
- `analysis_json_path`: path to the specforge analysis JSON
- `output_dir`: where to write your documents (`.claude/skills/[domain]`)

### 2. Read the Analysis JSON
Always start by reading `analysis_json_path`. Use the `packages` array to identify which packages belong to your domain. Filter by packages where `path` contains the domain name.

### 3. Explore Domain Code
Navigate to the domain's packages and read key files. Use Glob and Grep to find patterns. Read 3-5 representative files — don't read everything.

### 4. Write Documents Directly
Use the Write tool to create documents. Never return document contents to the orchestrator — write directly to disk.

### 5. Return Confirmation Only
Return a brief confirmation (10 lines max): what domain, what documents written, line counts.

## Quality Standards

- **Every finding needs a file path**: Use backticks. `internal/auth/handler.go` not "the auth handler".
- **Be prescriptive**: "Use X pattern" not "X pattern is used".
- **Severity for concerns**: Always `Severity: High`, `Severity: Medium`, or `Severity: Low`.
- **Mermaid for architecture**: Include a `graph TD` or `graph LR` Mermaid diagram in ARCHITECTURE.md.
- **Real examples**: Show actual code patterns from the codebase, not generic examples.
