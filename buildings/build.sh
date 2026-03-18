#!/bin/bash
set -e

PROJECT_DIR="$(dirname "$0" | xargs dirname)"
OUTPUT_DIR="${1:-.}"

mkdir -p "$OUTPUT_DIR"
rm -f "$OUTPUT_DIR/vmm.so" "$OUTPUT_DIR/vmm.dylib" "$OUTPUT_DIR/vmm.h"

cd "$PROJECT_DIR"

export CGO_ENABLED=1

go mod download
go mod tidy
go build -buildmode=c-shared -o "$OUTPUT_DIR/vmm.so" .

cat > "$OUTPUT_DIR/vmm.h" << 'EOF'
#ifndef VMM_H
#define VMM_H

#include <stdint.h>

typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;

extern Result VMCreate(const char* filename, int size, const char* typ, int stringLength);
extern Result VMOpen(const char* filename);
extern Result VMClose(int handle);
extern Result VMRead(int handle, int index);
extern Result VMWrite(int handle, int index, const char* value);
extern Result VMHelp(const char* filename, const char* helpText);

#endif
EOF

echo ""
echo "DLL saved to: $(cd "$OUTPUT_DIR" && pwd)/vmm.so"
echo "Header saved to: $(cd "$OUTPUT_DIR" && pwd)/vmm.h"
echo ""
read -p "Press Enter to exit..."
