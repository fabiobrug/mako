package shell

import "fmt"

// handleCompletion generates shell completion scripts
func handleCompletion(args []string) (string, error) {
	if len(args) == 0 {
		return "Usage: mako completion <bash|zsh|fish>\r\n", nil
	}

	var script string
	switch args[0] {
	case "bash":
		script = getBashCompletion()
	case "zsh":
		script = getZshCompletion()
	case "fish":
		script = getFishCompletion()
	default:
		return fmt.Sprintf("Unknown shell: %s (supported: bash, zsh, fish)\r\n", args[0]), nil
	}

	return script + "\r\n", nil
}

func getBashCompletion() string {
	return `_mako_completions() {
    local cur prev commands
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    commands="ask history stats help version config alias export import health update sync draw clear completion uninstall"
    
    if [ $COMP_CWORD -eq 1 ]; then
        COMPREPLY=($(compgen -W "${commands}" -- ${cur}))
        return 0
    fi
    
    case "${prev}" in
        history)
            COMPREPLY=($(compgen -W "semantic" -- ${cur}))
            ;;
        alias)
            COMPREPLY=($(compgen -W "save list delete run" -- ${cur}))
            ;;
        config)
            COMPREPLY=($(compgen -W "list get set reset" -- ${cur}))
            ;;
        update)
            COMPREPLY=($(compgen -W "check install" -- ${cur}))
            ;;
        completion)
            COMPREPLY=($(compgen -W "bash zsh fish" -- ${cur}))
            ;;
        export|import)
            COMPREPLY=($(compgen -f -- ${cur}))
            ;;
    esac
}

complete -F _mako_completions mako`
}

func getZshCompletion() string {
	return `#compdef mako

_mako() {
    local -a commands
    commands=(
        'ask:Generate command from natural language'
        'history:Search command history'
        'stats:Show usage statistics'
        'alias:Manage command aliases'
        'config:Manage configuration'
        'update:Check for updates'
        'export:Export command history'
        'import:Import command history'
        'health:Show system health'
        'sync:Sync bash history'
        'help:Show help'
        'version:Show version'
        'draw:Show shark art'
        'clear:Clear screen'
        'completion:Generate shell completion'
        'uninstall:Show uninstall instructions'
    )

    if (( CURRENT == 2 )); then
        _describe 'command' commands
        return
    fi

    case "$words[2]" in
        history)
            _arguments '2:mode:(semantic)'
            ;;
        alias)
            _arguments '2:action:(save list delete run)'
            ;;
        config)
            _arguments '2:action:(list get set reset)'
            ;;
        update)
            _arguments '2:action:(check install)'
            ;;
        completion)
            _arguments '2:shell:(bash zsh fish)'
            ;;
        export|import)
            _files
            ;;
    esac
}

_mako`
}

func getFishCompletion() string {
	return `# Mako completion for fish shell

# Commands
complete -c mako -n "__fish_use_subcommand" -a ask -d "Generate command from natural language"
complete -c mako -n "__fish_use_subcommand" -a history -d "Search command history"
complete -c mako -n "__fish_use_subcommand" -a stats -d "Show usage statistics"
complete -c mako -n "__fish_use_subcommand" -a alias -d "Manage command aliases"
complete -c mako -n "__fish_use_subcommand" -a config -d "Manage configuration"
complete -c mako -n "__fish_use_subcommand" -a update -d "Check for updates"
complete -c mako -n "__fish_use_subcommand" -a export -d "Export command history"
complete -c mako -n "__fish_use_subcommand" -a import -d "Import command history"
complete -c mako -n "__fish_use_subcommand" -a health -d "Show system health"
complete -c mako -n "__fish_use_subcommand" -a sync -d "Sync bash history"
complete -c mako -n "__fish_use_subcommand" -a help -d "Show help"
complete -c mako -n "__fish_use_subcommand" -a version -d "Show version"
complete -c mako -n "__fish_use_subcommand" -a completion -d "Generate shell completion"
complete -c mako -n "__fish_use_subcommand" -a uninstall -d "Show uninstall instructions"

# Subcommands
complete -c mako -n "__fish_seen_subcommand_from history" -a semantic -d "Semantic search"
complete -c mako -n "__fish_seen_subcommand_from alias" -a "save list delete run"
complete -c mako -n "__fish_seen_subcommand_from config" -a "list get set reset"
complete -c mako -n "__fish_seen_subcommand_from update" -a "check install"
complete -c mako -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"`
}
