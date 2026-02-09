# Changelog

All notable changes to Mako will be documented in this file.

## [0.2.0] - 2026-02-09

### Added - Week 8 Features

#### 🎯 Explain Command
- New "Explain what this does" option in command menu
- AI-powered command explanations before execution
- Shows what the command does, key flags/options, and potential warnings

#### ✏️ Edit Before Running
- New "Edit before running" option in command menu
- Simple line editor with pre-filled generated command
- Backspace support for corrections
- Edited commands are saved to history

#### 🔖 Alias System
- Complete alias management with 4 new commands:
  - `mako alias save <name> <command>` - Save command shortcuts
  - `mako alias list` - View all saved aliases
  - `mako alias run <name>` - Execute saved aliases
  - `mako alias delete <name>` - Remove aliases
- Aliases stored in `~/.mako/aliases.json`
- Aliased commands are tracked in history

#### 📊 Enhanced History Display
- Beautiful status icons: ✓ (success) / ✗ (failure)
- Timestamps in `[HH:MM:SS]` format
- Execution duration (ms or seconds)
- Applies to all history commands (regular, keyword, semantic search)

### Fixed
- Alias run command output formatting (proper line endings)
- Edit command now preserves prefilled text correctly

### Changed
- Updated help text to include new alias commands
- Improved visual consistency across all features

---

## [0.1.3] - Previous Release

### Features (Weeks 1-7)
- PTY-based shell wrapper with stream interception
- AI command generation with natural language
- Command history with vector embeddings
- Full-text search (FTS5) and semantic search
- Interactive menu system for command actions
- Context-aware suggestions using recent terminal output
- Safety guardrails for dangerous commands
- Auto-explain failed commands (error autopsy)
- Secret redaction from history
- Beautiful ocean-themed UI with shark ASCII art

### Technical
- Built with Go 1.25+
- Google Gemini API integration (gemini-2.5-flash, text-embedding-004)
- SQLite with FTS5 for full-text search
- PTY handling via creack/pty

---

## Version Numbering

Mako follows semantic versioning:
- **Major** (x.0.0): Breaking changes or major architectural updates
- **Minor** (0.x.0): New features, backward compatible
- **Patch** (0.0.x): Bug fixes and minor improvements
