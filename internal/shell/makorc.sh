# Mako Shell Configuration
export MAKO_ACTIVE=1

# Custom shark-themed prompt
PS1='\[\033[36m\]\w\[\033[0m\] \[\033[34m\]❯\[\033[0m\] '

# Mako command function - THIS IS CRITICAL
mako() {
    local cmd_file="$HOME/.mako/last_command.txt"
    
    # Write the command to file
    echo "$*" > "$cmd_file"
    
    # Send marker to stderr (intercepted by Mako)
    echo "<<<MAKO_EXECUTE>>>" >&2
}

# Startup message
echo " 🦈 Mako shell active - type 'mako help' for commands"
