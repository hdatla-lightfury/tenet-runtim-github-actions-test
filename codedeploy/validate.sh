#!/bin/bash
set -e

# --- Load constants ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONSTANTS_FILE="${CONSTANTS_FILE:-${SCRIPT_DIR}/constants.cnf}"
if [ ! -f "${CONSTANTS_FILE}" ]; then
  echo "[error] constants file not found at ${CONSTANTS_FILE}"
  exit 1
fi

source "${CONSTANTS_FILE}"

echo "Showing status for ${SERVICE_NAME}..."
sudo systemctl status --no-pager -l "${SERVICE_NAME}" || true
