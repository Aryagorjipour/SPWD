#!/bin/bash
# Linux Install Script for spwd

set -e  # Exit immediately if a command fails

echo "Installing spwd..."
echo

# Define the latest release URL from GitHub
LATEST_RELEASE=$(curl -s https://api.github.com/repos/Aryagorjipour/spwd/releases/latest | grep "browser_download_url" | cut -d '"' -f 4 | grep "linux")

if [[ -z "$LATEST_RELEASE" ]]; then
    echo "Error: Could not fetch the latest release. Please check the repository."
    exit 1
fi

echo "Downloading latest version from: $LATEST_RELEASE"

# Ensure download location is writable
TMP_FILE="/tmp/spwd"
sudo rm -f "$TMP_FILE"  # Remove any existing file to prevent conflicts

# Download the binary to /tmp/
curl -L -o "$TMP_FILE" "$LATEST_RELEASE"

# Verify the file was downloaded
if [[ ! -f "$TMP_FILE" ]]; then
    echo "Error: Failed to download spwd binary."
    exit 1
fi

# Make it executable
chmod +x "$TMP_FILE"

# Move to /usr/local/bin with sudo
echo "Moving executable to /usr/local/bin/"
sudo mv "$TMP_FILE" /usr/local/bin/spwd

# Ensure /etc/spwd/ directory exists
sudo mkdir -p /etc/spwd

# Define the URL for config.sample.json
CONFIG_URL="https://raw.githubusercontent.com/Aryagorjipour/spwd/main/config.sample.json"

# Download config.sample.json and create config.json
if [ ! -f "/etc/spwd/config.json" ]; then
    echo "Generating config.json..."
    sudo curl -sSL -o /etc/spwd/config.sample.json "$CONFIG_URL"

    if [ ! -f "/etc/spwd/config.sample.json" ]; then
        echo "Error: Failed to download config.sample.json"
        exit 1
    fi

    # Generate a 32-byte secret key and remove newlines
    SECRET_KEY=$(head -c 32 /dev/urandom | base64 | tr -d '\n')

    # Use jq to modify JSON properly
    sudo jq --arg key "$SECRET_KEY" '.secret_key = $key' /etc/spwd/config.sample.json | sudo tee /etc/spwd/config.json > /dev/null
fi

echo "Installation completed successfully!"
echo "You can now run 'spwd' from any terminal."
