#!/bin/bash
set -e

SERVICE="nakama"
sudo systemctl status --no-pager -l "${SERVICE}" || true
