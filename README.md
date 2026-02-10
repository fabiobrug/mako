# Mako - AI-Native Shell Orchestrator

<p align="center">
  <strong>Transform your terminal with AI-powered command assistance</strong>
</p>

<p align="center">
  <a href="https://github.com/fabiobrug/mako/releases"><img src="https://img.shields.io/github/v/release/fabiobrug/mako?include_prereleases&style=for-the-badge" alt="GitHub release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="MIT License"></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go" alt="Go version"></a>
</p>

**Mako** is an AI-native shell orchestrator that wraps around your existing shell (bash/zsh) to provide intelligent command assistance. Generate commands from natural language, search your history semantically with vector embeddings, and work faster with an AI that learns your preferences and understands context.

The tool intercepts terminal I/O through a PTY (pseudo-terminal) and routes commands to Google Gemini for natural language processing, while maintaining a high-performance searchable history with async embedding generation and LRU caching.

## Quick Start

Runtime: **Go 1.21**

```bash
# Clone and build
git clone https://github.com/fabiobrug/mako.git
cd mako

# Set up API key
echo "GEMINI_API_KEY=your_api_key_here" > .env

# Build binaries (fts5 tag required for full-text search)
go build -tags "fts5" -o mako cmd/mako/main.go
go build -o mako-menu cmd/mako-menu/main.go

# Start Mako
./mako

# Inside Mako shell
mako ask "find files larger than 100MB"
mako history --interactive
mako health
mako stats
```

