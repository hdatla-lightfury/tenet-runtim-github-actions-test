#!/usr/bin/env bash
# Centralized constants for Nakama setup, runtime, and post-deploy steps.

# --- System / service ---
SERVICE_NAME="nakama"
RUNTIME_USER="ec2-user"
RUNTIME_GROUP="ec2-user"

# --- Paths ---
RUNTIME_DIR="/home/ec2-user/tenet-runtime"
MODULES_DIR="${RUNTIME_DIR}/modules"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
NAKAMA_BIN="${RUNTIME_DIR}/nakama"
BACKEND_SO_NAME="backend.so"
BACKEND_SO_PATH="${RUNTIME_DIR}/${BACKEND_SO_NAME}"
LOCAL_YAML_FILE="${RUNTIME_DIR}/local.yml"

# --- Backup settings ---
BACKUP_BASE="/home/ec2-user"
BACKUP_PREFIX="tenet-runtime-backup"
BACKUP_KEEP=7                 # keep most recent N backups (0 = disable pruning)

# --- Versions ---
GO_VERSION="1.24.5"
NAKAMA_VERSION="3.28.0"
JQ_VERSION="1.8.1"
YQ_VERSION="v4.47.1"
YQ_BINARY="yq_linux_amd64"

# --- URLs ---
GO_TARBALL_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
NAKAMA_TARBALL_URL="https://github.com/heroiclabs/nakama/releases/download/v${NAKAMA_VERSION}/nakama-${NAKAMA_VERSION}-linux-amd64.tar.gz"
AWSCLI_ZIP_URL="https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip"
JQ_TARBALL_URL="https://github.com/jqlang/jq/releases/download/jq-${JQ_VERSION}/jq-${JQ_VERSION}.tar.gz"
YQ_TARBALL_URL="https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/${YQ_BINARY}.tar.gz"

# --- Secrets / DB ---
SECRET_ID="staging/tenet-runtime/aws-rds-postgres"
AWS_REGION="ap-south-1"
DB_NAME="postgres"

# --- Permission  (used by post-deploy fixups) ---
CHMOD_RECURSE_TARGET="${RUNTIME_DIR}/*"
CHMOD_RECURSE_MODE="777" 
CHOWN_SINGLE_TARGET="${BACKEND_SO_PATH}"
CHOWN_USER="${RUNTIME_USER}"
CHOWN_GROUP="${RUNTIME_GROUP}"

# --- Service wait loop ---
SERVICE_RETRY_COUNT=30    # how many attempts
SERVICE_RETRY_SLEEP=1     # seconds between attempts