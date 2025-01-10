#!/bin/bash

# Чтение переменных из .env
source "$(dirname "${BASH_SOURCE[0]}")/.env"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Start new session"
screen -dmS "$SESSION_NAME" bash -c "$SCRIPT_DIR/start.sh"