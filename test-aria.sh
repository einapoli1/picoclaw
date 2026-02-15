#!/bin/bash

# Test script for Aria agent

set -e

ARIA_URL="http://localhost:8092"

echo "🧪 Testing Aria Agent..."

# Test health endpoint
echo "📊 Checking health..."
curl -s "$ARIA_URL/health" | jq . || echo "Health check failed"

echo -e "\n💬 Testing chat..."
# Test chat endpoint
curl -s -X POST "$ARIA_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hey Aria, what do you think about building a minimal agent service instead of using a complex framework?"}' \
  | jq . || echo "Chat test failed"

echo -e "\n📚 Checking history..."
curl -s "$ARIA_URL/history" | jq . || echo "History check failed"

echo -e "\n✅ Test complete!"