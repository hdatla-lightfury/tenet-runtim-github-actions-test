# Nakama Go Backend Deployment Guide

This guide covers the complete setup process for running Nakama with a Go backend plugin using PM2 process manager on Linux.

## System Requirements

- **OS**: Linux (64-bit) - tested on Ubuntu/WSL2
- **Architecture**: x86_64 (AMD64)
- **Shell**: bash

## Version Information

| Component | Version | Notes |
|-----------|---------|-------|
| Nakama | v3.28.0 | Latest stable release |
| Go | 1.24.5 | Required for building plugins |
| Node.js | LTS (via nvm) | For PM2 process manager |
| PostgreSQL | 12.2+ | Database backend |
| nakama-common | v1.38.0 | Go runtime library |
| protobuf | v1.35.1 | Version-locked for compatibility |

## Installation Steps

### 1. Install Node.js with NVM

Install Node Version Manager and Node.js:

```bash
# Download and install nvm
curl -sL https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.0/install.sh -o install_nvm.sh
bash install_nvm.sh

# Restart terminal or source profile
source ~/.bashrc  # or ~/.zshrc

# Install latest LTS Node.js
nvm install --lts
nvm use --lts

# Verify installation
node --version
npm --version
```

### 2. Install PM2 Process Manager

```bash
# Install PM2 globally
npm install -g pm2

# Verify installation
pm2 --version
```


### 4. Download Nakama Binary

```bash
# Download Nakama v3.28.0
wget https://github.com/heroiclabs/nakama/releases/download/v3.28.0/nakama-3.28.0-linux-amd64.tar.gz

# Extract binary
tar -xzf nakama-3.28.0-linux-amd64.tar.gz

# Make executable
chmod +x nakama

# Verify installation
./nakama --version
```

### 5. Setup Go Dependencies

Ensure your project has the correct dependency versions:

**go.mod**:
```go
module github.com/titan/titan-runtime

go 1.24.5

require github.com/heroiclabs/nakama-common v1.38.0

require google.golang.org/protobuf v1.36.6 // indirect
```

### 6. Run Database Migrations

```bash
# Run Nakama database migrations
./nakama migrate up --database.address postgres:nakama:localdb@localhost:5432/nakama
```

## Project Structure

```
tenet-runtime/
├── nakama                    # Nakama binary (v3.28.0)
├── backend.so               # Built Go plugin
├── modules/                 # Nakama modules directory
│   └── backend.so          # Plugin copy for runtime
├── scripts/
│   └── restart-nakama.sh   # Deployment script
├── go.mod                  # Go dependencies
├── go.work                 # Version compatibility
├── local.yml               # Nakama configuration
├── main.go                 # Plugin entry point
└── modules/                # Go source modules
    ├── account/
    ├── common/
    └── utils/
```

## Deployment Script

The `scripts/restart-nakama.sh` script handles the complete build and deployment process:

```bash
#!/bin/bash

# Script to rebuild backend.so and restart Nakama using PM2
# Usage: ./scripts/restart-nakama.sh

set -e

# Configuration
NAKAMA_BINARY="./nakama"  # Path to downloaded nakama binary
NAKAMA_CONFIG="./local.yml"
DATABASE_ADDRESS="postgres:nakama:localdb@localhost:5432/nakama"
PM2_APP_NAME="nakama"
MODULES_DIR="./modules"

echo "🔨 Building backend.so..."

# Sync go.work dependencies and update vendor
go work sync
go mod vendor

# Build the backend plugin using vendored dependencies
go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

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
```

## Usage

### Initial Setup

```bash
# Make script executable
chmod +x scripts/restart-nakama.sh

# Run the deployment script
./scripts/restart-nakama.sh
```

### Development Workflow

```bash
# After making code changes, rebuild and restart:
./scripts/restart-nakama.sh
```

### PM2 Management Commands

```bash
# View running processes
pm2 list

# View logs
pm2 logs nakama

# Restart service
pm2 restart nakama

# Stop service
pm2 stop nakama

# Remove from PM2
pm2 delete nakama

# Save PM2 configuration
pm2 save

# Setup PM2 to start on boot
pm2 startup
```

## Verification

### Check Services

```bash
# Verify Nakama is running
pm2 list

# Check Nakama logs
pm2 logs nakama

# Test database connection
./nakama --database.address postgres:nakama:localdb@localhost:5432/nakama --check
```

### Access Points

- **Nakama Console**: http://localhost:7351 (admin/password)
- **API Endpoint**: http://localhost:7350
- **Socket Endpoint**: ws://localhost:7350/ws

## Troubleshooting

### Common Issues

**1. Protobuf Version Mismatch**
```
Error: plugin was built with a different version of package google.golang.org/protobuf/types/known/timestamppb
```
*Solution*: Ensure `go.work` file exists with protobuf version replacement.

**2. Runtime Path Error**
```
Error: mkdir ./backend.so: not a directory
```
*Solution*: Use `--runtime.path ./modules` (directory) not `--runtime.path ./backend.so` (file).

**3. Database Connection Issues**
```
Error: failed to connect to database
```
*Solution*: Verify PostgreSQL is running and credentials are correct.

### Debug Commands

```bash
# Check Go environment
go version
go env

# Verify dependencies
go mod verify
go work sync

# Test build without PM2
go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

# Test Nakama directly
./nakama --database.address postgres:nakama:localdb@localhost:5432/nakama --runtime.path ./modules
```

## Production Considerations

### Security
- Change default Nakama console credentials
- Use environment variables for database passwords
- Configure proper firewall rules
- Enable SSL/TLS in production

### Performance
- Use PM2 cluster mode for multiple instances
- Configure database connection pooling
- Monitor resource usage with PM2 monitoring
- Set up log rotation

### Monitoring
```bash
# PM2 monitoring
pm2 monit

# System resource monitoring
pm2 install pm2-server-monit
```

## Backup Strategy

```bash
# Database backup
pg_dump -U nakama -h localhost nakama > nakama_backup.sql

# Configuration backup
cp local.yml local.yml.backup
cp -r modules/ modules.backup/
```

This guide provides a complete, reproducible setup for running Nakama with Go backend plugins using PM2 process management. 