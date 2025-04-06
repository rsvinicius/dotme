package mocks

import (
	"os"
	"path/filepath"
	"testing"
)

// MockRepository creates a mock git repository in a temporary directory
// for testing the repository cloning functionality without using actual git operations
func MockRepository(t *testing.T) (string, string) {
	// Create a temporary directory for the mock repository
	repoDir, err := os.MkdirTemp("", "dotme-mock-repo-*")
	if err != nil {
		t.Fatalf("failed to create mock repository directory: %v", err)
	}

	// Create a .git directory to mimic a git repository
	gitDir := filepath.Join(repoDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Create some dotfiles
	dotfiles := []struct {
		name    string
		content string
		isDir   bool
	}{
		{".gitconfig", "git config content", false},
		{".vscode", "", true},
		{".bashrc", "bash config content", false},
	}

	for _, f := range dotfiles {
		path := filepath.Join(repoDir, f.name)
		if f.isDir {
			if err := os.Mkdir(path, 0755); err != nil {
				t.Fatalf("failed to create mock dotfile directory %s: %v", path, err)
			}
		} else {
			if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
				t.Fatalf("failed to create mock dotfile %s: %v", path, err)
			}
		}
	}

	// Create a file in the .vscode directory
	vscodePath := filepath.Join(repoDir, ".vscode", "settings.json")
	if err := os.WriteFile(vscodePath, []byte("vscode settings"), 0644); err != nil {
		t.Fatalf("failed to create mock file in .vscode directory: %v", err)
	}

	// Create some non-dotfiles
	nonDotfiles := []struct {
		name    string
		content string
		isDir   bool
	}{
		{"README.md", "readme content", false},
		{"src", "", true},
	}

	for _, f := range nonDotfiles {
		path := filepath.Join(repoDir, f.name)
		if f.isDir {
			if err := os.Mkdir(path, 0755); err != nil {
				t.Fatalf("failed to create mock non-dotfile directory %s: %v", path, err)
			}
		} else {
			if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
				t.Fatalf("failed to create mock non-dotfile %s: %v", path, err)
			}
		}
	}

	// Create a temporary directory for destination
	destDir, err := os.MkdirTemp("", "dotme-mock-dest-*")
	if err != nil {
		t.Fatalf("failed to create mock destination directory: %v", err)
	}

	return repoDir, destDir
}
