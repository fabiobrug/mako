# Contributing to Mako

Thank you for your interest in contributing to Mako! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mako.git
   cd mako
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/fabiobrug/mako.git
   ```

## Development Setup

### Requirements

- Go 1.21 or later
- SQLite3 with FTS5 support
- Git
- A Gemini API key for testing

### Building from Source

```bash
# Build main binary
go build -tags "fts5" -o mako cmd/mako/main.go

# Build menu binary
go build -o mako-menu cmd/mako-menu/main.go

# Run tests
go test -v -tags "fts5" ./...
```

### Project Structure

```
mako/
├── cmd/
│   ├── mako/           # Main binary
│   └── mako-menu/      # Interactive menu binary
├── internal/
│   ├── ai/             # AI/Gemini integration
│   ├── alias/          # Alias management
│   ├── cache/          # Embedding cache
│   ├── config/         # Configuration management
│   ├── context/        # Context extraction
│   ├── database/       # SQLite database
│   ├── export/         # Import/export
│   ├── health/         # Health checks
│   ├── parser/         # Command parsing
│   ├── safety/         # Safety validation
│   ├── shell/          # Shell commands
│   ├── stream/         # PTY stream interception
│   └── ui/             # User interface
├── scripts/            # Installation/uninstall scripts
├── packaging/          # Distribution files
│   ├── homebrew/       # Homebrew formula
│   └── completions/    # Shell completions
├── docs/               # Documentation
└── .github/
    └── workflows/      # CI/CD workflows
```

## Making Changes

### Branching

Create a feature branch from `dev`:

```bash
git checkout dev
git pull upstream dev
git checkout -b feature/your-feature-name
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `go vet` to catch common issues
- Keep functions focused and well-named
- Add comments for complex logic

### Terminal Output

**CRITICAL**: When adding or modifying output to the terminal:

1. **Always use `\r\n` for line endings** in PTY context:
   ```go
   output := strings.ReplaceAll(text, "\n", "\r\n")
   ```

2. **Use ANSI escape codes carefully**:
   ```go
   cyan := "\033[38;2;0;209;255m"
   reset := "\033[0m"
   fmt.Printf("%sText%s\r\n", cyan, reset)
   ```

3. **Clear properly when redrawing**:
   ```go
   tty.WriteString("\033[J")  // Clear from cursor to end
   ```

See `.cursorrules` for detailed terminal formatting guidelines.

### Testing

Add tests for new features:

```bash
# Run all tests
go test -v -tags "fts5" ./...

# Run specific package tests
go test -v -tags "fts5" ./internal/ai

# Run with coverage
go test -v -tags "fts5" -cover ./...
```

### Committing

Write clear commit messages:

```
type(scope): brief description

Longer description if needed.

Fixes #issue_number
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat(ai): add conversation context support

- Maintain conversation history across commands
- Store in ~/.mako/conversations/
- Add clear command to reset

Fixes #123
```

## Submitting Changes

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub:
   - Target the `dev` branch (not `main`)
   - Provide a clear description
   - Reference any related issues
   - Include screenshots for UI changes

3. **Wait for review**:
   - Address feedback promptly
   - Keep your PR updated with `dev`

## Pull Request Guidelines

### Good PRs

- Focus on a single feature/fix
- Include tests for new functionality
- Update documentation as needed
- Follow existing code style
- Keep changes minimal and focused

### PR Checklist

- [ ] Code follows project style
- [ ] Tests pass locally
- [ ] New tests added for new features
- [ ] Documentation updated
- [ ] Terminal output uses proper line endings
- [ ] No breaking changes (or documented)
- [ ] Commit messages are clear

## Areas to Contribute

### High Priority

- **Performance improvements**: Faster embedding generation, better caching
- **Test coverage**: Add tests for untested code
- **Documentation**: Improve README, add tutorials, fix typos
- **Bug fixes**: Check the issues page

### Medium Priority

- **New commands**: Add useful Mako commands
- **AI improvements**: Better prompts, context extraction
- **Shell support**: Add support for other shells (fish, etc.)
- **Export formats**: Add CSV, XML export options

### Good First Issues

Look for issues labeled `good-first-issue`:
- Documentation fixes
- Simple bug fixes
- Adding examples
- Improving error messages

## Reporting Bugs

Create an issue with:

1. **Description**: Clear description of the bug
2. **Steps to reproduce**: Exact steps to trigger the bug
3. **Expected behavior**: What should happen
4. **Actual behavior**: What actually happens
5. **Environment**:
   - OS and version
   - Mako version (`mako version`)
   - Go version (`go version`)
   - Shell (bash/zsh)
6. **Logs**: Any relevant error messages

Example:
```
**Description**
Semantic search returns no results even when commands exist.

**Steps to reproduce**
1. Add some git commands to history
2. Run: mako history semantic "git operations"
3. See "No commands found"

**Expected**
Should show git commands

**Actual**
Shows empty results

**Environment**
- OS: Ubuntu 22.04
- Mako: v1.0.0
- Go: 1.21.5
- Shell: bash 5.1.16

**Logs**
(paste relevant logs here)
```

## Feature Requests

Create an issue with:

1. **Description**: What feature you'd like
2. **Use case**: Why you need it
3. **Alternatives**: Other solutions you've considered
4. **Examples**: How it would work

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions
- No harassment or discrimination

## Questions?

- Open a [Discussion](https://github.com/fabiobrug/mako/discussions)
- Ask in the pull request
- Check existing issues and docs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Mako!
