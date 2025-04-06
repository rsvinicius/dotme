package fs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	fspackage "github.com/rsvinicius/dotme/internal/fs"
)

// TestIsDotfile tests the logic for identifying dotfiles
func TestIsDotfile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"empty string", "", false},
		{"dotfile", ".gitconfig", true},
		{"dotfolder", ".vscode", true},
		{"non-dotfile", "README.md", false},
		{"non-dotfolder", "src", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fspackage.IsDotfile(tt.filename); got != tt.want {
				t.Errorf("IsDotfile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestCopyFile tests the file copying functionality
func TestCopyFile(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "dotme-test-src-*")
	if err != nil {
		t.Fatalf("failed to create source temp dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "dotme-test-dst-*")
	if err != nil {
		t.Fatalf("failed to create destination temp dir: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create a test file in the source directory
	testContent := "test content"
	testFilename := ".testfile"
	srcPath := filepath.Join(srcDir, testFilename)
	if err := os.WriteFile(srcPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Copy the file
	dstPath := filepath.Join(dstDir, testFilename)
	if err := fspackage.CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify the file was copied correctly
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(dstContent) != testContent {
		t.Errorf("CopyFile did not copy content correctly, got %q, want %q", string(dstContent), testContent)
	}

	// Check file permissions
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		t.Fatalf("failed to stat source file: %v", err)
	}

	dstInfo, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("failed to stat destination file: %v", err)
	}

	if srcInfo.Mode() != dstInfo.Mode() {
		t.Errorf("CopyFile did not preserve file mode, got %v, want %v", dstInfo.Mode(), srcInfo.Mode())
	}
}

// TestCopyDir tests the directory copying functionality
func TestCopyDir(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "dotme-test-src-*")
	if err != nil {
		t.Fatalf("failed to create source temp dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "dotme-test-dst-*")
	if err != nil {
		t.Fatalf("failed to create destination temp dir: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create a directory structure in the source
	testDirPath := filepath.Join(srcDir, ".testdir")
	if err := os.Mkdir(testDirPath, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create files in the test directory
	testFiles := []struct {
		path    string
		content string
	}{
		{filepath.Join(testDirPath, "file1"), "content1"},
		{filepath.Join(testDirPath, "file2"), "content2"},
		{filepath.Join(testDirPath, ".dotfile"), "dotcontent"},
	}

	for _, tf := range testFiles {
		if err := os.WriteFile(tf.path, []byte(tf.content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", tf.path, err)
		}
	}

	// Create a subdirectory
	subDirPath := filepath.Join(testDirPath, "subdir")
	if err := os.Mkdir(subDirPath, 0755); err != nil {
		t.Fatalf("failed to create test subdirectory: %v", err)
	}

	// Create a file in the subdirectory
	subFilePath := filepath.Join(subDirPath, "subfile")
	if err := os.WriteFile(subFilePath, []byte("subcontent"), 0644); err != nil {
		t.Fatalf("failed to create test file in subdirectory: %v", err)
	}

	// Copy the directory
	dstDirPath := filepath.Join(dstDir, ".testdir")
	if err := fspackage.CopyDir(testDirPath, dstDirPath); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify the directory was copied correctly
	var srcFiles []string
	var dstFiles []string

	err = filepath.WalkDir(testDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(testDirPath, path)
			if err != nil {
				return err
			}
			srcFiles = append(srcFiles, relPath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk source directory: %v", err)
	}

	err = filepath.WalkDir(dstDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(dstDirPath, path)
			if err != nil {
				return err
			}
			dstFiles = append(dstFiles, relPath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk destination directory: %v", err)
	}

	// Sort for comparison
	sort.Strings(srcFiles)
	sort.Strings(dstFiles)

	if !reflect.DeepEqual(srcFiles, dstFiles) {
		t.Errorf("CopyDir did not copy files correctly, got %v, want %v", dstFiles, srcFiles)
	}

	// Verify content of copied files
	for _, tf := range testFiles {
		srcRelPath, err := filepath.Rel(testDirPath, tf.path)
		if err != nil {
			t.Fatalf("failed to get relative path: %v", err)
		}

		dstPath := filepath.Join(dstDirPath, srcRelPath)
		dstContent, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("failed to read destination file %s: %v", dstPath, err)
		}

		if string(dstContent) != tf.content {
			t.Errorf("CopyDir did not copy content correctly for %s, got %q, want %q", srcRelPath, string(dstContent), tf.content)
		}
	}

	// Check subdirectory content
	dstSubFilePath := filepath.Join(dstDirPath, "subdir", "subfile")
	dstSubContent, err := os.ReadFile(dstSubFilePath)
	if err != nil {
		t.Fatalf("failed to read destination subdirectory file: %v", err)
	}

	if string(dstSubContent) != "subcontent" {
		t.Errorf("CopyDir did not copy subdirectory content correctly, got %q, want %q", string(dstSubContent), "subcontent")
	}
}

// TestCopyDotFiles tests the dotfile filtering and copying logic
func TestCopyDotFiles(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "dotme-test-src-*")
	if err != nil {
		t.Fatalf("failed to create source temp dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "dotme-test-dst-*")
	if err != nil {
		t.Fatalf("failed to create destination temp dir: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create files and directories in the source
	files := []struct {
		name    string
		content string
		isDir   bool
	}{
		{".gitconfig", "git config content", false},
		{".vscode", "", true},
		{"README.md", "readme content", false},
		{"src", "", true},
		{".npmrc", "npm config content", false},
	}

	for _, f := range files {
		path := filepath.Join(srcDir, f.name)
		if f.isDir {
			if err := os.Mkdir(path, 0755); err != nil {
				t.Fatalf("failed to create test directory %s: %v", path, err)
			}
		} else {
			if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
				t.Fatalf("failed to create test file %s: %v", path, err)
			}
		}
	}

	// Create a file in the .vscode directory
	vscodePath := filepath.Join(srcDir, ".vscode", "settings.json")
	if err := os.WriteFile(vscodePath, []byte("vscode settings"), 0644); err != nil {
		t.Fatalf("failed to create test file in .vscode directory: %v", err)
	}

	// Copy dotfiles
	// Redirect stdout to avoid cluttering test output
	originalStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = originalStdout }()

	if err := fspackage.CopyDotFiles(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDotFiles failed: %v", err)
	}

	// Verify only dotfiles were copied
	entries, err := os.ReadDir(dstDir)
	if err != nil {
		t.Fatalf("failed to read destination directory: %v", err)
	}

	var copiedNames []string
	for _, entry := range entries {
		copiedNames = append(copiedNames, entry.Name())
	}

	expectedNames := []string{".gitconfig", ".npmrc", ".vscode"}
	sort.Strings(copiedNames)
	sort.Strings(expectedNames)

	if !reflect.DeepEqual(copiedNames, expectedNames) {
		t.Errorf("CopyDotFiles did not filter correctly, got %v, want %v", copiedNames, expectedNames)
	}

	// Check content of the copied files
	for _, f := range files {
		if !fspackage.IsDotfile(f.name) || f.isDir {
			continue
		}

		dstPath := filepath.Join(dstDir, f.name)
		content, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("failed to read copied file %s: %v", dstPath, err)
		}

		if string(content) != f.content {
			t.Errorf("CopyDotFiles did not copy content correctly for %s, got %q, want %q", f.name, string(content), f.content)
		}
	}

	// Check .vscode directory contents
	dstVscodePath := filepath.Join(dstDir, ".vscode", "settings.json")
	content, err := os.ReadFile(dstVscodePath)
	if err != nil {
		t.Fatalf("failed to read copied file in .vscode directory: %v", err)
	}

	if string(content) != "vscode settings" {
		t.Errorf("CopyDotFiles did not copy .vscode contents correctly, got %q, want %q", string(content), "vscode settings")
	}
}
