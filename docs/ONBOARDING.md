# Mako Onboarding System

## Overview

Mako now features a polished, interactive onboarding wizard that helps new users get started quickly. The wizard guides users through provider selection, API key configuration, and initial settings.

## Features

### 🎨 Beautiful Visual Design
- Colorful ASCII art welcome banner
- Provider-specific colors:
  - **Gemini**: Blue (#4285F4)
  - **Claude**: Orange (#FF8A4C)
  - **OpenAI**: Green (#10A37F)
  - **DeepSeek**: Purple (#8A2BE2)
  - **Ollama**: Yellow/Gold (#FFD700)
  - **OpenRouter**: Magenta (#FF1493)

### ⚙️ Multi-Provider Support
- Configure multiple AI providers during setup
- Easy switching between providers
- Automatic API key management
- Local-first option with Ollama

### 🔐 Security-Focused
- Secure API key storage in `~/.mako/.env`
- File permissions set to `0600` (read/write for owner only)
- API keys masked in display
- 100% private - no telemetry

### ✨ User-Friendly
- Skippable wizard (type 'skip' at start)
- Smart defaults (Gemini recommended)
- Validates API keys
- Detects existing environment variables
- Clear error messages and help

## User Flow

### First Run Experience

When you run Mako for the first time:

```
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║          ███╗   ███╗ █████╗ ██╗  ██╗ ██████╗        ║
║          ████╗ ████║██╔══██╗██║ ██╔╝██╔═══██╗       ║
║          ██╔████╔██║███████║█████╔╝ ██║   ██║       ║
║          ██║╚██╔╝██║██╔══██║██╔═██╗ ██║   ██║       ║
║          ██║ ╚═╝ ██║██║  ██║██║  ██╗╚██████╔╝       ║
║          ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝        ║
║                                                       ║
║           AI-Native Shell Orchestrator                ║
║              Welcome! Let's get started.              ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝

Press Enter to continue, or type 'skip' to use defaults:
```

### Step 1: Select AI Providers

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Select AI Providers
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

You can configure multiple providers and switch between them later.

  1. gemini ☁ Google's Gemini (Fast & Free tier available)
  2. anthropic ☁ Anthropic Claude (High quality reasoning)
  3. openai ☁ OpenAI GPT Models (Industry standard)
  4. deepseek ☁ DeepSeek (Cost-effective alternative)
  5. openrouter ☁ OpenRouter (Access to multiple models)
  6. ollama 🏠 Ollama (Local, private, free)

Enter provider numbers separated by commas (e.g., 1,2,3)
Or press Enter for Gemini (recommended):
```

### Step 2: Configure API Keys

For each selected cloud provider:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Configure gemini
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Google's Gemini (Fast & Free tier available)

Get your API key: https://aistudio.google.com/app/apikey

Enter your gemini API key (or press Enter to skip):
```

### Step 3: Select Default Provider

If multiple providers configured:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Select Default Provider
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Which provider would you like to use by default?

  1. gemini
  2. anthropic

Enter number (default: 1):
```

### Step 4: Additional Settings

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Additional Settings
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Safety Level (confirm dangerous commands):
  1. Low    - No confirmations
  2. Medium - Confirm destructive commands (recommended)
  3. High   - Confirm all commands

Enter number (default: 2):

Enable automatic updates? (Y/n):
```

### Step 5: Completion

```
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║              ✨ Setup Complete! ✨                  ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝

Default provider: gemini

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Quick Start Guide
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Try these commands:

  ▸ mako ask "list files in current directory"
  ▸ mako ask "find large files over 100MB"
  ▸ mako history
  ▸ mako stats

Manage providers:

  ▸ mako config providers       - View configured providers
  ▸ mako config switch <name>   - Switch active provider
  ▸ mako setup                  - Re-run this wizard

Documentation: https://github.com/fabiobrug/mako

Press Enter to start Mako...
```

## Commands

### Run/Re-run Setup Wizard

```bash
mako setup
```

Runs the interactive onboarding wizard. Can be used at any time to:
- Add new providers
- Update API keys
- Change default provider
- Adjust settings

### List Configured Providers

Inside Mako shell:

```bash
mako config providers
```

Output:
```
Configured AI Providers
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ● gemini       Google's Gemini
  ○ anthropic    Anthropic Claude
  ○ openai       OpenAI GPT Models
  ✕ deepseek     DeepSeek
  ✕ openrouter   OpenRouter
  ○ ollama       Ollama (Local)

● Active  ○ Configured  ✕ Not configured

Switch provider: mako config switch <provider>
Setup wizard:   mako setup
```

### Switch Active Provider

Inside Mako shell:

```bash
mako config switch anthropic
```

Switches to Claude/Anthropic:
```
✓ Switched to anthropic
```

Aliases supported:
- `claude` → `anthropic`
- `gpt` → `openai`

## Configuration Files

### ~/.mako/config.json

Main configuration file:

```json
{
  "version": "1.0",
  "llm_provider": "gemini",
  "llm_model": "",
  "llm_base_url": "",
  "theme": "ocean",
  "cache_size": 10000,
  "telemetry": false,
  "auto_update": true,
  "history_limit": 100000,
  "safety_level": "medium",
  "embedding_batch_size": 10
}
```

### ~/.mako/.env

Secure API key storage:

```bash
# Mako AI Provider API Keys
# Generated by Mako setup wizard

GEMINI_API_KEY=your_gemini_key_here
ANTHROPIC_API_KEY=your_anthropic_key_here
OPENAI_API_KEY=your_openai_key_here
```

**Security**: File permissions set to `0600` (owner read/write only)

## Implementation Details

### Architecture

```
apps/cli/
├── cmd/mako/main.go              # First-run detection & wizard trigger
├── internal/
│   ├── onboarding/
│   │   └── wizard.go             # Interactive wizard logic
│   ├── config/
│   │   ├── config.go             # Configuration management
│   │   └── firstrun.go           # First-run detection
│   └── shell/
│       └── config.go             # Provider management commands
```

### Key Functions

#### onboarding.RunWizard()
Main wizard entry point. Handles:
- Welcome screen
- Provider selection
- API key configuration
- Settings collection
- Configuration save

#### config.IsFirstRun()
Checks if `~/.mako/config.json` exists

#### handleProvidersList()
Displays all configured providers with status indicators

#### handleProviderSwitch(provider)
Switches active provider and validates API key availability

### Color Constants

```go
const (
	ColorGemini     = "\033[38;2;66;133;244m"  // Google Blue
	ColorClaude     = "\033[38;2;255;138;76m"  // Claude Orange
	ColorOpenAI     = "\033[38;2;16;163;127m"  // OpenAI Green
	ColorDeepSeek   = "\033[38;2;138;43;226m"  // Purple
	ColorOllama     = "\033[38;2;255;215;0m"   // Gold/Yellow
	ColorOpenRouter = "\033[38;2;255;20;147m"  // Deep Pink/Magenta
	ColorCyan       = "\033[38;2;0;209;255m"   // Mako Cyan
	ColorLightBlue  = "\033[38;2;93;173;226m"  // Mako Light Blue
	ColorGreen      = "\033[38;2;46;204;113m"  // Success Green
	ColorReset      = "\033[0m"
)
```

## Design Principles

### 1. **Skip-First Approach**
Users can skip wizard immediately if they want to configure manually later

### 2. **Smart Defaults**
- Gemini as default provider (free tier available)
- Medium safety level (balanced)
- Auto-update enabled

### 3. **Multi-Provider First**
Wizard encourages setting up multiple providers for flexibility

### 4. **Progressive Disclosure**
Only asks for information as needed (e.g., no API key for Ollama)

### 5. **Validation & Feedback**
- Checks API key length
- Detects existing environment variables
- Confirms successful configuration
- Warns about missing API keys

### 6. **Consistency**
- Uses Mako's color palette
- Matches overall CLI aesthetic
- Terminal-native interface

## Competitor Analysis

### Inspiration from Leading Tools

#### Warp
- **Adopted**: Beautiful welcome screen with ASCII art
- **Adopted**: Step-by-step guided flow
- **Improved**: Made skippable, added multi-provider setup

#### Claude Code Design
- **Adopted**: Clean, modern interface
- **Adopted**: Provider selection with clear descriptions
- **Improved**: Added color coding, status indicators

#### Mako Improvements
- **Multi-provider configuration** in single flow
- **Color-coded providers** for easy identification
- **Instant provider switching** without restart
- **Privacy-first** with local API key storage
- **Ollama integration** for local-only users

## Future Enhancements

### Planned Features
- [ ] API key validation (test connection during setup)
- [ ] Model selection per provider
- [ ] Import/export configuration
- [ ] Team configuration sharing (opt-in)
- [ ] Configuration profiles (work, personal, etc.)
- [ ] Interactive tutorial mode
- [ ] Health check after setup

### Nice-to-Have
- [ ] Video walkthrough embedded in terminal
- [ ] Example command demonstration
- [ ] Integration testing mode
- [ ] Cost estimation per provider
- [ ] Performance comparison

## Testing

### Manual Testing Checklist

- [ ] First run with no config
- [ ] Skip wizard and use defaults
- [ ] Select single provider
- [ ] Select multiple providers
- [ ] Configure with existing env vars
- [ ] Invalid API key handling
- [ ] Re-run wizard to add providers
- [ ] Switch between providers
- [ ] List configured providers
- [ ] Color display on different terminals

### Edge Cases

- [ ] No internet connection
- [ ] Invalid provider selection
- [ ] Empty API key input
- [ ] Ollama not installed
- [ ] Corrupted config file
- [ ] Permission issues with .mako directory

## Troubleshooting

### Wizard doesn't appear
```bash
# Manually trigger wizard
mako setup
```

### API key not working
```bash
# Check configuration
mako config list

# Edit .env file manually
nano ~/.mako/.env

# Verify provider
mako config providers
```

### Provider switch fails
```bash
# Check if provider has API key
cat ~/.mako/.env | grep PROVIDER_KEY

# Re-run setup to configure
mako setup
```

### Colors not displaying
Terminal may not support 24-bit color. Try:
```bash
# Check terminal capabilities
echo $TERM

# Use a modern terminal (iTerm2, Alacritty, etc.)
```

## Contributing

To enhance the onboarding system:

1. **UI/UX improvements**: Edit `internal/onboarding/wizard.go`
2. **Add providers**: Update provider list and color map
3. **New steps**: Add to wizard flow in `RunWizard()`
4. **Configuration options**: Extend `config.Config` struct

## Related Documentation

- [Adding Providers Guide](./ADDING_PROVIDERS.md)
- [Configuration Reference](./CONFIGURATION.md)
- [Growth Plan](../GROWTH.md)

---

**Status**: ✅ Implemented (Growth Plan Section 4)
**Impact**: ⭐⭐ High (User Onboarding)
**Version**: Added in v1.3.6
