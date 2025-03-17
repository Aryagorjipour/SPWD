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

:: Ensure C:\ProgramData\spwd\ exists
if not exist "C:\ProgramData\spwd\" mkdir "C:\ProgramData\spwd"

:: Ensure the database exists
if not exist "C:\ProgramData\spwd\passwords.db" (
    echo Creating database...
    type nul > "C:\ProgramData\spwd\passwords.db"
)

:: Download config.sample.json
powershell -Command "(New-Object System.Net.WebClient).DownloadFile('https://raw.githubusercontent.com/Aryagorjipour/spwd/main/config.sample.json', 'C:\ProgramData\spwd\config.sample.json')"

:: Verify if the file was downloaded
if not exist "C:\ProgramData\spwd\config.sample.json" (
    echo Error: Failed to download config.sample.json
    exit /b 1
)

:: Generate a 32-byte base64 secret key using PowerShell
for /f %%A in ('powershell -Command "[Convert]::ToBase64String((1..32 | % {Get-Random -Maximum 256}))"') do set SECRET_KEY=%%A

:: Create config.json by replacing the placeholder in config.sample.json
powershell -Command "(Get-Content 'C:\ProgramData\spwd\config.sample.json') -replace '\"secret_key\": *\".*\"', '\"secret_key\": \"%SECRET_KEY%\"' | Set-Content 'C:\ProgramData\spwd\config.json'"

echo Installation completed successfully!
echo You can now run 'spwd' from any command prompt.
pause
