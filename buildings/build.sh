#!/bin/bash
set -e
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_NAME="VirtualMemoryManagement"

# If no output dir specified, use default C# project output path
if [ -z "$1" ]; then
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_DIR="$( cd "${SCRIPT_DIR}/.." && pwd )"
    # Default output: CLI/TimpLaba2_VirtualMemory/TimpLaba2_VirtualMemory/bin/Debug/net10.0
    OUTPUT_DIR="${PROJECT_DIR}/CLI/TimpLaba2_VirtualMemory/TimpLaba2_VirtualMemory/bin/Debug/net10.0"
else
    OUTPUT_DIR="$1"
fi

DLL_NAME="vmm.so"
DLL_PATH="${OUTPUT_DIR}/${DLL_NAME}"
echo -e "${YELLOW}Building ${PROJECT_NAME} for Linux${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}Go version: ${GO_VERSION}${NC}"

# If PROJECT_DIR wasn't set from parameter, set it now
if [ -z "$PROJECT_DIR" ]; then
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_DIR="$( cd "${SCRIPT_DIR}/.." && pwd )"
fi

echo -e "${YELLOW}Project directory: ${PROJECT_DIR}${NC}"
echo -e "${YELLOW}Output directory: ${OUTPUT_DIR}${NC}"

echo -e "${YELLOW}Cleaning previous builds${NC}"
rm -f "${DLL_PATH}"
go clean -cache 2>/dev/null || true

echo -e "${YELLOW}Downloading dependencies${NC}"
cd "${PROJECT_DIR}"
go mod download
go mod tidy

echo -e "${YELLOW}Building DLL${NC}"
cd "${PROJECT_DIR}"

GOOS=linux GOARCH=amd64 go build \
    -tags dll \
    -buildmode=c-shared \
    -o "${DLL_PATH}" \
    .

# Check if build was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN} Build successful!${NC}"

    # Generate header file for C# interop
    HEADER_FILE="${OUTPUT_DIR}/vmm.h"
    echo -e "${YELLOW}Generating C header file...${NC}"

    # Simple header generation (basic structure)
    cat > "${HEADER_FILE}" << 'EOF'
#ifndef VMM_H
#define VMM_H

#include <stdint.h>

typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;

// Virtual Memory Manager Functions
int VMCreate(const char* filename, int size, const char* typ, int stringLength);
int VMOpen(const char* filename);
int VMClose(int handle);
Result VMRead(int handle, int index);
int VMWrite(int handle, int index, const char* value);
Result VMHelp(const char* filename);

#endif // VMM_H
EOF

    echo -e "${GREEN} Header file generated: ${HEADER_FILE}${NC}"

    # Check file size
    DLL_SIZE=$(ls -lh "${DLL_PATH}" | awk '{print $5}')
    echo -e "${GREEN} DLL size: ${DLL_SIZE}${NC}"
    echo -e "${GREEN} DLL path: ${DLL_PATH}${NC}"

    echo ""
    echo -e "${GREEN}Build completed successfully!${NC}"
else
    echo -e "${RED} Build failed!${NC}"
    exit 1
fi

