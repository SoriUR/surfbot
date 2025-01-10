#!/bin/bash

# Чтение переменных из .env
source "$(dirname "${BASH_SOURCE[0]}")/.env"

echo "Killing old session"
screen -X -S "$SESSION_NAME" quit