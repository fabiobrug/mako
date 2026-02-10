# Mako Installation Guide

This guide covers all methods to install Mako on your system.

## Quick Install (Recommended)

### One-Command Install

The fastest way to install Mako:

```bash
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/main/scripts/install.sh | bash
```

This will:
- Detect your OS and architecture
- Download the latest release
- Install to `/usr/local/bin` (or `~/.local/bin` if no sudo access)
- Set up the `~/.mako` directory
- Prompt for your Gemini API key
- Verify the installation

### Homebrew (macOS/Linux)

If you prefer using Homebrew:

```bash
brew tap fabiobrug/mako
brew install mako
```

After installation, set your API key:

```bash
mako config set api_key YOUR_API_KEY
```

Get your free API key at: https://ai.google.dev/

## Manual Installation

### From Pre-built Binaries

1. Download the appropriate binary for your platform from the [releases page](https://github.com/fabiobrug/mako/releases/latest):
   - `mako-linux-amd64` - Linux x86_64
   - `mako-linux-arm64` - Linux ARM64
   - `mako-darwin-amd64` - macOS Intel
   - `mako-darwin-arm64` - macOS Apple Silicon

2. Download the menu binary as well:
   - `mako-menu-linux-amd64`
   - `mako-menu-darwin-arm64`
   - etc.

3. Make them executable and move to your PATH:

```bash
chmod +x mako-* mako-menu-*
sudo mv mako-* /usr/local/bin/mako
sudo mv mako-menu-* /usr/local/bin/mako-menu
```

4. Create the Mako directory:

```bash
mkdir -p ~/.mako
```

5. Set your API key:

```bash
mako config set api_key YOUR_API_KEY
```

### From Source

Requirements:
- Go 1.21 or later
- SQLite3 with FTS5 support
- Git

Steps:

```bash
# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako

# Build the binaries
go build -tags "fts5" -o mako cmd/mako/main.go
go build -o mako-menu cmd/mako-menu/main.go

# Install to your PATH
sudo mv mako mako-menu /usr/local/bin/

# Create Mako directory
mkdir -p ~/.mako

# Set your API key
mako config set api_key YOUR_API_KEY
```

## Post-Installation Setup

### API Key

Mako requires a Google Gemini API key to work. Get yours for free:

1. Visit https://ai.google.dev/
2. Sign in with your Google account
3. Create a new API key
4. Set it in Mako:

```bash
mako config set api_key YOUR_API_KEY
```

Alternatively, you can set the `GEMINI_API_KEY` environment variable:

```bash
export GEMINI_API_KEY=your_api_key_here
```

### Shell Completions (Optional)

Enable tab completion for Mako commands:

**Bash:**
```bash
mako completion bash | sudo tee /etc/bash_completion.d/mako
source /etc/bash_completion.d/mako
```

**Zsh:**
```bash
mkdir -p ~/.zsh/completions
mako completion zsh > ~/.zsh/completions/_mako

# Add to ~/.zshrc:
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

**Fish:**
```bash
mako completion fish > ~/.config/fish/completions/mako.fish
```

### Man Page (Optional)

Install the man page for offline documentation:

```bash
sudo mkdir -p /usr/local/share/man/man1
sudo cp docs/man/mako.1 /usr/local/share/man/man1/
sudo mandb  # Update man database
```

Then view with:
```bash
man mako
```

## Verification

Verify your installation:

```bash
# Check version
mako version

# Check configuration
mako config list

# Test command generation
mako ask "list files"
```

## Updating

Keep Mako up to date:

```bash
# Check for updates
mako update check

# Install latest version
mako update install
```

Or reinstall via the installation script:

```bash
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/main/scripts/install.sh | bash
```

## Uninstallation

If you need to remove Mako:

```bash
# Using uninstall script
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/main/scripts/uninstall.sh | bash

# Or manually
sudo rm /usr/local/bin/mako /usr/local/bin/mako-menu
rm -rf ~/.mako
```

## Platform-Specific Notes

### Linux

- Mako works on all major distributions (Ubuntu, Debian, Fedora, Arch, etc.)
- Requires bash or zsh as your shell
- SQLite3 is usually pre-installed, but may need the FTS5 extension

### macOS

- Works on both Intel and Apple Silicon Macs
- Requires macOS 10.15 (Catalina) or later
- May need to allow the binary in System Preferences > Security & Privacy on first run

### Windows (WSL)

- Install Mako in WSL (Windows Subsystem for Linux)
- Use the Linux installation method
- Works with WSL 1 or WSL 2

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Make sure the binary is in your PATH:
```bash
echo $PATH
which mako
```

2. If installed to `~/.local/bin`, add it to PATH:
```bash
export PATH="$PATH:$HOME/.local/bin"
```

Add this to your `~/.bashrc` or `~/.zshrc` to make it permanent.

### Permission Denied

If you get permission errors:

```bash
chmod +x /usr/local/bin/mako /usr/local/bin/mako-menu
```

### Database Errors

If you see SQLite errors about FTS5:

```bash
# Rebuild with FTS5 support
go build -tags "fts5" -o mako cmd/mako/main.go
```

### API Key Issues

If the AI features don't work:

1. Verify your API key is set:
```bash
mako config get api_key
```

2. Test API connectivity:
```bash
mako health
```

3. Check your internet connection

## Getting Help

- Documentation: https://github.com/fabiobrug/mako
- Issues: https://github.com/fabiobrug/mako/issues
- Discussions: https://github.com/fabiobrug/mako/discussions
