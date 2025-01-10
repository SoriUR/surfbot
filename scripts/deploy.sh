#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Stopping session..."
bash "$SCRIPT_DIR/stop_session.sh"

echo "Starting new session..."
bash "$SCRIPT_DIR/start_session.sh"

echo "Deployment complete."