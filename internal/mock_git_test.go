package internal

import (
	"os"
	"path/filepath"
	"testing"
)

// mockRepository creates a mock git repository in a temporary directory
// for testing the repository cloning functionality without using actual git operations
func mockRepository(t *testing.T) (string, string) {
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

// TestProcessRepositoryWithMock tests ProcessRepository with a mock repository
func TestProcessRepositoryWithMock(t *testing.T) {
	// Create a mock repository and destination directory
	repoDir, destDir := mockRepository(t)
	defer os.RemoveAll(repoDir)
	defer os.RemoveAll(destDir)

	// Save current working directory and chdir to destination for testing
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer os.Chdir(origWd) // Restore working directory when done

	if err := os.Chdir(destDir); err != nil {
		t.Fatalf("failed to change to destination directory: %v", err)
	}

	// Redirect stdout to avoid cluttering test output
	originalStdout := os.Stdout
	nullFile, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("failed to open devNull: %v", err)
	}
	os.Stdout = nullFile
	defer func() {
		os.Stdout = originalStdout
		nullFile.Close()
	}()

	// This test uses a mock repository, but cannot test the ProcessRepository directly
	// because it relies on real git operations. Instead, we test the copyDotFiles function
	// which is the main functionality after cloning.
	if err := copyDotFiles(repoDir, destDir); err != nil {
		t.Fatalf("copyDotFiles failed: %v", err)
	}

	// Verify only dotfiles were copied
	entries, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatalf("failed to read destination directory: %v", err)
	}

	// Check for expected dotfiles
	hasDotGitconfig := false
	hasDotBashrc := false
	hasDotVscode := false
	hasReadme := false

	for _, entry := range entries {
		name := entry.Name()
		switch name {
		case ".gitconfig":
			hasDotGitconfig = true
		case ".bashrc":
			hasDotBashrc = true
		case ".vscode":
			hasDotVscode = true
		case "README.md":
			hasReadme = true
		}
	}

	if !hasDotGitconfig {
		t.Errorf("copyDotFiles did not copy .gitconfig")
	}
	if !hasDotBashrc {
		t.Errorf("copyDotFiles did not copy .bashrc")
	}
	if !hasDotVscode {
		t.Errorf("copyDotFiles did not copy .vscode directory")
	}
	if hasReadme {
		t.Errorf("copyDotFiles incorrectly copied README.md, which is not a dotfile")
	}

	// Check content of copied dotfiles
	gitconfigContent, err := os.ReadFile(filepath.Join(destDir, ".gitconfig"))
	if err != nil {
		t.Fatalf("failed to read copied .gitconfig: %v", err)
	}
	if string(gitconfigContent) != "git config content" {
		t.Errorf("copyDotFiles did not copy .gitconfig content correctly, got %q, want %q", string(gitconfigContent), "git config content")
	}

	// Check .vscode directory contents
	vscodeSettingsContent, err := os.ReadFile(filepath.Join(destDir, ".vscode", "settings.json"))
	if err != nil {
		t.Fatalf("failed to read copied .vscode/settings.json: %v", err)
	}
	if string(vscodeSettingsContent) != "vscode settings" {
		t.Errorf("copyDotFiles did not copy .vscode/settings.json content correctly, got %q, want %q", string(vscodeSettingsContent), "vscode settings")
	}
}
