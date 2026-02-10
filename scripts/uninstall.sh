#!/bin/bash
set -e

# Mako Uninstallation Script
# Usage: curl -sSL https://get-mako.sh/uninstall.sh | bash
# Or: ./scripts/uninstall.sh

# Color codes
CYAN='\033[38;2;0;209;255m'
LIGHT_BLUE='\033[38;2;93;173;226m'
DIM_BLUE='\033[38;2;120;150;180m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
RESET='\033[0m'

# Configuration
INSTALL_LOCATIONS=(
    "/usr/local/bin/mako"
    "/usr/local/bin/mako-menu"
    "$HOME/.local/bin/mako"
    "$HOME/.local/bin/mako-menu"
)
MAKO_DIR="$HOME/.mako"
COMPLETION_LOCATIONS=(
    "/etc/bash_completion.d/mako"
    "/usr/local/etc/bash_completion.d/mako"
    "$HOME/.local/share/bash-completion/completions/mako"
    "$HOME/.zsh/completions/_mako"
    "$HOME/.config/fish/completions/mako.fish"
)

# Print functions
print_header() {
    echo -e "\n${CYAN}  Mako Uninstaller${RESET}\n"
}

print_info() {
    echo -e "${LIGHT_BLUE}ℹ${RESET} $1"
}

print_success() {
    echo -e "${LIGHT_BLUE}✓${RESET} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${RESET} $1"
}

print_error() {
    echo -e "${RED}✗${RESET} $1" >&2
}

# Count commands in database
count_commands() {
    if [ -f "$MAKO_DIR/mako.db" ]; then
        if command -v sqlite3 >/dev/null 2>&1; then
            COUNT=$(sqlite3 "$MAKO_DIR/mako.db" "SELECT COUNT(*) FROM commands" 2>/dev/null || echo "0")
            echo "$COUNT"
        else
            echo "unknown"
        fi
    else
        echo "0"
    fi
}

# Offer to export history
offer_export() {
    local count=$1
    
    if [ "$count" = "0" ]; then
        return
    fi
    
    echo ""
    print_warning "Command history will be deleted ($count commands)"
    echo ""
    
    read -p "$(echo -e ${CYAN}Export history first? \[y/N\]: ${RESET})" -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BACKUP_FILE="$HOME/mako-backup-$(date +%Y-%m-%d-%H%M%S).json"
        
        if command -v mako >/dev/null 2>&1; then
            mako export > "$BACKUP_FILE" 2>/dev/null || true
            if [ -f "$BACKUP_FILE" ]; then
                print_success "Exported to $BACKUP_FILE"
            else
                print_error "Failed to export history"
            fi
        else
            print_error "Cannot export: mako command not found"
        fi
    fi
}

# Remove binaries
remove_binaries() {
    echo ""
    print_info "Removing binaries..."
    
    local removed=0
    
    for location in "${INSTALL_LOCATIONS[@]}"; do
        if [ -f "$location" ]; then
            if rm "$location" 2>/dev/null; then
                print_success "Removed $location"
                ((removed++))
            elif command -v sudo >/dev/null 2>&1; then
                if sudo rm "$location" 2>/dev/null; then
                    print_success "Removed $location (with sudo)"
                    ((removed++))
                fi
            fi
        fi
    done
    
    if [ $removed -eq 0 ]; then
        print_warning "No binaries found to remove"
    fi
}

# Remove configuration
remove_config() {
    print_info "Removing configuration..."
    
    if [ -d "$MAKO_DIR" ]; then
        rm -rf "$MAKO_DIR"
        print_success "Removed $MAKO_DIR"
    else
        print_warning "Configuration directory not found"
    fi
}

# Remove shell completions
remove_completions() {
    print_info "Removing shell completions..."
    
    local removed=0
    
    for location in "${COMPLETION_LOCATIONS[@]}"; do
        if [ -f "$location" ]; then
            if rm "$location" 2>/dev/null; then
                print_success "Removed $location"
                ((removed++))
            elif command -v sudo >/dev/null 2>&1; then
                if sudo rm "$location" 2>/dev/null; then
                    print_success "Removed $location (with sudo)"
                    ((removed++))
                fi
            fi
        fi
    done
    
    if [ $removed -eq 0 ]; then
        print_warning "No completion files found"
    fi
}

# Show what will be removed
show_removal_plan() {
    echo ""
    echo -e "${DIM_BLUE}This will remove:${RESET}"
    echo -e "  • Mako binaries (mako, mako-menu)"
    echo -e "  • Configuration directory (~/.mako/)"
    echo -e "  • Shell completion files"
    echo ""
}

# Confirm uninstallation
confirm_uninstall() {
    echo ""
    read -p "$(echo -e ${CYAN}Proceed with uninstall? \[y/N\]: ${RESET})" -n 1 -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Uninstall cancelled"
        exit 0
    fi
}

# Show goodbye message
show_goodbye() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo -e "${LIGHT_BLUE}  Mako has been uninstalled${RESET}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo ""
    echo -e "${DIM_BLUE}We're sad to see you go!${RESET}"
    echo ""
    echo -e "${DIM_BLUE}Feedback:${RESET} ${CYAN}https://github.com/fabiobrug/mako/issues${RESET}"
    echo ""
}

# Main uninstallation flow
main() {
    print_header
    
    # Check if Mako is installed
    if ! command -v mako >/dev/null 2>&1 && [ ! -d "$MAKO_DIR" ]; then
        print_warning "Mako doesn't appear to be installed"
        exit 0
    fi
    
    show_removal_plan
    
    # Count and offer to export
    COMMAND_COUNT=$(count_commands)
    offer_export "$COMMAND_COUNT"
    
    confirm_uninstall
    
    echo ""
    remove_binaries
    remove_completions
    remove_config
    
    show_goodbye
}

main "$@"
