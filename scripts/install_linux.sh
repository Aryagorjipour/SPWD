#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.

echo "Installing spwd..."
echo

# Check if Go is installed (not required for pre-built binary)
if ! command -v go &> /dev/null; then
    echo "Warning: Go is not installed. This script will install the pre-built binary instead."
fi

# Define the latest release URL from GitHub
LATEST_RELEASE=$(curl -s https://api.github.com/repos/Aryagorjipour/spwd/releases/latest | grep "browser_download_url" | cut -d '"' -f 4 | grep "linux")

if [[ -z "$LATEST_RELEASE" ]]; then
    echo "Error: Could not fetch the latest release. Please check the repository."
    exit 1
fi

echo "Downloading latest version from: $LATEST_RELEASE"
curl -L -o spwd "$LATEST_RELEASE"

# Make it executable
chmod +x spwd

# Move to /usr/local/bin for system-wide usage
echo "Moving executable to /usr/local/bin/"
sudo mv spwd /usr/local/bin/spwd

echo "Installation completed successfully!"
echo "You can now run 'spwd' from any terminal."
