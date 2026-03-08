package pm

import (
	"fmt"
	"net/url"
	"os"

	"specforge/internal/ai"
	"specforge/internal/config"
	"specforge/internal/jira"

	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pm",
		Short: "Product Manager commands",
		Long:  "Manage requirements and sync from Confluence/Jira",
	}

	cmd.AddCommand(syncCmd(cfg))
	return cmd
}

func syncCmd(cfg *config.Config) *cobra.Command {
	var sourceType string
	var sourceURL string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync requirements from Confluence or Jira and generate specs",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔄 Starting Spec Sync...")

			parsedURL, err := url.Parse(sourceURL)
			if err != nil {
				return fmt.Errorf("invalid URL: %w", err)
			}

			var content string

			switch sourceType {
			case "confluence":
				client := jira.NewConfluenceClient(cfg.Atlassian)
				content, err = client.GetPageContent(parsedURL.Path)
				if err != nil {
					return fmt.Errorf("failed to fetch Confluence page: %w", err)
				}
			case "jira":
				client := jira.NewJiraClient(cfg.Atlassian)
				issueKey := parsedURL.Query().Get("issueKey")
				if issueKey == "" {
					return fmt.Errorf("missing issueKey query parameter")
				}
				content, err = client.GetIssueDescription(issueKey)
				if err != nil {
					return fmt.Errorf("failed to fetch Jira issue: %w", err)
				}
			default:
				return fmt.Errorf("unsupported source type: %s (use confluence or jira)", sourceType)
			}

			fmt.Printf("✅ Fetched %d characters from %s\n", len(content), sourceType)

			if outputDir == "" {
				outputDir = "."
			}

			aiClient := ai.NewAIClient(cfg.AI)

			fmt.Println("🤖 Generating PROJECT.md...")
			projectSpec, err := aiClient.GenerateSpec(content, "project")
			if err != nil {
				fmt.Printf("⚠️  AI generation failed: %v\n", err)
				fmt.Println("   Saving raw content instead...")
				projectSpec = "# Project\n\n" + content
			}

			projectPath := outputDir + "/PROJECT.md"
			if err := os.WriteFile(projectPath, []byte(projectSpec), 0644); err != nil {
				return fmt.Errorf("failed to write PROJECT.md: %w", err)
			}
			fmt.Printf("   ✓ Saved to %s\n", projectPath)

			fmt.Println("🤖 Generating REQUIREMENTS.md...")
			requirementsSpec, err := aiClient.GenerateSpec(content, "requirements")
			if err != nil {
				fmt.Printf("⚠️  AI generation failed: %v\n", err)
				fmt.Println("   Saving raw content instead...")
				requirementsSpec = "# Requirements\n\n" + content
			}

			requirementsPath := outputDir + "/REQUIREMENTS.md"
			if err := os.WriteFile(requirementsPath, []byte(requirementsSpec), 0644); err != nil {
				return fmt.Errorf("failed to write REQUIREMENTS.md: %w", err)
			}
			fmt.Printf("   ✓ Saved to %s\n", requirementsPath)

			fmt.Println("\n✅ Spec files generated!")
			fmt.Println("   - PROJECT.md")
			fmt.Println("   - REQUIREMENTS.md")

			return nil
		},
	}

	cmd.Flags().StringVar(&sourceType, "type", "confluence", "Source type: confluence or jira")
	cmd.Flags().StringVar(&sourceURL, "url", "", "URL to sync from")
	cmd.Flags().StringVar(&outputDir, "output", ".", "Output directory for generated files")
	cmd.MarkFlagRequired("url")

	return cmd
}
