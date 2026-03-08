package em

import (
	"fmt"
	"os"

	"specforge/internal/ai"
	"specforge/internal/config"
	"specforge/internal/jira"

	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "em",
		Short: "Engineering Manager commands",
		Long:  "Architect roadmaps, verify features, and manage bugs",
	}

	cmd.AddCommand(architectCmd(cfg))
	cmd.AddCommand(verifyCmd(cfg))
	cmd.AddCommand(bugCmd(cfg))

	return cmd
}

func architectCmd(cfg *config.Config) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "architect",
		Short: "Analyze requirements and create domain roadmaps",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🏗️ Starting Architecture Analysis...")

			if outputDir == "" {
				outputDir = "."
			}

			requirementsPath := outputDir + "/REQUIREMENTS.md"
			content, err := os.ReadFile(requirementsPath)
			if err != nil {
				return fmt.Errorf("failed to read REQUIREMENTS.md: %w", err)
			}

			fmt.Println("📖 Read REQUIREMENTS.md")
			fmt.Println("🎨 Generating architecture with Mermaid diagrams...")

			aiClient := ai.NewAIClient(cfg.AI)

			architectureSpec, err := aiClient.GenerateSpec(string(content), "architecture")
			if err != nil {
				fmt.Printf("⚠️  AI generation failed: %v\n", err)
				architectureSpec = "# Architecture\n\n" + string(content)
			}

			archPath := outputDir + "/ARCHITECTURE.md"
			if err := os.WriteFile(archPath, []byte(architectureSpec), 0644); err != nil {
				return fmt.Errorf("failed to write ARCHITECTURE.md: %w", err)
			}
			fmt.Printf("   ✓ Saved to %s\n", archPath)

			fmt.Println("📋 Generating ROADMAP.md...")
			roadmapSpec, err := aiClient.GenerateSpec(string(content), "roadmap")
			if err != nil {
				fmt.Printf("⚠️  AI generation failed: %v\n", err)
				roadmapSpec = "# Roadmap\n\n" + string(content)
			}

			roadmapPath := outputDir + "/ROADMAP.md"
			if err := os.WriteFile(roadmapPath, []byte(roadmapSpec), 0644); err != nil {
				return fmt.Errorf("failed to write ROADMAP.md: %w", err)
			}
			fmt.Printf("   ✓ Saved to %s\n", roadmapPath)

			fmt.Println("\n✅ Architecture generated!")
			fmt.Println("   - ARCHITECTURE.md (with Mermaid.js diagrams)")
			fmt.Println("   - ROADMAP.md")

			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "output", ".", "Directory containing REQUIREMENTS.md")

	return cmd
}

func verifyCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Run UAT verification checklist",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🧪 Starting UAT Verification...")
			fmt.Println("📋 Loading acceptance criteria from REQUIREMENTS.md...")

			if err := RunInteractiveUAT(); err != nil {
				return fmt.Errorf("UAT failed: %w", err)
			}

			return nil
		},
	}
}

func bugCmd(cfg *config.Config) *cobra.Command {
	var feedback string

	cmd := &cobra.Command{
		Use:   "bug",
		Short: "Create Jira bug from UAT feedback",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🐞 Creating Jira Bug...")

			if feedback == "" {
				return fmt.Errorf("feedback is required (use --feedback flag)")
			}

			fmt.Printf("📝 Feedback: %s\n", feedback)

			client := jira.NewJiraClient(cfg.Atlassian)

			description := fmt.Sprintf("h3. Bug Report\n\n*Description:*\n%s\n\n*Reported from:* UAT Verification\n*Reporter:* %s", feedback, cfg.Atlassian.Email)

			issueKey, err := client.CreateBug(feedback, description, "Bug")
			if err != nil {
				return fmt.Errorf("failed to create Jira ticket: %w", err)
			}

			fmt.Printf("✅ Bug ticket created: %s\n", issueKey)
			fmt.Printf("🔗 View at: https://%s.atlassian.net/browse/%s\n", cfg.Atlassian.Domain, issueKey)

			return nil
		},
	}

	cmd.Flags().StringVar(&feedback, "feedback", "", "Natural language description of the bug")
	cmd.MarkFlagRequired("feedback")

	return cmd
}
