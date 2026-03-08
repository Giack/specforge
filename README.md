# SpecForge

Enterprise Spec-Driven Development Tool - CLI written in Golang

## Overview

SpecForge bridges business requirements (Confluence/Jira) with AI execution engines (Claude Code/OpenCode) using Spec-Driven Development principles.

## Installation

```bash
# Build from source
make build

# Or install globally
make install
```

## Configuration

Copy `config.example.yaml` to `~/.specforge/config.yaml` and fill in your credentials:

```yaml
atlassian:
  domain: your-company
  email: your-email@company.com
  api_token: your-api-token
  project_key: PROJ

# Version Control System (GitHub, GitLab, or Bitbucket)
vcs:
  provider: github  # github, gitlab, bitbucket

  github:
    token: ghp_your_token
    owner: your-org
    repo: your-repo

  gitlab:
    domain: gitlab.com
    token: glpat_your_token
    group: your-group

  bitbucket:
    domain: bitbucket.org
    username: your-username
    app_password: your-app_password
    workspace: your-workspace

ai:
  provider: claude
  model: sonnet-4-20250514
```

## Getting Started

### 1. Initialize SpecForge Commands

Install the Claude Code/OpenCode commands:

```bash
make init
```

### 2. PM Workflow - Sync Requirements

```bash
# Sync from Confluence
specforge pm sync --type confluence --url "https://your-company.atlassian.net/wiki/spaces/PROJ/pages/123456789"

# Sync from Jira
specforge pm sync --type jira --url "https://your-company.atlassian.net/browse/PROJ-123?issueKey=PROJ-123"
```

### 3. EM Workflow - Architect & Verify

```bash
# Analyze requirements and create roadmaps
specforge em architect

# Run UAT verification
specforge em verify

# Create bug from UAT feedback
specforge em bug --feedback "Login button doesn't show spinner"
```

### 4. Dev Workflow

Run these commands **inside Claude Code or OpenCode**:

```bash
# Start the workflow
/specforge:dev-discuss

# Create execution plans
/specforge:dev-plan

# Execute with parallel sub-agents
/specforge:dev-execute

# Fix a Jira bug
/specforge:dev-fix PROJ-123
```

## Architecture

- **Core CLI**: Golang binary handling business logic
- **Skill Injector**: `specforge init` installs commands into Claude Code
- **Sub-agents**: Execution uses fresh context per task to prevent context rot

## Commands

| Command | Description |
|---------|-------------|
| `specforge init` | Install Claude Code commands |
| `specforge pm sync` | Sync requirements from Confluence/Jira |
| `specforge em architect` | Create domain roadmaps |
| `specforge em verify` | Run UAT checklist |
| `specforge em bug` | Create Jira bug from feedback |
| `/specforge:dev-discuss` | Capture implementation context |
| `/specforge:dev-plan` | Create atomic execution plans |
| `/specforge:dev-execute` | Execute plans with sub-agents |
| `/specforge:dev-fix` | Fix Jira bug with AI debugger |
| `specforge dev pr` | Create Pull Request (GitHub/GitLab/Bitbucket) |
