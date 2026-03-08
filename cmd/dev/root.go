package dev

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"specforge/internal/ai"
	"specforge/internal/config"
	"specforge/internal/jira"

	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Developer commands",
		Long:  "Execute specs, discuss context, and manage fixes",
	}

	cmd.AddCommand(startCmd(cfg))
	cmd.AddCommand(discussCmd(cfg))
	cmd.AddCommand(planCmd(cfg))
	cmd.AddCommand(executeCmd(cfg))
	cmd.AddCommand(fixCmd(cfg))
	cmd.AddCommand(prCmd(cfg))
	cmd.AddCommand(subAgentCmd(cfg))
	cmd.AddCommand(waveCmd(cfg))

	return cmd
}

func startCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the full Dev workflow (discuss -> plan -> execute)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🚀 Starting Dev Workflow...")
			fmt.Println()
			fmt.Println("This will guide you through:")
			fmt.Println("  1. /specforge:dev-discuss - Capture implementation context")
			fmt.Println("  2. /specforge:dev-plan    - Create atomic execution plans")
			fmt.Println("  3. /specforge:dev-execute - Run plans with AI agents")
			fmt.Println()
			fmt.Println("Run these commands inside Claude Code or OpenCode:")
			fmt.Println("  /specforge:dev-discuss")
			fmt.Println("  /specforge:dev-plan")
			fmt.Println("  /specforge:dev-execute")

			return nil
		},
	}
}

func discussCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "discuss",
		Short: "Capture implementation context and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("💬 Starting Discuss Phase...")
			fmt.Println("This command should be run inside Claude Code/OpenCode:")
			fmt.Println("  /specforge:dev-discuss")

			return runInteractiveDiscuss()
		},
	}
}

func runInteractiveDiscuss() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n📝 Implementation Context - Answer these questions:")
	fmt.Println(strings.Repeat("─", 50))

	answers := map[string]string{}

	questions := []string{
		"Technology stack (e.g., React, Go, PostgreSQL)",
		"Framework preferences (e.g., Next.js, Gin, Express)",
		"State management approach",
		"API design preferences (REST, GraphQL, gRPC)",
		"Testing strategy (unit, integration, e2e)",
		"Coding conventions or style guide",
		"Any specific libraries to use or avoid",
	}

	for i, question := range questions {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(questions), question)
		fmt.Print("Answer: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		answers[question] = strings.TrimSpace(input)
	}

	content := "# Implementation Context\n\n"
	content += "This document captures the implementation preferences for the current development cycle.\n\n"

	for question, answer := range answers {
		if answer != "" {
			content += fmt.Sprintf("## %s\n\n%s\n\n", question, answer)
		}
	}

	if err := os.WriteFile("CONTEXT.md", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write CONTEXT.md: %w", err)
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("✅ CONTEXT.md created!")
	fmt.Println("   This file will be used by the planner and executors.")

	return nil
}

func planCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Create atomic execution plans from roadmap",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📋 Starting Plan Phase...")
			fmt.Println("This command should be run inside Claude Code/OpenCode:")
			fmt.Println("  /specforge:dev-plan")

			return generatePlans()
		},
	}
}

func generatePlans() error {
	roadmapContent, err := os.ReadFile("ROADMAP.md")
	if err != nil {
		return fmt.Errorf("failed to read ROADMAP.md: %w", err)
	}

	contextContent, _ := os.ReadFile("CONTEXT.md")

	content := fmt.Sprintf(`# Execution Plans

This file contains the atomic execution plans for the current phase.

## Context
%s

## Roadmap
%s

---

## Tasks

### Task 1: [Task Name]
- **Files to modify**: 
- **Description**: 
- **Tests to write**: 
- **Verification**: 

### Task 2: [Task Name]
- **Files to modify**: 
- **Description**: 
- **Tests to write**: 
- **Verification**: 

---

Each task will be executed in a fresh context to prevent context rot.
`, string(contextContent), string(roadmapContent))

	if err := os.WriteFile("PLANS.md", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write PLANS.md: %w", err)
	}

	fmt.Println("✅ PLANS.md created!")
	fmt.Println("   Run 'specforge dev execute' to start execution.")

	return nil
}

func executeCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "execute",
		Short: "Execute plans using sub-agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("⚡ Starting Execute Phase...")
			fmt.Println("This command should be run inside Claude Code/OpenCode:")
			fmt.Println("  /specforge:dev-execute")

			return executePlans()
		},
	}
}

