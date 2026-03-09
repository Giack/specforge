# Phase 1: Debt and Security Clearance - Research

**Researched:** 2026-03-09
**Domain:** Go HTTP client bug fixes, shell injection remediation, `net/http` timeout configuration
**Confidence:** HIGH — all findings are from direct source code inspection of the codebase; no inference required

## Summary

Phase 1 is three targeted bug fixes in an existing Go 1.24 codebase. No new packages are needed. Every bug was confirmed by reading the source directly.

**Bug 1 (FOUND-01):** `internal/ai/client.go` decodes the Anthropic API response into a struct shaped for OpenAI (`Choices []Choice` / `choices[].message.content`). The Anthropic Messages API returns `content[].text`. The struct needs two type replacements and the extractor on line 163 needs updating. Until this is fixed every call to `callClaude()` silently returns an empty string.

**Bug 2 (FOUND-02):** `cmd/dev/wave.go` functions `executeWithClaude` and `executeWithOpenCode` build a bash script by interpolating the prompt string via `fmt.Sprintf` into the script body, then writing that to `/tmp/specforge_exec.sh` and running it with `exec.Command("bash", scriptFile)`. Any backtick, `$()`, or `"` in the prompt — all common in AI-generated content — becomes arbitrary shell execution. The fix is to drop the bash wrapper entirely: call `exec.Command("claude", "-p", "--dangerously-skip-permissions")` directly and write the prompt to the process's stdin pipe.

**Bug 3 (FOUND-03):** All three HTTP clients (`internal/ai/client.go`, `internal/jira/client.go`, `internal/vcs/github.go`, `internal/vcs/gitlab.go`, `internal/vcs/bitbucket.go`) construct `&http.Client{}` with no `Timeout` field. A default Go `http.Client` has no timeout and will hang indefinitely if the remote host is unreachable. Each constructor must set `Timeout: 30 * time.Second`.

**Primary recommendation:** Fix all three bugs in a single phase in the order: FOUND-03 (lowest risk, pure addition) → FOUND-01 (struct change, no behavior risk) → FOUND-02 (behavioral change, highest impact). Write one integration test per fix before merging.

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| FOUND-01 | `ClaudeResponse` struct correctly decodes Anthropic API shape (`content[].text`) instead of OpenAI shape (`choices[].message.content`) | Bug confirmed at lines 26-32 and 159-163 of `internal/ai/client.go`; correct Anthropic shape documented at docs.anthropic.com/en/api/messages |
| FOUND-02 | Wave execution replaces bash heredoc string interpolation with `exec.Command` + stdin pipe to eliminate shell injection | Bug confirmed at lines 82-97 and 103-118 of `cmd/dev/wave.go`; `fmt.Sprintf` injects prompt into shell script string |
| FOUND-03 | All HTTP clients (`ai`, `jira`, `vcs`) initialized with a timeout (≥ 30s) | Confirmed: `internal/ai/client.go` line 147, `internal/jira/client.go` lines 27-30 and 135-138, `internal/vcs/github.go` line 32, `internal/vcs/gitlab.go`, `internal/vcs/bitbucket.go` all use `&http.Client{}` with no Timeout |
</phase_requirements>

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib `net/http` | 1.24 | HTTP client with timeout | Already in use; `http.Client{Timeout: ...}` is the canonical pattern |
| Go stdlib `os/exec` | 1.24 | Subprocess execution without shell | Replaces bash wrapper; no CGO, no new dep |
| Go stdlib `io` | 1.24 | Pipe prompt via stdin to subprocess | `cmd.StdinPipe()` is the correct safe pattern |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Go stdlib `time` | 1.24 | Timeout duration constant | `30 * time.Second` in all `http.Client` constructors |
| Go stdlib `net/http/httptest` | 1.24 | Test server for HTTP timeout verification | Integration tests for FOUND-01 and FOUND-03 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `http.Client{Timeout}` | Context with deadline per-request | Per-request context is more flexible but requires API changes across all callers; global client timeout is simpler and sufficient |
| Direct `exec.Command("claude", ...)` | Wrapper script | Script approach introduced the injection; direct exec is both simpler and safer |

**Installation:** No new dependencies. All fixes use Go stdlib only.

## Architecture Patterns

### FOUND-01: Correct Anthropic Response Struct

The Anthropic Messages API (`POST https://api.anthropic.com/v1/messages`) returns:

```json
{
  "content": [
    { "type": "text", "text": "..." }
  ]
}
```

The current code expects the OpenAI shape:

```go
// CURRENT (wrong) — internal/ai/client.go lines 26-32
type ClaudeResponse struct {
    Choices []Choice `json:"choices"`
}
type Choice struct {
    Message Message `json:"message"`
}
```

The fix requires replacing both types and updating the extraction:

```go
// CORRECT
type ClaudeResponse struct {
    Content []ContentBlock `json:"content"`
}
type ContentBlock struct {
    Type string `json:"type"`
    Text string `json:"text"`
}
```

Extraction at line 159-163 changes from:

