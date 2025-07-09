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
${binary} -c "${CONFIG_FILE}"