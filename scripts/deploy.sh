#!/bin/bash

echo "Killing old session"
screen -X -S "$SESSION_NAME" quit

echo "Start new session"
screen -dmS "$SESSION_NAME" bash -c "/var/www/surbot/scripts/run.sh"
