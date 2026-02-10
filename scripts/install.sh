#!/bin/bash
set -e

# Mako Installation Script
# Usage: curl -sSL https://get-mako.sh | bash

# Color codes
CYAN='\033[38;2;0;209;255m'
LIGHT_BLUE='\033[38;2;93;173;226m'
DIM_BLUE='\033[38;2;120;150;180m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
RESET='\033[0m'

# Configuration
REPO="fabiobrug/mako"
INSTALL_DIR="/usr/local/bin"
MAKO_DIR="$HOME/.mako"
VERSION="${MAKO_VERSION:-latest}"

# Print functions
print_header() {
    echo -e "\n${CYAN}  Mako - AI-Native Shell Orchestrator${RESET}\n"
}

print_step() {
    echo -e "${LIGHT_BLUE}▸${RESET} $1"
}

print_success() {
    echo -e "${GREEN}✓${RESET} $1"
}

print_error() {
    echo -e "${RED}✗${RESET} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}⚠${RESET} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64 | arm64)
            ARCH="arm64"
            ;;
        armv7l | armv6l)
            ARCH="arm"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    case $OS in
        linux | darwin)
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    print_success "Detected: $OS $ARCH"
}

# Get latest release version
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            print_error "Failed to get latest version"
            exit 1
        fi
    fi
    print_success "Version: $VERSION"
}

# Download binaries
download_binaries() {
    print_step "Downloading Mako $VERSION..."
    
    BINARY_NAME="mako-${OS}-${ARCH}"
    MENU_BINARY_NAME="mako-menu-${OS}-${ARCH}"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"
    MENU_DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$MENU_BINARY_NAME"
    
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    # Download main binary
    if ! curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/mako"; then
        print_error "Failed to download mako from $DOWNLOAD_URL"
        print_error "This might be because pre-built binaries are not available yet."
        print_error "Please build from source: https://github.com/$REPO"
        exit 1
    fi
    
    # Download menu binary
    if ! curl -sL "$MENU_DOWNLOAD_URL" -o "$TMP_DIR/mako-menu"; then
        print_warning "Failed to download mako-menu, trying to continue..."
    fi
    
    chmod +x "$TMP_DIR/mako"
    if [ -f "$TMP_DIR/mako-menu" ]; then
        chmod +x "$TMP_DIR/mako-menu"
    fi
    
    print_success "Downloaded binaries"
}

# Install binaries
install_binaries() {
    print_step "Installing to $INSTALL_DIR..."
    
    # Try to install to /usr/local/bin first
    if [ -w "$INSTALL_DIR" ]; then
        cp "$TMP_DIR/mako" "$INSTALL_DIR/mako"
        if [ -f "$TMP_DIR/mako-menu" ]; then
            cp "$TMP_DIR/mako-menu" "$INSTALL_DIR/mako-menu"
        fi
    elif command -v sudo >/dev/null 2>&1; then
        sudo cp "$TMP_DIR/mako" "$INSTALL_DIR/mako"
        if [ -f "$TMP_DIR/mako-menu" ]; then
            sudo cp "$TMP_DIR/mako-menu" "$INSTALL_DIR/mako-menu"
        fi
    else
        # Fallback to ~/.local/bin
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        cp "$TMP_DIR/mako" "$INSTALL_DIR/mako"
        if [ -f "$TMP_DIR/mako-menu" ]; then
            cp "$TMP_DIR/mako-menu" "$INSTALL_DIR/mako-menu"
        fi
        
        # Add to PATH if not already there
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            print_warning "Add $INSTALL_DIR to your PATH:"
            echo -e "  ${CYAN}export PATH=\"\$PATH:$INSTALL_DIR\"${RESET}"
            echo ""
            
            # Try to add to shell config
            if [ -f "$HOME/.bashrc" ]; then
                echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
                print_success "Added to ~/.bashrc"
            elif [ -f "$HOME/.zshrc" ]; then
                echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc"
                print_success "Added to ~/.zshrc"
            fi
        fi
    fi
    
    print_success "Installed binaries to $INSTALL_DIR"
}

# Setup Mako directory
setup_mako_dir() {
    print_step "Setting up $MAKO_DIR..."
    
    mkdir -p "$MAKO_DIR"
    
    print_success "Created Mako directory"
}

# Setup API key
setup_api_key() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo -e "${LIGHT_BLUE}  Gemini API Key Setup${RESET}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo ""
    echo -e "${DIM_BLUE}Mako uses Google's Gemini API (free tier available)${RESET}"
    echo -e "${DIM_BLUE}Get your key: ${CYAN}https://ai.google.dev/${RESET}"
    echo ""
    
    read -p "$(echo -e ${CYAN}Enter API key \(or press Enter to skip\): ${RESET})" API_KEY
    echo ""
    
    if [ -n "$API_KEY" ]; then
        echo "{\"api_key\":\"$API_KEY\"}" > "$MAKO_DIR/config.json"
        print_success "API key saved!"
    else
        print_warning "You can set it later with: ${CYAN}mako config set api_key YOUR_KEY${RESET}"
    fi
}

# Verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    if command -v mako >/dev/null 2>&1; then
        VERSION_OUTPUT=$(mako version 2>&1 || true)
        print_success "Mako installed successfully!"
    else
        print_error "Installation failed: mako command not found"
        exit 1
    fi
}

# Show completion message
show_completion() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo -e "${GREEN}  Mako installed successfully!${RESET}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
    echo ""
    echo -e "${LIGHT_BLUE}Quick start:${RESET}"
    echo -e "  ${CYAN}mako${RESET}                    Start Mako shell"
    echo -e "  ${CYAN}mako ask \"list files\"${RESET}   Generate commands"
    echo -e "  ${CYAN}mako help${RESET}               See all commands"
    echo ""
    echo -e "${LIGHT_BLUE}Documentation:${RESET} ${CYAN}https://github.com/$REPO${RESET}"
    echo ""
    
    # Shell completion hint
    echo -e "${DIM_BLUE}Optional: Enable shell completion:${RESET}"
    echo -e "  ${CYAN}mako completion bash | sudo tee /etc/bash_completion.d/mako${RESET}"
    echo -e "  ${CYAN}mako completion zsh > ~/.zsh/completions/_mako${RESET}"
    echo -e "  ${CYAN}mako completion fish > ~/.config/fish/completions/mako.fish${RESET}"
    echo ""
}

# Main installation flow
main() {
    print_header
    
    print_step "Installing Mako..."
    echo ""
    
    detect_platform
    get_latest_version
    download_binaries
    install_binaries
    setup_mako_dir
    setup_api_key
    verify_installation
    show_completion
}

main "$@"
