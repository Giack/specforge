# SpecForge

Codebase analyzer and Claude Code plugin for Spec-Driven Development.

## What it does today

**CLI** (`specforge analyze`): Walks a codebase and outputs structured JSON of
packages, types, functions, and imports. Supports Go, TypeScript, Python, Kotlin
with auto-detection.

**Claude Code plugin** (`/specforge:map`): Interactive command that runs
`specforge analyze`, detects architectural domains, and dispatches parallel
mapper agents to generate Skills 2.0 artifacts per domain.

## Installation

### From source

```bash
git clone https://github.com/Giack/specforge.git
cd specforge
make build          # → bin/specforge
make install        # → installs to GOPATH/bin
```

### Binary (once releases are published)

```bash
curl -sSL https://raw.githubusercontent.com/Giack/specforge/main/install.sh | bash
```

## CLI Usage

```bash
# Analyze current directory (auto-detect language)
specforge analyze

# Analyze specific path, force Go, exclude vendor
specforge analyze ./myproject --lang go --exclude vendor,testdata

# Output as Markdown summary instead of JSON
specforge analyze --format markdown
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--lang` | `auto` | `go` \| `ts` \| `kotlin` \| `python` \| `auto` |
| `--format` | `json` | `json` \| `markdown` |
| `--exclude` | — | Comma-separated dirs to exclude |

## Claude Code Plugin: /specforge:map

Run inside Claude Code to map a codebase into Skills 2.0 domain artifacts:

```
/specforge:map [path]
```

The command will:
1. Ask you for language, domains, and exclude dirs
2. Run `specforge analyze` to gather structure
3. Detect or confirm architectural domains
4. Spawn parallel mapper agents (tech, arch, quality, concerns) per domain
5. Write `.claude/skills/[domain]/*.md` artifacts
6. Enrich or create `CLAUDE.md` with architecture context

**Requires `specforge` binary on PATH.** Install from source if not present.

## Architecture

- **`cmd/specforge/`** — Cobra root command + version
- **`cmd/analyze/`** — `analyze` subcommand
- **`internal/analyzer/`** — multi-language analysis (detect, go, ts, python, kotlin)
- **`internal/mapper/`** — AST mapper + snapshot logic for Go
- **`commands/map.md`** — `/specforge:map` Claude Code plugin command
- **`agents/`** — mapper agent definitions (tech, arch, quality, concerns, synthesizer)
- **`skills/`** — specforge-mapper-workflow and specforge-output-format skills

## Building

```bash
GOROOT=/usr/local/go go build ./...
GOROOT=/usr/local/go go test ./...
```

See `CLAUDE.md` for Apple Silicon Go toolchain notes.

