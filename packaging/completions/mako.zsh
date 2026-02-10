#compdef mako

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
            _arguments '2:mode:(semantic --failed --success --interactive)'
            ;;
        alias)
            _arguments '2:action:(save list delete run export import)'
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

_mako
