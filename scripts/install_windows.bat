@echo off
REM Windows Install Script for spwd

echo Installing spwd...
echo.

REM Check if Go is installed
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Go is not installed. Please install Go from https://golang.org/dl/.
    exit /b
)

REM Clone the repository
echo Cloning repository...
git clone https://github.com/Aryagorjipour/spwd.git
cd spwd

REM Build the project
echo Building the project...
go build -o spwd.exe .

REM Move the executable to a directory in the PATH
echo Moving executable to PATH...
move spwd.exe C:\Windows\System32

echo Installation completed. You can now use the spwd command in any terminal.

pause
