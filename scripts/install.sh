#!/bin/bash

BINDIR="$HOME/.bin"
mkdir -p "$BINDIR"
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    DOWNLOAD_URL="https://github.com/v1ctorio/termpet/releases/latest/download/termpet_Linux-x86_64"
    elif [[ "$ARCH" == "aarch64" ]]; then
    DOWNLOAD_URL="https://github.com/v1ctorio/termpet/releases/latest/download/termpet_Linux-arm64"
else
    echo "Unsupported arch"
    exit 1
fi

EXEPATH="$BINDIR/termpet"


echo "Downloading termpet to $EXEPATH"
curl -L "$DOWNLOAD_URL" -o "$EXEPATH"
chmod +x "$EXEPATH"

if [[ ":$PATH:" != *":$BINDIR:"* ]]; then
    if [[ -f "$HOME/.zshrc" ]]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [[ -f "$HOME/.bashrc" ]]; then
        SHELL_CONFIG="$HOME/.bashrc"
    else
        echo "Could not find shell configuration file"
        exit 1
    fi

    echo "export PATH=\"$BINDIR:\$PATH\"" >> "$SHELL_CONFIG"
    echo "Added $BINDIR to PATH in $SHELL_CONFIG"
fi

echo "Termpet successfully installed! Restart your terminal or run 'source $SHELL_CONFIG' to use it."
echo "Just write 'termpet' to get started!"