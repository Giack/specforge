package jira

import (
	"testing"
	"time"

	"specforge/internal/config"
)

const minTimeout = 30 * time.Second

// TestJiraClientTimeout verifies that the http.Client embedded in JiraClient
// has a Timeout of at least 30 s so network calls cannot hang indefinitely.
//
// RED state: NewJiraClient currently uses &http.Client{} with Timeout=0.
// This test MUST fail until FOUND-03 is fixed.
func TestJiraClientTimeout(t *testing.T) {
	t.Parallel()

	cfg := config.AtlassianConfig{
		Domain:     "example",
		Email:      "test@example.com",
		APIToken:   "token",
		ProjectKey: "TEST",
	}

	c := NewJiraClient(cfg)

	if c.client == nil {
		t.Fatal("JiraClient.client is nil; want a non-nil *http.Client")
	}
	if c.client.Timeout < minTimeout {
		t.Errorf("JiraClient.client.Timeout = %v; want >= %v (FOUND-03: no HTTP timeout)", c.client.Timeout, minTimeout)
	}
}

// TestConfluenceClientTimeout verifies that the http.Client embedded in
// ConfluenceClient has a Timeout of at least 30 s.
//
// RED state: NewConfluenceClient currently uses &http.Client{} with Timeout=0.
// This test MUST fail until FOUND-03 is fixed.
func TestConfluenceClientTimeout(t *testing.T) {
	t.Parallel()

	cfg := config.AtlassianConfig{
		Domain:     "example",
		Email:      "test@example.com",
		APIToken:   "token",
		ProjectKey: "TEST",
	}

	c := NewConfluenceClient(cfg)

	if c.client == nil {
		t.Fatal("ConfluenceClient.client is nil; want a non-nil *http.Client")
	}
	if c.client.Timeout < minTimeout {
		t.Errorf("ConfluenceClient.client.Timeout = %v; want >= %v (FOUND-03: no HTTP timeout)", c.client.Timeout, minTimeout)
	}
}
