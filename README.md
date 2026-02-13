<div align="center">

# Mako

[Installation](#installation) ✦ [Documentation](apps/cli/README.md) ✦ [Features](#features) ✦ [Contributing](#contributing) ✦ [Discord](#) ✦ [Twitter/X](#)

Transform your terminal with AI-powered command assistance.

<br>

[![CI Status](https://img.shields.io/github/actions/workflow/status/fabiobrug/mako/test.yml?branch=dev&style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)](https://github.com/fabiobrug/mako/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/fabiobrug/mako?include_prereleases&style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)](https://github.com/fabiobrug/mako/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge&labelColor=F0F0E8&color=1d4ed8)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&labelColor=F0F0E8&color=1d4ed8)](https://go.dev/)

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

Mako works on **Linux** and **macOS**. Choose from multiple AI providers including local models (Ollama) or cloud services (OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter).

---

## Choose Your Installation Method

### Option 1: One-Line Install (Recommended)

Fast installation with optional environment variable configuration:

```bash
# Basic installation
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash

# Or install with provider configuration
LLM_PROVIDER=openai LLM_MODEL=gpt-4o-mini LLM_API_KEY=sk-your-key \
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash
```

**After installation, configure your AI provider:**

```bash
# Start Mako
mako

# Inside Mako shell, configure your provider:
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# Or for Ollama (local):
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# View all settings:
mako config list
```

**Supported configuration keys:**
- `llm_provider` - AI provider (openai, anthropic, gemini, deepseek, openrouter, ollama)
- `llm_model` - Model name (provider-specific)
- `llm_base_url` - Base URL (optional, for custom endpoints)
- `api_key` - Your API key (not required for Ollama)

**Check your configuration:**
```bash
# View current settings
mako config list

# Check provider and API status
mako health
```

---

### Option 2: From Source with .env File

Clone the repository and configure via `.env` file:

```bash
# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Copy and edit configuration
cp .env.example .env
nano .env  # Edit with your provider settings

# Build
make build

# Install (optional, requires sudo)
make install

# Or run directly
./mako
```

**Example `.env` configuration:**

```bash
# OpenAI
LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-key

# Or Ollama (local, free)
LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434
```

---

#### Verify Installation

```bash
# Start Mako shell
mako

# Inside Mako shell, try:
mako ask "find files larger than 100MB"
mako history
mako health    # Check configuration status
mako help
```

### Configuration Priority

Mako checks for configuration in this order:

1. **Environment variables** (`.env` file in `apps/cli/`)
2. **Config file** (`~/.mako/config.json`) - set via `mako config set`
3. **Default values** (Gemini provider)

You can use either method, or combine them. Environment variables take precedence.

---

### AI Provider Configuration

Mako supports multiple AI providers. Configure your preferred provider using environment variables or CLI commands:

#### Quick Setup

```bash
# Navigate to CLI directory
cd apps/cli

# Copy the example configuration
cp .env.example .env

# Edit the file and set your preferred provider
nano .env
```

#### Supported Providers

| Provider | Type | Cost | Best For |
|----------|------|------|----------|
| **Ollama** | Local | Free | Privacy, offline use, no API costs |
| **OpenAI** | Cloud | Paid | Best quality, GPT-4o models |
| **Anthropic** | Cloud | Paid | Claude models, great reasoning |
| **Google Gemini** | Cloud | Free tier available | Default option, good balance |
| **OpenRouter** | Cloud | Paid | Access to multiple models |
| **DeepSeek** | Cloud | Paid | Cost-effective alternative |

#### Example Configurations

**Ollama (Local, Free)**
```bash
LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434
```

**OpenAI**
```bash
LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-api-key-here
```

**Anthropic (Claude)**
```bash
LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-your-key
```

**Google Gemini**
```bash
LLM_PROVIDER=gemini
LLM_MODEL=gemini-2.5-flash
LLM_API_KEY=your-gemini-key
```

#### Using Ollama (Local Models)

Ollama allows you to run AI models locally on your machine:

```bash
# Install Ollama
curl https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2

# Navigate to CLI directory and configure Mako
cd apps/cli
echo "LLM_PROVIDER=ollama" > .env
echo "LLM_MODEL=llama3.2" >> .env
echo "LLM_API_BASE=http://localhost:11434" >> .env

# Start Mako
./mako
```

Benefits of Ollama:
- ✅ Completely free
- ✅ Works offline
- ✅ Privacy - data never leaves your machine
- ✅ No API rate limits

### How It Works

1. **Start Mako** - Wraps around your bash/zsh shell
2. **Natural Language** - Type `mako ask "compress this video"` 
3. **AI Generation** - Your configured AI provider generates the appropriate shell command
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
| AI Providers | OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Ollama |
| Embeddings | Provider-specific (Gemini, OpenAI, Ollama) |
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

Made with care by [Fabio Brugnara](https://github.com/fabiobrug)

---

<div align="center">

**Star the repo** to support development and get notified of new releases.

[![GitHub Stars](https://img.shields.io/github/stars/fabiobrug/mako?style=social)](https://github.com/fabiobrug/mako)

</div>
