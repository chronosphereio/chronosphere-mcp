#!/bin/bash
set -euo pipefail

# Check if there are any changes
if ! git diff --quiet; then
    echo "Error: Code generation produced changes that are not checked in!"
    echo "Please run 'make all-gen' locally and commit the changes."
    git diff --name-status
    exit 1
else
    echo "Success: No code generation changes detected. Generated code is up-to-date."
fi