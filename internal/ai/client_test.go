package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"specforge/internal/config"
)

// minimalAIConfig returns a minimal AIConfig for white-box tests.
func minimalAIConfig() config.AIConfig {
	return config.AIConfig{
		Provider:  "claude",
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 100,
	}
}

// setAPIKey sets ANTHROPIC_API_KEY for the duration of a test and restores the
// original value on cleanup. Compatible with t.Parallel().
func setAPIKey(t *testing.T, value string) {
	t.Helper()
	orig, had := os.LookupEnv("ANTHROPIC_API_KEY")
	if err := os.Setenv("ANTHROPIC_API_KEY", value); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if had {
			_ = os.Setenv("ANTHROPIC_API_KEY", orig)
		} else {
			_ = os.Unsetenv("ANTHROPIC_API_KEY")
		}
	})
}

// TestCallClaudeResponseDecoding verifies that callClaude correctly parses the
// Anthropic API response format {"content":[{"type":"text","text":"..."}]}.
//
// RED state: the current production code reads result.Choices (OpenAI format)
// so it always returns a "no response from Claude" error — this test MUST fail
// until FOUND-01 is fixed.
func TestCallClaudeResponseDecoding(t *testing.T) {
	t.Parallel()

	// Spin up a local server returning an Anthropic-shaped response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": "hello world"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Override the package-level URL so callClaude hits our test server.
	origURL := claudeBaseURL
	claudeBaseURL = ts.URL
	t.Cleanup(func() { claudeBaseURL = origURL })

	setAPIKey(t, "test-key")

	c := NewAIClient(minimalAIConfig())
	response, err := c.callClaude("test prompt")
	if err != nil {
		t.Fatalf("callClaude returned error: %v", err)
	}
	if len(response) == 0 {
		t.Fatal("callClaude returned empty response; want non-empty (FOUND-01: struct uses choices[] instead of content[])")
	}
	if !strings.Contains(response, "hello world") {
		t.Errorf("callClaude response = %q; want it to contain %q", response, "hello world")
	}
}

// TestAIClientTimeout verifies that the HTTP client used by callClaude
// respects a configurable timeout so slow servers cannot hang the process.
//
// RED state: the current production code creates &http.Client{} with Timeout=0
// (no timeout). This test injects a 1 s timeout via newHTTPClient, starts a
// server that sleeps 5 s, and expects callClaude to fail with a
// timeout/deadline error within 3 s.
func TestAIClientTimeout(t *testing.T) {
	t.Parallel()

	// Server that deliberately takes longer than the injected client timeout.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	origURL := claudeBaseURL
	claudeBaseURL = ts.URL
	t.Cleanup(func() { claudeBaseURL = origURL })

	origFactory := newHTTPClient
	newHTTPClient = func() *http.Client { return &http.Client{Timeout: 1 * time.Second} }
	t.Cleanup(func() { newHTTPClient = origFactory })

	setAPIKey(t, "test-key")

	c := NewAIClient(minimalAIConfig())

	done := make(chan error, 1)
	go func() {
		_, err := c.callClaude("test prompt")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("callClaude succeeded but should have returned a timeout error (FOUND-03)")
		}
		errStr := strings.ToLower(err.Error())
		if !strings.Contains(errStr, "timeout") && !strings.Contains(errStr, "deadline") && !strings.Contains(errStr, "context") {
			t.Errorf("expected a timeout/deadline error; got: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("callClaude did not respect the HTTP timeout — still blocking after 3 s (FOUND-03: no timeout on http.Client)")
	}
}
