package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
)

// ProcessRepository handles cloning the repository and copying dotfiles
func ProcessRepository(repoURL string) error {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "dotme-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("🔄 Cloning repository: %s\n", repoURL)

	// Clone the repository
	r, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Get repository information for detailed output
	ref, err := r.Head()
	if err == nil {
		fmt.Printf("✅ Repository cloned, using branch: %s\n", ref.Name().Short())
	}

	// Get current working directory
	destDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	fmt.Println("📋 Scanning for dotfiles...")

	// Process files from the temporary directory
	return copyDotFiles(tempDir, destDir)
}

// copyDotFiles copies all dotfiles from source to destination directory
func copyDotFiles(srcDir, destDir string) error {
	var copied, ignored int
	var copiedItems []string
	var ignoredItems []string

	// Read all files and directories from the source
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Process each entry
	for _, entry := range entries {
		name := entry.Name()

		// Skip .git directory
		if name == ".git" {
			continue
		}

		// Only process files/directories that start with a dot
		if !strings.HasPrefix(name, ".") {
			ignored++
			ignoredItems = append(ignoredItems, name)
			continue
		}

		srcPath := filepath.Join(srcDir, name)
		destPath := filepath.Join(destDir, name)

		if entry.IsDir() {
			// Process directories recursively
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
			copied++
			copiedItems = append(copiedItems, name+"/")
		} else {
			// Copy files directly
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
			copied++
			copiedItems = append(copiedItems, name)
		}
	}

	// Display summary
	fmt.Printf("\n📦 Summary:\n")
	fmt.Printf("✅ Copied %d items:\n", copied)
	for _, item := range copiedItems {
		fmt.Printf("   - %s\n", item)
	}

	fmt.Printf("\n❌ Ignored %d items:\n", ignored)
	for _, item := range ignoredItems {
		fmt.Printf("   - %s\n", item)
	}

	fmt.Printf("\n🎉 Done! Your dotfiles have been applied successfully.\n")

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a file from source to destination
func copyFile(src, dst string) error {
	// Check if destination file exists
	if _, err := os.Stat(dst); err == nil {
		fmt.Printf("⚠️  Warning: %s already exists, overwriting\n", dst)
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	// Get source file mode
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info %s: %w", src, err)
	}

	// Create destination file
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content from %s to %s: %w", src, dst, err)
	}

	fmt.Printf("📄 Copied: %s\n", dst)
	return nil
}
