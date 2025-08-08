package patterns

import (
	"testing"

	"github.com/rsvinicius/dotme/internal/patterns"
)

func TestIsDotfile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"dotfile", ".gitconfig", true},
		{"dotfile with extension", ".vimrc", true},
		{"dotfolder", ".vscode", true},
		{"regular file", "README.md", false},
		{"regular folder", "src", false},
		{"empty string", "", false},
		{"just dot", ".", true},
		{"double dot", "..", true},
		{"hidden file with path", ".config/nvim", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := patterns.IsDotfile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsDotfile(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestParsePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", nil},
		{"single pattern", ".gitconfig", []string{".gitconfig"}},
		{"multiple patterns", ".gitconfig,.vimrc,.bashrc", []string{".gitconfig", ".vimrc", ".bashrc"}},
		{"patterns with spaces", ".gitconfig, .vimrc , .bashrc", []string{".gitconfig", ".vimrc", ".bashrc"}},
		{"patterns with empty parts", ".gitconfig,,.vimrc,", []string{".gitconfig", ".vimrc"}},
		{"glob patterns", ".git*,.vim*", []string{".git*", ".vim*"}},
		{"mixed patterns", ".DS_Store,.git*,README.md", []string{".DS_Store", ".git*", "README.md"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := patterns.ParsePatterns(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParsePatterns(%q) returned %d patterns, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, pattern := range result {
				if pattern != tt.expected[i] {
					t.Errorf("ParsePatterns(%q)[%d] = %q, want %q", tt.input, i, pattern, tt.expected[i])
				}
			}
		})
	}
}

func TestFilterOptions_ShouldInclude(t *testing.T) {
	tests := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		filename        string
		expected        bool
	}{
		// Default behavior (no patterns)
		{"no patterns - dotfile", nil, nil, ".gitconfig", true},
		{"no patterns - regular file", nil, nil, "README.md", false},

		// Include patterns only
		{"include specific - match", []string{".gitconfig"}, nil, ".gitconfig", true},
		{"include specific - no match", []string{".gitconfig"}, nil, ".vimrc", false},
		{"include glob - match", []string{".git*"}, nil, ".gitconfig", true},
		{"include glob - no match", []string{".git*"}, nil, ".vimrc", false},
		{"include multiple - first match", []string{".gitconfig", ".vimrc"}, nil, ".gitconfig", true},
		{"include multiple - second match", []string{".gitconfig", ".vimrc"}, nil, ".vimrc", true},
		{"include multiple - no match", []string{".gitconfig", ".vimrc"}, nil, ".bashrc", false},

		// Exclude patterns only
		{"exclude specific - match", nil, []string{".DS_Store"}, ".DS_Store", false},
		{"exclude specific - no match", nil, []string{".DS_Store"}, ".gitconfig", true},
		{"exclude specific - regular file", nil, []string{".DS_Store"}, "README.md", false}, // Not a dotfile
		{"exclude glob - match", nil, []string{".DS_*"}, ".DS_Store", false},
		{"exclude glob - no match", nil, []string{".DS_*"}, ".gitconfig", true},
		{"exclude multiple - first match", nil, []string{".DS_Store", ".Trash"}, ".DS_Store", false},
		{"exclude multiple - second match", nil, []string{".DS_Store", ".Trash"}, ".Trash", false},
		{"exclude multiple - no match", nil, []string{".DS_Store", ".Trash"}, ".gitconfig", true},

		// Both include and exclude patterns
		{"both - include match, exclude no match", []string{".git*"}, []string{".DS_Store"}, ".gitconfig", true},
		{"both - include match, exclude match", []string{".git*"}, []string{".git*"}, ".gitconfig", false},
		{"both - include no match", []string{".vim*"}, []string{".DS_Store"}, ".gitconfig", false},
		{"both - complex case", []string{".git*", ".vim*"}, []string{".DS_Store"}, ".gitconfig", true},

		// Edge cases
		{"include empty pattern", []string{""}, nil, ".gitconfig", false},
		{"exclude empty pattern", nil, []string{""}, ".gitconfig", true},
		{"glob with brackets", []string{".git[ci]*"}, nil, ".gitconfig", true},
		{"glob with question mark - no match", []string{".vim?"}, nil, ".vimrc", false}, // .vimrc has 6 chars, ? matches 1
		{"glob with question mark - match", []string{".vim?"}, nil, ".vimr", true}, // .vimr has 5 chars, ? matches 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &patterns.FilterOptions{
				IncludePatterns: tt.includePatterns,
				ExcludePatterns: tt.excludePatterns,
			}
			result := filter.ShouldInclude(tt.filename)
			if result != tt.expected {
				t.Errorf("ShouldInclude(%q) with include=%v, exclude=%v = %v, want %v",
					tt.filename, tt.includePatterns, tt.excludePatterns, result, tt.expected)
			}
		})
	}
}