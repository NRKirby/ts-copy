# ts-copy

A simple Go program that recursively finds files and copies them to a Tailscale machine using `tailscale cp`.

## Overview

ts-copy searches the current directory for files (configurable extensions) and copies them to a specified Tailscale machine using concurrent workers.

## Setup

1. Build the program:

   ```bash
   go build
   ```

2. Create config directory and file:

   ```bash
   mkdir -p ~/.ts-copy && touch ~/.ts-copy/config.yaml
   ```

3. Edit `~/.ts-copy/config.yaml`:
   ```yaml
   extensions:
     - ".mp3"
     - ".flac"
     - ".wav"
   targetTsMachine: your-tailscale-machine-name
   ```

## Usage

Run from any directory containing audio files:

```bash
# Dry run to see what would be copied
./ts-copy --dry-run

# Actually copy the files
./ts-copy
```

The program will:

- Recursively search the current directory for files matching configured extensions
- Copy each file to the specified Tailscale machine using up to 5 concurrent workers
- Execute: `sudo tailscale cp <file> <targetTsMachine>:`

## Options

- `--dry-run`: Show what files would be copied without executing commands

## Future Enhancements

Features not yet implemented but planned:

- **GoReleaser CI/CD**: Automated binary builds and releases from git tags
- **Tailscale status check**: Detect and warn if Tailscale is not running before attempting file transfers
- **Progress indicators**: Show transfer progress for large files
- **Retry logic**: Automatically retry failed transfers with exponential backoff
- **Selective copying**: Skip files that already exist on the target machine
- **Custom destination paths**: Specify target directory on remote machine instead of default home directory
