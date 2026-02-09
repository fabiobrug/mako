# 🦈 Mako

**AI-Native Shell Orchestrator**

Mako is an intelligent shell wrapper that brings AI-powered command assistance to your terminal. Generate commands from natural language, search your history semantically, and let AI help you work faster.

## Features

- **Natural Language Commands**: Type `mako ask "find large files"` → get the right command
- **Semantic History Search**: Find past commands by meaning, not just keywords
- **Smart Command Tracking**: Automatic history with vector embeddings
- **Beautiful Interface**: Shark-themed prompt with ocean colors
- **Zero Friction**: Works alongside your normal shell workflow

## Installation

### Prerequisites

- Go 1.25+
- SQLite3
- Gemini API key ([Get one here](https://ai.google.dev/))

### Build

```bash
# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako

# Set up environment
echo "GEMINI_API_KEY=your_api_key_here" > .env

# Build (requires fts5 tag for SQLite full-text search)
go build -tags "fts5" -o mako cmd/mako/main.go

# Build menu binary (required for interactive menus)
go build -o mako-menu cmd/mako-menu/main.go
```

## Usage

### Start Mako Shell

```bash
./mako
```

### Commands

**Inside Mako shell:**

```bash
# Generate command from natural language
mako ask "list all files larger than 100MB"

# Show recent command history
mako history

# Search history by keyword
mako history grep

# Search history by meaning/intent
mako history semantic "compress video"

# Show statistics
mako stats

# Get help
mako help
```

**Regular shell commands work normally** - they're automatically saved with AI embeddings for future semantic search.

## How It Works

1. Mako wraps your shell (bash/zsh) using a PTY (pseudo-terminal)
2. Commands you run are captured and stored in SQLite with vector embeddings
3. When you use `mako ask`, it generates shell commands using Gemini AI
4. Semantic search uses embeddings to find similar past commands by meaning

## Architecture

```
┌─────────────┐
│    User     │
└──────┬──────┘
       │
       ↓
┌─────────────────────────────────┐
│   Mako Shell Wrapper (PTY)      │
│   - Stream Interception         │
│   - Command Detection           │
└────────┬──────────────────┬─────┘
         │                  │
         ↓                  ↓
┌─────────────────┐  ┌──────────────────┐
│  Gemini API     │  │  SQLite Database │
│  - Commands     │  │  - History       │
│  - Embeddings   │  │  - FTS5 Search   │
└─────────────────┘  └──────────────────┘
```

## Project Status

**Current Version**: v0.1.2

**Completed Features** (Weeks 1-5):
- ✅ PTY-based shell wrapper
- ✅ AI command generation
- ✅ Command history with embeddings
- ✅ Full-text and semantic search
- ✅ Interactive menu system

**Planned** (Week 6+):
- 🔄 Context-aware suggestions (using recent terminal output)
- 🔄 Safety guardrails (dangerous command detection)
- 🔄 Error autopsy (auto-explain failed commands)
- 🔄 Secret redaction from history

## Tech Stack

- **Language**: Go 1.25+
- **AI**: Google Gemini API (`gemini-2.5-flash`, `text-embedding-004`)
- **Database**: SQLite with FTS5 (full-text search)
- **Terminal**: PTY via `creack/pty`
- **Dependencies**: See `go.mod`

## Development

```bash
# Run tests
go test ./...

# Run specific test
go test ./internal/buffer -v

# Build both binaries
go build -tags "fts5" -o mako cmd/mako/main.go && \
go build -o mako-menu cmd/mako-menu/main.go

# Clean build artifacts
rm -f mako mako-menu
```

## Configuration

Mako stores data in `~/.mako/`:
- `history.db` - SQLite database with command history
- `last_command.txt` - IPC file for command passing (temporary)
- `pause_input` - Signal file for menu system (temporary)

## Contributing

Contributions welcome! Please:
1. Test your changes thoroughly
2. Follow existing code style
3. Include relevant tests
4. Update documentation if needed

## License

MIT License - see LICENSE file for details

## Author

Created by [@fabiobrug](https://github.com/fabiobrug)

## Acknowledgments

- Built with [Google Gemini API](https://ai.google.dev/)
- Terminal handling via [creack/pty](https://github.com/creack/pty)
- Shark ASCII art for the aesthetic 🦈
