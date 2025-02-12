@echo off
:: Windows Install Script for spwd

echo Installing spwd...
echo.

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

:: Generate config.json if it doesn't exist
if not exist "C:\ProgramData\spwd\" mkdir "C:\ProgramData\spwd"
if not exist "C:\ProgramData\spwd\config.json" (
    echo Generating config.json...
    powershell -Command "$secretKey = [System.Convert]::ToBase64String((1..32 | % { Get-Random -Minimum 0 -Maximum 256 }))"
    copy config.sample.json C:\ProgramData\spwd\config.json
    powershell -Command "(Get-Content C:\ProgramData\spwd\config.json) -replace 'GENERATE_ON_INSTALL', $secretKey | Set-Content C:\ProgramData\spwd\config.json"
)

echo Installation completed successfully!
echo You can now run 'spwd' from any command prompt.
pause
