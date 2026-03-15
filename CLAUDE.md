# SpecForge — Project Instructions

## Go Toolchain

This project runs on Apple Silicon (arm64). The PATH `go` binary is the x86_64 v1.24 system install whose GOROOT dynamically resolves to the arm64 mise installation, causing `go: no such tool "compile"`.

**Always prefix all `go` commands with `GOROOT=/usr/local/go`:**

```bash
GOROOT=/usr/local/go go build ./...
GOROOT=/usr/local/go go test ./...
GOROOT=/usr/local/go go test ./internal/mapper/... ./cmd/map/...
```

The `GOROOT=/usr/local/go` prefix is the standard workaround in this project for all agents.

## Module

- Module: `specforge`
- Go version: `go 1.24.0` (go.mod)
- Package structure: `cmd/` (Cobra commands), `internal/` (libraries)

## Test Commands

```bash
# Quick: mapper + map command only
GOROOT=/usr/local/go go test ./internal/mapper/... ./cmd/map/...

# Full suite with race detector
GOROOT=/usr/local/go go test -v -race ./...
```

## Code Conventions

- No `fmt.Println` / `fmt.Printf` in `cmd/map/` — progress goes to `os.Stderr` only (stdout reserved for MCP)
- White-box tests (same package) to access unexported functions
- Package-level var injection (e.g., `execCommand`, `newHTTPClient`) for testability
- `t.Parallel()` only on tests that don't mutate package-level vars
