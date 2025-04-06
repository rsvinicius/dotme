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

### 1. Update the Changelog

Follow the instructions above to update the `CHANGELOG.md` file.

### 2. Create a Git Tag and Push

```bash
# Create an annotated tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# Push the tag
git push origin vX.Y.Z
```

### 3. Automated Release Process

Once you push a tag that starts with 'v', the GitHub Actions release workflow will automatically:

1. Build release binaries for all supported platforms
2. Create a GitHub Release with the same name as the tag
3. Upload the compiled binaries to the release
4. Generate checksums for verification

You can monitor the release process in the Actions tab of the GitHub repository.

### 4. Verify the Release

After the workflow completes:
1. Go to the repository's Releases page on GitHub
2. Verify that the new release appears with the correct version
3. Check that all binary assets are attached
4. Review the generated changelog

## CI/CD Setup

This project uses GitHub Actions for continuous integration and delivery:

### CI Workflow (.github/workflows/ci.yml)

The CI workflow runs on every push to main/master and pull requests. It includes:
- Running tests (with and without race detector)
- Linting the code with golangci-lint
- Building the application
- Cross-platform builds for multiple platforms

### Release Workflow (.github/workflows/release.yml)

The release workflow runs when a version tag (v*) is pushed. It uses GoReleaser to:
- Build binaries for multiple platforms
- Create release archives
- Generate checksums
- Publish the GitHub release

### GoReleaser Configuration (.goreleaser.yaml)

Contains the configuration for the automated release process, including:
- Build settings for different platforms
- Archive formats
- Changelog generation
- Release settings

## Handling Issues and Pull Requests

- Respond to issues and pull requests promptly
- Label issues appropriately
- Review pull requests with care, ensuring they follow coding standards and pass CI checks
- Merge pull requests that meet the project's quality standards
- Update the changelog when merging significant changes