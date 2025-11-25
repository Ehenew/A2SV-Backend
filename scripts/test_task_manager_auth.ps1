$BaseUrl = "http://localhost:8081"

function Write-Color($text, $color) {
    Write-Host $text -ForegroundColor $color
}

function Test-Endpoint {
    param (
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{},
        [object]$Body = $null,
        [string]$Description,
        [int]$ExpectedStatus = 200
    )

    Write-Host "Testing: $Description" -NoNewline
    
    try {
        $params = @{
            Method = $Method
            Uri = $Url
            ErrorAction = "Stop"
        }
        if ($Headers.Count -gt 0) { $params.Headers = $Headers }
        if ($Body) { 
            $params.Body = ($Body | ConvertTo-Json -Depth 10) 
            $params.ContentType = "application/json"
        }

        $response = Invoke-RestMethod @params
        
        # Invoke-RestMethod throws on 4xx/5xx, so if we are here, it's likely 2xx
        if ($ExpectedStatus -ge 200 -and $ExpectedStatus -lt 300) {
            Write-Color " [PASS]" "Green"
            return $response
        } else {
            Write-Color " [FAIL] (Expected $ExpectedStatus, got success)" "Red"
            return $null
        }
    }
    catch {
        $ex = $_.Exception
        $status = 0
        if ($ex.Response) {
            $status = [int]$ex.Response.StatusCode
        }

        if ($status -eq $ExpectedStatus) {
            Write-Color " [PASS] (Got expected error $status)" "Green"
            return $null
        } else {
            Write-Color " [FAIL] (Expected $ExpectedStatus, got $status)" "Red"
            Write-Host "Error: $($ex.Message)"
            if ($ex.Response) {
                $reader = New-Object System.IO.StreamReader($ex.Response.GetResponseStream())
                Write-Host "Response Body: $($reader.ReadToEnd())"
            }
            return $null
        }
    }
}

# --- Main Test Flow ---

$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$adminUser = "admin_$timestamp"
$regularUser = "user_$timestamp"
$password = "password123"

Write-Color "`n--- 1. Register Admin User ---" "Cyan"
# The first user registered in the system becomes admin. 
# Note: If the DB is not empty, this might fail to be admin if users exist.
# But for the script, we assume we can register.
$registerResponse = Test-Endpoint -Method "POST" -Url "$BaseUrl/register" -Body @{ username = $adminUser; password = $password } -Description "Register Admin User" -ExpectedStatus 201

Write-Color "`n--- 2. Login Admin User ---" "Cyan"
$loginResponse = Test-Endpoint -Method "POST" -Url "$BaseUrl/login" -Body @{ username = $adminUser; password = $password } -Description "Login Admin User"
$adminToken = $loginResponse.token
$adminHeaders = @{ Authorization = "Bearer $adminToken" }
Write-Host "Admin Token acquired."

Write-Color "`n--- 2b. Check Admin Role ---" "Cyan"
$me = Test-Endpoint -Method "GET" -Url "$BaseUrl/me" -Headers $adminHeaders -Description "Check Admin Role"
Write-Host "Role: $($me.role)"

Write-Color "`n--- 3. Register Regular User ---" "Cyan"
Test-Endpoint -Method "POST" -Url "$BaseUrl/register" -Body @{ username = $regularUser; password = $password } -Description "Register Regular User" -ExpectedStatus 201

Write-Color "`n--- 4. Login Regular User ---" "Cyan"
$userLoginResponse = Test-Endpoint -Method "POST" -Url "$BaseUrl/login" -Body @{ username = $regularUser; password = $password } -Description "Login Regular User"
$userToken = $userLoginResponse.token
$userId = $userLoginResponse.user_id
$userHeaders = @{ Authorization = "Bearer $userToken" }
Write-Host "User Token acquired. User ID: $userId"

Write-Color "`n--- 5. Admin: Create Task ---" "Cyan"
$taskBody = @{
    title = "Admin Task"
    description = "Created by Admin"
    due_date = (Get-Date).AddDays(1).ToString("yyyy-MM-ddTHH:mm:ssZ")
    status = "Pending"
}
$task = Test-Endpoint -Method "POST" -Url "$BaseUrl/tasks" -Headers $adminHeaders -Body $taskBody -Description "Admin creates task" -ExpectedStatus 201
$taskId = $task.id

Write-Color "`n--- 6. User: Try Create Task (Should Fail) ---" "Cyan"
Test-Endpoint -Method "POST" -Url "$BaseUrl/tasks" -Headers $userHeaders -Body $taskBody -Description "User tries to create task" -ExpectedStatus 403

Write-Color "`n--- 7. User: Get Tasks ---" "Cyan"
Test-Endpoint -Method "GET" -Url "$BaseUrl/tasks" -Headers $userHeaders -Description "User gets all tasks"

Write-Color "`n--- 8. Admin: Promote User ---" "Cyan"
Test-Endpoint -Method "POST" -Url "$BaseUrl/promote" -Headers $adminHeaders -Body @{ user_id = $userId } -Description "Admin promotes user"

Write-Color "`n--- 8b. User Re-Login (to get new token with admin role) ---" "Cyan"
$userLoginResponse = Test-Endpoint -Method "POST" -Url "$BaseUrl/login" -Body @{ username = $regularUser; password = $password } -Description "User re-logins"
$userToken = $userLoginResponse.token
$userHeaders = @{ Authorization = "Bearer $userToken" }
Write-Host "New User Token acquired."

Write-Color "`n--- 9. Promoted User: Create Task (Should Pass) ---" "Cyan"
$userTaskBody = @{
    title = "User Task"
    description = "Created by Promoted User"
    due_date = (Get-Date).AddDays(1).ToString("yyyy-MM-ddTHH:mm:ssZ")
    status = "Pending"
}
Test-Endpoint -Method "POST" -Url "$BaseUrl/tasks" -Headers $userHeaders -Body $userTaskBody -Description "Promoted user creates task" -ExpectedStatus 201

Write-Color "`n--- 10. Admin: Delete Task ---" "Cyan"
if ($taskId) {
    Test-Endpoint -Method "DELETE" -Url "$BaseUrl/tasks/$taskId" -Headers $adminHeaders -Description "Admin deletes task"
}

Write-Color "`n--- Test Complete ---" "Green"
