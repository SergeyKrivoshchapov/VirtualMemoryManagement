#!/bin/bash
set -e

OUTPUT_DIR="${1:-.}"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS_TYPE="linux"
    SCRIPT_NAME="build.sh"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    OS_TYPE="windows"
    SCRIPT_NAME="build.ps1"
else
    OS_TYPE="unknown"
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Universal Build Script for VirtualMemoryManagement"
echo "Detected OS: $OS_TYPE"
echo "Output directory: $OUTPUT_DIR"
echo ""

case $OS_TYPE in
    linux|macos)
        echo "Running Linux build script"
        bash "$SCRIPT_DIR/build.sh" "$OUTPUT_DIR"
        ;;
    windows)
        echo "Running Windows build script"
        powershell -NoProfile -ExecutionPolicy Bypass -File "$SCRIPT_DIR/build.ps1" -OutputDir "$OUTPUT_DIR"
        ;;
    *)
        echo "Error: Unknown operating system: $OSTYPE"
        exit 1
        ;;
esac

exit $?

