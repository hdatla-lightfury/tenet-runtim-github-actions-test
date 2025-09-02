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

# --- Permissions fixups (use constants) ---
sudo chmod -R "${CHMOD_RECURSE_MODE}" ${CHMOD_RECURSE_TARGET}
sudo chown "${CHOWN_USER}:${CHOWN_GROUP}" "${CHOWN_SINGLE_TARGET}"

echo "[debug] Fetching DB credentials from Secrets Manager: ${SECRET_ID}"

# Fetch the secret JSON
secret_str="$(aws secretsmanager get-secret-value \
  --secret-id "${SECRET_ID}" \
  ${AWS_REGION:+--region "${AWS_REGION}"} \
  --query 'SecretString' \
  --output text)"

if [[ -z "$secret_str" || "$secret_str" == "null" ]]; then
  echo "[error] SecretString not found for ${SECRET_ID}"
  exit 1
fi

# Parse fields and build DATABASE_ADDRESS
username="$(echo "$secret_str" | jq -r '.username')"
password="$(echo "$secret_str" | jq -r '.password')"
host="$(echo "$secret_str" | jq -r '.host')"
port="$(echo "$secret_str" | jq -r '.port')"
dbname="${DB_NAME}"

if [[ -z "$username" || -z "$password" || -z "$host" || -z "$port" || -z "$dbname" ]]; then
  echo "[error] One or more required fields missing in secret JSON"
  exit 2
fi

DATABASE_ADDRESS="${username}:${password}@${host}:${port}/${dbname}"
echo "[debug] Built DATABASE_ADDRESS"

# Codedeploy looks for folders in source directory; local.yml must already exist
if [[ ! -f "${LOCAL_YAML_FILE}" ]]; then
  echo "[error] local.yml not found at ${LOCAL_YAML_FILE}. Exiting."
  exit 1
fi

echo "[debug] Updating $(basename "${LOCAL_YAML_FILE}") with yq"
DATABASE_ADDRESS="${DATABASE_ADDRESS}" yq -i '
  .database.address = [ env(DATABASE_ADDRESS) ]
' "${LOCAL_YAML_FILE}"

# Ensure modules dir exists (from constants)
mkdir -p "${MODULES_DIR}"

# Copy backend module
if [[ ! -f "${BACKEND_SO_PATH}" ]]; then
  echo "[error] ${BACKEND_SO_NAME} not found at ${BACKEND_SO_PATH}"
  exit 1
fi

echo "[debug] Copying ${BACKEND_SO_NAME} into ${MODULES_DIR}"
cp -f "${BACKEND_SO_PATH}" "${MODULES_DIR}/"

# Run database migrations
echo "🔄 Running database migrations..."
"${NAKAMA_BIN}" migrate up --database.address "${DATABASE_ADDRESS}"