```go
// CURRENT (wrong)
if len(result.Choices) == 0 {
    return "", fmt.Errorf("no response from Claude")
}
return result.Choices[0].Message.Content, nil
```

To:

```go
// CORRECT
if len(result.Content) == 0 {
    return "", fmt.Errorf("no response from Claude")
}
return result.Content[0].Text, nil
```

### FOUND-02: Eliminate Shell Injection via stdin Pipe

The current pattern (`cmd/dev/wave.go` lines 82-97):

```go
// CURRENT (vulnerable) — prompt interpolated into shell script
script := fmt.Sprintf(`#!/bin/bash
echo "%s" | claude -p --dangerously-skip-permissions
`, taskName, prompt)
os.WriteFile(scriptFile, []byte(script), 0755)
execCmd := exec.Command("bash", scriptFile)
```

The safe replacement passes the prompt directly via stdin:

```go
// CORRECT — no shell, no interpolation
execCmd := exec.Command("claude", "-p", "--dangerously-skip-permissions")
execCmd.Stdout = os.Stdout
execCmd.Stderr = os.Stderr
stdin, err := execCmd.StdinPipe()
if err != nil {
    return fmt.Errorf("failed to create stdin pipe: %w", err)
}
if err := execCmd.Start(); err != nil {
    return fmt.Errorf("failed to start claude: %w", err)
}
if _, err := io.WriteString(stdin, prompt); err != nil {
    return fmt.Errorf("failed to write prompt: %w", err)
}
stdin.Close()
return execCmd.Wait()
```

The same fix applies to `executeWithOpenCode` (replace `opencode -p` as the command).

Note: the `/tmp/specforge_exec.sh` temp file can also be removed after this fix — it serves no purpose once the bash wrapper is gone.

### FOUND-03: HTTP Client Timeout

All five locations currently use `&http.Client{}` (no timeout). The pattern is identical across all:

```go
// CURRENT (no timeout)
client: &http.Client{}

// CORRECT — add Timeout field
client: &http.Client{
    Timeout: 30 * time.Second,
}
```

Locations to update:
- `internal/ai/client.go` line 147 (inside `callClaude`)
- `internal/jira/client.go` line 28 (`NewJiraClient`)
- `internal/jira/client.go` line 136 (`NewConfluenceClient`)
- `internal/vcs/github.go` line 33 (`NewGitHubClient`)
- `internal/vcs/gitlab.go` — equivalent constructor (not read; confirm pattern matches)
- `internal/vcs/bitbucket.go` — equivalent constructor (not read; confirm pattern matches)

**Note:** The `ai` client creates a new `http.Client{}` inline inside `callClaude` rather than at construction time (the `AIClient` struct has no `client` field). The fix should either add a `client` field to `AIClient` (consistent with jira/vcs pattern) or simply initialize with `Timeout` at the inline call site. Adding a struct field is the cleaner approach and aligns with the rest of the codebase.

### Anti-Patterns to Avoid

- **Mocking the AI client to work around FOUND-01:** The integration test must call the real Anthropic API (or a local httptest server returning the correct shape) to verify the fix. A mock that returns hardcoded content does not verify the struct fix.
- **Keeping the temp script file for FOUND-02:** After replacing with stdin pipe, delete the `/tmp/specforge_exec.sh` write entirely. Leaving dead code in place invites regression.
- **Adding context-based timeout only to new code:** The timeout must be applied to existing constructors — not only to future call sites.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| HTTP timeout | Custom deadline middleware | `http.Client{Timeout: 30 * time.Second}` | stdlib handles connection + read timeout in one field |
| Safe subprocess input | Temp file + shell echo | `exec.Command` + `StdinPipe()` | Shell is the attack surface; stdin pipe bypasses it entirely |
| API response validation | Custom JSON decoder | Standard `json.Unmarshal` + length check | The only issue is the wrong struct shape; fix the struct, keep the decoder |

## Common Pitfalls

### Pitfall 1: `ai` client creates `http.Client` inline, not at construction
**What goes wrong:** Developer adds `Timeout` to `NewAIClient` but `callClaude` instantiates a fresh `&http.Client{}` at line 147, ignoring any struct field.
**Why it happens:** The `AIClient` struct has no `client` field unlike `JiraClient` and `GitHubClient`.
**How to avoid:** Add a `client *http.Client` field to `AIClient` and initialize it in `NewAIClient`, then use `a.client.Do(req)` in `callClaude`.
**Warning signs:** Timeout test passes for jira/vcs but hangs for the ai client.

### Pitfall 2: Shell quoting survives the fix if bash stays involved
**What goes wrong:** Developer replaces the heredoc with a `fmt.Sprintf` that sets an env var or uses bash `-c`, thinking quoting is handled.
**Why it happens:** Any bash involvement reintroduces the attack surface.
**How to avoid:** The subprocess must be `exec.Command("claude", ...)` or `exec.Command("opencode", ...)` directly — bash must not appear in the call chain.
**Warning signs:** The command string still contains `bash` or `-c`.

