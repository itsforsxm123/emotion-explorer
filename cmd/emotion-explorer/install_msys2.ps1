#Requires -RunAsAdministrator

<#
.SYNOPSIS
Downloads and silently installs the base MSYS2 system, adds the expected
MinGW-w64 bin directory to the System PATH, and guides the user through
the manual pacman steps by opening the MSYS2 shell.

.DESCRIPTION
This script automates the initial setup for MSYS2 and MinGW-w64 required by Fyne.
It performs the following steps:
1. Downloads the latest known MSYS2 installer.
2. Runs the installer silently to the specified directory (default C:\msys64).
3. Appends the standard MinGW-w64 bin path (e.g., C:\msys64\mingw64\bin) to the
   System PATH environment variable if it's not already present.
4. Guides the user through the required 'pacman' update and toolchain install steps
   by opening the 'MSYS2 MSYS' shell and pausing for user confirmation.

!!! IMPORTANT !!!
You MUST follow the instructions and execute the 'pacman' commands displayed
within the separate 'MSYS2 MSYS' shell window when prompted by this script.

.PARAMETER InstallDir
The directory where MSYS2 should be installed. Defaults to 'C:\msys64'.

.PARAMETER InstallerUrl
The direct download URL for the MSYS2 installer executable. Defaults to a
recent known version, but you might want to update this from msys2.org.

.EXAMPLE
.\Install-Msys2MinGW-Guided.ps1
Downloads and installs MSYS2 to C:\msys64, updates the System PATH, and
guides through the pacman steps.

.NOTES
- Requires running PowerShell as Administrator.
- You MUST interact with the separate 'MSYS2 MSYS' shell window when prompted.
- You MUST restart any open terminals or IDEs after completing ALL steps.
#>
param(
    [Parameter(Mandatory=$false)]
    [string]$InstallDir = "C:\msys64",

    [Parameter(Mandatory=$false)]
    [string]$InstallerUrl = "https://github.com/msys2/msys2-installer/releases/download/2024-01-13/msys2-x86_64-20240113.exe" # <-- Check msys2.org for latest stable URL if needed
)

# --- Script Body ---

Write-Host "Starting MSYS2 Guided Installation and PATH Setup..." -ForegroundColor Cyan
Write-Host "Installation Directory: $InstallDir"
Write-Host "Installer URL: $InstallerUrl"
Write-Warning "This script requires Administrator privileges."

$ErrorActionPreference = 'Stop' # Exit script on any error
$DownloadPath = Join-Path $env:TEMP "msys2-installer-$($PID).exe"
$MingwBinPath = Join-Path $InstallDir "mingw64\bin"
$MsysShellPath = Join-Path $InstallDir "msys2.exe" # Common path to the shell launcher

# --- Step 1: Download ---
Write-Host "`n[Step 1/5] Downloading MSYS2 installer..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $InstallerUrl -OutFile $DownloadPath
    Write-Host "Download complete: $DownloadPath"
} catch {
    Write-Error "Failed to download MSYS2 installer from '$InstallerUrl': $($_.Exception.Message)"
    exit 1
}

# --- Step 2: Silent Install ---
Write-Host "`n[Step 2/5] Running MSYS2 installer silently to '$InstallDir'..." -ForegroundColor Yellow
Write-Host "(This may take a few moments...)"
try {
    $process = Start-Process -FilePath $DownloadPath -ArgumentList "/S /D=$InstallDir" -Wait -PassThru
    if ($process.ExitCode -ne 0) {
        throw "MSYS2 Installer exited with code $($process.ExitCode)."
    }
    Write-Host "MSYS2 base installation seems complete."
} catch {
    Write-Error "MSYS2 installation failed: $($_.Exception.Message)"
    Remove-Item $DownloadPath -ErrorAction SilentlyContinue
    exit 1
} finally {
     if (Test-Path $DownloadPath) {
        Remove-Item $DownloadPath -Force -ErrorAction SilentlyContinue
        Write-Host "Cleaned up installer file."
     }
}

