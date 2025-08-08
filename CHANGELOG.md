# Changelog

All notable changes to the `dotme` project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.0] - 2025-01-27

### Added
- Pattern-based filtering support:
  - `--include` flag to specify patterns for files to include
  - `--exclude` flag to specify patterns for files to exclude
  - Support for glob patterns (*, ?, [abc], etc.)
  - Comma-separated pattern lists for multiple patterns
- Default pattern configuration:
  - `dotme config set-default-patterns` command to set default include/exclude patterns
  - `dotme config show` command to display current configuration
  - Persistent storage of default patterns in configuration file
  - Automatic use of default patterns when no command-line patterns are specified
- Enhanced filtering logic:
  - Flexible pattern matching with glob support
  - Combination of include and exclude patterns
  - Clear display of active filters in output summary
- Comprehensive test coverage:
  - Unit tests for pattern matching functionality
  - Tests for filtering logic with various pattern combinations
  - Tests for configuration management
  - Integration tests for file copying with patterns

### Changed
- Updated `CopyDotFiles` function to accept filter options parameter
- Enhanced command-line interface with new pattern flags
- Improved output to show active filters when patterns are used
- Extended configuration structure to support default patterns

### Fixed
- Pattern matching edge cases with invalid glob patterns
- Proper handling of empty pattern strings

## [v0.2.0] - 2025-04-15
### Added
- Repository alias functionality:
  - Save repositories with custom aliases for quick access (`-s, --save` flag)
  - Use saved repositories by their aliases (`-a, --alias` flag)
  - List all saved aliases (`list-aliases` command)
  - Remove saved aliases (`remove-alias` command)
  - Persistent storage of aliases in user's home directory
- Unit tests for alias functionality
- Unit tests for core functionality:
  - Tests for dotfile filtering logic
  - Tests for file copying functionality
  - Tests for directory copying functionality
  - Mock tests for repository processing
- GitHub Actions workflows for CI/CD:
  - Continuous integration with testing, linting, and building
  - Cross-platform binary building
  - Automated release process with GoReleaser
- Version command to display build information

### Changed
- Improved project structure and organization:
  - Separated implementation code by responsibility (fs, git, alias)
  - Moved all tests to dedicated test directory
  - Clear separation between tests and implementation code
  - Better code modularity with focused packages

## [v0.1.0] - 2025-04-04

### Added
- Initial release of `dotme` utility
- Command-line interface using Cobra
- Support for cloning Git repositories
- Filtering files and directories to only copy dotfiles
- Cross-platform support (Linux, macOS, Windows)
- Detailed terminal output showing what was copied and ignored
- Automatic cleanup of temporary directories
- MIT License
- Documentation (README.md, CONTRIBUTING.md)