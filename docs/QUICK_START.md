# Mako Quick Start Guide

Get started with Mako in under 2 minutes!

## Installation

### One Command (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/main/scripts/install.sh | bash
```

This will:
1. Download Mako for your platform
2. Install to `/usr/local/bin`
3. Set up configuration
4. Prompt for your API key

### Alternative: Homebrew

```bash
brew tap fabiobrug/mako
brew install mako
```

## Get Your API Key

1. Visit https://ai.google.dev/
2. Sign in with Google
3. Create a new API key (free)
4. Set it in Mako:

```bash
mako config set api_key YOUR_API_KEY
```

## First Steps

### 1. Start Mako

```bash
mako
```

This starts the Mako shell - your bash/zsh with AI superpowers!

### 2. Generate Your First Command

```bash
mako ask "list all files larger than 10MB"
```

Mako will generate: `find . -type f -size +10M`

### 3. Search Your History

```bash
# Search by keyword
mako history docker

# Search by meaning
mako history semantic "git operations"
```

### 4. Save Useful Commands

```bash
# Save a command as an alias
mako alias save deploy "git push && ssh server './deploy.sh'"

# Run it later
mako alias run deploy
```

## Essential Commands

### Command Generation
```bash
mako ask "your question in natural language"
```

### History
```bash
mako history                    # Recent commands
mako history semantic "query"   # AI-powered search
mako history --failed           # Only failed commands
```

### Configuration
```bash
mako config list                # View all settings
mako config set theme ocean     # Change a setting
```

### Aliases
```bash
mako alias save <name> <cmd>    # Save alias
mako alias list                 # List all aliases
mako alias run <name>           # Run alias
```

### System
```bash
mako stats                      # Usage statistics
mako health                     # System health check
mako update check               # Check for updates
```

## Pro Tips

### 1. Context-Aware Commands

Mako sees recent terminal output:
```bash
ls  # Shows files
mako ask "delete all .log files"  # Knows what's in current directory
```

### 2. Complex Queries

Ask for multi-step operations:
```bash
mako ask "find all python files modified today and count lines"
```

### 3. Learn Commands

Use Mako to learn new tools:
```bash
mako ask "how do I use rsync to backup with progress"
```

### 4. Save Workflows

Complex workflows become simple aliases:
```bash
mako alias save backup "rsync -avz --progress ~/Documents /backup/"
mako alias run backup
```

### 5. Interactive History

Browse history with arrow keys:
```bash
mako history --interactive
```

## Configuration

Customize Mako to your needs:

```bash
# Set API key
mako config set api_key sk-...

# Adjust cache size
mako config set cache_size 20000

# Change safety level
mako config set safety_level high

# Disable auto-update
mako config set auto_update false
```

View all options:
```bash
mako config list
```

## Shell Completion

Enable tab completion:

**Bash:**
```bash
mako completion bash | sudo tee /etc/bash_completion.d/mako
```

**Zsh:**
```bash
mako completion zsh > ~/.zsh/completions/_mako
```

**Fish:**
```bash
mako completion fish > ~/.config/fish/completions/mako.fish
```

## Troubleshooting

### Command not found

Add to PATH:
```bash
export PATH="$PATH:$HOME/.local/bin"
```

### API not working

Check configuration:
```bash
mako config get api_key
mako health
```

### Slow semantic search

First search generates embeddings (slow). Subsequent searches are fast due to caching.

## Getting Help

### In Mako
```bash
mako help
```

### Man Page
```bash
man mako
```

### Online
- Documentation: https://github.com/fabiobrug/mako
- Issues: https://github.com/fabiobrug/mako/issues
- Full guide: https://github.com/fabiobrug/mako/blob/main/docs/INSTALL.md

## What's Next?

1. **Explore aliases** - Save your common commands
2. **Try semantic search** - Find commands by what they do
3. **Check health** - `mako health` for system status
4. **Export history** - `mako export > backup.json`
5. **Share with team** - Import/export workflow automation

## Cheat Sheet

| Command | Description |
|---------|-------------|
| `mako` | Start Mako shell |
| `mako ask "<query>"` | Generate command |
| `mako history` | Show history |
| `mako history semantic "<query>"` | AI search |
| `mako alias save <name> <cmd>` | Save alias |
| `mako alias run <name>` | Run alias |
| `mako stats` | Usage stats |
| `mako health` | Health check |
| `mako config list` | Show config |
| `mako update check` | Check updates |
| `mako help` | Full help |

---

Happy hacking with Mako!
