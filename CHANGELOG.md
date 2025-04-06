# Changelog

All notable changes to the `dotme` project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Unit tests for core functionality:
  - Tests for dotfile filtering logic
  - Tests for file copying functionality
  - Tests for directory copying functionality
  - Mock tests for repository processing

### Changed
- Improved project structure and organization:
  - Separated implementation code by responsibility (fs, git)
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