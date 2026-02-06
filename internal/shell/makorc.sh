# Mako Shell Configuration
export MAKO_ACTIVE=1

# Custom Mako PS1 (prompt)
PS1='\[\033[1;36m\]\u\[\033[0m\]@\[\033[1;35m\]\h\[\033[0m\]:\[\033[1;34m\]\w\[\033[0m\] \[\033[1;32m\]=>\[\033[0m\] '

# Optional: Add Mako info to the right side
# PROMPT_COMMAND='echo -ne "\033[s\033[0;$((COLUMNS-15))H[\033[1;32mMako\033[0m]\033[u"'

echo " Mako shell active - type 'mako help' for commands"
