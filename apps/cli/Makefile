.PHONY: build clean install test test-coverage test-coverage-html release quick

# Build both binaries
build:
	@echo "Building Mako..."
	@go build -o mako ./cmd/mako
	@go build -o mako-menu ./cmd/mako-menu
	@echo "Build complete: mako, mako-menu"

# Clean build artifacts
clean:
	@rm -f mako mako-menu
	@rm -rf dist coverage.out coverage.html
	@echo "Cleaned build artifacts"

# Install to /usr/local/bin (requires sudo)
install: build
	@echo "Installing Mako..."
	@sudo cp mako /usr/local/bin/
	@sudo cp mako-menu /usr/local/bin/
	@echo "Installed to /usr/local/bin/"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -tags "fts5" ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -tags "fts5" -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo ""
	@echo "Coverage report saved to coverage.out"
	@echo "Run 'make test-coverage-html' to view HTML report"

# Generate HTML coverage report
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"
	@echo "Open coverage.html in your browser to view"

# Run tests for a specific package
test-pkg:
	@go test -v -tags "fts5" ./$(PKG)/...

# Quick rebuild (just mako)
quick:
	@go build -o mako ./cmd/mako
	@echo "Mako rebuilt"

# Build release binaries for all platforms
release:
	@./scripts/build-release.sh $(VERSION)
