#!/bin/bash

set -e

# --- Load constants ---
# Expect constants.sh in the same directory as this script, or set $CONSTANTS_FILE externally.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONSTANTS_FILE="${CONSTANTS_FILE:-${SCRIPT_DIR}/constants.cnf}"

if [ ! -f "${CONSTANTS_FILE}" ]; then
  echo "[error] constants file not found at ${CONSTANTS_FILE}"
  exit 1
fi

source "${CONSTANTS_FILE}"

echo "Setting up Nakama Go Backend dependencies..."

# --- Timestamped backup of tenet-runtime ---
TS="$(date +%Y%m%d-%H%M%S)"
SRC="${RUNTIME_DIR}"
DST="${BACKUP_BASE}/${BACKUP_PREFIX}-${TS}"

if [ -d "$SRC" ] && [ "$(ls -A "$SRC" 2>/dev/null)" ]; then
  echo "[debug] Backing up $SRC to $DST ..."
  sudo mkdir -p "$DST"
  sudo cp -a "$SRC"/. "$DST"/
  sudo chown -R ${RUNTIME_USER}:${RUNTIME_GROUP} "$DST"
  echo "[debug] Backup complete: $DST"
else
  echo "[debug] $SRC missing or empty; skipping backup."
fi

# Prune old backups if BACKUP_KEEP > 0
if [ "${BACKUP_KEEP}" -gt 0 ]; then
  echo "[debug] Pruning old backups, keeping most recent ${BACKUP_KEEP}..."
  # list newest first, skip first N, delete the rest
  ls -1dt "${BACKUP_BASE}/${BACKUP_PREFIX}-"* 2>/dev/null | tail -n +$((BACKUP_KEEP+1)) | xargs -r sudo rm -rf
fi

# --- Temp workspace ---
TEMP_DIR="$(mktemp -d)"
trap "rm -rf '$TEMP_DIR'" EXIT

sudo mkdir -p "${RUNTIME_DIR}"
sudo chown -R ${RUNTIME_USER}:${RUNTIME_GROUP} "${RUNTIME_DIR}"

# --- Stop Nakama if running to avoid 'Text file busy' on replace ---
if systemctl list-unit-files | grep -q "^${SERVICE_NAME}\.service"; then
  echo "[debug] Stopping ${SERVICE_NAME}.service if running..."
  sudo systemctl stop "${SERVICE_NAME}" || true
fi

# --- Install Go ---
echo "Installing Go"
if ! command -v go &> /dev/null || [[ $(go version 2>/dev/null) != *"go${GO_VERSION}"* ]]; then
  echo " [debug] Installing Go ${GO_VERSION}..."
  cd "$TEMP_DIR"
  wget -q "${GO_TARBALL_URL}" -O go.linux-amd64.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf go.linux-amd64.tar.gz
  cd - > /dev/null

  # Add Go to PATH if not already there
  if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
  fi
fi

# --- Download and install Nakama binary ---
if [ ! -f "${RUNTIME_DIR}/nakama" ]; then
  echo " [debug] Installing Nakama v${NAKAMA_VERSION}..."
  cd "$TEMP_DIR"
  wget -q "${NAKAMA_TARBALL_URL}" -O nakama.tar.gz
  tar -xzf nakama.tar.gz
  cd - > /dev/null
  cp "$TEMP_DIR/nakama" "${RUNTIME_DIR}/nakama"
  chmod +x "${RUNTIME_DIR}/nakama"
fi

# --- Install AWS CLI v2 for secrets manager ---
if ! command -v aws &> /dev/null; then
  echo "[debug] Installing AWS CLI v2..."
  cd "$TEMP_DIR"
  wget "${AWSCLI_ZIP_URL}" -O "awscliv2.zip"
  unzip -q awscliv2.zip
  sudo ./aws/install
  cd - > /dev/null
fi

# --- Install jq for command line JSON parsing ---
if ! command -v jq &> /dev/null; then
  echo "[debug] Installing jq ${JQ_VERSION} from source..."
  cd "$TEMP_DIR"
  wget -q "${JQ_TARBALL_URL}" -O jq.tar.gz
  tar -xzf jq.tar.gz
  cd "jq-${JQ_VERSION}"
  ./configure
  make
  sudo make install
  cd - > /dev/null
fi

# --- Ensure yq exists (mikefarah/yq). Install if missing. ---
if ! command -v yq >/dev/null 2>&1; then
  echo "[debug] Installing yq (${YQ_VERSION})"
  cd "$TEMP_DIR"
  wget -q "${YQ_TARBALL_URL}" -O - | tar xz
  sudo mv "${YQ_BINARY}" /usr/local/bin/yq
  cd - >/dev/null
fi

mkdir -p "${MODULES_DIR}"

# --- systemd unit for Nakama ---
echo "[debug] Writing systemd unit for Nakama"
sudo tee "${SERVICE_FILE}" >/dev/null <<UNIT
[Unit]
Description=Nakama Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${RUNTIME_USER}
Group=${RUNTIME_GROUP}
WorkingDirectory=${RUNTIME_DIR}
# Adjust flags as you need; --config should point to your local.yml
ExecStart=${RUNTIME_DIR}/nakama --config ${RUNTIME_DIR}/local.yml
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
sudo systemctl enable "${SERVICE_NAME}"
echo "[debug] ${SERVICE_NAME}.service installed and enabled (not started yet)"
