#!/bin/bash

set -e

PM2_APP_NAME="nakama"

echo "Stopping Nakama via PM2..."
pm2 stop "$PM2_APP_NAME" || echo "Nakama not running in PM2"

echo "PM2 Status:"
pm2 list || true

echo "Stop completed."
