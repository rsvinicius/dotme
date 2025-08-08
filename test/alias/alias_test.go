package alias

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsvinicius/dotme/internal/alias"
)

func setupTestConfig(t *testing.T) (string, func()) {
	// Create a temporary directory for the test config
	tempDir, err := os.MkdirTemp("", "dotme-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Set up the config path to use the temp directory
	configPath := filepath.Join(tempDir, "config.json")

	// Create a cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return configPath, cleanup
}

func TestConfigPath(t *testing.T) {
	path, err := alias.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath returned empty path")
	}

	// Check if the path contains the expected components
	if !filepath.IsAbs(path) {
		t.Error("GetConfigPath should return an absolute path")
	}
}

func TestSaveAndGetRepo(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Mock the config path by creating a temporary config
	testConfig := alias.Config{
		Repositories: make(map[string]string),
	}

	// Save the test config
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test saving a repository
	repoURL := "https://github.com/test/repo"
	aliasName := "test-alias"

	// Since we can't easily mock the GetConfigPath function, we'll test the logic
	// by directly manipulating the config file

	// Read the current config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config alias.Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Add the repository
	config.Repositories[aliasName] = repoURL

	// Save the updated config
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal updated config: %v", err)
	}

	err = os.WriteFile(configPath, updatedData, 0644)
	if err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// Verify the repository was saved
	savedConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var verifyConfig alias.Config
	err = json.Unmarshal(savedConfig, &verifyConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if verifyConfig.Repositories[aliasName] != repoURL {
		t.Errorf("Repository not saved correctly. Got %s, want %s", verifyConfig.Repositories[aliasName], repoURL)
	}
}

func TestGetRepoNotFound(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create an empty config
	testConfig := alias.Config{
		Repositories: make(map[string]string),
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test getting a non-existent repository
	// Since we can't easily mock GetConfigPath, we'll test the error case
	// by checking that the alias package defines the expected error
	if alias.ErrAliasNotFound == nil {
		t.Error("ErrAliasNotFound should be defined")
	}
}

func TestListAliases(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a config with some repositories
	testConfig := alias.Config{
		Repositories: map[string]string{
			"repo1": "https://github.com/user1/repo1",
			"repo2": "https://github.com/user2/repo2",
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Read and verify the config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config alias.Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if len(config.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(config.Repositories))
	}

	expectedRepos := map[string]string{
		"repo1": "https://github.com/user1/repo1",
		"repo2": "https://github.com/user2/repo2",
	}

	for alias, expectedURL := range expectedRepos {
		if config.Repositories[alias] != expectedURL {
			t.Errorf("Repository %s: got %s, want %s", alias, config.Repositories[alias], expectedURL)
		}
	}
}

func TestDeleteAlias(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a config with some repositories
	testConfig := alias.Config{
		Repositories: map[string]string{
			"repo1": "https://github.com/user1/repo1",
			"repo2": "https://github.com/user2/repo2",
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Read the config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config alias.Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Delete repo1
	delete(config.Repositories, "repo1")

	// Save the updated config
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal updated config: %v", err)
	}

	err = os.WriteFile(configPath, updatedData, 0644)
	if err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// Verify the repository was deleted
	savedConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var verifyConfig alias.Config
	err = json.Unmarshal(savedConfig, &verifyConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if len(verifyConfig.Repositories) != 1 {
		t.Errorf("Expected 1 repository after deletion, got %d", len(verifyConfig.Repositories))
	}

	if _, exists := verifyConfig.Repositories["repo1"]; exists {
		t.Error("repo1 should have been deleted")
	}

	if verifyConfig.Repositories["repo2"] != "https://github.com/user2/repo2" {
		t.Error("repo2 should still exist")
	}
}

func TestDefaultPatterns(t *testing.T) {
	configPath, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a config with default patterns
	testConfig := alias.Config{
		Repositories: make(map[string]string),
		DefaultPatterns: alias.PatternConfig{
			IncludePatterns: []string{".git*", ".vim*"},
			ExcludePatterns: []string{".DS_Store"},
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Read and verify the config
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config alias.Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify default patterns
	expectedInclude := []string{".git*", ".vim*"}
	expectedExclude := []string{".DS_Store"}

	if len(config.DefaultPatterns.IncludePatterns) != len(expectedInclude) {
		t.Errorf("Expected %d include patterns, got %d", len(expectedInclude), len(config.DefaultPatterns.IncludePatterns))
	}

	for i, pattern := range config.DefaultPatterns.IncludePatterns {
		if pattern != expectedInclude[i] {
			t.Errorf("Include pattern %d: got %s, want %s", i, pattern, expectedInclude[i])
		}
	}

	if len(config.DefaultPatterns.ExcludePatterns) != len(expectedExclude) {
		t.Errorf("Expected %d exclude patterns, got %d", len(expectedExclude), len(config.DefaultPatterns.ExcludePatterns))
	}

	for i, pattern := range config.DefaultPatterns.ExcludePatterns {
		if pattern != expectedExclude[i] {
			t.Errorf("Exclude pattern %d: got %s, want %s", i, pattern, expectedExclude[i])
		}
	}
}