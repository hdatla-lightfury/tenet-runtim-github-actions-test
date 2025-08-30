#!/bin/bash
set -euo pipefail

# --- Load constants ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONSTANTS_FILE="${CONSTANTS_FILE:-${SCRIPT_DIR}/constants.sh}"
if [ ! -f "${CONSTANTS_FILE}" ]; then
  echo "[error] constants file not found at ${CONSTANTS_FILE}"
  exit 1
fi

source "${CONSTANTS_FILE}"

echo "Stopping ${SERVICE_NAME} via systemd..."
if sudo systemctl is-active --quiet "${SERVICE_NAME}"; then
  sudo systemctl stop "${SERVICE_NAME}"
  echo "Waiting up to $((SERVICE_RETRY_COUNT * SERVICE_RETRY_SLEEP))s for clean stop..."
  for i in $(seq 1 "${SERVICE_RETRY_COUNT}"); do
    if ! sudo systemctl is-active --quiet "${SERVICE_NAME}"; then
      echo "Service ${SERVICE_NAME} stopped."
      break
    fi
    sleep "${SERVICE_RETRY_SLEEP}"
  done
else
  echo "Service ${SERVICE_NAME} is not active."
fi

echo "Status after stop:"
sudo systemctl status --no-pager -l "${SERVICE_NAME}" || true
echo "Stop completed."
