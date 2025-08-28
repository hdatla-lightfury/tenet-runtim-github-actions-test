#!/bin/bash

set -e

echo "Setting up Nakama Go Backend dependencies..."

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

sudo mkdir -p /home/ec2-user/tenet-runtime
sudo chown -R ec2-user:ec2-user /home/ec2-user/tenet-runtime



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

mkdir -p /home/ec2-user/tenet-runtime/modules


# --- systemd unit for Nakama ---
echo "[debug] Writing systemd unit for Nakama"
sudo tee /etc/systemd/system/nakama.service >/dev/null <<'UNIT'
[Unit]
Description=Nakama Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=ec2-user
Group=ec2-user
WorkingDirectory=/home/ec2-user/tenet-runtime
# Adjust flags as you need; --config should point to your local.yml
ExecStart=/home/ec2-user/tenet-runtime/nakama --config /home/ec2-user/tenet-runtime/local.yml
Restart=on-failure
RestartSec=5
# Increase if your startup needs more time
StartLimitBurst=3
StartLimitIntervalSec=60

[Install]
WantedBy=multi-user.target
UNIT

# Reload systemd and enable service (autostart on boot)
sudo systemctl daemon-reload
sudo systemctl enable nakama
echo "[debug] nakama.service installed and enabled (not started yet)"