#!/bin/bash

set -e

set -e

sudo chmod -R 666 /home/ec2-user/tenet-runtime/*

SECRET_ID="staging/tenet-runtime/aws-rds-postgres"
AWS_REGION="ap-south-1"

echo "[debug] Fetching DB credentials from Secrets Manager: ${SECRET_ID}"

# Fetch the secret JSON
secret_str="$(aws secretsmanager get-secret-value \
  --secret-id "${SECRET_ID}" \
  ${AWS_REGION:+--region "$AWS_REGION"} \
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
dbname="postgres"

if [[ -z "$username" || -z "$password" || -z "$host" || -z "$port" || -z "$dbname" ]]; then
  echo "[error] One or more required fields missing in secret JSON"
  exit 2
fi

DATABASE_ADDRESS="${username}:${password}@${host}:${port}/${dbname}"
echo "[debug] Built DATABASE_ADDRESS"


# Codedeploy looks for all the folders in source directory, not relative to this folder.
LOCAL_YAML_FILE="/home/ec2-user/tenet-runtime/local.yml"

if [[ ! -f "${LOCAL_YAML_FILE}" ]]; then
  echo "[error] local.yml not found in parent dir,  Exiting."
  exit 1
fi

echo "[debug] Updating $(basename "$LOCAL_YAML_FILE") with yq"
DATABASE_ADDRESS="${DATABASE_ADDRESS}" yq -i '
  .database.address = [ env(DATABASE_ADDRESS) ]
' "${LOCAL_YAML_FILE}"


# Run database migrations
echo "🔄 Running database migrations..."
./nakama migrate up --database.address ${DATABASE_ADDRESS}


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
