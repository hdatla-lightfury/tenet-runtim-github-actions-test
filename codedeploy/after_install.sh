#!/bin/bash

set -e

SECRET_ID="staging/tenet-runtime"
AWS_REGION=ap-south-1

echo "[debug] Fetching DATABASE_ADDRESS from Secrets Manager: ${SECRET_ID}"

DATABASE_ADDRESS="$(aws secretsmanager get-secret-value \
  --secret-id "${SECRET_ID}" \
  ${AWS_REGION:+--region "$AWS_REGION"} \
  --query 'SecretString' \
  --output text \
  | jq -r '.DATABASE_ADDRESS')"

if [[ -z "$DATABASE_ADDRESS" || "$DATABASE_ADDRESS" == "null" ]]; then
  echo "[error] DATABASE_ADDRESS not found in secret ${SECRET_ID}"
  exit 1
fi

LOCAL_YAML_FILE="./local.yml"

if [[ ! -f "${LOCAL_YAML_FILE}" ]]; then
  echo "[error] local.yml not found in parent dir,  Exiting."
  exit 1
fi

echo "[debug] Updating $(basename "$LOCAL_YAML_FILE") with yq"
DATABASE_ADDRESS="${DATABASE_ADDRESS}" yq -i '
  .database.address = [ env(DATABASE_ADDRESS) ]
' "${LOCAL_YAML_FILE}"

MODULES_DIR="./modules"

if [[ ! -f "./backend.so" ]]; then
  echo "[error] backend.so not found in ."
  exit 1
fi

echo "[debug] Copying backend.so into ${MODULES_DIR}"
cp -f "./backend.so" "${MODULES_DIR}/backend.so"

# Run database migrations
echo "🔄 Running database migrations..."
./nakama migrate up --database.address ${DATABASE_ADDRESS}
