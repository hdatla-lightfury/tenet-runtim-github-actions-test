#!/bin/bash
set -euo pipefail

SERVICE="nakama"
echo "Starting ${SERVICE} via systemd..."
sudo systemctl daemon-reload
sudo systemctl restart "${SERVICE}"

# Optionally wait until it’s 'active'
for i in {1..30}; do
  if sudo systemctl is-active --quiet "${SERVICE}"; then
    echo "Service ${SERVICE} is active."
    break
  fi
  sleep 1
done

sudo systemctl status --no-pager -l "${SERVICE}" || true
