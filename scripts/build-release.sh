#!/bin/bash
set -e

# Mako Release Builder
# Builds binaries for all supported platforms

VERSION=${1:-"dev"}
OUTPUT_DIR="dist"

# Color codes
CYAN='\033[38;2;0;209;255m'
LIGHT_BLUE='\033[38;2;93;173;226m'
GREEN='\033[0;32m'
RED='\033[0;31m'
RESET='\033[0m'

print_step() {
    echo -e "${LIGHT_BLUE}▸${RESET} $1"
}

print_success() {
    echo -e "${GREEN}✓${RESET} $1"
}

print_error() {
    echo -e "${RED}✗${RESET} $1" >&2
}

# Clean and create output directory
print_step "Cleaning output directory..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Build configurations
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

echo ""
print_step "Building Mako $VERSION for all platforms..."
echo ""

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    output_name_main="mako-${GOOS}-${GOARCH}"
    output_name_menu="mako-menu-${GOOS}-${GOARCH}"
    
    print_step "Building $GOOS/$GOARCH..."
    
    # Build main binary
    if GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT_DIR/$output_name_main" ./cmd/mako; then
        print_success "Built $output_name_main ($(du -h "$OUTPUT_DIR/$output_name_main" | cut -f1))"
    else
        print_error "Failed to build $output_name_main"
        exit 1
    fi
    
    # Build menu binary
    if GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT_DIR/$output_name_menu" ./cmd/mako-menu; then
        print_success "Built $output_name_menu ($(du -h "$OUTPUT_DIR/$output_name_menu" | cut -f1))"
    else
        print_error "Failed to build $output_name_menu"
        exit 1
    fi
    
    echo ""
done

print_success "All binaries built successfully!"
echo ""
print_step "Binaries available in $OUTPUT_DIR/"
ls -lh "$OUTPUT_DIR"
echo ""

# Create checksums
print_step "Generating checksums..."
cd "$OUTPUT_DIR"
sha256sum * > SHA256SUMS
cd ..
print_success "Checksums generated: $OUTPUT_DIR/SHA256SUMS"

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${GREEN}  Release build complete!${RESET}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${LIGHT_BLUE}Next steps:${RESET}"
echo -e "  1. Create a git tag: ${CYAN}git tag -a v$VERSION -m \"Release v$VERSION\"${RESET}"
echo -e "  2. Push the tag: ${CYAN}git push origin v$VERSION${RESET}"
echo -e "  3. Upload binaries to GitHub release"
echo ""
echo -e "${LIGHT_BLUE}Using GitHub CLI:${RESET}"
echo -e "  ${CYAN}gh release create v$VERSION dist/* --title \"v$VERSION\" --notes \"Release notes here\"${RESET}"
echo ""
