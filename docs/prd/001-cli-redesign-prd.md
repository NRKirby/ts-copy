# CLI Interface Redesign - Product Requirements Document

**Date:** 2025-06-29  
**Status:** Approved  
**Author:** Nick Kirby

## Problem Statement

The current ts-copy tool requires users to create and maintain a configuration file at `~/.ts-copy/config.yaml` before they can use the tool. This creates several pain points:

- **Setup friction**: Users must create a config file before first use
- **Portability issues**: Config file ties the tool to a specific machine/user setup
- **Inflexibility**: Changing target machines or extensions requires editing config files
- **Discoverability**: New users don't immediately understand what the tool does without reading documentation

## Goals

1. **Eliminate config file dependency** - Make the tool completely self-contained
2. **Improve portability** - Tool should work immediately after installation without setup
3. **Enhance flexibility** - Allow easy switching between target machines and file types per invocation
4. **Maintain simplicity** - Keep the tool easy to use for common cases
5. **Preserve existing functionality** - Don't break core file transfer capabilities

## Requirements

### Core CLI Interface

**Target Machine Parameter**

- Target Tailscale machine should be a required command-line parameter
- Should be the first positional argument for clarity: `tscp <target-machine>`
- **Validation**: Show clear error message if no target machine is specified

**File Extensions Parameter**

- Accept file extensions as repeatable command-line flags
- Use `--ext` or `-e` flag that can be specified multiple times
- **Minimum requirement**: At least one `--ext` flag must be provided
- **Validation**: Show clear error message if no extensions are specified
- Example: `tscp my-server --ext .mp3 --ext .flac --ext .pdf`

**Existing Functionality**

- Maintain `--dry-run` flag for preview mode
- Keep concurrent workers (5 workers) for performance
- Preserve recursive directory scanning
- Maintain error handling and reporting

### CLI Interface Specification

```bash
tscp <target-machine> --ext <extension> [--ext <extension>...] [options]
```

**Arguments:**

- `<target-machine>`: Required. Name of the Tailscale machine to copy files to

**Options:**

- `--ext <extension>`, `-e <extension>`: File extension to copy (repeatable, required)
- `--dry-run`: Show what would be copied without executing commands
- `--help`, `-h`: Show help message and usage examples

### Usage Examples

```bash
# Copy audio files to my-server (using --ext)
tscp my-server --ext .mp3 --ext .flac --ext .wav

# Copy audio files using short form (-e)
tscp my-server -e .mp3 -e .flac -e .wav

# Copy documents with dry-run
tscp my-server --ext .pdf --ext .docx --dry-run

# Mixed usage of short and long form
tscp my-server -e .zip --ext .tar.gz --dry-run

# Copy single file type
tscp my-server -e .zip

# Error case - no extensions specified
tscp my-server
# Output: Error: At least one file extension must be specified using --ext or -e flag

# Error case - no target machine specified
tscp
# Output: Error: Target machine is required as the first argument

# Error case - no target machine or extensions
tscp
# Output: Error: Target machine is required as the first argument
```

## Success Metrics

1. **Zero-config operation**: Tool works immediately after installation without any setup
2. **Improved UX**: Users can understand and use the tool from `--help` output alone
3. **Flexibility**: Users can easily change target machines and file types per invocation
4. **Backward compatibility**: Existing core functionality remains unchanged

## Implementation Notes

- Remove `loadConfig()` function and config file handling
- Refactor argument parsing to handle positional target machine argument
- Implement repeatable `--ext` flag parsing
- **Add validation**: Ensure target machine is provided as first argument, show helpful error if missing
- **Add validation**: Ensure at least one `--ext` flag is provided, show helpful error if missing
- Update help text to show clear usage examples
- Remove references to config file from documentation
- Maintain existing file discovery and transfer logic

## Out of Scope

- Preset extension groups (--audio, --video, etc.) - keep initial implementation simple
- Configuration file migration - clean break from config approach
- Advanced filtering options - focus on core use case

## Dependencies

- Go standard library `flag` package for argument parsing
- Existing Tailscale integration remains unchanged
- No new external dependencies required
