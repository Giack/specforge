package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"specforge/internal/config"
)

// claudeBaseURL is the Anthropic API endpoint. Tests can override this to point
// at a local httptest.Server.
var claudeBaseURL = "https://api.anthropic.com/v1/messages"

// httpClientTimeout is the default timeout applied to every HTTP client created
// by callClaude. Tests can override this to a shorter value.
var httpClientTimeout = 30 * time.Second

// newHTTPClient is a factory used by callClaude to create an *http.Client.
// Tests can replace this to inject a client with a shorter timeout or other
// behaviour without modifying production logic.
var newHTTPClient = func() *http.Client {
	return &http.Client{Timeout: httpClientTimeout}
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type AIClient struct {
	config config.AIConfig
}

func NewAIClient(cfg config.AIConfig) *AIClient {
	return &AIClient{config: cfg}
}

func (a *AIClient) GenerateSpec(content string, specType string) (string, error) {
	prompt := a.buildPrompt(content, specType)

	switch a.config.Provider {
	case "claude":
		return a.callClaude(prompt)
	case "opencode":
		return a.callOpenCode(prompt)
	default:
		return a.callClaude(prompt)
	}
}

func (a *AIClient) buildPrompt(content string, specType string) string {
	switch specType {
	case "project":
		return fmt.Sprintf(`Transform this content into a PROJECT.md file for SpecForge.

The PROJECT.md should include:
- Project Vision (what are we building and why)
- Problem Statement (the pain point we're solving)
- Target Users (who is this for)
- Success Metrics (how do we measure success)

Content:
%s

Output only the PROJECT.md content, no explanations.`, content)

	case "requirements":
		return fmt.Sprintf(`Transform this content into a REQUIREMENTS.md file for SpecForge.

The REQUIREMENTS.md should include:
- Functional Requirements (what the system must do)
- Non-Functional Requirements (performance, security, etc.)
- User Stories (who, what, why)
- Acceptance Criteria (how do we verify it works)
- Phase Mapping (which phase implements what)

Content:
%s

Output only the REQUIREMENTS.md content, no explanations.`, content)

	case "roadmap":
		return fmt.Sprintf(`Transform this into a ROADMAP.md for SpecForge.

The ROADMAP.md should include:
- Phases (logical groupings of work)
- Each phase should have clear deliverables
- Dependencies between phases
- Timeline estimation

Content:
%s

Output only the ROADMAP.md content, no explanations.`, content)

	case "architecture":
		return fmt.Sprintf(`Create an ARCHITECTURE.md with Mermaid.js diagrams.

Include:
- System context diagram
- Component diagram showing microservices/frontend
- Data flow diagram
- API contracts

Content:
%s

Output the ARCHITECTURE.md with Mermaid diagrams embedded.`, content)

	default:
		return content
	}
}

func (a *AIClient) callClaude(prompt string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	reqBody := ClaudeRequest{
		Model:     a.config.Model,
		MaxTokens: a.config.MaxTokens,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", claudeBaseURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	req.Header.Add("Content-Type", "application/json")

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}

	return result.Choices[0].Message.Content, nil
}

func (a *AIClient) callOpenCode(prompt string) (string, error) {
	cmd := exec.Command("opencode", "--print", prompt)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("opencode failed: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func (a *AIClient) GenerateMermaidDiagram(architecture string) (string, error) {
	prompt := fmt.Sprintf(`Create a Mermaid.js diagram showing the architecture.

Requirements:
- Use Mermaid.js syntax
- Include microservices, databases, and frontend
- Show data flow with arrows

Architecture description:
%s

Output only the Mermaid code block, no explanations.`, architecture)

	return a.GenerateSpec(prompt, "mermaid")
}

func (a *AIClient) AnalyzeAndFix(bugDescription, codeContext string) (string, error) {
	prompt := fmt.Sprintf(`You are a expert debugger. Analyze this bug and create a fix plan.

## Bug Description
%s

## Code Context
%s

## Your Task
1. Analyze the bug description and code context
2. Identify the root cause
3. Create a detailed fix plan with:
   - Root cause analysis
   - Files to modify
   - Code changes needed
   - Tests to add/verify

## Output Format
Return a FIX-PLAN.md with:
- ## Root Cause
- ## Fix Steps (numbered)
- ## Files to Modify
- ## Test Verification

Be specific and actionable.`, bugDescription, codeContext)

	return a.GenerateSpec(prompt, "debug")
}

func (a *AIClient) ExecuteFix(fixPlan string) (string, error) {
	prompt := fmt.Sprintf(`Execute the following fix plan. 

## Fix Plan
%s

## Your Task
1. Read the fix plan carefully
2. Make the necessary code changes
3. Write or update tests
4. Run tests to verify the fix
5. Create a git commit with the fix

## Constraints
- Follow the existing code style
- Write tests that verify the fix
- Commit message should reference the bug

Execute the fix now.`, fixPlan)

	return a.GenerateSpec(prompt, "fix")
}
