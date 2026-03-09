package vcs

import (
	"testing"
	"time"

	"specforge/internal/config"
)

const minTimeout = 30 * time.Second

// TestGitHubClientTimeout verifies that the http.Client embedded in
// GitHubClient has a Timeout of at least 30 s.
//
// RED state: NewGitHubClient currently uses &http.Client{} with Timeout=0.
// This test MUST fail until FOUND-03 is fixed.
func TestGitHubClientTimeout(t *testing.T) {
	t.Parallel()

	cfg := config.GitHubConfig{
		Token: "test-token",
		Owner: "test-owner",
		Repo:  "test-repo",
	}

	c := NewGitHubClient(cfg)

	if c.client == nil {
		t.Fatal("GitHubClient.client is nil; want a non-nil *http.Client")
	}
	if c.client.Timeout < minTimeout {
		t.Errorf("GitHubClient.client.Timeout = %v; want >= %v (FOUND-03: no HTTP timeout)", c.client.Timeout, minTimeout)
	}
}

// TestGitLabClientTimeout verifies that the http.Client embedded in
// GitLabClient has a Timeout of at least 30 s.
//
// RED state: NewGitLabClient currently uses &http.Client{} with Timeout=0.
// This test MUST fail until FOUND-03 is fixed.
func TestGitLabClientTimeout(t *testing.T) {
	t.Parallel()

	cfg := config.GitLabConfig{
		Domain: "gitlab.example.com",
		Token:  "test-token",
	}

	c := NewGitLabClient(cfg)

	if c.client == nil {
		t.Fatal("GitLabClient.client is nil; want a non-nil *http.Client")
	}
	if c.client.Timeout < minTimeout {
		t.Errorf("GitLabClient.client.Timeout = %v; want >= %v (FOUND-03: no HTTP timeout)", c.client.Timeout, minTimeout)
	}
}

// TestBitbucketClientTimeout verifies that the http.Client embedded in
// BitbucketClient has a Timeout of at least 30 s.
//
// RED state: NewBitbucketClient currently uses &http.Client{} with Timeout=0.
// This test MUST fail until FOUND-03 is fixed.
func TestBitbucketClientTimeout(t *testing.T) {
	t.Parallel()

	cfg := config.BitbucketConfig{
		Domain:      "bitbucket.org",
		Username:    "test-user",
		AppPassword: "test-pass",
	}

	c := NewBitbucketClient(cfg)

	if c.client == nil {
		t.Fatal("BitbucketClient.client is nil; want a non-nil *http.Client")
	}
	if c.client.Timeout < minTimeout {
		t.Errorf("BitbucketClient.client.Timeout = %v; want >= %v (FOUND-03: no HTTP timeout)", c.client.Timeout, minTimeout)
	}
}
