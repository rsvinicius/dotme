package alias

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Errors defined for alias operations
var (
	ErrAliasNotFound      = errors.New("alias not found")
	ErrAliasAlreadyExists = errors.New("alias already exists")
)

// Config represents the structure of the configuration file
type Config struct {
	Repositories    map[string]string `json:"repositories"`     // Maps alias to repository URL
	DefaultPatterns PatternConfig     `json:"default_patterns"` // Default include/exclude patterns
}

// PatternConfig holds the default pattern configuration
type PatternConfig struct {
	IncludePatterns []string `json:"include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create the .dotme directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".dotme")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// loadConfig loads the configuration from the file
func loadConfig() (Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return Config{}, err
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return an empty config if the file doesn't exist
		return Config{
			Repositories: make(map[string]string),
		}, nil
	}

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize the map if it's nil
	if config.Repositories == nil {
		config.Repositories = make(map[string]string)
	}

	return config, nil
}

// saveConfig saves the configuration to the file
func saveConfig(config Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write the config to the file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveRepo saves a repository URL with the given alias
func SaveRepo(repoURL, alias string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// Check if the alias already exists
	if _, exists := config.Repositories[alias]; exists {
		return ErrAliasAlreadyExists
	}

	// Save the repository
	config.Repositories[alias] = repoURL

	return saveConfig(config)
}

// GetRepo retrieves a repository URL by its alias
func GetRepo(alias string) (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}

	// Check if the alias exists
	repoURL, exists := config.Repositories[alias]
	if !exists {
		return "", ErrAliasNotFound
	}

	return repoURL, nil
}

// ListAliases returns a map of all saved aliases and their repository URLs
func ListAliases() (map[string]string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Return a copy of the repositories map
	result := make(map[string]string, len(config.Repositories))
	for alias, repoURL := range config.Repositories {
		result[alias] = repoURL
	}

	return result, nil
}

// DeleteAlias removes an alias from the configuration
func DeleteAlias(alias string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// Check if the alias exists
	if _, exists := config.Repositories[alias]; !exists {
		return ErrAliasNotFound
	}

	// Delete the alias
	delete(config.Repositories, alias)

	return saveConfig(config)
}

// GetDefaultPatterns returns the default pattern configuration
func GetDefaultPatterns() (PatternConfig, error) {
	config, err := loadConfig()
	if err != nil {
		return PatternConfig{}, err
	}

	return config.DefaultPatterns, nil
}

// SetDefaultPatterns saves the default pattern configuration
func SetDefaultPatterns(patterns PatternConfig) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.DefaultPatterns = patterns
	return saveConfig(config)
}