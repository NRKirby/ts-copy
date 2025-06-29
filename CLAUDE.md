# ts-copy Project Overview

## Purpose
ts-copy is a Go CLI application that recursively searches for files matching configured extensions in the current directory and copies them to a specified Tailscale machine using concurrent workers. It's designed to simplify bulk file transfers over Tailscale networks.

## Architecture & Key Components

### Core Application (`main.go`)
- **Language**: Go 1.23.8
- **Main functionality**: File discovery, concurrent copying using goroutines
- **Configuration**: YAML-based config at `~/.ts-copy/config.yaml`
- **Concurrency**: Uses 5 worker goroutines for parallel file transfers
- **Commands executed**: `sudo tailscale cp <file> <targetMachine>:`

### Configuration Structure
```yaml
extensions:
  - ".mp3"
  - ".flac" 
  - ".wav"
  - ".pdf"
  - ".txt"
targetTsMachine: "your-tailscale-machine-name"
```

### Key Features
- **Recursive file discovery**: Walks directory tree to find matching files
- **Configurable file extensions**: Supports any file types via config
- **Dry-run mode**: `--dry-run` flag to preview operations without execution
- **Concurrent transfers**: 5 parallel workers for efficient copying
- **Error handling**: Graceful error reporting for failed transfers

## Project Structure
```
ts-copy/
├── main.go             # Main application logic
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── justfile            # Task automation (contains dry-run test)
├── LICENSE             # MIT License
├── README.md           # User documentation
├── .gitignore          # Excludes test files, binaries, and IDE files
├── .github/
│   └── workflows/
│       └── release.yml # GitHub Actions release workflow
├── .goreleaser.yml     # GoReleaser build configuration
└── test/               # Test directory (files excluded from git)
```

## Dependencies
- `gopkg.in/yaml.v3` - YAML configuration parsing
- Standard Go libraries: `flag`, `os`, `filepath`, `sync`, `os/exec`

## Development & Build

### Local Development
- **Build**: `go build` produces single binary
- **Test**: `just dry-run` runs dry-run mode in test directory

### CI/CD Release Process
**GitHub Actions Workflow** (`.github/workflows/release.yml`):
- **Trigger**: Automatic on version tags matching `v*` pattern (e.g., `v1.0.0`)
- **Runner**: Ubuntu latest
- **Go Version**: 1.21 with module caching
- **Steps**:
  1. Checkout code
  2. Setup Go environment
  3. Run tests: `go test ./... -v`
  4. Use GoReleaser for build and release

**GoReleaser Configuration** (`.goreleaser.yml`):
- **Binary Name**: `tscp` (changed from `ts-copy`)
- **Build ID**: `ts-copy`
- **Platforms**: Linux, macOS, Windows
- **Architectures**: ARM64, AMD64 (excludes Windows ARM64)
- **Archive Formats**: tar.gz (Unix), zip (Windows)
- **Archive Naming**: `{ProjectName}_{OS}_{Arch}` format
- **Changelog**: Excludes commits prefixed with `docs:`, `test:`, `chore:`
- **Release Settings**: Auto-prerelease detection, draft disabled, make latest enabled

**Release Process**:
1. Push version tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions triggers automatically
3. GoReleaser builds cross-platform binaries
4. Creates GitHub release with binaries and changelog
5. Binaries available as `tscp` executable

## Current Limitations & Future Enhancements
As documented in README.md, planned improvements include:
- Tailscale status checking
- Progress indicators for large files
- Retry logic with exponential backoff
- Selective copying (skip existing files)
- Custom destination paths

## Important Notes for Development
- Uses `sudo tailscale cp` - requires elevated privileges
- Config file is mandatory - app will exit if not found
- File extensions are case-insensitive
- No current mechanism to resume interrupted transfers
- No validation of target machine availability before starting
- Test files matching configured extensions are gitignored - test files exist locally but aren't committed to the repository

## License
MIT License (Copyright 2025 Nick Kirby)