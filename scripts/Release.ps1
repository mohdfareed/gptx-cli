#!/usr/bin/env pwsh
param (
  [Parameter(Mandatory = $true, HelpMessage = "the version tag, e.g. v1.0.0")]
  [string] $Version
)

Write-Host "Creating release: $Version"
Write-Host "Press Enter to continue"
Read-Host

git tag $Version
git push origin $Version
