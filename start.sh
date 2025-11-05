#!/bin/bash
set -e

# Load env vars
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

echo "Starting tanggalan-api server..."
exec ./server
