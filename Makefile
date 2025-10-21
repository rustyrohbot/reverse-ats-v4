.PHONY: generate migrate-up migrate-down run build clean install-tools css css-watch

# Generate sqlc and templ code
generate:
	sqlc generate
	templ generate

# Build CSS with Tailwind
css:
	npm run build:css

# Watch CSS for changes
css-watch:
	npm run watch:css

# Database migrations
migrate-up:
	goose -dir migrations sqlite3 ./data.db up

migrate-down:
	goose -dir migrations sqlite3 ./data.db down

migrate-status:
	goose -dir migrations sqlite3 ./data.db status

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Cross-platform builds
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/reverse-ats-linux-amd64 cmd/server/main.go

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/reverse-ats-linux-arm64 cmd/server/main.go

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/reverse-ats-windows-amd64.exe cmd/server/main.go

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o bin/reverse-ats-windows-arm64.exe cmd/server/main.go

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/reverse-ats-darwin-amd64 cmd/server/main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/reverse-ats-darwin-arm64 cmd/server/main.go

# Build all platforms
build-all: build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64 build-darwin-amd64 build-darwin-arm64

# Clean generated files and binaries
clean:
	rm -rf bin/
	rm -rf internal/db/
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
	@echo "  generate           - Generate sqlc and templ code"
	@echo "  css                - Build CSS with Tailwind"
	@echo "  migrate-up         - Run database migrations"
	@echo "  clean              - Clean generated files and binaries"

# Install development tools
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/a-h/templ/cmd/templ@latest
	npm install