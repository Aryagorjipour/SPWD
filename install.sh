#!/bin/bash
# Linux Install Script for spwd

set -e  # Exit immediately if a command exits with a non-zero status.

echo "Installing spwd..."
echo

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

# Ensure /etc/spwd/ directory exists
sudo mkdir -p /etc/spwd

# Download config.sample.json from GitHub and create config.json
CONFIG_URL="https://raw.githubusercontent.com/Aryagorjipour/spwd/main/config.sample.json"

if [ ! -f "/etc/spwd/config.json" ]; then
    echo "Generating config.json..."
    sudo curl -sSL -o /etc/spwd/config.sample.json "$CONFIG_URL"

    if [ ! -f "/etc/spwd/config.sample.json" ]; then
        echo "Error: Failed to download config.sample.json"
        exit 1
    fi

    SECRET_KEY=$(head -c 32 /dev/urandom | base64)
    sudo cp /etc/spwd/config.sample.json /etc/spwd/config.json
    sudo sed -i "s/GENERATE_ON_INSTALL/$SECRET_KEY/" /etc/spwd/config.json
fi

echo "Installation completed successfully!"
echo "You can now run 'spwd' from any terminal."
