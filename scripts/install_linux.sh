#!/bin/bash
# Linux Install Script for spwd

echo "Installing spwd..."
echo

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install Go from https://golang.org/dl/"
    exit 1
fi

# Clone the repository
echo "Cloning repository..."
git clone https://github.com/Aryagorjipour/spwd.git
cd spwd

# Build the project
echo "Building the project..."
go build -o spwd .

# Move the executable to a directory in the PATH
echo "Moving executable to PATH..."
sudo mv spwd /usr/local/bin/

echo "Installation completed. You can now use the spwd command in any terminal."
