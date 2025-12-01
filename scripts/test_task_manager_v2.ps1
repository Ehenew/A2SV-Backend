$ErrorActionPreference = "Stop"
$ServerPort = 8081
$ServerUrl = "http://localhost:$ServerPort"
$ProcessName = "task_manager_v2_test"
$ProjectDir = "$PSScriptRoot/../task_manager_v2"

Write-Host "Building the application..."
Push-Location $ProjectDir
try {
    go mod tidy
    go build -o $ProcessName.exe ./Delivery
}
finally {
    Pop-Location
}

Write-Host "Starting server on port $ServerPort..."
$env:PORT = $ServerPort
# Using the same credentials as before
$env:MONGODB_URI = "mongodb+srv://ehenewamogne:sxSsjuKDhw7Sl0v3@cluster0.y8khhhs.mongodb.net/?appName=Cluster0"
$env:JWT_SECRET = "clean-arch-secret-v2"

$ExePath = Join-Path $ProjectDir "$ProcessName.exe"
$ServerProcess = Start-Process -FilePath $ExePath -PassThru -NoNewWindow

# Wait for server to start
Start-Sleep -Seconds 5

try {
    Write-Host "Running tests..."
    & "$PSScriptRoot/test_task_manager_auth.ps1"
}
finally {
    Write-Host "Stopping server..."
    if ($ServerProcess) {
        Stop-Process -Id $ServerProcess.Id -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path $ExePath) {
        Remove-Item $ExePath -ErrorAction SilentlyContinue
    }
}
