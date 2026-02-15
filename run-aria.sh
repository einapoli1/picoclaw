#!/bin/bash

# Aria Agent Runner Script
# A minimal local collaborator agent for Eva

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤖 Starting Aria - Eva's Collaborator Agent${NC}"

# Check for ANTHROPIC_API_KEY
if [ -z "$ANTHROPIC_API_KEY" ]; then
    echo -e "${RED}❌ ANTHROPIC_API_KEY environment variable is not set${NC}"
    echo -e "${YELLOW}💡 You can set it by:${NC}"
    echo "   export ANTHROPIC_API_KEY='your-api-key-here'"
    echo "   or add it to your ~/.zshrc or ~/.bashrc"
    echo ""
    echo -e "${YELLOW}🔗 Get your API key at: https://console.anthropic.com/keys${NC}"
    exit 1
fi

# Change to aria directory
cd "$(dirname "$0")/aria"

# Check if Go modules are initialized
if [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}📦 Installing Go dependencies...${NC}"
    go mod tidy
fi

echo -e "${GREEN}🚀 Starting Aria on port 8092...${NC}"
echo -e "${BLUE}📋 Personality: Sharp, creative, slightly contrarian thinker${NC}"
echo -e "${BLUE}🎯 Role: Collaborative thinking partner for Eva${NC}"
echo ""
echo -e "${YELLOW}Available endpoints:${NC}"
echo "  POST http://localhost:8092/chat     - Send a message to Aria"  
echo "  GET  http://localhost:8092/health   - Health check"
echo "  GET  http://localhost:8092/history  - Conversation history"
echo "  POST http://localhost:8092/clear    - Clear conversation"
echo ""
echo -e "${GREEN}Press Ctrl+C to stop${NC}"

# Start the Go service
exec go run main.go