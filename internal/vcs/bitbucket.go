package vcs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"specforge/internal/config"
)

type BitbucketClient struct {
	config config.BitbucketConfig
	client *http.Client
}

type BitbucketPRRequest struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type BitbucketPRResponse struct {
	ID  int    `json:"id"`
	URL string `json:"links,omitempty"`
}

func NewBitbucketClient(cfg config.BitbucketConfig) *BitbucketClient {
	return &BitbucketClient{
		config: cfg,
		client: &http.Client{},
	}
}

func (b *BitbucketClient) GetProviderName() string {
	return "Bitbucket"
}

func (b *BitbucketClient) CreatePullRequest(repoSlug, title, sourceBranch, targetBranch string) (string, error) {
	workspace := b.config.Workspace
	if workspace == "" {
		parts := strings.Split(b.config.Domain, "/")
		workspace = parts[len(parts)-1]
	}

	apiURL := fmt.Sprintf("https://api.%s/2.0/repositories/%s/%s/pullrequests", b.config.Domain, workspace, repoSlug)

	pr := BitbucketPRRequest{
		Title:  title,
		Source: sourceBranch,
		Target: targetBranch,
	}

	body, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(b.config.Username + ":" + b.config.AppPassword))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result BitbucketPRResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s/pullrequests/%d", b.config.Domain, workspace, repoSlug, result.ID), nil
}
