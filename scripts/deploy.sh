#!/bin/bash

SESSION_NAME = "surfbot"

echo "Killing old session"
screen -X -S "$SESSION_NAME" quit

echo "Start new session"
screen -dmS "$SESSION_NAME" bash -c "/var/www/surfbot/scripts/start.sh"
