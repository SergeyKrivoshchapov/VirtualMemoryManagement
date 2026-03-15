param(
    [string]$OutputDir = ".",
    [string]$Configuration = "Release"
)

$ErrorActionPreference = "Stop"
function Write-Success {
    Write-Host $args -ForegroundColor Green
}
function Write-Error-Custom {
    Write-Host $args -ForegroundColor Red
}
function Write-Warning-Custom {
    Write-Host $args -ForegroundColor Yellow
}

$ProjectName = "VirtualMemoryManagement"
$DLLName = "vmm.dll"
$DLLPath = Join-Path $OutputDir $DLLName

Write-Host "Building $ProjectName for Windows" -ForegroundColor Yellow
try {
    $GoVersion = go version
    Write-Success "Go found: $GoVersion"
} catch {
    Write-Error-Custom "Go is not installed or not in PATH"
    exit 1
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir = $ScriptDir

Write-Warning-Custom "Project directory: $ProjectDir"
Write-Warning-Custom "Output directory: $OutputDir"
Write-Warning-Custom "Configuration: $Configuration"

if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
    Write-Success "Created output directory: $OutputDir"
}

Write-Warning-Custom "Cleaning previous builds"
Remove-Item -Force -ErrorAction SilentlyContinue $DLLPath
Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path $OutputDir "*.h")
Push-Location $ProjectDir
go clean -cache 2>$null
Pop-Location

Write-Warning-Custom "Downloading dependencies"
Push-Location $ProjectDir
go mod download
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Error-Custom "Failed to download dependencies"
    Pop-Location
    exit 1
}
Pop-Location

Write-Warning-Custom "Building DLL for Windows"
Push-Location $ProjectDir

$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

try {
    go build `
        -tags dll `
        -buildmode=c-shared `
        -o "$DLLPath" `
        .

    if ($LASTEXITCODE -ne 0) {
        throw "Build failed with exit code $LASTEXITCODE"
    }
} catch {
    Write-Error-Custom "Build failed: $_"
    Pop-Location
    exit 1
}

Write-Success "Build successful!"

Write-Warning-Custom "Generating C header file..."
$HeaderPath = Join-Path $OutputDir "vmm.h"

$HeaderContent = @"
#ifndef VMM_H
#define VMM_H

#include <stdint.h>

/* Result structure for API responses */
typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;

/* Virtual Memory Manager Functions */

/* Create a new virtual array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMCreate(const char* filename, int size, const char* typ, int stringLength);

/* Open an existing virtual array file
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMOpen(const char* filename);

/* Close a virtual array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMClose(int handle);

/* Read a value from the array
   Returns: Result struct with data and error code */
extern Result __cdecl VMRead(int handle, int index);

/* Write a value to the array
   Returns: 1 on success, 0 on failure */
extern int __cdecl VMWrite(int handle, int index, const char* value);

/* Get help information
   Returns: Result struct with help text */
extern Result __cdecl VMHelp(const char* filename);

#endif /* VMM_H */
"@

$HeaderContent | Out-File -FilePath $HeaderPath -Encoding ASCII
Write-Success "Header file generated: $HeaderPath"

# Check file size
$DLLItem = Get-Item $DLLPath -ErrorAction SilentlyContinue
if ($DLLItem) {
    $DLLSize = "{0:N2} MB" -f ($DLLItem.Length / 1MB)
    Write-Success "DLL size: $DLLSize"
    Write-Success "DLL path: $DLLPath"
} else {
    Write-Error-Custom "DLL file not found after build"
    Pop-Location
    exit 1
}

Pop-Location

Write-Host ""
Write-Host "Build completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Success "Next steps:"
Write-Host "1. Copy $DLLName to your C# project"
Write-Host "2. Create C# P/Invoke declarations from $HeaderPath"
Write-Host "3. Reference the DLL in your C# project"

