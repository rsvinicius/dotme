package git

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
)

// CloneRepository clones a Git repository to a local temporary directory
func CloneRepository(repoURL string) (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "dotme-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	fmt.Printf("🔄 Cloning repository: %s\n", repoURL)

	// Clone the repository
	r, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Get repository information for detailed output
	ref, err := r.Head()
	if err == nil {
		fmt.Printf("✅ Repository cloned, using branch: %s\n", ref.Name().Short())
	}

	return tempDir, nil
}
