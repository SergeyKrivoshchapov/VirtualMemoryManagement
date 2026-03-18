#!/bin/bash
set -e

OUTPUT_DIR="${1:-.}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    powershell -NoProfile -ExecutionPolicy Bypass -File "$SCRIPT_DIR/buildings/build.ps1"
else
    bash "$SCRIPT_DIR/buildings/build.sh" "$OUTPUT_DIR"
fi

