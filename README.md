# dotme

[![CI](https://github.com/rsvinicius/dotme/actions/workflows/ci.yml/badge.svg)](https://github.com/rsvinicius/dotme/actions/workflows/ci.yml)
[![Release](https://github.com/rsvinicius/dotme/actions/workflows/release.yml/badge.svg)](https://github.com/rsvinicius/dotme/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsvinicius/dotme)](https://goreportcard.com/report/github.com/rsvinicius/dotme)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool to apply dotfiles from a Git repository to your current working directory.

## 🚀 Overview

`dotme` is a simple utility that helps you quickly set up your development environment by applying dotfiles from a Git repository. It only copies files and folders that start with a dot (`.`) from the root of the specified repository into your current directory.

## ✨ Features

- Apply dotfiles from any Git repository with a single command
- Only copies files and folders that start with a dot (`.`) at the repository root
- Recursively copies contents of dotfiles folders
- Cross-platform (Linux, macOS, Windows)
- Save repositories with aliases for quick access
- Clear terminal output with information about what was copied and ignored
- Automatically cleans up temporary files after execution
- Comprehensive test suite with high code coverage
- Continuous integration and automated releases

## 📋 Installation

### From Binary Releases

Download the latest release archive for your platform from the [Releases](https://github.com/rsvinicius/dotme/releases) page and extract it.

#### Linux

```bash
# Download the tar.gz file (choose amd64 or arm64 according to your architecture)
curl -L https://github.com/rsvinicius/dotme/releases/download/v0.2.0/dotme_0.2.0_linux_amd64.tar.gz -o dotme.tar.gz

# Extract the archive
tar -xzf dotme.tar.gz

# Move the binary to a directory in your PATH (optional)
sudo mv dotme /usr/local/bin/
```

#### macOS

```bash
# Download the tar.gz file (choose amd64 or arm64 according to your architecture)
curl -L https://github.com/rsvinicius/dotme/releases/download/v0.2.0/dotme_0.2.0_darwin_amd64.tar.gz -o dotme.tar.gz

# Extract the archive
tar -xzf dotme.tar.gz

# Move the binary to a directory in your PATH (optional)
sudo mv dotme /usr/local/bin/
```

#### Windows

1. Download the [Windows zip file](https://github.com/rsvinicius/dotme/releases/download/v0.2.0/dotme_0.2.0_windows_amd64.zip).
2. Extract the zip file.
3. Rename the extracted file to `dotme.exe` (optional).
4. Add it to your PATH or place it in a directory that's already in your PATH.

### From Source

```bash
# Clone the repository
git clone https://github.com/rsvinicius/dotme.git

# Change to the project directory
cd dotme

# Build the binary
go build -o dotme

# Install the binary to your PATH (optional)
go install
```

## 🔧 Usage

```bash
# Apply dotfiles from a Git repository
dotme https://github.com/your-username/dotfiles

# Save a repository with an alias
dotme -s my-dotfiles https://github.com/your-username/dotfiles

# Apply dotfiles using a saved alias
dotme -a my-dotfiles

# List all saved repository aliases
dotme list-aliases

# Remove a saved repository alias
dotme remove-alias my-dotfiles

# Display version information
dotme version
```

### Example

```bash
# Apply dotfiles from a repository
dotme https://github.com/your-username/dotfiles
```

This will:
1. Clone the repository to a temporary directory
2. Copy only files and folders that start with a dot (`.`) from the root of the repository to your current directory
3. Show what was copied and what was ignored
4. Clean up the temporary directory

### Repository Aliases

You can save frequently used repositories with aliases for easier access:

```bash
# Save a repository with an alias
dotme -s work-dotfiles https://github.com/your-company/dotfiles

# Apply dotfiles using the saved alias
dotme -a work-dotfiles

# List all saved aliases
dotme list-aliases

# Remove an alias you no longer need
dotme remove-alias work-dotfiles
```

Aliases are stored in `~/.dotme/config.json` and can be used across sessions.

## ⚙️ How It Works

`dotme` performs the following steps:
1. Clones the specified Git repository to a temporary directory
2. Scans the root of the cloned repository for files and folders that start with a dot (`.`)
3. Copies those files/folders to your current working directory
   - For folders, it recursively copies all contents (regardless of whether the inner files start with a dot)
4. Displays a summary of what was copied and what was ignored
5. Cleans up the temporary directory

## 🧪 Development and Testing

The project follows a structured organization with clear separation of concerns:

```
dotme/
├── cmd/                    # Command-line interface
├── internal/               # Implementation code
│   ├── alias/              # Repository alias management
│   ├── fs/                 # File system operations
│   ├── git/                # Git repository operations
│   └── dotfiles.go         # Integration layer
└── test/                   # Test code
    ├── alias/              # Alias tests
    ├── fs/                 # File system tests
    └── mocks/              # Mock implementations
```

All tests can be run using:
```bash
go test ./...
```

## 🔄 CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **CI pipeline** runs on every push and pull request, performing tests, linting, and builds
- **Release workflow** automatically creates releases when version tags are pushed
- **Cross-platform builds** ensure compatibility across Linux, macOS, and Windows

## 🤝 Contributing

Contributions are welcome! Please check out our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📦 Versioning

We use [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/rsvinicius/dotme/tags).