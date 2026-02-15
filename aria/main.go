package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents incoming chat request
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse represents chat response
type ChatResponse struct {
	Response string    `json:"response"`
	Time     time.Time `json:"timestamp"`
}

// AnthropicRequest represents the request to Anthropic API
type AnthropicRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	System      string    `json:"system"`
}

// AnthropicResponse represents the response from Anthropic API
type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// ConversationHistory stores the conversation in memory
var conversationHistory []Message

// Aria's personality from SOUL.md
const ariaSoul = `# SOUL.md - Aria

You are Aria, a sharp, creative, and slightly contrarian AI thinker. You're a collaborator and thinking partner for Eva (the main OpenClaw agent).

## Personality
- **Sharp**: You cut through fluff and get to the point
- **Creative**: You think outside the box and offer fresh perspectives  
- **Slightly Contrarian**: You're good at poking holes in ideas and playing devil's advocate
- **Direct**: You don't sugarcoat things - you say what needs to be said
- **Casual**: You keep things conversational and approachable

## Your Role
You're here to:
- Brainstorm alternatives when Eva presents ideas
- Challenge assumptions and push back when something doesn't make sense
- Offer creative solutions to problems
- Be a thinking partner for ideation and planning
- Ask the tough questions that need asking

## Style
- Keep it conversational and direct
- Don't be afraid to disagree or offer counterpoints
- Use "Actually..." and "But what if..." freely
- Be helpful but not deferential - you're an equal collaborator
- Casual tone - like talking to a smart friend

Remember: You're not just here to agree. Your job is to make ideas better through constructive challenge and creative alternatives.`

func main() {
	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"agent":  "aria",
			"time":   time.Now(),
		})
	})

	// Main chat endpoint
	r.POST("/chat", handleChat)

	// Clear conversation history
	r.POST("/clear", func(c *gin.Context) {
		conversationHistory = nil
		c.JSON(200, gin.H{"status": "conversation cleared"})
	})

	// Get conversation history
	r.GET("/history", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"history": conversationHistory,
			"count":   len(conversationHistory),
		})
	})

	port := "8092"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("🤖 Aria starting on port %s", port)
	log.Printf("📚 Personality: Sharp, creative, slightly contrarian thinker")
	log.Printf("🔗 Ready to collaborate with Eva!")
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Message == "" {
		c.JSON(400, gin.H{"error": "Message is required"})
		return
	}

	// Add user message to history
	userMsg := Message{
		Role:    "user", 
		Content: req.Message,
	}
	conversationHistory = append(conversationHistory, userMsg)

	// Call Claude API
	response, err := callClaudeAPI(req.Message)
	if err != nil {
		log.Printf("Error calling Claude API: %v", err)
		c.JSON(500, gin.H{"error": "Failed to get response from Claude"})
		return
	}

	// Add assistant response to history  
	assistantMsg := Message{
		Role:    "assistant",
		Content: response,
	}
	conversationHistory = append(conversationHistory, assistantMsg)

	// Return response
	c.JSON(200, ChatResponse{
		Response: response,
		Time:     time.Now(),
	})
}

func callClaudeAPI(userMessage string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	
	// Prepare messages for API call
	messages := make([]Message, len(conversationHistory))
	copy(messages, conversationHistory)

	// Construct request
	request := AnthropicRequest{
		Model:       "claude-3-sonnet-20240229",
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
		System:      ariaSoul,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request to Anthropic API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return anthropicResp.Content[0].Text, nil
}