#!/bin/bash
set -e

echo "Installing twigo..."
cd "$(dirname "$0")"

# Create ~/.local/bin if it doesn't exist
mkdir -p ~/.local/bin

# Install executable
install -m 755 twigo ~/.local/bin/

# Create systemd user directory if it doesn't exist
mkdir -p ~/.config/systemd/user

# Create temporary service file for editing
tmp_service=$(mktemp)
cp ../twigo.service "$tmp_service"

# Open editor for token configuration
echo "Please configure your WEBHOOK_TOKEN and GOTIFY_TOKEN in the editor"
${EDITOR:-nano} "$tmp_service"

# Check if tokens were modified
if grep -q "webhook-token\|gotify-token" "$tmp_service"; then
    echo "Error: Tokens not configured. Aborting installation."
    rm "$tmp_service"
    exit 1
fi

# Install service file in user systemd directory
install -m 644 "$tmp_service" ~/.config/systemd/user/twigo.service
rm "$tmp_service"

# Reload systemd user daemon and enable service
systemctl --user daemon-reload
systemctl --user enable twigo.service
systemctl --user restart twigo.service

echo "Installation complete!"
echo "You can check the service status with: systemctl --user status twigo" 