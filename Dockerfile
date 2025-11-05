# syntax=docker/dockerfile:1

### STAGE 1: Build application and goose ###
FROM golang:1.23.2 AS builder

WORKDIR /app

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build app binary
# RUN go build -o server ./cmd/api/main.go
# Build dengan flags untuk static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o server ./cmd/api/main.go

### STAGE 2: Migration-only image (optional but useful) ###
FROM golang:1.23.2 AS migrator

WORKDIR /app

# Install goose again
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy all needed files
COPY . .

# Mark scripts executable
RUN chmod +x ./scripts/*.sh

ENTRYPOINT ["/bin/bash", "./scripts/init_db.sh"]

### STAGE 3: Final runtime image ###
FROM debian:bullseye AS runtime

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/server /app/server

# Copy env and start script
COPY --from=builder /app/start.sh /app/start.sh
COPY --from=builder /app/.env /app/.env

# App might need access to SQL files at runtime (e.g. embedded mode)
COPY --from=builder /app/sql /app/sql

RUN chmod +x /app/start.sh

EXPOSE 3000

ENTRYPOINT ["/app/start.sh"]

