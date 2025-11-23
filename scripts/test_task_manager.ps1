Param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$CleanStart
)

Write-Host "--- Task Manager API Automated Test ---" -ForegroundColor Cyan
Write-Host "(Verbose count fix enabled)" -ForegroundColor DarkGray

function New-JsonBody {
    param(
        [string]$Title,
        [string]$Description,
        [string]$DueDate,
        [string]$Status
    )
    return @{ title = $Title; description = $Description; due_date = $DueDate; status = $Status } | ConvertTo-Json -Compress
}

function Create-Task {
    param([string]$Title, [string]$Description, [string]$DueDate, [string]$Status = "pending")
    $body = New-JsonBody -Title $Title -Description $Description -DueDate $DueDate -Status $Status
    $resp = Invoke-RestMethod -Method Post -Uri "$BaseUrl/tasks" -ContentType "application/json" -Body $body -ErrorAction Stop
    if (-not $resp.id) { throw "Create failed: no id returned" }
    Write-Host "Created Task ID:" $resp.id -ForegroundColor Green
    return $resp.id
}

function List-Tasks {
    $raw = Invoke-RestMethod -Method Get -Uri "$BaseUrl/tasks" -ErrorAction Stop
    # Normalize to array; if API returned [] we get $null, treat as empty
    if ($null -eq $raw) { $tasks = @() }
    elseif ($raw -is [System.Array]) { $tasks = $raw }
    else { $tasks = @($raw) }
    # Remove null / empty placeholder objects
    $tasks = $tasks | Where-Object { $_ -ne $null -and ($_.id -or $_.title -or $_.status) }
    $count = $tasks.Count
    Write-Host "Total tasks:" $count -ForegroundColor Yellow
    if ($count -gt 0) {
        Write-Host (ConvertTo-Json -InputObject $tasks -Depth 6) -ForegroundColor DarkYellow
    } else {
        Write-Host "(empty list)" -ForegroundColor DarkYellow
    }
    return $tasks
}

function Get-Task {
    param([string]$Id)
    $resp = Invoke-RestMethod -Method Get -Uri "$BaseUrl/tasks/$Id" -ErrorAction Stop
    Write-Host "Fetched Task ($Id): status=$($resp.status)" -ForegroundColor Yellow
    return $resp
}

function Update-Task {
    param([string]$Id, [hashtable]$Fields)
    $json = $Fields | ConvertTo-Json -Compress
    $resp = Invoke-RestMethod -Method Put -Uri "$BaseUrl/tasks/$Id" -ContentType "application/json" -Body $json -ErrorAction Stop
    Write-Host "Update response:" ($resp.message) -ForegroundColor Green
}

function Delete-Task {
    param([string]$Id)
    $resp = Invoke-RestMethod -Method Delete -Uri "$BaseUrl/tasks/$Id" -ErrorAction Stop
    Write-Host "Delete response:" ($resp.message) -ForegroundColor Green
}

# 0. Initial list before any action (optionally clean existing tasks)
$initial = List-Tasks
if ($CleanStart -and $initial.Count -gt 0) {
    Write-Host "Cleaning existing tasks..." -ForegroundColor Magenta
    $initial | ForEach-Object {
        if ($_.id) {
            try { Invoke-RestMethod -Method Delete -Uri "$BaseUrl/tasks/$($_.id)" -ErrorAction Stop | Out-Null } catch {}
        }
    }
    Write-Host "Cleanup done. Listing again:" -ForegroundColor Magenta
    List-Tasks | Out-Null
}

# 1. Create
$id = Create-Task -Title "Automated" -Description "Script flow" -DueDate "2025-04-01T00:00:00Z"

# 2. List
$list1 = List-Tasks

# 3. Get by ID
$one = Get-Task -Id $id

# 4. Update status
Update-Task -Id $id -Fields @{ status = "completed" }
$updated = Get-Task -Id $id

# 5. Delete
Delete-Task -Id $id

# 6. Verify deletion
try {
    $gone = Get-Task -Id $id
    Write-Host "Unexpected: task still exists" -ForegroundColor Red
} catch {
    Write-Host "Confirmed deletion (not found)." -ForegroundColor Green
}

# 7. Final list (should reflect deletion)
$list2 = List-Tasks

# Extra explicit verification block
if ($list2.Count -eq 0) {
    Write-Host "Final verification: collection empty." -ForegroundColor Green
} else {
    Write-Host "Final verification: tasks remain:" -ForegroundColor Red
    $list2 | ForEach-Object { Write-Host " - ID: $($_.id) Title: $($_.title) Status: $($_.status)" -ForegroundColor Red }
}

Write-Host "--- Test Flow Complete ---" -ForegroundColor Cyan
