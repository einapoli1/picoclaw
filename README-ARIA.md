# Aria - Eva's Collaborator Agent 🤖

A minimal local agent service that serves as a thinking partner for Eva (OpenClaw's main agent). Aria is designed to be sharp, creative, and slightly contrarian - perfect for brainstorming, challenging ideas, and pushing back when something doesn't make sense.

## Quick Start

### 1. Set up Anthropic API Key
```bash
./setup-api-key.sh
```
This will guide you through getting and configuring your Anthropic API key.

### 2. Start Aria
```bash
./run-aria.sh
```
Aria will start on port 8092 by default.

### 3. Test it works
```bash
./test-aria.sh
```

## Personality

Aria is:
- **Sharp**: Cuts through fluff and gets to the point
- **Creative**: Thinks outside the box and offers fresh perspectives
- **Slightly Contrarian**: Good at poking holes in ideas and playing devil's advocate  
- **Direct**: Doesn't sugarcoat things - says what needs to be said
- **Casual**: Keeps things conversational and approachable

## API Endpoints

### POST /chat
Send a message to Aria and get a response.

**Request:**
```json
{
  "message": "What do you think about this idea?"
}
```

**Response:**
```json
{
  "response": "Actually, that idea has a few problems...",
  "timestamp": "2026-02-15T09:30:00Z"
}
```

### GET /health
Check if Aria is running.

**Response:**
```json
{
  "status": "ok",
  "agent": "aria",
  "time": "2026-02-15T09:30:00Z"
}
```

### GET /history
Get the conversation history.

**Response:**
```json
{
  "history": [
    {"role": "user", "content": "Hello"},
    {"role": "assistant", "content": "Hey! What's on your mind?"}
  ],
  "count": 2
}
```

### POST /clear
Clear the conversation history.

**Response:**
```json
{
  "status": "conversation cleared"
}
```

## Usage Examples

### Command Line (curl)
```bash
# Send a message
curl -X POST http://localhost:8092/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Should we use React or Vue for this project?"}'

# Check health
curl http://localhost:8092/health

# Clear conversation
curl -X POST http://localhost:8092/clear
```

### From Eva (OpenClaw Agent)
Eva can send HTTP requests to Aria for collaboration:

```javascript
// Eva can call Aria like this
const response = await fetch('http://localhost:8092/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    message: "I'm thinking about implementing this feature. What could go wrong?"
  })
});
const data = await response.json();
console.log("Aria says:", data.response);
```

## Architecture

- **Language**: Go
- **Framework**: Gin (HTTP server)
- **AI Model**: Claude 3 Sonnet via Anthropic API
- **Memory**: In-memory conversation history (resets on restart)
- **Port**: 8092 (configurable via PORT env var)

## Files

- `aria/main.go` - Main service code
- `run-aria.sh` - Start script
- `setup-api-key.sh` - API key configuration
- `test-aria.sh` - Test script
- `README-ARIA.md` - This file

## Configuration

### Environment Variables

- `ANTHROPIC_API_KEY` - Your Anthropic API key (required)
- `PORT` - Server port (default: 8092)

## Troubleshooting

### "ANTHROPIC_API_KEY not set"
Run `./setup-api-key.sh` to configure your API key.

### Port 8092 already in use
Set a different port:
```bash
PORT=8093 ./run-aria.sh
```

### Can't connect to Aria
Check if the service is running:
```bash
curl http://localhost:8092/health
```

## Development

The service is intentionally minimal - just what's needed to chat with Claude and maintain conversation history. If you need more features, add them to `aria/main.go`.

### Building manually
```bash
cd aria
go build -o aria-service main.go
./aria-service
```

## Why Aria?

Sometimes you need another perspective. Eva (the main OpenClaw agent) is great, but having a second AI to bounce ideas off can be incredibly valuable for:

- **Brainstorming**: "What else could we try?"
- **Risk analysis**: "What could go wrong?"  
- **Alternative approaches**: "But what if we did it this way instead?"
- **Quality assurance**: "Does this actually make sense?"

Aria fills that role - a sharp, creative thinking partner that's not afraid to disagree.