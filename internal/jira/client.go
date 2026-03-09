package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"specforge/internal/config"
)

type JiraClient struct {
	config config.AtlassianConfig
	client *http.Client
}

type JiraIssue struct {
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
	} `json:"fields"`
}

func NewJiraClient(cfg config.AtlassianConfig) *JiraClient {
	return &JiraClient{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (j *JiraClient) GetIssueDescription(issueKey string) (string, error) {
	apiURL := fmt.Sprintf("https://%s.atlassian.net/rest/api/3/issue/%s", j.config.Domain, issueKey)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(j.config.Email + ":" + j.config.APIToken))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Accept", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var issue JiraIssue
	if err := json.Unmarshal(body, &issue); err != nil {
		return "", err
	}

	return issue.Fields.Description, nil
}

type CreateIssueRequest struct {
	Fields struct {
		Project     map[string]string `json:"project"`
		IssueType   string            `json:"issuetype"`
		Summary     string            `json:"summary"`
		Description string            `json:"description"`
	} `json:"fields"`
}

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (j *JiraClient) CreateBug(summary, description, issueType string) (string, error) {
	apiURL := fmt.Sprintf("https://%s.atlassian.net/rest/api/3/issue/", j.config.Domain)

	reqBody := CreateIssueRequest{}
	reqBody.Fields.Project = map[string]string{"key": j.config.ProjectKey}
	reqBody.Fields.IssueType = issueType
	reqBody.Fields.Summary = summary
	reqBody.Fields.Description = description

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(j.config.Email + ":" + j.config.APIToken))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := j.client.Do(req)
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

	var result CreateIssueResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.Key, nil
}

type ConfluenceClient struct {
	config config.AtlassianConfig
	client *http.Client
}

func NewConfluenceClient(cfg config.AtlassianConfig) *ConfluenceClient {
	return &ConfluenceClient{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *ConfluenceClient) GetPageContent(pageIDOrSlug string) (string, error) {
	pageID := strings.TrimPrefix(pageIDOrSlug, "/")

	apiURL := fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api/content/%s?expand=body.storage", c.config.Domain, pageID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.config.Email + ":" + c.config.APIToken))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	bodyStorage := result["body"].(map[string]interface{})["storage"].(map[string]interface{})["value"].(string)

	return bodyStorage, nil
}
