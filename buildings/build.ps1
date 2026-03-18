$ProjectDir = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
$OutputDir = "."
$DLLPath = Join-Path $OutputDir "vmm.dll"
$HeaderPath = Join-Path $OutputDir "vmm.h"

if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
}

Remove-Item -Force -ErrorAction SilentlyContinue $DLLPath
Remove-Item -Force -ErrorAction SilentlyContinue $HeaderPath

Push-Location $ProjectDir

$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

go mod download
go mod tidy

go build -tags dll -buildmode=c-shared -o "$DLLPath" .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    Pop-Location
    Read-Host "Press Enter to exit"
    exit 1
}

$HeaderContent = @"
#ifndef VMM_H
#define VMM_H

#include <stdint.h>

typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;

extern Result __cdecl VMCreate(const char* filename, int size, const char* typ, int stringLength);
extern Result __cdecl VMOpen(const char* filename);
extern Result __cdecl VMClose(int handle);
extern Result __cdecl VMRead(int handle, int index);
extern Result __cdecl VMWrite(int handle, int index, const char* value);
extern Result __cdecl VMHelp(const char* filename);

#endif
"@

$HeaderContent | Out-File -FilePath $HeaderPath -Encoding ASCII

Pop-Location

Write-Host ""
Write-Host "DLL saved to: $(Resolve-Path $DLLPath)"
Write-Host "Header saved to: $(Resolve-Path $HeaderPath)"
Write-Host ""
Write-Host "Press Enter to exit..." -ForegroundColor Cyan
Read-Host
