#!/usr/bin/env pwsh
<# .SYNOPSIS #>
Param(
  [string]$Output = ".bin",
  [switch]$Help
)
$ErrorActionPreference = "Stop"

# MARK: Init ==================================================================

# Help
if ($Help) {
  Write-Host "Usage: Build.ps1 [Output <path=.bin>]"
  exit 0
}

# Arguments
$AppPath = Join-Path -Path $PWD -ChildPath "cmd/gptx"
New-Item -ItemType Directory -Path $Output -ErrorAction SilentlyContinue | Out-Null
$BinPath = Resolve-Path -Path $Output

# Store current environment
$DevOS = $env:GOOS
$DevArch = $env:GOARCH

# Download dependencies and clean bin
go build -o "$BinPath/_" "$AppPath"
Remove-Item -Path $BinPath -Recurse -Force
New-Item -ItemType Directory -Path $Output | Out-Null

# MARK: Build =================================================================

function Build-Binary {
  [CmdletBinding()]
  Param(
    [Parameter(Mandatory = $true)][string]$Platform,
    [Parameter(Mandatory = $true)][string]$Architecture,
    [Parameter(Mandatory = $true)][string]$Id
  )

  $ExeName = Split-Path -Leaf $AppPath
  $OutputExe = Join-Path -Path $Output -ChildPath $ExeName
  if ($Platform -eq 'windows') {
    $OutputExe = "$OutputExe.exe"
  }

  $ArchiveName = "$ExeName-$Id.zip"
  $ArchivePath = Join-Path -Path $Output -ChildPath $ArchiveName

  Write-Host "Building for $Platform $Architecture..."
  $env:GOOS = $Platform
  $env:GOARCH = $Architecture
  go build -C $Output $AppPath

  if ($Platform -ne 'windows') {
    chmod +x "$OutputExe" # make it executable
  }
  Compress-Archive -Path $OutputExe -DestinationPath $ArchivePath -Force
  Remove-Item -Path $OutputExe -Force
  Write-Host "-> Packaged: $ArchivePath"
}

# MARK: Targets ===============================================================

# Release builds
Build-Binary -Platform 'linux'   -Architecture 'arm64' -Id 'linux-arm'
Build-Binary -Platform 'linux'   -Architecture 'amd64' -Id 'linux-x64'
Build-Binary -Platform 'darwin'  -Architecture 'arm64' -Id 'macos-arm'
Build-Binary -Platform 'darwin'  -Architecture 'amd64' -Id 'macos-x64'
Build-Binary -Platform 'windows' -Architecture 'arm64' -Id 'win-arm'
Build-Binary -Platform 'windows' -Architecture 'amd64' -Id 'win-x64'

# Development build
Write-Host "Building for development..."
$env:GOOS = $DevOS
$env:GOARCH = $DevArch
go build -C $BinPath -tags=dev $AppPath
Write-Host "-> Dev pkg at: $BinPath/$(Split-Path -Leaf $AppPath)"

# Update documentation
# echo "updating docs..."
# echo "# gptx" > "cmd/cli.md" # root/cmd/cli.md
# echo > "cmd/cli.md"
# if [ "$GOOS" = "windows" ]; then
#   go run -C $BIN/gptx.exe -h >> "cmd/cli.md"
# fi


Write-Host "updating docs..."
"# gptx" | Out-File "cmd/cli.md" -Encoding utf8
"" | Add-Content "cmd/cli.md"
if ($env:GOOS -eq "windows") {
  & go run -C "$env:BIN/gptx.exe" -h | Add-Content "cmd/cli.md"
}
else {
  & go run -C "$env:BIN/gptx" -h | Add-Content "cmd/cli.md"
}
