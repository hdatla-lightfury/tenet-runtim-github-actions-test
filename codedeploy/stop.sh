#!/bin/bash

set -e

PM2_APP_NAME="nakama"

# Check if pm2 is installed
if command -v pm2 &> /dev/null; then
  pm2 stop "$PM2_APP_NAME" || echo "Nakama not running in PM2"

  echo "PM2 Status:"
  pm2 list || true
else
  echo "PM2 not found. Skipping stop step."
fi

echo "Stop completed."
