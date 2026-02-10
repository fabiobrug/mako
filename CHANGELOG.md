# Changelog

All notable changes to Mako will be documented in this file.

## [0.3.0] - 2026-02-09

### Added - Week 9 Major Feature Expansion

#### 🎨 Enhanced Command Explanation
- **Suggest alternatives** option in menu - AI generates 2-3 alternative ways to accomplish the same goal
- Security warnings integrated into explanations
- Comparison of different approaches with trade-offs

#### ⌨️ Simple Line Editor
- Edit before running with pre-filled command
- Backspace support for corrections
- Simple and reliable terminal handling
- No complex key handling to avoid terminal issues

#### 🔖 Advanced Alias System
- **Parameter support**: Use $1, $2, ... $n, $@, $# in alias commands
  - Example: `mako alias save rm-safe "rm -i $1"` then `mako alias run rm-safe file.txt`
- **Tags/Categories**: Organize aliases with --tags flag
  - Example: `mako alias save deploy "kubectl apply -f app.yaml" --tags kubernetes,production`
  - Filter: `mako alias list --tag kubernetes`
- **Import/Export**: Share aliases between systems
  - `mako alias export ~/my-aliases.json`
  - `mako alias import ~/my-aliases.json`
- Backward compatible with v0.2.0 alias format

#### 📊 Enhanced History System
- **Exit code filtering**:
  - `mako history --failed` - Show only failed commands
  - `mako history --success` - Show only successful commands
  - Works with keyword and semantic search
- **Output preview**: See first 60 characters of command output in history display
- **Interactive history browser**: `mako history --interactive`
  - Browse recent commands with arrow keys
  - Select and re-run commands
  - Copy commands to clipboard
  - View full output of past commands
  - Supports --failed and --success filters

### Technical Changes
- Database method `GetCommandsByExitCode()` for filtered queries
- Alias structure extended with `AliasInfo` type supporting tags
- Menu system integration for interactive history browsing

### Backward Compatibility
- Old alias format automatically migrated to new format
- Import supports both old and new alias formats
- All existing commands and features remain unchanged

---

## [0.2.0] - 2026-02-09

### Added - Week 8 Features

####  Explain Command
- New "Explain what this does" option in command menu
- AI-powered command explanations before execution
- Shows what the command does, key flags/options, and potential warnings

####  Edit Before Running
- New "Edit before running" option in command menu
- Simple line editor with pre-filled generated command
- Backspace support for corrections
- Edited commands are saved to history

####  Alias System
- Complete alias management with 4 new commands:
  - `mako alias save <name> <command>` - Save command shortcuts
  - `mako alias list` - View all saved aliases
  - `mako alias run <name>` - Execute saved aliases
  - `mako alias delete <name>` - Remove aliases
- Aliases stored in `~/.mako/aliases.json`
- Aliased commands are tracked in history

####  Enhanced History Display
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
