# Changelog

All notable changes to Mako will be documented in this file.

## [0.4.0] - 2026-02-10

### Added - Week 10 Advanced AI Features & Intelligence

#### 🧠 Smart Context Switching
- **Project type detection** - Automatically detects Go, Node, Python, Rust, Java, Ruby, PHP, Elixir projects
- **Framework awareness** - Identifies Django, Flask, Rails, Laravel, Next.js, React, Angular, Vue
- **Smart command suggestions** - Context-aware test/build/run commands based on project type
  - Example: `mako ask "test"` → `go test ./...` (Go) / `npm test` (Node) / `pytest` (Python)
- New file: `internal/context/project.go` with comprehensive project detection

#### 💬 Multi-Turn Conversations (GAME CHANGER!)
- **Conversation memory** - Mako remembers last 5 exchanges for context-aware refinements
- **Auto-timeout** - Conversations auto-clear after 5 minutes of inactivity
- **New command**: `mako clear` - Manually clear conversation history
- **Contextual refinements** - Build upon previous commands without repeating context
  - Example:
    ```bash
    mako ask "find large files"
    → find . -type f -size +100M
    
    mako ask "only PDFs"
    → find . -type f -name "*.pdf" -size +100M
    ```
- Conversation state stored in `~/.mako/conversation.json`
- New file: `internal/ai/conversation.go` for conversation management

#### ⚡ Enhanced Command Composition
- **Pipeline intelligence** - Better understanding of complex multi-stage commands
- **Operator awareness** - AI understands `|`, `&&`, `||`, `;` operators
- **Pipeline validation** - Syntax checking before command suggestion
- **Examples in prompt** - AI learns from command composition patterns
- New validation functions in `internal/parser/command.go`

#### 🎯 Personalization & Learning
- **Pattern learning** - Tracks commonly used flags and options per command
- **Smart suggestions** - Suggests your preferred flags after 3+ uses
  - Example: After using `ls -lah` repeatedly, `mako ask "list files"` → `ls -lah`
- **Preference hints** - Learned patterns injected into AI context
- Preferences stored in `~/.mako/preferences.json`
- New file: `internal/ai/personalization.go` for preference management

### Technical Changes
- Enhanced `SystemContext` with `Project` and `Preferences` fields
- Improved AI prompts with conversation history and learned preferences
- JSON-based storage for conversations and preferences (lightweight, fast)
- Pipeline complexity scoring for better command generation

### Performance
- Minimal overhead: <20ms additional latency from all new features
- Conversation: ~5-10ms (file I/O)
- Personalization: ~2-5ms (file I/O)
- Project detection: ~1-3ms (filesystem checks)

### Backward Compatibility
- All new features work transparently with existing commands
- Graceful degradation if preference/conversation files don't exist
- No breaking changes to existing functionality

---

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
