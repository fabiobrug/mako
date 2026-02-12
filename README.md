<div align="center">

# Mako

[Installation](#installation) ✦ [Documentation](apps/cli/README.md) ✦ [Features](#features) ✦ [Contributing](#contributing) ✦ [Discord](#) ✦ [Twitter/X](#)

Transform your terminal with AI-powered command assistance.

<br>

![CI Status](https://img.shields.io/github/actions/workflow/status/fabiobrug/mako/test.yml?branch=dev&style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)
![Release](https://img.shields.io/github/v/release/fabiobrug/mako?include_prereleases&style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)
![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&labelColor=F0F0E8&color=1d4ed8)

</div>

<br>

> [!IMPORTANT]
>
> Mako is in active development. New features are being added continuously, and we welcome contributions from the community. If you have suggestions or feature requests, please open an issue on GitHub.

## What is Mako?

**Mako** is an AI-native shell orchestrator that wraps around your existing shell (bash/zsh) to provide intelligent command assistance. Generate commands from natural language, search your history semantically, and work faster with an AI that understands context.

Unlike traditional command-line tools, Mako intercepts terminal I/O through a PTY (pseudo-terminal) and routes commands to AI for natural language processing, making your command-line experience more intuitive and productive.

## Getting Started

<a id="installation"></a>

### Installation

Mako works on **Linux** and **macOS**. You'll need a Gemini API key from [Google AI Studio](https://ai.google.dev/).

#### One-Line Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash
```

#### From Source

```bash
# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Build
make build

# Install (requires sudo)
make install
```

#### Verify Installation

```bash
# Start Mako shell
mako

# Inside Mako shell, try:
mako ask "find files larger than 100MB"
mako history
mako help
```

### How It Works

1. **Start Mako** - Wraps around your bash/zsh shell
2. **Natural Language** - Type `mako ask "compress this video"` 
3. **AI Generation** - Gemini generates the appropriate shell command
4. **Review & Execute** - Review the command before running it
5. **Learn & Improve** - Mako learns your preferences over time

<a id="features"></a>

## Key Features

### Natural Language Commands

Generate shell commands from plain English. No need to remember complex syntax.

```bash
mako ask "find all PDF files modified in the last week"
# Generates: find . -name "*.pdf" -mtime -7
```

### Semantic History Search

Search your command history by meaning, not exact text. Find that command you ran months ago.

```bash
mako history semantic "backup database"
# Finds commands like: pg_dump -U postgres mydb > backup.sql
```

### Context-Aware AI

Mako understands your current directory, recent output, and command patterns to provide better suggestions.

### Safety First

Detects potentially dangerous commands before execution. Prevents accidental data loss.

```bash
# Mako will warn you about:
rm -rf /
sudo dd if=/dev/zero of=/dev/sda
```

### Usage Analytics

Track your command patterns, most-used commands, and efficiency over time.

```bash
mako stats
```

### Shell Integration

Works seamlessly with your existing bash or zsh configuration. No need to change your workflow.

### Smart History

SQLite database with FTS5 for lightning-fast full-text and semantic search across your entire command history.

## Tech Stack

This is a monorepo containing multiple applications:

```
mako/
├── apps/
│   ├── cli/          # Go-based CLI application
│   └── landing/      # Next.js landing page
├── packages/         # Shared packages (future)
├── docs/            # Documentation
└── scripts/         # Build and deployment scripts
```

### CLI Application

| Component | Technology |
|-----------|------------|
| Language | Go 1.24+ |
| PTY Handling | creack/pty |
| Database | SQLite with FTS5 (modernc.org/sqlite) |
| AI Provider | Google Gemini API (gemini-2.0-flash-exp) |
| Embeddings | text-embedding-004 |
| Build System | Make |

### Landing Page

| Component | Technology |
|-----------|------------|
| Framework | Next.js 16 (App Router) |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Build Tool | Turbopack |

## Development

### Prerequisites

| Tool | Version | Installation |
|------|---------|--------------|
| Go | 1.24+ | [golang.org](https://golang.org) |
| Node.js | 18+ | [nodejs.org](https://nodejs.org) |
| Make | Latest | Usually pre-installed on Unix systems |

### Building the CLI

```bash
cd apps/cli
make build
./mako
```

See the [CLI Documentation](apps/cli/README.md) for detailed development instructions.

### Running the Landing Page

```bash
cd apps/landing
npm install
npm run dev
```

Visit http://localhost:3000

### Project Structure

**CLI Application** (`apps/cli/`)
- `cmd/mako/` - Main application entry point
- `cmd/mako-menu/` - Interactive menu TUI
- `internal/ai/` - Gemini API integration
- `internal/database/` - SQLite and search
- `internal/shell/` - Command routing and execution
- `internal/stream/` - PTY stream interception

**Landing Page** (`apps/landing/`)
- `app/` - Next.js pages and layouts
- `components/` - React components
- `public/` - Static assets

## Documentation

- [Installation Guide](docs/INSTALL.md)
- [Quick Start Guide](docs/QUICK_START.md)
- [Contributing Guide](docs/CONTRIBUTING.md)
- [CLI Documentation](apps/cli/README.md)
- [Man Page](docs/man/mako.1)

<a id="contributing"></a>

## Contributing

We welcome contributions from everyone! Whether you're fixing bugs, adding features, improving documentation, or sharing ideas.

### How to Contribute

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/mako.git
cd mako

# Create a new branch
git checkout -b feature/your-feature-name

# Make your changes
cd apps/cli
make build
make test

# Commit and push
git add .
git commit -m "Add your feature"
git push origin feature/your-feature-name

# Open a Pull Request on GitHub
```

Please read our [Contributing Guide](docs/CONTRIBUTING.md) for detailed information.

### Development Roadmap

- Visual output analysis and suggestions
- Command history synchronization across devices
- Plugin system for custom commands
- Multi-shell support (fish, PowerShell)
- Team collaboration features
- Advanced safety rules customization

## Community

Join the conversation and get help:

- **Discord** - [Coming Soon]
- **GitHub Issues** - [Report bugs and request features](https://github.com/fabiobrug/mako/issues)
- **GitHub Discussions** - [Ask questions and share ideas](https://github.com/fabiobrug/mako/discussions)

## License

Mako is open source software licensed under the [MIT License](LICENSE).

## Creators

Made with care by [Fabio Brug](https://github.com/fabiobrug)

---

<div align="center">

**Star the repo** to support development and get notified of new releases.

[![GitHub Stars](https://img.shields.io/github/stars/fabiobrug/mako?style=social)](https://github.com/fabiobrug/mako)

</div>
