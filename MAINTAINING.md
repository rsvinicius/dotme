# Maintaining dotme

This document provides instructions for maintainers of the `dotme` project, covering versioning, releases, and changelog updates.

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version when making incompatible API changes (e.g., 1.0.0)
- **MINOR** version when adding functionality in a backwards compatible manner (e.g., 0.1.0)
- **PATCH** version when making backwards compatible bug fixes (e.g., 0.0.1)

## Updating the Changelog

1. Open the `CHANGELOG.md` file
2. Add a new section for the new version at the top of the file
3. Use the following format:

```markdown
## [vX.Y.Z] - YYYY-MM-DD

### Added
- Feature A
- Feature B

### Changed
- Change A
- Change B

### Fixed
- Bug fix A
- Bug fix B

### Removed
- Removed feature A
```

4. Make sure all significant changes are documented

## Creating a New Release

### 1. Update Version

Update version references in the code:

- Main CLI version string in the `cmd/root.go` file

### 2. Update the Changelog

Follow the instructions above to update the `CHANGELOG.md` file.

### 3. Create a Git Tag

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag
git push origin v1.0.0
```

### 4. Build Release Binaries

Build binaries for all supported platforms:

```bash
# Create release directory
mkdir -p release

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o release/dotme-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o release/dotme-linux-arm64 .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o release/dotme-macos-amd64 .
GOOS=darwin GOARCH=arm64 go build -o release/dotme-macos-arm64 .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o release/dotme-windows-amd64.exe .
```

### 5. Create a GitHub Release

1. Go to the repository on GitHub
2. Click on "Releases"
3. Click "Draft a new release"
4. Select the tag you just created
5. Add a title and description (you can copy from the CHANGELOG)
6. Attach the binary files from the `release` directory
7. Click "Publish release"

## Handling Issues and Pull Requests

- Respond to issues and pull requests promptly
- Label issues appropriately
- Review pull requests with care, ensuring they follow coding standards
- Merge pull requests that meet the project's quality standards
- Update the changelog when merging significant changes