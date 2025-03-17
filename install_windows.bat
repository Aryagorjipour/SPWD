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

:: Move the file to a location in the PATH (System32)
set INSTALL_DIR=C:\Windows\System32
echo Moving executable to %INSTALL_DIR%
move spwd.exe %INSTALL_DIR% >nul 2>&1

:: Download config.sample.json next to spwd.exe
echo Downloading config.sample.json...
powershell -Command "(New-Object System.Net.WebClient).DownloadFile('https://raw.githubusercontent.com/Aryagorjipour/spwd/main/config.sample.json', '%INSTALL_DIR%\config.sample.json')"

:: Verify if the file was downloaded
if not exist "%INSTALL_DIR%\config.sample.json" (
    echo Error: Failed to download config.sample.json
    exit /b 1
)

:: Rename config.sample.json to config.json
echo Renaming config.sample.json to config.json...
rename "%INSTALL_DIR%\config.sample.json" "config.json"

:: Replace "secret_key": "GENERATE_ON_INSTALL" with a static key
echo Updating config.json with static secret key...
powershell -Command "(Get-Content '%INSTALL_DIR%\config.json') -replace '\"secret_key\": *\"GENERATE_ON_INSTALL\"', '\"secret_key\": \"STATIC_SECRET_KEY_HERE\"' | Set-Content '%INSTALL_DIR%\config.json'"

echo Installation completed successfully!
echo You can now run 'spwd' from any command prompt.
pause
