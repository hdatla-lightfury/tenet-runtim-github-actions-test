#!/bin/bash

# Script to rebuild backend.so and restart Nakama using PM2
# Usage: ./scripts/restart-nakama.sh

set -e

# Configuration
NAKAMA_BINARY="./nakama"  # Path to downloaded nakama binary
NAKAMA_CONFIG="./local.yml"
DATABASE_ADDRESS="postgres:localdb@localhost:5432/nakama"
PM2_APP_NAME="nakama"
MODULES_DIR="./modules"

echo "🔨 Building backend.so..."

# Sync go.work dependencies and update vendor
go mod vendor

# Build the backend plugin using vendored dependencies
go build --trimpath --mod=vendor  --buildmode=plugin -o ./backend.so

echo "✅ Backend built successfully!"

# Create modules directory and copy backend.so
echo "📁 Setting up modules directory..."
mkdir -p $MODULES_DIR
cp ./backend.so $MODULES_DIR/

# Stop Nakama using PM2
echo "⏹️  Stopping Nakama via PM2..."
pm2 stop $PM2_APP_NAME || echo "Nakama not running in PM2"

# Wait a moment
sleep 2

# Start/Restart Nakama with PM2
echo "🚀 Starting Nakama with PM2..."
if [ -f "$NAKAMA_CONFIG" ]; then
    # Start with config file
    pm2 start $NAKAMA_BINARY --name $PM2_APP_NAME -- --config $NAKAMA_CONFIG --database.address $DATABASE_ADDRESS --runtime.path $MODULES_DIR
else
    # Start with minimal config
    pm2 start $NAKAMA_BINARY --name $PM2_APP_NAME -- --database.address $DATABASE_ADDRESS --runtime.path $MODULES_DIR
fi

# Show PM2 status
echo "📋 PM2 Status:"
pm2 list

echo "✅ Nakama restarted successfully with PM2!"
echo "🌐 Console available at: http://localhost:7351"
echo "📋 View logs: pm2 logs $PM2_APP_NAME" 