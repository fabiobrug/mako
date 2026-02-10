_mako_completions() {
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
            COMPREPLY=($(compgen -W "semantic --failed --success --interactive" -- ${cur}))
            ;;
        alias)
            COMPREPLY=($(compgen -W "save list delete run export import" -- ${cur}))
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

complete -F _mako_completions mako
