#!/bin/bash
# run-supervisor.sh

echo "Starting Game Factory Supervisor Wrapper..."

# Ensure we use the local Go and TinyGo installations
export PATH="/home/luke/.local/go/bin:/home/luke/.local/tinygo/bin:$PATH"

while true; do
  node supervisor.js
  EXIT_CODE=$?
  if [ $EXIT_CODE -eq 0 ]; then
    echo "Supervisor exited cleanly. Stopping."
    exit 0
  elif [ $EXIT_CODE -eq 3 ]; then
    echo "Supervisor requested a reboot (self-update detected). Reloading..."
    sleep 1
  else
    echo "Supervisor crashed with exit code $EXIT_CODE. Restarting in 10 seconds..."
    sleep 10
  fi
done
