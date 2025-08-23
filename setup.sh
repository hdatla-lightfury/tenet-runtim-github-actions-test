#!/bin/bash

# Nakama Go Backend Setup Script
# Installs all dependencies for Nakama runtime plugin

set -e

echo "🚀 Setting up Nakama Go Backend dependencies..."

# Create temporary directory for downloads
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Install Node.js with NVM
echo "📦 Installing Node.js with NVM..."
if ! command -v nvm &> /dev/null; then
    cd $TEMP_DIR
    curl -sL https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.0/install.sh -o install_nvm.sh
    bash install_nvm.sh
    cd - > /dev/null
    source ~/.bashrc
fi

# Install Node.js LTS
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm install --lts
nvm use --lts

# Install PM2
echo "⚙️  Installing PM2..."
npm install -g pm2

# Install Go 1.24.5
echo "🐹 Installing Go 1.24.5..."
if ! command -v go &> /dev/null || [[ $(go version) != *"go1.24.5"* ]]; then
    cd $TEMP_DIR
    wget -q https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz
    cd - > /dev/null
    
    # Add Go to PATH if not already there
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        export PATH=$PATH:/usr/local/go/bin
    fi
fi


# Download Nakama binary
echo "⚡ Downloading Nakama v3.28.0..."
if [ ! -f "./nakama" ]; then
    cd $TEMP_DIR
    wget -q https://github.com/heroiclabs/nakama/releases/download/v3.28.0/nakama-3.28.0-linux-amd64.tar.gz
    tar -xzf nakama-3.28.0-linux-amd64.tar.gz
    cd - > /dev/null
    cp $TEMP_DIR/nakama ./nakama
    chmod +x ./nakama
fi

# Run database migrations
echo "🔄 Running database migrations..."
./nakama migrate up --database.address postgres:localdb@localhost:5432/nakama

# Create modules directory
mkdir -p modules

./scripts/restart-nakama.sh

echo "✅ Setup complete!"
echo "📋 Next steps:"
echo "   1. Run: ./scripts/restart-nakama.sh"
echo "   2. Access console: http://localhost:7351"
echo "   3. View logs: pm2 logs nakama" 