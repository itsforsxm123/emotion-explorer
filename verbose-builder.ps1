<#
.SYNOPSIS
Runs 'go build' with verbose command execution output (-x flag).
.DESCRIPTION
This script executes 'go build -x ./...' in the current directory.
The -x flag prints the commands executed by the Go tool, which is useful
for seeing progress, especially during CGO compilation. It also measures
the build time.
.EXAMPLE
.\verbose-build.ps1
#>

Write-Host "Starting verbose Go build (go build -x ./...)" -ForegroundColor Cyan
Write-Host "This will show the commands being executed, including C compiler calls."
Write-Host "The first build involving CGO can take several minutes."
Write-Host "-------------------------------------------------------------"

# Measure the time taken for the build command
$buildDuration = Measure-Command {
    # Execute the verbose build command
    go build -x ./...

    # Capture the exit code of the last command (go build)
    # Note: $LASTEXITCODE needs to be checked immediately after the command
    $script:buildExitCode = $LASTEXITCODE
}

Write-Host "-------------------------------------------------------------"

# Check the exit code captured from the go build command
if ($script:buildExitCode -eq 0) {
    Write-Host "Verbose build completed successfully." -ForegroundColor Green
} else {
    Write-Host "Verbose build FAILED with exit code $script:buildExitCode." -ForegroundColor Red
    Write-Host "Review the output above for specific error messages."
}

Write-Host "Total build time: $($buildDuration.TotalSeconds) seconds." -ForegroundColor Cyan

# Optional: Pause at the end if running by double-clicking
# Read-Host "Press Enter to exit"