# --- Step 3: Update System PATH ---
Write-Host "`n[Step 3/5] Attempting to add '$MingwBinPath' to System PATH..." -ForegroundColor Yellow
try {
    $pathScope = 'Machine'
    $currentPath = [Environment]::GetEnvironmentVariable('Path', $pathScope)
    $pathElements = $currentPath -split ';' -ne ''

    if ($pathElements -contains $MingwBinPath) {
        Write-Host "'$MingwBinPath' already found in System PATH. No changes needed."
    } else {
        Write-Host "'$MingwBinPath' not found in System PATH. Appending..."
        $newPath = ($pathElements + $MingwBinPath) -join ';'
        [Environment]::SetEnvironmentVariable('Path', $newPath, $pathScope)
        Write-Host "Successfully appended '$MingwBinPath' to System PATH."
        Write-Warning "Remember to restart terminals/IDEs AFTER all steps are complete."
    }
} catch {
    Write-Error "Failed to update System PATH. Error: $($_.Exception.Message)"
    Write-Warning "You may need to add '$MingwBinPath' to your PATH manually."
}

# --- Step 4: Guided Pacman Update ---
Write-Host "`n[Step 4/5] Manual Step: Update MSYS2 Packages via Pacman" -ForegroundColor Yellow
Write-Host "----------------------------------------------------------"
Write-Host "The script will now open the 'MSYS2 MSYS' shell window."
Write-Host "In THAT WINDOW, you need to run the update commands."
Write-Host ""
Write-Warning "Command to run in MSYS2 shell:"
Write-Host "   pacman -Syu" -ForegroundColor Green
Write-Host ""
Write-Host "IMPORTANT NOTES for 'pacman -Syu':"
Write-Host " - If it asks to close the terminal ([Y/n]), press 'Y', Enter."
Write-Host " - CLOSE the MSYS2 window manually."
Write-Host " - Re-open 'MSYS2 MSYS' from the Start Menu."
Write-Host " - Run this command in the new MSYS2 window: pacman -Su" -ForegroundColor Green
Write-Host " - Repeat 'pacman -Syu' and 'pacman -Su' until MSYS2 says 'there is nothing to do'."
Write-Host ""

# Try to launch the MSYS2 shell
if (Test-Path $MsysShellPath) {
    Write-Host "Launching MSYS2 MSYS shell..."
    Start-Process $MsysShellPath
} else {
    Write-Warning "Could not find '$MsysShellPath'. Please open 'MSYS2 MSYS' manually from the Start Menu."
}

# Pause PowerShell script and wait for user confirmation
Read-Host "--> Perform the 'pacman -Syu'/'pacman -Su' update steps in the MSYS2 window now. Press Enter in THIS (PowerShell) window when ALL MSYS2 updates are complete"
Write-Host "Update step confirmation received."

# --- Step 5: Guided Toolchain Install ---
Write-Host "`n[Step 5/5] Manual Step: Install MinGW Toolchain via Pacman" -ForegroundColor Yellow
Write-Host "------------------------------------------------------------"
Write-Host "Now, you need to install the compiler toolchain."
Write-Host "Make sure the 'MSYS2 MSYS' shell window is still open (or reopen it)."
Write-Host ""
Write-Warning "Command to run in MSYS2 shell:"
Write-Host "   pacman -S --needed base-devel mingw-w64-x86_64-toolchain" -ForegroundColor Green
Write-Host ""
Write-Host "IMPORTANT NOTES for toolchain install:"
Write-Host " - When asked 'Enter a selection (default=all)', just press Enter."
Write-Host " - When asked 'Proceed with installation? [Y/n]', press 'Y' and Enter."
Write-Host ""

# Pause PowerShell script and wait for user confirmation
Read-Host "--> Perform the toolchain installation in the MSYS2 window now. Press Enter in THIS (PowerShell) window when the toolchain installation is complete"
Write-Host "Toolchain installation confirmation received."

# --- Final Instructions ---
Write-Host "`n------------------------------------------------------------------" -ForegroundColor Cyan
Write-Host "Guided MSYS2 Installation and Setup Complete!" -ForegroundColor Cyan
Write-Host "------------------------------------------------------------------" -ForegroundColor Cyan
Write-Host "FINAL MANDATORY STEPS:" -ForegroundColor Yellow
Write-Host "1. CLOSE and REOPEN any PowerShell/CMD terminals and your code editor (VS Code)."
Write-Host "2. Open a NEW terminal and verify the C compiler by typing: gcc --version"
Write-Host "   (You should see version information, not an error)."
Write-Host "3. You should now be able to build Fyne projects."
Write-Host "------------------------------------------------------------------" -ForegroundColor Yellow