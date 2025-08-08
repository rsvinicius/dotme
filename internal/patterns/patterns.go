package patterns

import (
	"path/filepath"
	"strings"
)

// FilterOptions contains the filtering options for dotfiles
type FilterOptions struct {
	IncludePatterns []string
	ExcludePatterns []string
}

// ShouldInclude determines if a file should be included based on the filter options
func (f *FilterOptions) ShouldInclude(filename string) bool {
	// If no patterns are specified, use default behavior (include all dotfiles)
	if len(f.IncludePatterns) == 0 && len(f.ExcludePatterns) == 0 {
		return IsDotfile(filename)
	}

	// If include patterns are specified, file must match at least one
	if len(f.IncludePatterns) > 0 {
		included := false
		for _, pattern := range f.IncludePatterns {
			if matchesPattern(filename, pattern) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	} else {
		// No include patterns, so include all dotfiles by default
		if !IsDotfile(filename) {
			return false
		}
	}

	// If exclude patterns are specified, file must not match any
	if len(f.ExcludePatterns) > 0 {
		for _, pattern := range f.ExcludePatterns {
			if matchesPattern(filename, pattern) {
				return false
			}
		}
	}

	return true
}

// IsDotfile checks if a file or directory name starts with a dot
func IsDotfile(name string) bool {
	return strings.HasPrefix(name, ".")
}

// matchesPattern checks if a filename matches a glob pattern
func matchesPattern(filename, pattern string) bool {
	// Use filepath.Match for glob pattern matching
	matched, err := filepath.Match(pattern, filename)
	if err != nil {
		// If pattern is invalid, fall back to exact string matching
		return filename == pattern
	}
	return matched
}

// ParsePatterns parses a comma-separated string of patterns into a slice
func ParsePatterns(patterns string) []string {
	if patterns == "" {
		return nil
	}
	
	var result []string
	parts := strings.Split(patterns, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}