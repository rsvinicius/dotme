package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rsvinicius/dotme/internal/patterns"
)

// CopyDotFiles copies dotfiles from source to destination directory based on filter options
func CopyDotFiles(srcDir, destDir string, filterOptions *patterns.FilterOptions) error {
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

		// Check if the file should be included based on filter options
		if !filterOptions.ShouldInclude(name) {
			ignored++
			ignoredItems = append(ignoredItems, name)
			continue
		}

		srcPath := filepath.Join(srcDir, name)
		destPath := filepath.Join(destDir, name)

		if entry.IsDir() {
			// Process directories recursively
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
			copied++
			copiedItems = append(copiedItems, name+"/")
		} else {
			// Copy files directly
			if err := CopyFile(srcPath, destPath); err != nil {
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

	// Display active filters if any
	if len(filterOptions.IncludePatterns) > 0 || len(filterOptions.ExcludePatterns) > 0 {
		fmt.Printf("\n🔍 Active filters:\n")
		if len(filterOptions.IncludePatterns) > 0 {
			fmt.Printf("   Include patterns: %v\n", filterOptions.IncludePatterns)
		}
		if len(filterOptions.ExcludePatterns) > 0 {
			fmt.Printf("   Exclude patterns: %v\n", filterOptions.ExcludePatterns)
		}
	}

	fmt.Printf("\n🎉 Done! Your dotfiles have been applied successfully.\n")

	return nil
}

// CopyDir recursively copies a directory
func CopyDir(src, dst string) error {
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
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
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