func executePlans() error {
	plansContent, err := os.ReadFile("PLANS.md")
	if err != nil {
		return fmt.Errorf("failed to read PLANS.md: %w", err)
	}

	fmt.Println("\n📦 Loading execution plans...")
	fmt.Println(string(plansContent))

	fmt.Println("\n⚠️  Execute phase should be run inside Claude Code/OpenCode")
	fmt.Println("   to leverage sub-agent parallelization.")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  /specforge:dev-discuss  - Setup context")
	fmt.Println("  /specforge:dev-plan    - Create plans")
	fmt.Println("  /specforge:dev-execute - Run with AI agents")

	return nil
}

func fixCmd(cfg *config.Config) *cobra.Command {
	var jiraID string
	var analyzeOnly bool

	cmd := &cobra.Command{
		Use:   "fix",
		Short: "Fix a Jira bug using AI debugger",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jiraID = args[0]
			return runBugFix(cfg, jiraID, analyzeOnly)
		},
	}

	cmd.Flags().BoolVar(&analyzeOnly, "analyze", false, "Only analyze the bug, don't fix")

	return cmd
}

func runBugFix(cfg *config.Config, jiraID string, analyzeOnly bool) error {
	fmt.Printf("🔧 Starting Fix for Jira %s...\n", jiraID)

	jiraClient := jira.NewJiraClient(cfg.Atlassian)

	fmt.Println("📥 Fetching bug details from Jira...")
	bugDescription, err := jiraClient.GetIssueDescription(jiraID)
	if err != nil {
		return fmt.Errorf("failed to fetch Jira issue: %w", err)
	}

	fmt.Println("✅ Fetched bug details")

	codeContext := ""
	if files := getRelevantFiles(); len(files) > 0 {
		fmt.Println("📂 Reading relevant code files...")
		for _, f := range files {
			if content, err := os.ReadFile(f); err == nil {
				codeContext += fmt.Sprintf("\n\n=== %s ===\n%s", f, string(content))
			}
		}
	}

	aiClient := ai.NewAIClient(cfg.AI)

	fmt.Println("🤖 Analyzing bug and generating fix plan...")
	fixPlan, err := aiClient.AnalyzeAndFix(bugDescription, codeContext)
	if err != nil {
		fmt.Printf("⚠️  AI analysis failed: %v\n", err)
		fixPlan = fmt.Sprintf("# Fix Plan for %s\n\n## Bug Description\n%s\n\n## Code Context\n%s\n\n(No AI analysis available)", jiraID, bugDescription, codeContext)
	}

	fixPlanPath := fmt.Sprintf("FIX-%s.md", jiraID)
	if err := os.WriteFile(fixPlanPath, []byte(fixPlan), 0644); err != nil {
		return fmt.Errorf("failed to write fix plan: %w", err)
	}

	fmt.Printf("   ✓ Saved to %s\n", fixPlanPath)

	if analyzeOnly {
		fmt.Println("\n✅ Analysis complete! Review FIX-" + jiraID + ".md for the fix plan.")
		return nil
	}

	fmt.Println("🤖 Executing fix...")

	executeFix, err := aiClient.ExecuteFix(fixPlan)
	if err != nil {
		fmt.Printf("⚠️  AI execution failed: %v\n", err)
		fmt.Println("\n✅ Fix plan created! Execute manually inside Claude Code/OpenCode:")
		fmt.Printf("   /specforge:dev-fix %s\n", jiraID)
		return nil
	}

	summaryPath := fmt.Sprintf("FIX-%s-SUMMARY.md", jiraID)
	if err := os.WriteFile(summaryPath, []byte(executeFix), 0644); err != nil {
		return fmt.Errorf("failed to write fix summary: %w", err)
	}

	fmt.Printf("   ✓ Fix executed, summary saved to %s\n", summaryPath)

	fmt.Println("\n✅ Bug fix completed!")
	fmt.Printf("   Review: %s\n", summaryPath)
	fmt.Println("   Commit your changes and update the Jira ticket.")

	return nil
}

func getRelevantFiles() []string {
	var files []string

	entries, err := os.ReadDir(".")
	if err != nil {
		return files
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, ".ts") ||
			strings.HasSuffix(name, ".tsx") ||
			strings.HasSuffix(name, ".js") ||
			strings.HasSuffix(name, ".jsx") {
			files = append(files, name)
		}
	}

	return files
}
