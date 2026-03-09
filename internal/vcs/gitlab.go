package vcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"specforge/internal/config"
)

type GitLabClient struct {
	config config.GitLabConfig
	client *http.Client
}

type GitLabMRRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
}

type GitLabMRResponse struct {
	IID    int    `json:"iid"`
	WebURL string `json:"web_url"`
	State  string `json:"state"`
}

func NewGitLabClient(cfg config.GitLabConfig) *GitLabClient {
	return &GitLabClient{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GitLabClient) GetProviderName() string {
	return "GitLab"
}

func (g *GitLabClient) CreatePullRequest(repoSlug, title, sourceBranch, targetBranch string) (string, error) {
	projectID := url.QueryEscape(repoSlug)
	apiURL := fmt.Sprintf("https://%s/api/v4/projects/%s/merge_requests", g.config.Domain, projectID)

	mr := GitLabMRRequest{
		Title:        title,
		Description:  "Created via SpecForge",
		SourceBranch: sourceBranch,
		TargetBranch: targetBranch,
	}

	body, err := json.Marshal(mr)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("PRIVATE-TOKEN", g.config.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("GitLab API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result GitLabMRResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.WebURL, nil
}