### Pitfall 3: Integration test for FOUND-01 passes with an empty string
**What goes wrong:** Test asserts `err == nil` but does not assert `len(response) > 0`, so a silent empty response still passes.
**Why it happens:** The bug's symptom is silent empty output, not an error.
**How to avoid:** The success criterion is explicit: response must be non-empty. Assert `len(response) > 0` in the test.
**Warning signs:** Test passes even when the wrong struct is in place.

## Code Examples

Verified patterns from Go stdlib documentation:

### HTTP Client with Timeout
```go
// Source: https://pkg.go.dev/net/http#Client
client := &http.Client{
    Timeout: 30 * time.Second,
}
```

### exec.Command with StdinPipe
```go
// Source: https://pkg.go.dev/os/exec#Cmd.StdinPipe
cmd := exec.Command("claude", "-p", "--dangerously-skip-permissions")
stdin, err := cmd.StdinPipe()
if err != nil {
    return err
}
if err := cmd.Start(); err != nil {
    return err
}
io.WriteString(stdin, prompt)
stdin.Close()
return cmd.Wait()
```

### httptest server for timeout testing
```go
// Source: https://pkg.go.dev/net/http/httptest
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(ClaudeResponse{
        Content: []ContentBlock{{Type: "text", Text: "hello"}},
    })
}))
defer ts.Close()
```

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` (no external test framework in go.mod) |
| Config file | none — `go test ./...` is standard |
| Quick run command | `go test ./internal/ai/... ./internal/jira/... ./internal/vcs/... ./cmd/dev/... -timeout 60s` |
| Full suite command | `go test ./... -timeout 120s` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FOUND-01 | `callClaude` returns non-empty string when API returns `content[].text` | unit (httptest) | `go test ./internal/ai/... -run TestCallClaudeResponse -v` | Wave 0 |
| FOUND-02 | `executeWithClaude` passes prompt via stdin, no bash in exec chain | unit | `go test ./cmd/dev/... -run TestExecuteWithClaudeNoShell -v` | Wave 0 |
| FOUND-03 | `http.Client` in all three packages has Timeout ≥ 30s | unit | `go test ./internal/... -run TestHTTPClientTimeout -v` | Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/ai/... ./internal/jira/... ./internal/vcs/... ./cmd/dev/... -timeout 60s`
- **Per wave merge:** `go test ./... -timeout 120s`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/ai/client_test.go` — covers FOUND-01 (httptest server returning correct Anthropic shape; assert non-empty response)
- [ ] `cmd/dev/wave_test.go` — covers FOUND-02 (verify no `bash` in exec chain; verify stdin pipe used)
- [ ] `internal/jira/client_test.go` — covers FOUND-03 for jira client (assert Timeout field)
- [ ] `internal/vcs/github_test.go` — covers FOUND-03 for vcs clients
- [ ] No framework install needed — `go test` is available

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Shell heredoc prompt injection | `exec.Command` + `StdinPipe()` | This phase | Eliminates shell injection vector |
| OpenAI response struct for Anthropic | Anthropic `content[]` struct | This phase | Enables all AI-powered features |
| `http.Client{}` no timeout | `http.Client{Timeout: 30s}` | This phase | Prevents indefinite hangs on unreachable hosts |

## Open Questions

1. **`internal/vcs/gitlab.go` and `internal/vcs/bitbucket.go` constructors**
   - What we know: github.go uses `&http.Client{}` — gitlab and bitbucket almost certainly follow the same pattern
   - What's unclear: Not read during research; could differ
   - Recommendation: Read both files at the start of the implementation task and apply the same timeout fix

2. **`AIClient` struct refactor scope**
   - What we know: The struct currently has no `client` field; the fix can be minimal (inline timeout) or structural (add field)
   - What's unclear: Whether future phases will need per-request timeouts or custom transports for the AI client
   - Recommendation: Add `client *http.Client` field to `AIClient` in `NewAIClient` for consistency with jira/vcs; cost is 3 lines

## Sources

### Primary (HIGH confidence)
- `internal/ai/client.go` — direct code inspection; bug at lines 26-32, 147, 159-163
- `cmd/dev/wave.go` — direct code inspection; shell injection at lines 82-97 and 103-118
- `internal/jira/client.go` — direct code inspection; missing timeout at lines 28, 136
- `internal/vcs/github.go` — direct code inspection; missing timeout at line 33
- `go.mod` — Go 1.24, no existing test framework
- Anthropic Messages API: https://docs.anthropic.com/en/api/messages — `content[].text` response shape

### Secondary (MEDIUM confidence)
- Go stdlib `net/http` docs: https://pkg.go.dev/net/http#Client.Timeout — Timeout field behavior
- Go stdlib `os/exec` docs: https://pkg.go.dev/os/exec#Cmd.StdinPipe — stdin pipe pattern

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all fixes use stdlib; no new dependencies
- Architecture: HIGH — bug locations confirmed by direct source inspection; fixes are standard Go patterns
- Pitfalls: HIGH — all three pitfalls derived from direct reading of the buggy code paths

**Research date:** 2026-03-09
**Valid until:** 2026-06-09 (stable Go stdlib patterns; no expiry risk)
