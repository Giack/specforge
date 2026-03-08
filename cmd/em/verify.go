package em

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunInteractiveUAT() error {
	criteria := []string{
		"User can log in with email and password",
		"Login shows loading spinner during authentication",
		"Invalid credentials show error message",
		"Successful login redirects to dashboard",
		"Session persists on page refresh",
	}

	passed := 0
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🧪 UAT Verification - Check each criteria manually")
	fmt.Println(strings.Repeat("─", 50))

	for i, criterion := range criteria {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(criteria), criterion)
		fmt.Print("Passed? (y/n): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			passed++
			fmt.Println("  ✓ PASSED")
		} else {
			fmt.Println("  ✗ FAILED")
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Printf("✅ UAT Results: %d/%d passed\n", passed, len(criteria))

	if passed == len(criteria) {
		fmt.Println("🎉 All tests passed! Feature is ready for release.")
	} else {
		fmt.Printf("⚠️  %d test(s) failed. Run 'specforge em bug' to report issues.\n", len(criteria)-passed)
	}

	return nil
}
