# ts-copy

A Go CLI tool that efficiently copies multiple files to Tailscale machines using concurrent workers. Instead of running many sequential `tailscale file cp` commands, ts-copy transfers files in parallel for significantly better performance.

## Overview

ts-copy finds files matching your specified extensions and copies them to a Tailscale machine using up to 5 concurrent workers. This parallel approach is much faster than running individual `tailscale file cp` commands sequentially.

## Installation

1. **Download the binary:**

   - Go to the [Releases page](https://github.com/NRKirby/ts-copy/releases)
   - Download the appropriate binary for your platform (Linux, macOS, or Windows)
   - Extract the archive

2. **Install the binary:**

   ```bash
   # Move to a directory in your PATH
   sudo mv tscp /usr/local/bin/tscp

   # Or for macOS with Homebrew
   mv tscp /opt/homebrew/bin/tscp

   # Make executable (if needed)
   chmod +x /usr/local/bin/tscp
   ```

That's it! No configuration files needed - the tool is ready to use.

## For Developers

If you want to build from source:

```bash
git clone https://github.com/NRKirby/ts-copy.git
cd ts-copy
go build
```

## Usage

Run from any directory containing files you want to transfer:

```bash
# Copy audio files to my-server
tscp my-server --ext .mp3 --ext .flac --ext .wav

# Copy documents using short form flags
tscp my-server -e .pdf -e .docx -e .txt

# Dry run to see what would be copied
tscp my-server --ext .zip --dry-run

# Copy single file type
tscp my-server -e .mp4
```

## Command Line Arguments

```
Usage: tscp <target-machine> [options]

Arguments:
  <target-machine>    Name of the Tailscale machine to copy files to

Options:
  --ext <extension>   File extension to copy (repeatable)
  -e <extension>      File extension to copy (repeatable, short form)
  --dry-run          Show what would be copied without executing commands
  --help, -h         Show help message and usage examples
```

The program will:

- Search the current directory (recursively) for files matching the specified extensions
- Copy files to the specified Tailscale machine using up to 5 concurrent workers for optimal performance
- Execute: `sudo tailscale file cp <file> <targetMachine>:` (note: uses the standard tailscale command)

## Future Enhancements

Features not yet implemented but planned:

- **Tailscale status check**: Detect and warn if Tailscale is not running before attempting file transfers
- **Progress indicators**: Show transfer progress for large files
- **Retry logic**: Automatically retry failed transfers with exponential backoff
- **Selective copying**: Skip files that already exist on the target machine
- **Custom destination paths**: Specify target directory on remote machine instead of default home directory
