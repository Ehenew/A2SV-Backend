Param(
    [string]$BaseUrl = "http://localhost:8080"
)

Write-Host "--- Go-Auth API Automated Test ---" -ForegroundColor Cyan

# 1. Register
Write-Host "`n1. Registering new user..." -ForegroundColor Yellow
$regBody = @{
    email = "user@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $regResp = Invoke-RestMethod -Method Post -Uri "$BaseUrl/register" -ContentType "application/json" -Body $regBody -ErrorAction Stop
    Write-Host "Success: $($regResp.message)" -ForegroundColor Green
} catch {
    Write-Host "Registration failed (might already exist): $($_.Exception.Message)" -ForegroundColor Red
}

# 2. Login
Write-Host "`n2. Logging in..." -ForegroundColor Yellow
try {
    $loginResp = Invoke-RestMethod -Method Post -Uri "$BaseUrl/login" -ContentType "application/json" -Body $regBody -ErrorAction Stop
    $token = $loginResp.token
    if ($token) {
        Write-Host "Login Successful! Token received." -ForegroundColor Green
        Write-Host "Token: $token" -ForegroundColor DarkGray
    } else {
        Write-Host "Login failed: No token returned." -ForegroundColor Red
        exit
    }
} catch {
    Write-Host "Login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 3. Access Secure Route (No Token)
Write-Host "`n3. Testing Secure Route (No Token)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Method Get -Uri "$BaseUrl/secure" -ErrorAction Stop
    Write-Host "Failed: Should have been unauthorized!" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode -eq [System.Net.HttpStatusCode]::Unauthorized) {
        Write-Host "Success: Access denied as expected (401)." -ForegroundColor Green
    } else {
        Write-Host "Unexpected error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 4. Access Secure Route (With Token)
Write-Host "`n4. Testing Secure Route (With Token)..." -ForegroundColor Yellow
try {
    $headers = @{ Authorization = "Bearer $token" }
    $secureResp = Invoke-RestMethod -Method Get -Uri "$BaseUrl/secure" -Headers $headers -ErrorAction Stop
    Write-Host "Success: $($secureResp.message)" -ForegroundColor Green
} catch {
    Write-Host "Failed to access secure route: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n--- Test Complete ---" -ForegroundColor Cyan
