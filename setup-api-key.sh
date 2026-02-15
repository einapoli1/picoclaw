#!/bin/bash

# Setup script to configure ANTHROPIC_API_KEY for Aria

set -e

echo "🔧 Setting up Anthropic API Key for Aria"
echo ""

# Check if key is already set
if [ ! -z "$ANTHROPIC_API_KEY" ]; then
    echo "✅ ANTHROPIC_API_KEY is already configured in your environment"
    echo ""
    echo "Current key starts with: ${ANTHROPIC_API_KEY:0:8}..."
    echo ""
    read -p "Do you want to use a different key? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Using existing key."
        exit 0
    fi
fi

echo "🔑 You need an Anthropic API key to run Aria."
echo "📋 Get one at: https://console.anthropic.com/keys"
echo ""

# Prompt for API key
read -p "Enter your Anthropic API key: " api_key

if [ -z "$api_key" ]; then
    echo "❌ No API key provided. Exiting."
    exit 1
fi

# Validate key format (basic check)
if [[ ! $api_key =~ ^sk-ant-api03- ]]; then
    echo "⚠️  Warning: The key doesn't look like a valid Anthropic API key"
    echo "   Anthropic keys typically start with 'sk-ant-api03-'"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Determine shell config file
if [ -f ~/.zshrc ]; then
    SHELL_CONFIG=~/.zshrc
    SHELL_NAME="zsh"
elif [ -f ~/.bashrc ]; then
    SHELL_CONFIG=~/.bashrc
    SHELL_NAME="bash"
elif [ -f ~/.bash_profile ]; then
    SHELL_CONFIG=~/.bash_profile
    SHELL_NAME="bash"
else
    echo "❌ Could not find ~/.zshrc, ~/.bashrc, or ~/.bash_profile"
    echo "Please manually add this line to your shell configuration:"
    echo "export ANTHROPIC_API_KEY='$api_key'"
    exit 1
fi

# Check if already in config
if grep -q "ANTHROPIC_API_KEY" "$SHELL_CONFIG" 2>/dev/null; then
    echo "⚠️  ANTHROPIC_API_KEY is already in $SHELL_CONFIG"
    echo ""
    read -p "Replace it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Remove existing line
        sed -i.bak '/ANTHROPIC_API_KEY/d' "$SHELL_CONFIG"
        echo "🗑️  Removed old ANTHROPIC_API_KEY from $SHELL_CONFIG"
    else
        echo "Keeping existing configuration."
        exit 0
    fi
fi

# Add to shell config
echo "export ANTHROPIC_API_KEY='$api_key'" >> "$SHELL_CONFIG"
echo "✅ Added ANTHROPIC_API_KEY to $SHELL_CONFIG"
echo ""

# Set for current session
export ANTHROPIC_API_KEY="$api_key"
echo "🔄 Set ANTHROPIC_API_KEY for current session"
echo ""

echo "🎉 Setup complete!"
echo ""
echo "To use the key in new terminal sessions:"
echo "  source $SHELL_CONFIG"
echo "  # or restart your terminal"
echo ""
echo "Ready to start Aria! Run:"
echo "  ./run-aria.sh"