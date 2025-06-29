# ts-copy Project Overview

## Purpose
ts-copy is a Go CLI application that recursively searches for audio files in the current directory and copies them to a specified Tailscale machine using concurrent workers. It's designed to simplify bulk file transfers over Tailscale networks.

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
├── main.go           # Main application logic
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── justfile         # Task automation (contains dry-run test)
├── LICENSE          # MIT License
├── README.md        # User documentation
├── .gitignore       # Excludes audio files, binaries, and IDE files
└── test/           # Test directory (audio files excluded from git)
```

## Dependencies
- `gopkg.in/yaml.v3` - YAML configuration parsing
- Standard Go libraries: `flag`, `os`, `filepath`, `sync`, `os/exec`

## Development & Build
- **Build**: `go build` produces single binary
- **Test**: `just dry-run` runs dry-run mode in test directory
- **Release**: GitHub releases provide pre-built binaries for multiple platforms

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
- Audio files (*.mp3, *.flac, *.wav) are gitignored for testing purposes - test files exist locally but aren't committed to the repository

## License
MIT License (Copyright 2025 Nick Kirby)