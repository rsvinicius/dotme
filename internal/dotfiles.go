package internal

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal/fs"
	"github.com/rsvinicius/dotme/internal/git"
)

// ProcessRepository handles cloning the repository and copying dotfiles
func ProcessRepository(repoURL string) error {
	// Clone the repository into a temporary directory
	tempDir, err := git.CloneRepository(repoURL)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Get current working directory
	destDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	fmt.Println("📋 Scanning for dotfiles...")

	// Process files from the temporary directory
	return fs.CopyDotFiles(tempDir, destDir)
}
