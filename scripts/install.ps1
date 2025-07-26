$binDir = Join-Path $env:USERPROFILE ".bin"
if (!(Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

$exePath = Join-Path $binDir "termpet.exe"

try {
    Invoke-WebRequest -Uri "https://github.com/v1ctorio/termpet/releases/latest/download/termpet_Windows-x86_64.exe" -OutFile $exePath
}
catch {
    Write-Error "Download failed: $_"
}

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$binDir*") {
    $newPath = $currentPath + ";$binDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "Added $binDir to PATH"
}

Write-Host "Termpet succesfully installed! Just write `termpet` to get started using it!"
