.PHONY: build clean install test

# Build both binaries
build:
	@echo "Building Mako..."
	@go build -o mako ./cmd/mako
	@go build -o mako-menu ./cmd/mako-menu
	@echo "Build complete: mako, mako-menu"

# Clean build artifacts
clean:
	@rm -f mako mako-menu
	@echo "Cleaned build artifacts"

# Install to /usr/local/bin (requires sudo)
install: build
	@echo "Installing Mako..."
	@sudo cp mako /usr/local/bin/
	@sudo cp mako-menu /usr/local/bin/
	@echo "Installed to /usr/local/bin/"

# Run tests
test:
	@go test -v ./...

# Quick rebuild (just mako)
quick:
	@go build -o mako ./cmd/mako
	@echo "Mako rebuilt"
