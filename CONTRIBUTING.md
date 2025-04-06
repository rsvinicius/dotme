# Contributing to dotme

Thank you for considering contributing to `dotme`! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with the following information:
- A clear title and description
- Steps to reproduce the bug
- Expected behavior
- Actual behavior
- Screenshots if applicable
- Your environment (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:
- A clear title and description
- As much detail as possible about the suggested enhancement
- The rationale for the enhancement
- Examples of how the enhancement would be used

### Pull Requests

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and ensure they pass
5. Commit your changes with clear commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Project Structure

The project is organized into the following directories:

```
dotme/
├── cmd/                    # Command-line application entrypoints
│   └── root.go
├── internal/               # Private application code
│   ├── fs/                 # Filesystem operations
│   │   └── copy.go
│   ├── git/                # Git operations
│   │   └── repository.go
│   └── dotfiles.go         # Integration layer
├── test/                   # Test code
│   ├── fs/                 # Filesystem tests
│   │   └── copy_test.go
│   └── mocks/              # Mock implementations
│       └── git_mock.go
├── .github/workflows/      # CI/CD workflows
│   ├── ci.yml
│   └── release.yml
└── ...
```

## Development Setup

1. Ensure you have Go 1.18 or later installed
2. Fork and clone the repository
3. Install dependencies with `go mod download`
4. Build the project with `go build`
5. Run tests with `go test ./...`

## Continuous Integration

This project uses GitHub Actions for continuous integration:

- **CI Workflow**: Runs tests, linting, and builds on every push and pull request
- **Release Workflow**: Automatically builds and publishes releases when version tags are pushed

When contributing, ensure your code passes all CI checks. You can view the workflow definitions in the `.github/workflows` directory.

## Coding Guidelines

- Follow standard Go coding conventions
- Format your code with `go fmt`
- Use `golint` and `go vet` to check for issues
- Write tests for new functionality
- Place implementation code in the appropriate package under `internal/`
- Place tests in the corresponding package under `test/`

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Routine tasks, maintenance, etc.

## Versioning

This project follows [Semantic Versioning](https://semver.org/).

## Questions?

Feel free to open an issue if you have questions about contributing.

Thank you for your contributions!