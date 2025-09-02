#!/bin/bash
set -euo pipefail

# --- Load constants ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONSTANTS_FILE="${CONSTANTS_FILE:-${SCRIPT_DIR}/constants.cnf}"
if [ ! -f "${CONSTANTS_FILE}" ]; then
  echo "[error] constants file not found at ${CONSTANTS_FILE}"
  exit 1
fi

source "${CONSTANTS_FILE}"

echo "Starting ${SERVICE_NAME} via systemd..."
sudo systemctl daemon-reload
sudo systemctl restart "${SERVICE_NAME}"

# Optionally wait until it’s 'active'
for i in {1..30}; do
  if sudo systemctl is-active --quiet "${SERVICE_NAME}"; then
    echo "Service ${SERVICE_NAME} is active."
    break
  fi
  sleep 1
done
