#!/bin/bash

set -eoux pipefail

binary=$1

if [ ! -f "${CONFIG_FILE}" ]; then
  echo "Config file ${CONFIG_FILE} not found";
  exit 1;
fi

if [ -z "${CHRONOSPHERE_ORG_NAME:-}" ]; then
  echo "Environment variable CHRONOSPHERE_ORG_NAME is not set";
  exit 1;
fi

echo "Starting MCP server..."
ARGS=""
if [ -z "${CHRONOSPHERE_API_TOKEN:-}" ]; then
  if [ ! -f .chronosphere_api_token ]; then
    echo "Either CHRONOSPHERE_API_TOKEN needs to be set or .chronosphere_api_token file must contain the chronosphere api token";
    exit 1;
  fi
  ARGS=(--api-token-filename .chronosphere_api_token);
fi

${binary} -c "${CONFIG_FILE}" "${ARGS[@]}" --verbose