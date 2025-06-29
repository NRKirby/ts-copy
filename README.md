# ts-copy

A simple Go program that recursively finds files matching configured extensions and copies them to a Tailscale machine using `tailscale cp`.

## Overview

ts-copy searches the current directory for files matching configurable extensions and copies them to a specified Tailscale machine using concurrent workers.

## Installation

1. **Download the binary:**

   - Go to the [Releases page](https://github.com/NRKirby/ts-copy/releases)
   - Download the appropriate binary for your platform (Linux, macOS, or Windows)
   - Extract the archive

2. **Install the binary:**

   ```bash
   # Move to a directory in your PATH
   sudo mv ts-copy /usr/local/bin/ts-copy

   # Or for macOS with Homebrew
   mv ts-copy /opt/homebrew/bin/ts-copy

   # Make executable (if needed)
   chmod +x /usr/local/bin/ts-copy
   ```

3. **Create config directory and file:**

   ```bash
   mkdir -p ~/.ts-copy && touch ~/.ts-copy/config.yaml
   ```

4. **Edit `~/.ts-copy/config.yaml`:**
   ```yaml
   extensions:
     - ".mp3"
     - ".flac"
     - ".wav"
     - ".pdf"
     - ".txt"
   targetTsMachine: your-tailscale-machine-name
   ```

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
# Dry run to see what would be copied
ts-copy --dry-run

# Actually copy the files
ts-copy
```

The program will:

- Recursively search the current directory for files matching the configured extensions
- Copy each file to the specified Tailscale machine using up to 5 concurrent workers
- Execute: `sudo tailscale cp <file> <targetTsMachine>:`

## Options

- `--dry-run`: Show what files would be copied without executing commands

## Future Enhancements

Features not yet implemented but planned:

- **Tailscale status check**: Detect and warn if Tailscale is not running before attempting file transfers
- **Progress indicators**: Show transfer progress for large files
- **Retry logic**: Automatically retry failed transfers with exponential backoff
- **Selective copying**: Skip files that already exist on the target machine
- **Custom destination paths**: Specify target directory on remote machine instead of default home directory
