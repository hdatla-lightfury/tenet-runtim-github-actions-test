#!/bin/bash

set -e

echo "Setting up Nakama Go Backend dependencies..."

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

sudo mkdir -p /home/ec2-user/tenet-runtime
sudo chown -R ec2-user:ec2-user /home/ec2-user/tenet-runtime


# Install Node.js with NVM
if ! command -v nvm &> /dev/null; then
    echo "Installing Node.js with NVM..."
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
echo "Installing PM2..."
npm install -g pm2

# Install Go 1.24.5
GO_VERSION=1.24.5
echo "Installing Go"
if ! command -v go &> /dev/null || [[ $(go version) != *"go${GO_VERSION}"* ]]; then
    echo " [debug] Installing Go ${GO_VERSION}..."
    cd $TEMP_DIR
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    cd - > /dev/null
    
    # Add Go to PATH if not already there
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        export PATH=$PATH:/usr/local/go/bin
    fi
fi

# Download and install Nakama binary
NAKAMA_VERSION=3.28.0
if [ ! -f "./nakama" ]; then
    echo " [debug] Installing Nakama v${NAKAMA_VERSION}..."
    cd $TEMP_DIR
    wget -q https://github.com/heroiclabs/nakama/releases/download/v${NAKAMA_VERSION}/nakama-${NAKAMA_VERSION}-linux-amd64.tar.gz  -O nakama.tar.gz
    tar -xzf nakama.tar.gz
    cd - > /dev/null
    cp $TEMP_DIR/nakama /home/ec2-user/tenet-runtime/nakama
    chmod +x /home/ec2-user/tenet-runtime/nakama
fi

# Install AWS CLI v2 for secrets manager

if ! command -v aws &> /dev/null; then
    echo "[debug] Installing AWS CLI v2..."
    cd $TEMP_DIR
    wget "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip -q awscliv2.zip
    sudo ./aws/install
    cd - > /dev/null
fi

# Install jq for command line JSON parsing(for AWS secrets manager)
if ! command -v jq &> /dev/null; then
     echo "[debug] Installing jq 1.8.1 from source..."
    cd $TEMP_DIR

    # Download and extract jq source tarball
    wget -q https://github.com/jqlang/jq/releases/download/jq-1.8.1/jq-1.8.1.tar.gz
    tar -xzf jq-1.8.1.tar.gz
    cd jq-1.8.1

    # Build and install
    ./configure
    make
    sudo make install

    cd - > /dev/null
fi

# Ensure yq exists (mikefarah/yq). Install if missing. To edit local.yml file
if ! command -v yq >/dev/null 2>&1; then
  echo "[debug] Installing yq (mikefarah) from source"
  VERSION="v4.47.1"   # pin to a specific version
  BINARY="yq_linux_amd64"

  cd $TMP_DIR

  wget -q "https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY}.tar.gz" -O - \
    | tar xz

  sudo mv ${BINARY} /usr/local/bin/yq
  cd - >/dev/null
fi
