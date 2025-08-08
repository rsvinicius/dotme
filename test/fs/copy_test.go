package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rsvinicius/dotme/internal/fs"
	"github.com/rsvinicius/dotme/internal/patterns"
)

// TestCopyFile tests the file copying functionality
func TestCopyFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotme-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a source file
	srcFile := filepath.Join(tempDir, "source.txt")
	content := "Hello, World!"
	err = os.WriteFile(srcFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy the file
	destFile := filepath.Join(tempDir, "dest.txt")
	err = fs.CopyFile(srcFile, destFile)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify the destination file exists and has the correct content
	destContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(destContent) != content {
		t.Errorf("Destination file content = %q, want %q", string(destContent), content)
	}
}

// TestCopyDir tests the directory copying functionality
func TestCopyDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotme-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create files in the source directory
	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
	}

	for file, content := range files {
		filePath := filepath.Join(srcDir, file)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create source file %s: %v", file, err)
		}
	}

	// Copy the directory
	destDir := filepath.Join(tempDir, "dest")
	err = fs.CopyDir(srcDir, destDir)
	if err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify all files were copied correctly
	for file, expectedContent := range files {
		destFile := filepath.Join(destDir, file)
		content, err := os.ReadFile(destFile)
		if err != nil {
			t.Errorf("Failed to read destination file %s: %v", file, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("File %s content = %q, want %q", file, string(content), expectedContent)
		}
	}
}

// TestCopyDotFiles tests the dotfile filtering and copying logic with patterns
func TestCopyDotFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotme-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory
	srcDir := filepath.Join(tempDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create destination directory
	destDir := filepath.Join(tempDir, "dest")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	// Create test files and directories
	testFiles := map[string]string{
		".gitconfig":    "git config content",
		".vimrc":        "vim config content",
		".DS_Store":     "mac metadata",
		"README.md":     "readme content",
		"regular.txt":   "regular file content",
	}

	for file, content := range testFiles {
		filePath := filepath.Join(srcDir, file)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create a dotfile directory
	dotDir := filepath.Join(srcDir, ".config")
	err = os.MkdirAll(dotDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .config directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(dotDir, "config.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	tests := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		expectedFiles   []string
		unexpectedFiles []string
	}{
		{
			name:            "default behavior - all dotfiles",
			includePatterns: nil,
			excludePatterns: nil,
			expectedFiles:   []string{".gitconfig", ".vimrc", ".DS_Store", ".config/config.json"},
			unexpectedFiles: []string{"README.md", "regular.txt"},
		},
		{
			name:            "include specific files",
			includePatterns: []string{".gitconfig", ".vimrc"},
			excludePatterns: nil,
			expectedFiles:   []string{".gitconfig", ".vimrc"},
			unexpectedFiles: []string{".DS_Store", ".config/config.json", "README.md", "regular.txt"},
		},
		{
			name:            "exclude specific files",
			includePatterns: nil,
			excludePatterns: []string{".DS_Store"},
			expectedFiles:   []string{".gitconfig", ".vimrc", ".config/config.json"},
			unexpectedFiles: []string{".DS_Store", "README.md", "regular.txt"},
		},
		{
			name:            "include with glob pattern",
			includePatterns: []string{".git*"},
			excludePatterns: nil,
			expectedFiles:   []string{".gitconfig"},
			unexpectedFiles: []string{".vimrc", ".DS_Store", ".config/config.json", "README.md", "regular.txt"},
		},
		{
			name:            "exclude with glob pattern",
			includePatterns: nil,
			excludePatterns: []string{".DS_*"},
			expectedFiles:   []string{".gitconfig", ".vimrc", ".config/config.json"},
			unexpectedFiles: []string{".DS_Store", "README.md", "regular.txt"},
		},
		{
			name:            "both include and exclude",
			includePatterns: []string{".git*", ".vim*"},
			excludePatterns: []string{".DS_Store"},
			expectedFiles:   []string{".gitconfig", ".vimrc"},
			unexpectedFiles: []string{".DS_Store", ".config/config.json", "README.md", "regular.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean destination directory
			os.RemoveAll(destDir)
			err = os.MkdirAll(destDir, 0755)
			if err != nil {
				t.Fatalf("Failed to recreate destination directory: %v", err)
			}

			// Create filter options
			filterOptions := &patterns.FilterOptions{
				IncludePatterns: tt.includePatterns,
				ExcludePatterns: tt.excludePatterns,
			}

			// Copy dotfiles
			err = fs.CopyDotFiles(srcDir, destDir, filterOptions)
			if err != nil {
				t.Fatalf("CopyDotFiles failed: %v", err)
			}

			// Check expected files exist
			for _, expectedFile := range tt.expectedFiles {
				destFile := filepath.Join(destDir, expectedFile)
				if _, err := os.Stat(destFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not copied", expectedFile)
				}
			}

			// Check unexpected files don't exist
			for _, unexpectedFile := range tt.unexpectedFiles {
				destFile := filepath.Join(destDir, unexpectedFile)
				if _, err := os.Stat(destFile); err == nil {
					t.Errorf("Unexpected file %s was copied", unexpectedFile)
				}
			}
		})
	}
}