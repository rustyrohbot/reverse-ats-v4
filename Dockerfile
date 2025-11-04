# Multi-stage build for smaller final image
# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

# Install build dependencies (gcc and musl-dev needed for CGO/SQLite)
RUN apk add --no-cache git make npm gcc musl-dev

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install templ for template generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Add Go bin to PATH and generate templates
RUN export PATH="$(go env GOPATH)/bin:$PATH" && templ generate

# Build Tailwind CSS
RUN npm install
RUN npx tailwindcss -i ./static/input.css -o ./static/output.css --minify

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server cmd/server/main.go

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies (SQLite needs glibc)
RUN apk add --no-cache ca-certificates libc6-compat

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/server /app/server

# Copy static files and migrations
COPY --from=builder /app/static /app/static
COPY --from=builder /app/pb_migrations /app/pb_migrations

# Create directory for PocketBase data
RUN mkdir -p /app/pb_data

# Expose the port
EXPOSE 5627

# Set environment variables
ENV REVERSE_ATS_PORT=5627

# Run the application
CMD ["/app/server"]
