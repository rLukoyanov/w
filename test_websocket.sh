#!/bin/bash

# WebSocket test script

# Create a test user and get token
echo "📝 Creating test user..."
RESPONSE=$(curl -s -X POST http://localhost:8090/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "wstest",
    "email": "wstest@test.com",
    "password": "password123"
  }')

TOKEN=$(echo $RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token. Response:"
  echo $RESPONSE | jq
  
  # Try logging in instead
  echo ""
  echo "🔑 Trying to login..."
  RESPONSE=$(curl -s -X POST http://localhost:8090/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "email": "wstest@test.com",
      "password": "password123"
    }')
  
  TOKEN=$(echo $RESPONSE | jq -r '.token')
  
  if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Failed to login. Response:"
    echo $RESPONSE | jq
    exit 1
  fi
fi

echo "✅ Got token: ${TOKEN:0:50}..."
echo ""

# Create a server
echo "🏢 Creating server..."
SERVER_RESPONSE=$(curl -s -X POST http://localhost:8090/api/servers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Test Server"}')

SERVER_ID=$(echo $SERVER_RESPONSE | jq -r '.id')
echo "✅ Server created: $SERVER_ID"
echo ""

# Create a channel
echo "📺 Creating channel..."
CHANNEL_RESPONSE=$(curl -s -X POST http://localhost:8090/api/servers/$SERVER_ID/channels \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "general", "type": "text"}')

CHANNEL_ID=$(echo $CHANNEL_RESPONSE | jq -r '.id')
echo "✅ Channel created: $CHANNEL_ID"
echo ""

echo "🔌 WebSocket endpoint ready: ws://localhost:8090/ws?token=$TOKEN"
echo ""
echo "You can test with a WebSocket client like wscat:"
echo "  npm install -g wscat"
echo "  wscat -c 'ws://localhost:8090/ws?token=$TOKEN'"
echo ""
echo "Or send a test message:"
echo "  {\"type\":\"MESSAGE_CREATE\",\"data\":{\"channel_id\":\"$CHANNEL_ID\",\"content\":\"Hello WebSocket!\"}}"
