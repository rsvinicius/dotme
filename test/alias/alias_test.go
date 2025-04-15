package alias_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rsvinicius/dotme/internal/alias"
)

func setupTestConfig(t *testing.T) (string, func()) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "dotme-alias-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	// Store original HOME environment variable
	origHome := os.Getenv("HOME")

	// Set HOME to the temp directory for the test
	os.Setenv("HOME", tempDir)

	// Create a cleanup function
	cleanup := func() {
		os.Setenv("HOME", origHome)
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestSaveAndGetRepo(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	testCases := []struct {
		name      string
		repoURL   string
		aliasName string
		wantErr   bool
	}{
		{
			name:      "valid alias",
			repoURL:   "https://github.com/user/repo",
			aliasName: "test-repo",
			wantErr:   false,
		},
		{
			name:      "duplicate alias",
			repoURL:   "https://github.com/user/repo2",
			aliasName: "test-repo",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Try to save the repo
			err := alias.SaveRepo(tc.repoURL, tc.aliasName)

			// Check error result against expectation
			if (err != nil) != tc.wantErr {
				t.Errorf("SaveRepo() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// If we don't expect an error, check if we can retrieve the repo
			if !tc.wantErr {
				gotURL, err := alias.GetRepo(tc.aliasName)
				if err != nil {
					t.Errorf("GetRepo() error = %v", err)
					return
				}
				if gotURL != tc.repoURL {
					t.Errorf("GetRepo() = %v, want %v", gotURL, tc.repoURL)
				}
			}
		})
	}
}

func TestGetRepoNotFound(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	_, err := alias.GetRepo("nonexistent")
	if err != alias.ErrAliasNotFound {
		t.Errorf("GetRepo() error = %v, want %v", err, alias.ErrAliasNotFound)
	}
}

func TestListAliases(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Save some test repositories
	testRepos := map[string]string{
		"alias1": "https://github.com/user/repo1",
		"alias2": "https://github.com/user/repo2",
		"alias3": "https://github.com/user/repo3",
	}

	for a, url := range testRepos {
		if err := alias.SaveRepo(url, a); err != nil {
			t.Fatalf("SaveRepo() error = %v", err)
		}
	}

	// List all aliases
	got, err := alias.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}

	// Check that all test repos are in the list
	for a, url := range testRepos {
		if got[a] != url {
			t.Errorf("ListAliases() = %v, want %v for alias %s", got[a], url, a)
		}
	}

	// Check that no extra repos are in the list
	if len(got) != len(testRepos) {
		t.Errorf("ListAliases() returned %d aliases, want %d", len(got), len(testRepos))
	}
}

func TestDeleteAlias(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Save a test repository
	aliasName := "test-repo"
	repoURL := "https://github.com/user/repo"

	if err := alias.SaveRepo(repoURL, aliasName); err != nil {
		t.Fatalf("SaveRepo() error = %v", err)
	}

	// Delete the alias
	if err := alias.DeleteAlias(aliasName); err != nil {
		t.Errorf("DeleteAlias() error = %v", err)
	}

	// Check that the alias is gone
	_, err := alias.GetRepo(aliasName)
	if err != alias.ErrAliasNotFound {
		t.Errorf("GetRepo() after deletion, error = %v, want %v", err, alias.ErrAliasNotFound)
	}

	// Try to delete nonexistent alias
	err = alias.DeleteAlias("nonexistent")
	if err != alias.ErrAliasNotFound {
		t.Errorf("DeleteAlias() error = %v, want %v", err, alias.ErrAliasNotFound)
	}
}

func TestConfigPath(t *testing.T) {
	tempDir, cleanup := setupTestConfig(t)
	defer cleanup()

	path, err := alias.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}

	// Check that the path is under the temp directory
	if !filepath.IsAbs(path) {
		t.Errorf("GetConfigPath() = %v, want absolute path", path)
	}

	// Check that the path contains the expected filename
	if filepath.Base(path) != "config.json" {
		t.Errorf("GetConfigPath() filename = %v, want config.json", filepath.Base(path))
	}

	// Check that the directory exists
	configDir := filepath.Dir(path)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("GetConfigPath() did not create directory %v", configDir)
	}

	// Check that the directory is in the temp (home) directory
	// Replacing deprecated filepath.HasPrefix with filepath.Rel
	rel, err := filepath.Rel(tempDir, configDir)
	if err != nil {
		t.Errorf("Failed to get relative path: %v", err)
	}
	if rel == ".." || strings.HasPrefix(rel, "../") {
		t.Errorf("GetConfigPath() directory = %v is not under temp dir %v", configDir, tempDir)
	}
}
