# SpecForge - Enterprise Spec-Driven Development Tool

## Vision
SpecForge is an Enterprise-grade CLI tool built in Golang that implements Spec-Driven Development across the entire product lifecycle (PM -> EM -> Dev -> QA). It acts as a bridge between high-level business requirements (Confluence/Jira) and native AI execution engines (Claude Code/OpenCode) without suffering from context rot.

## Architecture & Integration Strategy
*   **Golang Core CLI (`specforge`)**: Handles business logic, TUI, Atlassian/Bitbucket API integrations, and multi-repo management.
*   **Native AI Plugin**: SpecForge acts as a "Skill Injector". Instead of wrapping LLM calls externally, `specforge init` auto-installs custom commands and sub-agent workflows directly into Claude Code (`~/.claude/commands/specforge/`) or OpenCode.
*   **Context Rot Prevention**: By leveraging native sub-agents (`Task` tools), every atomic plan generated from the roadmap is executed in a fresh context window. The agent writes the code, writes the tests, commits, and dies, returning only a success signal to the main loop.

## Key Workflows

### 1. The Spec Engine (PM Flow)
*   **Goal**: Translate informal documents into machine-readable specs.
*   **Command**: `specforge pm sync --confluence <URL>`
*   **Action**: Extracts content from Confluence/Jira, structures it into `PROJECT.md` and `REQUIREMENTS.md`, and stores them in a Centralized Spec Repo (or Confluence itself as Single Source of Truth).

### 2. The Architect Engine (EM Flow)
*   **Goal**: Break down requirements into actionable microservice/frontend roadmaps.
*   **Command**: `specforge em architect`
*   **Action**: 
    *   AI analyzes requirements and splits them across domains (Hybrid Monorepo/Polyrepo support).
    *   Generates `ARCHITECTURE.md` (with Mermaid.js diagrams).
    *   Generates a `ROADMAP.md` for each targeted repository.

### 3. The Execution Engine (Dev Flow)
*   **Goal**: Execute specs, write tested code, and open PRs.
*   **Execution Environment**: Inside Claude Code / OpenCode using native SpecForge skills.
*   **Commands**:
    *   `/specforge:dev-discuss`: Adds dev constraints (e.g., "Use React Query") to `CONTEXT.md`.
    *   `/specforge:dev-plan`: Creates atomic XML `PLAN.md`s. **Constraint:** Every plan must mandate writing Unit/Integration tests.
    *   `/specforge:dev-execute`: Launches parallel sub-agents (Waves). Each agent executes a plan in a fresh context, runs local tests, and makes an atomic git commit.
    *   *Post-Execution*: Automatically pushes to Bitbucket and opens a PR assigned to reviewers.

### 4. Verification & Bug Tracking (UAT/QA Flow)
*   **Goal**: Manual testing and seamless Jira integration.
*   **Commands**:
    *   `specforge em verify`: TUI-based interactive checklist based on acceptance criteria.
    *   `specforge em bug`: If UAT fails, takes natural language feedback, analyzes recent commits, and automatically opens a Jira ticket via Atlassian APIs.
    *   `/specforge:dev-fix <JIRA-ID>` (Inside AI CLI): Launches `gsd-debugger` sub-agent to find the root cause in isolation and generate a fix plan.

## Tech Stack
*   **Language**: Golang
*   **CLI & Config**: `spf13/cobra`, `spf13/viper`
*   **TUI**: `charmbracelet/bubbletea`, `charmbracelet/huh`
*   **APIs**: `go-jira` (Atlassian REST/GraphQL), VCS (GitHub/GitLab/Bitbucket API)
*   **AI Integration**: Native Claude Code / OpenCode custom tools/skills injection.