Get your Gemini API key: [Google AI Studio](https://ai.google.dev/)

## What's New in v1.0.0 🎉

**First Stable Release - Production Ready!**

Week 12 brought professional distribution and ease of use:
- **One-Command Install**: `curl -sSL https://get-mako.sh | bash`
- **Auto-Updates**: Stay current with `mako update`
- **Configuration**: Easy setup with `mako config`
- **Professional Docs**: Man page, installation guide, shell completions
- **Package Managers**: Homebrew support for easy installation

Plus all the powerful features from Weeks 1-11:
- AI-powered command generation from natural language
- Semantic history search with vector embeddings
- Async embedding generation (20x faster saves)
- LRU caching with 80%+ hit rate
- Command aliases and personalization
- Health diagnostics and performance metrics

See [CHANGELOG.md](CHANGELOG.md) for complete details.

## Features

### Core Intelligence
- **Natural Language to Commands** - Type `mako ask "compress this video"` and get the right command
- **Multi-Turn Conversations** - Mako remembers last 5 exchanges for context-aware refinements
- **Smart Context Switching** - Auto-detects project types (Go, Node, Python, etc.) and suggests appropriate commands
- **Personalization & Learning** - Learns your preferred flags and options after 3+ uses

### High-Performance Search
- **Async Embedding Generation** - Commands save in <10ms with background worker pool
- **LRU Embedding Cache** - 80%+ hit rate, 10,000 entry capacity with persistent storage
- **Two-Phase Semantic Search** - FTS5 keyword filter + vector similarity ranking, scales to 100k+ commands
- **Interactive History Browser** - Browse, re-run, and view full output of past commands

### Advanced Features
- **Database Deduplication** - SHA256 hash-based duplicate detection, 30-50% size reduction
- **Export/Import System** - JSON-based backup and sharing with conflict resolution (`skip`, `merge`, `overwrite`)
- **Batch History Sync** - Incremental sync from bash history with timestamp tracking
- **Health Check Diagnostics** - Comprehensive system health with performance metrics and actionable tips
- **Smart Aliases** - Create parameterized aliases with `$1`, `$2`, `$@`, `$#` support and tag organization
- **Error Autopsy** - Automatically explain why commands failed with AI analysis
- **Secret Redaction** - Sensitive data automatically removed from history
- **Enhanced Command Composition** - AI understands pipes, `&&`, `||`, `;` operators with pipeline validation
- **Project-Aware Commands** - Auto-detects Go, Node, Python, Rust, Django, React, and suggests context-appropriate commands

## How It Works

```
User Input → PTY Master → Bash Shell → PTY Slave → Stream Interceptor → Output
 ↑ ↓
 └─────────── Command Detection (<<<MAKO_EXECUTE>>>) ─────────────┘
                                 ↓
                          Gemini API + SQLite
```

1. Mako creates a PTY wrapper around your shell (bash/zsh)
2. All terminal I/O is intercepted and monitored
3. Shell function writes `mako` commands to `~/.mako/last_command.txt`
4. Marker `<<<MAKO_EXECUTE>>>` triggers command routing
5. AI generates commands via Gemini API
6. Commands stored in SQLite with vector embeddings
7. Semantic search uses embeddings to find similar commands by intent

## Core Commands

**Inside Mako shell:**

```bash
# Natural language command generation (with conversation context)
mako ask "list all docker containers"
mako ask "only running ones"    # Mako remembers previous context!
mako clear                      # Clear conversation history (auto-expires after 5 min)

# History management
mako history                    # Show recent commands
mako history <keyword>          # Search by keyword  
mako history semantic "backup"  # Search by meaning
mako history --failed           # Show only failed commands
mako history --success          # Show only successful commands
mako history --interactive      # Browse history interactively

# Aliases with parameters
mako alias save <name> <cmd> [--tags tag1,tag2]
mako alias list [--tag <tag>]
mako alias run <name> [args]    # $1, $2, $@, $# supported
mako alias delete <name>
mako alias export <file>
mako alias import <file>

# Performance & diagnostics
mako health                     # System health check and performance metrics
mako stats                      # Show usage statistics
mako sync                       # Manually sync bash history

# Export/Import
mako export [--last N] [--dir path] > file.json
mako import [--merge|--skip|--overwrite] file.json

# Utilities
mako help                       # Display help information
mako version                    # Show Mako version
mako draw                       # Display shark ASCII art
```

**Regular shell commands** work normally and are automatically tracked with AI embeddings processed asynchronously.

## Usage Examples

### Multi-Turn Conversations

Mako remembers your last 5 exchanges for context-aware refinements:

```bash
$ mako ask "find large files"
→ find . -type f -size +100M

$ mako ask "only PDFs"
→ find . -type f -name "*.pdf" -size +100M

$ mako ask "sort by size"
→ find . -type f -name "*.pdf" -size +100M -exec ls -lh {} \; | sort -k5 -h
```

### Project-Aware Commands

Mako detects your project type and suggests appropriate commands:

```bash
# In a Go project
$ mako ask "test"
→ go test ./...

# In a Node project
$ mako ask "test"
→ npm test

# In a Python project with pytest
$ mako ask "test"
→ pytest
```

### Personalization

Mako learns your preferences after 3+ uses:

```bash
# After using `ls -lah` multiple times
$ mako ask "list files"
→ ls -lah  # Mako learned you prefer -lah flags
```

### Interactive History

Browse and re-run commands interactively:

```bash
$ mako history --interactive
# Browse with arrow keys, select to run, copy, or view full output

$ mako history --failed
# Show only failed commands for debugging

$ mako history semantic "compress video"
# Find commands by meaning, not exact text
```

## Installation

### Prerequisites

- Go 1.25.6 or higher
- SQLite3 with FTS5 support
- Gemini API key from [Google AI Studio](https://ai.google.dev/)

### From Source

```bash
git clone https://github.com/fabiobrug/mako.git
cd mako

# Set up environment variables
cat > .env <<EOF
GEMINI_API_KEY=your_api_key_here
EOF

# Build with FTS5 tag (required for full-text search)
go build -tags "fts5" -o mako cmd/mako/main.go

# Build menu binary (required for interactive menus)
go build -o mako-menu cmd/mako-menu/main.go

# Run
./mako
```

### Build Requirements

- The main binary requires the `fts5` build tag for SQLite full-text search
- The menu binary (`mako-menu`) must be in the same directory as `mako`
- Both binaries are required for full functionality
- Database automatically migrates on first run for existing installations

## Architecture

### PTY Flow

```
User Input → PTY Master → Bash Shell → PTY Slave → Stream Interceptor → Output
```

The PTY (pseudo-terminal) layer allows Mako to intercept all terminal I/O while maintaining compatibility with existing shell features like job control, signals, and raw terminal mode.

### Command Interception

1. User types `mako ask "natural language"`
2. Shell function writes command to `~/.mako/last_command.txt`
3. Shell function prints marker: `<<<MAKO_EXECUTE>>>`
4. Stream interceptor detects marker
5. Reads command file, routes to `internal/shell/commands.go`
6. AI processes request and returns generated command
7. User can edit, explain, or execute the command

### Interactive Menu System

- **Main binary**: `mako` - Shell orchestrator with PTY management
- **Menu binary**: `mako-menu` - Standalone TUI for user choices
- **Communication**: Pause file (`~/.mako/pause_input`) stops PTY input during menu display
- **Direct I/O**: Both binaries use `/dev/tty` for direct terminal access

## Project Structure

```
mako/
├── cmd/
│   ├── mako/           # Main shell orchestrator (PTY wrapper)
│   │   └── main.go
│   └── mako-menu/      # Interactive menu TUI
│       └── main.go
├── internal/
│   ├── ai/             # AI & embeddings
│   │   ├── gemini.go           # Gemini API client
│   │   ├── embeddings.go       # Vector embeddings
│   │   ├── conversation.go     # Multi-turn conversations
│   │   ├── personalization.go  # Preference learning
│   │   └── context.go          # System context
│   ├── database/       # Database operations
│   │   ├── db.go              # Core SQLite operations
│   │   ├── async.go           # Async embedding worker pool
│   │   └── sync.go            # Bash history sync
│   ├── cache/          # Performance optimization
│   │   └── embedding.go       # LRU embedding cache
│   ├── export/         # Import/export system
│   │   ├── format.go          # JSON schema
│   │   ├── export.go          # Export functionality
│   │   └── import.go          # Import with conflict resolution
│   ├── health/         # Diagnostics
│   │   └── health.go          # Health check system
│   ├── context/        # Context detection
│   │   └── project.go         # Project type detection
│   ├── parser/         # Command analysis
│   │   └── command.go         # Pipeline validation
│   ├── safety/         # Security
│   │   └── validator.go       # Safety validation & secret redaction
│   ├── alias/          # Alias system
│   │   └── alias.go           # Alias management with parameters
│   ├── shell/          # Command execution
│   │   └── commands.go        # Command routing & handling
│   ├── stream/         # PTY management
│   │   └── interceptor.go     # Stream interception
│   ├── buffer/         # Output buffering
│   │   └── ringbuffer.go      # Ring buffer for recent output
│   └── ui/             # User interface
│       └── menu.go            # Menu utilities
└── .env                # API keys (gitignored)
```

## Configuration

Mako stores data in `~/.mako/`:

- `mako.db` - SQLite database with command history and embeddings
- `conversation.json` - Multi-turn conversation history (auto-expires after 5 min)
- `preferences.json` - Learned command preferences
- `aliases.json` - Saved command aliases with tags
- `last_command.txt` - IPC file for command passing (temporary)
- `pause_input` - Signal file to pause PTY input during menus (temporary)

No configuration file is required. The tool works out of the box with sensible defaults.

### Database Schema

The database automatically migrates to the latest schema on startup:
- Command deduplication with SHA256 hashing
- Async embedding status tracking (`pending` → `processing` → `completed`)
- Embedding cache table with hit count tracking
- Sync metadata for incremental bash history sync
- Optimized indexes for performance at scale

## Technology Stack

### Core Technologies

- **Language**: Go 1.25.6
- **AI Provider**: Google Gemini API
  - `gemini-2.5-flash` for command generation, explanations, and error analysis
  - `text-embedding-004` for semantic search with 768-dimensional vectors
- **Database**: SQLite with FTS5 extension for hybrid search
- **Terminal**: PTY via `creack/pty`

### Key Dependencies

- `creack/pty` v1.1.24 - PTY (pseudo-terminal) handling
- `mattn/go-sqlite3` v1.14.33 - SQLite driver with FTS5
- `atotto/clipboard` v0.1.4 - Clipboard operations
- `joho/godotenv` v1.5.1 - Environment variable management
- `golang.org/x/term` v0.39.0 - Terminal control

Full dependency list: `go.mod`

### Build Command

```bash
go build -tags "fts5" -o mako cmd/mako/main.go
```

The `fts5` tag enables SQLite full-text search capabilities for hybrid FTS5 + vector search.

### Performance Characteristics

- **Command save**: <10ms (async embedding generation)
- **Semantic search**: <100ms (100k+ commands with two-phase FTS5 filtering)
- **Cache hit rate**: 80%+ (typical usage with LRU cache)
- **Database size**: 30-50% reduction with deduplication
- **Startup time**: <100ms with cache preloading

## Development Status

**Current Version**: v1.0.0 (2026-02-10) - **First Stable Release!**

### Implemented Features

**Week 12 - Production Ready** (v1.0.0):
- One-command installation (`curl -sSL https://get-mako.sh | bash`)
- Configuration management system (`mako config`)
- Auto-update mechanism (`mako update`)
- First-run setup wizard
- Shell completions (bash/zsh/fish)
- Professional man page
- Homebrew formula
- GitHub Actions CI/CD
- Clean uninstall script

**Week 11 - Performance & Scale**:
- Async embedding generation with worker pool architecture (<10ms command saves)
- LRU embedding cache with 10,000 entry capacity (80%+ hit rate)
- Database deduplication with SHA256 hashing (30-50% size reduction)
- Two-phase semantic search (FTS5 + vector similarity, <100ms for 100k commands)
- Export/import system with JSON format and conflict resolution
- Batch history sync with incremental updates
- Health check diagnostics with performance metrics

**Week 10 - Advanced AI Features**:
- Multi-turn conversations with 5-minute auto-timeout
- Smart context switching with project type detection (Go, Node, Python, Rust, etc.)
- Enhanced command composition with pipeline intelligence
- Personalization & learning (learns preferred flags after 3+ uses)

**Week 9 - Feature Expansion**:
- AI-powered command alternatives and explanations
- Simple line editor for command modification
- Alias parameters with `$1`, `$2`, `$@`, `$#` support
- Alias tags and categories
- Import/export aliases for sharing
- History filters (`--failed`, `--success`)
- Output preview in history
- Interactive history browser

**Week 8 - Foundation**:
- Command explanation before execution
- Edit before running workflow
- Alias system with parameterization
- Enhanced history display with status icons and timestamps

**Core Foundation** (v0.1.x):
- PTY-based shell wrapper with stream interception
- AI command generation via Gemini 2.5 Flash
- Command history with vector embeddings
- Full-text and semantic search (FTS5 + embeddings)
- Interactive menu system with keyboard navigation
- Context-aware suggestions using recent output
- Safety guardrails for destructive commands
- Error autopsy with automatic explanations
- Secret redaction from history

### Roadmap

**Planned Enhancements**:
- Command templating system with variables
- Multi-line command editing support
- Plugin system for custom extensions
- Enhanced secret detection patterns
- Configurable history retention policies
- Additional AI models support
- Shell scripting assistance

## Development

### Building

```bash
# Build both binaries
go build -tags "fts5" -o mako cmd/mako/main.go && \
go build -o mako-menu cmd/mako-menu/main.go

# Clean build artifacts
rm -f mako mako-menu
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/database -v
go test ./internal/ai -v

# Run with coverage
go test -cover ./...
```

### Code Style

- Clean, working code over extensive comments
- Follow standard Go conventions
- Keep functions focused and testable
- Prefer explicit error handling

## Terminal Formatting

Mako handles terminal output carefully to ensure proper display:

- PTY output requires `\r\n` line endings (not just `\n`)
- All command output is converted: `strings.ReplaceAll(output, "\n", "\r\n")`
- Menu drawing uses ANSI escape sequences (`\033[K`, `\033[J`)
- Direct `/dev/tty` access for reliable I/O during menus

See `.cursorrules` for detailed terminal handling guidelines.

## Getting Help

- **Documentation**: See [CHANGELOG.md](CHANGELOG.md) for detailed feature documentation
- **Issues**: Open an issue on [GitHub](https://github.com/fabiobrug/mako/issues)
- **Health Check**: Run `mako health` to diagnose common problems
- **Debugging**: Check `~/.mako/` directory for logs and state files

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Follow existing code style and conventions
5. Test thoroughly (especially terminal I/O and PTY behavior)
6. Update documentation as needed
7. Submit a pull request

### Testing Guidelines

- Test all terminal formatting changes in actual PTY environment
- Verify menu navigation works with arrow keys
- Check that command output displays correctly
- Ensure no regression in existing features
- Test async embedding generation under load
- Verify cache hit rates with `mako health`

### Development Workflow

```bash
# Run tests
go test ./...

# Build and test locally
go build -tags "fts5" -o mako cmd/mako/main.go && \
go build -o mako-menu cmd/mako-menu/main.go

# Check health after changes
./mako
mako health
```

## Security Considerations

- **Secret Redaction**: Sensitive patterns (API keys, passwords, tokens) automatically removed from history
- **Destructive Command Safety**: Warnings and confirmations for potentially dangerous operations (`rm -rf`, `dd`, etc.)
- **Critical Command Blocking**: Extremely dangerous commands are blocked entirely
- **API Key Storage**: `.env` file is gitignored and never committed to version control
- **Local-First**: All data stays on your machine, AI queries sent only to Gemini API
- **Health Monitoring**: Run `mako health` to check for security and configuration issues
- **Conversation Privacy**: Conversation history auto-expires after 5 minutes of inactivity

### Secret Patterns Detected

Mako automatically redacts:
- API keys and tokens (AWS, GitHub, etc.)
- Password patterns in commands
- SSH private keys
- Database connection strings with credentials
- Generic `key=value` and `token=value` patterns

## Troubleshooting

### Menu not appearing
- Ensure `mako-menu` binary is in the same directory as `mako`
- Check file permissions: both binaries must be executable (`chmod +x mako mako-menu`)

### Output formatting issues
- Terminal line endings: Mako converts `\n` to `\r\n` automatically
- If staircase effect appears, verify stream interceptor is processing output
- Menu duplication: Old menu not clearing before redraw (this is a known issue that's been fixed in v0.5.0)

### Database errors
- Ensure build includes `-tags "fts5"` for SQLite FTS5 support
- Check permissions on `~/.mako/` directory
- Database automatically migrates on startup - if migration fails, check error logs
- For corrupted databases, backup and delete `~/.mako/mako.db` to start fresh

### API errors
- Verify `GEMINI_API_KEY` is set in `.env` file
- Check API key is valid at [Google AI Studio](https://ai.google.dev/)
- Rate limiting: Gemini API has rate limits - check `mako health` for status

### Performance issues
- Run `mako health` to check system status
- Cache hit rate below 50%: May need cache warm-up time
- Slow semantic search: Run `mako stats` to check command count
- Database size growing: Run export/import cycle to repack database

### Conversation not working
- Conversation history stored in `~/.mako/conversation.json`
- Auto-expires after 5 minutes of inactivity
- Use `mako clear` to manually reset conversation

### History sync issues
- Bash history format must be timestamped or plain
- Default path: `~/.bash_history`
- Check sync metadata: SQLite table `sync_metadata`
- Manual sync: `mako sync`

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Created by [Fabio Brug](https://github.com/fabiobrug)

## Acknowledgments

- Built with [Google Gemini API](https://ai.google.dev/) for AI-powered command generation and semantic embeddings
- Terminal handling via [creack/pty](https://github.com/creack/pty)
- Database powered by SQLite with FTS5 full-text search
- Clipboard integration via [atotto/clipboard](https://github.com/atotto/clipboard)
- Inspired by the need for better command-line AI integration with high performance at scale
