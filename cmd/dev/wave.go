package dev

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"specforge/internal/config"

	"github.com/spf13/cobra"
)

// execCommand is the factory used to create exec.Cmd values. Tests override
// this to capture command arguments without running real subprocesses.
var execCommand = exec.Command

func subAgentCmd(cfg *config.Config) *cobra.Command {
	var phase int

	cmd := &cobra.Command{
		Use:    "subagent",
		Short:  "Execute a single task using a sub-agent",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			return executeSubAgent(cfg, taskName, phase)
		},
	}

	cmd.Flags().IntVar(&phase, "phase", 1, "Phase number")

	return cmd
}

func executeSubAgent(cfg *config.Config, taskName string, phase int) error {
	fmt.Printf("🤖 Starting sub-agent for task: %s\n", taskName)

	planFile := fmt.Sprintf("phase-%d-%s-PLAN.md", phase, taskName)
	planContent, err := os.ReadFile(planFile)
	if err != nil {
		return fmt.Errorf("failed to read plan file: %w", err)
	}

	contextContent, _ := os.ReadFile("CONTEXT.md")

	prompt := fmt.Sprintf(`You are executing a task in a fresh context to prevent context rot.

## Task: %s
%s

## Context
%s

## Your Job
1. Read the plan carefully
2. Implement the changes
3. Write/update tests
4. Verify the implementation
5. Create a git commit

## Constraints
- Work in isolation - only use files mentioned in the plan
- Write tests that verify your implementation
- Commit with a descriptive message

Execute now!`, taskName, string(planContent), string(contextContent))

	switch cfg.AI.Provider {
	case "claude":
		return executeWithClaude(prompt, taskName)
	case "opencode":
		return executeWithOpenCode(prompt, taskName)
	default:
		return executeWithClaude(prompt, taskName)
	}
}

func executeWithClaude(prompt, taskName string) error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	fmt.Printf("Running task %s with Claude...\n", taskName)

	execCmd := execCommand("claude", "-p", "--dangerously-skip-permissions")
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
}

func executeWithOpenCode(prompt, taskName string) error {
	fmt.Printf("Running task %s with OpenCode...\n", taskName)

	execCmd := execCommand("opencode", "-p")
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("failed to start opencode: %w", err)
	}

	if _, err := io.WriteString(stdin, prompt); err != nil {
		return fmt.Errorf("failed to write prompt: %w", err)
	}
	stdin.Close()

	return execCmd.Wait()
}

func waveCmd(cfg *config.Config) *cobra.Command {
	var phase int

	cmd := &cobra.Command{
		Use:    "wave",
		Short:  "Execute tasks in waves (parallel execution)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeWave(cfg, phase)
		},
	}

	cmd.Flags().IntVar(&phase, "phase", 1, "Phase number")

	return cmd
}

func executeWave(cfg *config.Config, phase int) error {
	fmt.Printf("⚡ Starting wave execution for phase %d...\n", phase)

	taskDir := fmt.Sprintf(".planning/phase-%d/", phase)

	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return fmt.Errorf("failed to read task directory: %w", err)
	}

	var tasks []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "-PLAN.md") {
			taskName := strings.TrimSuffix(entry.Name(), "-PLAN.md")
			tasks = append(tasks, taskName)
		}
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found in wave")
		return nil
	}

	fmt.Printf("Found %d tasks to execute\n", len(tasks))

	for i, task := range tasks {
		fmt.Printf("\n[%d/%d] Executing: %s\n", i+1, len(tasks), task)
		if err := executeSubAgent(cfg, task, phase); err != nil {
			fmt.Printf("⚠️  Task %s failed: %v\n", task, err)
			continue
		}
		fmt.Printf("✅ Task %s completed\n", task)
	}

	fmt.Println("\n✅ Wave execution complete!")

	return nil
}
