# ts-copy Project Overview

## CRITICAL DEVELOPMENT DIRECTIVE

**MANDATORY DOCUMENTATION UPDATES AFTER ANY CHANGE**: When completing ANY feature, bug fix, or architectural change, you MUST IMMEDIATELY:

1. **Update README.md** - Ensure usage examples, installation instructions, and feature descriptions are current
2. **Update CLAUDE.md** - Ensure architecture descriptions, configuration details, and development notes reflect changes
3. **Verify consistency** - Both files must accurately reflect the current implementation

**This is non-negotiable. No work is complete until documentation is updated to match the implementation.**

## Purpose
ts-copy is a Go CLI application that recursively searches for files matching configured extensions in the current directory and copies them to a specified Tailscale machine using concurrent workers. It's designed to simplify bulk file transfers over Tailscale networks.

## Architecture & Key Components

### Core Application (`main.go`)
- **Language**: Go 1.23.8
- **Main functionality**: File discovery, concurrent copying using goroutines
- **Configuration**: CLI arguments only (no config files)
- **Concurrency**: Uses 5 worker goroutines for parallel file transfers
- **Commands executed**: `sudo tailscale cp <file> <targetMachine>:`

### CLI Usage
```bash
tscp <target-machine> --ext .mp3 --ext .flac --dry-run
tscp my-server -e .pdf -e .docx
```

### Key Features
- **Recursive file discovery**: Walks directory tree to find matching files
- **Configurable file extensions**: Supports any file types via CLI flags
- **Dry-run mode**: `--dry-run` flag to preview operations without execution
- **Concurrent transfers**: 5 parallel workers for efficient copying
- **Error handling**: Graceful error reporting for failed transfers

## Project Structure
```
ts-copy/
├── CLAUDE.md           # Project instructions and architecture docs
├── LICENSE             # MIT License
├── README.md           # User documentation
├── main.go             # Main application logic (legacy)
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── justfile            # Task automation (contains dry-run test)
├── cmd/
│   └── tscp/
│       └── main.go     # Main CLI entry point
├── internal/
│   ├── discovery/
│   │   └── files.go    # File discovery logic
│   ├── transfer/
│   │   └── tailscale.go # Tailscale transfer operations
│   └── worker/
│       └── pool.go     # Worker pool implementation
├── docs/
│   ├── prd/            # Product Requirements Documents
│   │   └── 001-cli-redesign-prd.md
│   └── adr/            # Architecture Decision Records
│       ├── 000-project-structure.md
│       └── 001-testing-strategy.md
└── test/               # Test directory (files excluded from git)
    ├── 01 - Leviticus - The Burial.flac
    └── Drs, Dogger - Let You Down (Original Mix).flac
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
- CLI arguments are mandatory - app will exit if target machine or extensions not provided
- File extensions are case-insensitive
- No current mechanism to resume interrupted transfers
- No validation of target machine availability before starting
- Test files matching configured extensions are gitignored - test files exist locally but aren't committed to the repository

## Development Workflow

### Feature Development Process
When building new features, we follow this structured approach:

1. **Product Requirements Document (PRD)**
   - Create PRD in `docs/prd/XXX-feature-name-prd.md` (where XXX is a 3-digit number for ordering)
   - Include: Problem Statement, Goals, Requirements, Success Metrics
   - Use this as the single source of truth for feature scope

2. **PRD Approval**
   - **CRITICAL**: PRD must be explicitly approved by Nick before implementation begins
   - Update PRD status from "Draft" to "Approved" once approved
   - No implementation work should start without this approval

3. **Architecture Decision Records (ADR)**
   - Create ADR in `docs/adr/XXX-decision-name.md` (where XXX is a 3-digit number for ordering)
   - Include: Context, Decision, Status, Consequences
   - Use for documenting significant technical and architectural decisions
   - Status options: Proposed, Accepted, Deprecated, Superseded

4. **Implementation Planning**
   - Break down requirements into actionable tasks
   - Use TodoWrite tool to track implementation progress
   - Plan technical approach and architecture changes
   - Reference relevant ADRs for architectural decisions

5. **Development & Testing**
   - Implement feature following existing code conventions
   - Update documentation to reflect changes
   - Ensure tests pass and code builds successfully

6. **Documentation Updates**
   - Update README.md with new functionality
   - Update CLAUDE.md if architectural changes are made
   - Update any relevant configuration examples

This workflow ensures clear requirements, proper planning, and maintainable code.

### Documentation Standards

**Formatting Requirements:**
- **Code blocks**: Ensure proper indentation and alignment in all documentation
- **Command examples**: Use consistent spacing and alignment for readability
- **Options lists**: Align descriptions consistently across all documentation
- **IMPORTANT**: Pay special attention to indentation - this has been a recurring issue

**Quality Check:**
- Review all documentation formatting before committing
- Ensure code blocks render correctly in markdown
- Verify command examples are properly formatted and aligned

## Critical Development Directive

**MANDATORY DOCUMENTATION UPDATES**: When completing ANY feature, bug fix, or architectural change, you MUST:

1. **Check README.md** - Update usage examples, installation instructions, feature descriptions
2. **Check CLAUDE.md** - Update architecture descriptions, configuration details, development notes
3. **Verify consistency** - Ensure both files accurately reflect the current implementation

**Failure to update documentation is considered incomplete work.** No feature is complete until documentation matches the implementation.

## License
MIT License (Copyright 2025 Nick Kirby)