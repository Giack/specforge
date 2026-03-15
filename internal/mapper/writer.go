package mapper

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteDocument writes content to a file named docName inside outputDir.
// It creates outputDir (with mode 0755) if it does not exist, and writes
// the file with mode 0644.
func WriteDocument(outputDir, docName, content string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	path := filepath.Join(outputDir, docName)
	return os.WriteFile(path, []byte(content), 0644)
}
