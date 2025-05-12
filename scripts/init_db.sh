#!/bin/sh
set -e

echo "🧱 Initializing SQLite DB..."
go run ./cmd/initdb
