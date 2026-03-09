package vcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"specforge/internal/config"
)

type GitHubClient struct {
	config config.GitHubConfig
	client *http.Client
}

type GitHubPRRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

type GitHubPRResponse struct {
	Number  int    `json:"number"`
	HTMLURL string `json:"html_url"`
	State   string `json:"state"`
}

func NewGitHubClient(cfg config.GitHubConfig) *GitHubClient {
	return &GitHubClient{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *GitHubClient) GetProviderName() string {
	return "GitHub"
}

func (g *GitHubClient) CreatePullRequest(repoSlug, title, sourceBranch, targetBranch string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", g.config.Owner, repoSlug)

	pr := GitHubPRRequest{
		Title: title,
		Body:  "Created via SpecForge",
		Head:  sourceBranch,
		Base:  targetBranch,
	}

	body, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+g.config.Token)
	req.Header.Add("Accept", "application/vnd.github+json")
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
		return "", fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result GitHubPRResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.HTMLURL, nil
}
