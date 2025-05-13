#!/bin/bash
set -e

# Load .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

DB_FILE="${DB_PATH:-./data/dev.db}"

goose -dir sql/schema sqlite3 "$DB_FILE" "$@"
