package internal

import (
	"fmt"
	"os"

	"github.com/rsvinicius/dotme/internal/alias"
	"github.com/rsvinicius/dotme/internal/fs"
	"github.com/rsvinicius/dotme/internal/git"
	"github.com/rsvinicius/dotme/internal/patterns"
)

// ProcessRepository handles cloning the repository and copying dotfiles
func ProcessRepository(repoURL string, includePatterns, excludePatterns []string) error {
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

	// Create filter options
	filterOptions := &patterns.FilterOptions{
		IncludePatterns: includePatterns,
		ExcludePatterns: excludePatterns,
	}

	// If no patterns provided via command line, try to load defaults from config
	if len(includePatterns) == 0 && len(excludePatterns) == 0 {
		defaultPatterns, err := alias.GetDefaultPatterns()
		if err == nil {
			filterOptions.IncludePatterns = defaultPatterns.IncludePatterns
			filterOptions.ExcludePatterns = defaultPatterns.ExcludePatterns
		}
		// If error loading defaults, continue with empty patterns (default behavior)
	}

	// Process files from the temporary directory
	return fs.CopyDotFiles(tempDir, destDir, filterOptions)
}