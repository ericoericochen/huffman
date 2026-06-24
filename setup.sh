#!/bin/bash

echo "Building hfcode..."
echo ""

# get the directory of the script
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$DIR/cli"
BUILD_DIR="$DIR/build"
BUILD_OUTPUT_PATH="$BUILD_DIR/hfcode"

# build the cli binary
go build -o $BUILD_OUTPUT_PATH $CLI_DIR

# print export command for user to add command to their PATH
echo "✅ Done!"
echo "To add hfcode to your PATH, add the command to your shell (~/.zshrc, ~/.bashrc)"
echo "export PATH=$BUILD_DIR:\$PATH"