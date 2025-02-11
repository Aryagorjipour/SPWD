@echo off
:: Windows Install Script for spwd

echo Installing spwd...
echo.

:: Check if Go is installed (not required for pre-built binary)
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Go is not installed. This script will install the pre-built binary instead.
)

:: Define the latest release URL from GitHub
for /f "tokens=*" %%i in ('powershell -Command "(Invoke-WebRequest -Uri 'https://api.github.com/repos/Aryagorjipour/spwd/releases/latest' -UseBasicParsing).Content | ConvertFrom-Json | Select-Object -ExpandProperty assets | Where-Object { $_.browser_download_url -match 'windows' } | Select-Object -ExpandProperty browser_download_url"') do set LATEST_RELEASE=%%i

if "%LATEST_RELEASE%"=="" (
    echo Error: Could not fetch the latest release. Please check the repository.
    exit /b 1
)

echo Downloading latest version from: %LATEST_RELEASE%
powershell -Command "(New-Object System.Net.WebClient).DownloadFile('%LATEST_RELEASE%', 'spwd.exe')"

:: Move the file to a location in the PATH
echo Moving executable to C:\Windows\System32
move spwd.exe C:\Windows\System32 >nul 2>&1

echo Installation completed successfully!
echo You can now run 'spwd' from any command prompt.
pause
