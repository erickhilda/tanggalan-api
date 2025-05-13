#!/bin/bash
set -e

# Load environment variables from .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Use DB_PATH from .env, fallback default
DB_FILE="${DB_PATH:-./data/dev.db}"

echo "Creating DB if not exist: $DB_FILE"
mkdir -p ./data
touch $DB_FILE

echo "Running goose migration..."
goose -dir sql/schema sqlite3 $DB_FILE up
