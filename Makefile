.PHONY: generate run build clean install-tools css css-watch build-all build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64 build-darwin-amd64 build-darwin-arm64 build-linux-amd64-only build-linux-arm64-only build-windows-amd64-only build-windows-arm64-only build-darwin-amd64-only build-darwin-arm64-only

# Generate templ code
generate:
	templ generate

# Build CSS with Tailwind
css:
	npm run build:css

# Watch CSS for changes
css-watch:
	npm run watch:css

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build: generate css
	go build -o bin/server cmd/server/main.go

# Cross-platform builds
build-linux-amd64: generate css
	GOOS=linux GOARCH=amd64 go build -o bin/reverse-ats-linux-amd64 cmd/server/main.go

build-linux-arm64: generate css
	GOOS=linux GOARCH=arm64 go build -o bin/reverse-ats-linux-arm64 cmd/server/main.go

build-windows-amd64: generate css
	GOOS=windows GOARCH=amd64 go build -o bin/reverse-ats-windows-amd64.exe cmd/server/main.go

build-windows-arm64: generate css
	GOOS=windows GOARCH=arm64 go build -o bin/reverse-ats-windows-arm64.exe cmd/server/main.go

build-darwin-amd64: generate css
	GOOS=darwin GOARCH=amd64 go build -o bin/reverse-ats-darwin-amd64 cmd/server/main.go

build-darwin-arm64: generate css
	GOOS=darwin GOARCH=arm64 go build -o bin/reverse-ats-darwin-arm64 cmd/server/main.go

# Build all platforms (generate and css run once before all builds)
build-all: generate css
	$(MAKE) build-linux-amd64-only build-linux-arm64-only build-windows-amd64-only build-windows-arm64-only build-darwin-amd64-only build-darwin-arm64-only

build-linux-amd64-only:
	GOOS=linux GOARCH=amd64 go build -o bin/reverse-ats-linux-amd64 cmd/server/main.go

build-linux-arm64-only:
	GOOS=linux GOARCH=arm64 go build -o bin/reverse-ats-linux-arm64 cmd/server/main.go

build-windows-amd64-only:
	GOOS=windows GOARCH=amd64 go build -o bin/reverse-ats-windows-amd64.exe cmd/server/main.go

build-windows-arm64-only:
	GOOS=windows GOARCH=arm64 go build -o bin/reverse-ats-windows-arm64.exe cmd/server/main.go

build-darwin-amd64-only:
	GOOS=darwin GOARCH=amd64 go build -o bin/reverse-ats-darwin-amd64 cmd/server/main.go

build-darwin-arm64-only:
	GOOS=darwin GOARCH=arm64 go build -o bin/reverse-ats-darwin-arm64 cmd/server/main.go

# Clean generated files and binaries
clean:
	rm -rf bin/
	rm -rf pb_data/
	rm -f *_templ.go
	rm -f internal/templates/*_templ.go
	rm -f static/output.css

# Show available build targets
help:
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-linux-amd64  - Build for Linux x86_64"
	@echo "  build-linux-arm64  - Build for Linux ARM64"
	@echo "  build-windows-amd64 - Build for Windows x86_64"
	@echo "  build-windows-arm64 - Build for Windows ARM64"
	@echo "  build-darwin-amd64 - Build for macOS x86_64"
	@echo "  build-darwin-arm64 - Build for macOS ARM64 (Apple Silicon)"
	@echo "  build-all          - Build for all platforms"
	@echo "  run                - Run the development server"
	@echo "  generate           - Generate templ code"
	@echo "  css                - Build CSS with Tailwind"
	@echo "  css-watch          - Watch and rebuild CSS"
	@echo "  clean              - Clean generated files and binaries"

# Install development tools
install-tools:
	go install github.com/a-h/templ/cmd/templ@latest
	npm install