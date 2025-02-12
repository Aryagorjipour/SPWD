#!/bin/bash
# Linux Install Script for spwd

set -e  # Exit if any command fails

echo "Installing spwd..."
echo

# Define installation directory
INSTALL_DIR="/usr/local/bin"
TMP_FILE="/tmp/spwd"
sudo rm -f "$TMP_FILE"

# Get latest release URL
LATEST_RELEASE=$(curl -s https://api.github.com/repos/Aryagorjipour/spwd/releases/latest | grep "browser_download_url" | cut -d '"' -f 4 | grep "linux")
if [[ -z "$LATEST_RELEASE" ]]; then
    echo "Error: Could not fetch the latest release. Please check the repository."
    exit 1
fi

echo "Downloading latest version from: $LATEST_RELEASE"
curl -L -o "$TMP_FILE" "$LATEST_RELEASE"
chmod +x "$TMP_FILE"

# Move to /usr/local/bin
echo "Moving executable to $INSTALL_DIR/"
sudo mv "$TMP_FILE" "$INSTALL_DIR/spwd"

# Ensure config.json exists next to spwd
CONFIG_PATH="$INSTALL_DIR/config.json"
CONFIG_SAMPLE_URL="https://raw.githubusercontent.com/Aryagorjipour/spwd/main/config.sample.json"

if [ ! -f "$CONFIG_PATH" ]; then
    echo "Generating config.json..."
    sudo curl -sSL -o "$INSTALL_DIR/config.sample.json" "$CONFIG_SAMPLE_URL"

    if [ ! -f "$INSTALL_DIR/config.sample.json" ]; then
        echo "Error: Failed to download config.sample.json"
        exit 1
    fi

    SECRET_KEY=$(head -c 32 /dev/urandom | base64 | tr -d '\n')
    sudo jq --arg key "$SECRET_KEY" '.secret_key = $key' "$INSTALL_DIR/config.sample.json" | sudo tee "$CONFIG_PATH" > /dev/null
    sudo rm -f "$INSTALL_DIR/config.sample.json"
fi

echo "Installation completed successfully!"
echo "You can now run 'spwd' from any terminal."
