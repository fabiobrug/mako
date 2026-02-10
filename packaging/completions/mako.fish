# Mako completion for fish shell

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
complete -c mako -n "__fish_seen_subcommand_from history" -a "semantic --failed --success --interactive"
complete -c mako -n "__fish_seen_subcommand_from alias" -a "save list delete run export import"
complete -c mako -n "__fish_seen_subcommand_from config" -a "list get set reset"
complete -c mako -n "__fish_seen_subcommand_from update" -a "check install"
complete -c mako -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
