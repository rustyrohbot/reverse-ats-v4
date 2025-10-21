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

# Clean generated files and binaries
clean:
	rm -rf bin/
	rm -rf internal/db/
	rm -f *_templ.go
	rm -f internal/templates/*_templ.go
	rm -f static/output.css

# Install development tools
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/a-h/templ/cmd/templ@latest
