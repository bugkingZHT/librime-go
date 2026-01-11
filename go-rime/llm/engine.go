package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Engine wraps OpenAI API for code completion
type Engine struct {
	client *openai.Client
}

// Config holds LLM engine configuration
type Config struct {
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
	Timeout     time.Duration
	Provider    string // "openai" or "dashscope"
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	// Try DashScope first, then fall back to OpenAI
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	baseURL := ""
	model := "qwen-turbo"
	provider := "dashscope"
	
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		baseURL = os.Getenv("OPENAI_BASE_URL")
		model = "gpt-4o-mini"
		provider = "openai"
	}
	
	if provider == "dashscope" && baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	
	return &Config{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       model,
		MaxTokens:   50,
		Temperature: 0.3,
		Timeout:     3 * time.Second,
		Provider:    provider,
	}
}

// NewEngine creates a new LLM engine
func NewEngine(config *Config) (*Engine, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key not found (set DASHSCOPE_API_KEY or OPENAI_API_KEY)")
	}

	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}

	client := openai.NewClient(opts...)

	return &Engine{
		client: &client,
	}, nil
}

// Complete generates code completion using the hole-filling prompt
func (e *Engine) Complete(ctx context.Context, prefix string, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Build the prompt
	prompt := buildPrompt(prefix)

	// Call OpenAI API
	completion, err := e.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(config.Model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		MaxCompletionTokens: openai.Int(int64(config.MaxTokens)),
		Temperature:         openai.Float(config.Temperature),
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	result := completion.Choices[0].Message.Content
	return cleanCompletion(result), nil
}

// buildPrompt constructs the hole-filling prompt
func buildPrompt(prefix string) string {
	return fmt.Sprintf(`You are a HOLE FILLER. Your ONLY task is to complete the missing code part marked by FILL_HERE.
Follow these RULES strictly:
1. Do NOT add any explanations, comments, or code block markers (e.g., `+"```"+`js).
2. Do NOT modify the given prefix and suffix code.
3. Keep the completion concise (≤ 50 chars), fit for input method candidate list.
4. Comply with Go syntax.
5. Output ONLY the completion text, nothing else.

## EXAMPLE
### QUERY
function sum(a, b) {
  {{FILL_HERE}}
}
### CORRECT COMPLETION
return a + b;

## YOUR TASK
<QUERY>
%s{{FILL_HERE}}
</QUERY>

Your completion (only the missing part):`, prefix)
}

// cleanCompletion removes unnecessary formatting from the completion
func cleanCompletion(completion string) string {
	// Remove leading/trailing whitespace
	result := strings.TrimSpace(completion)

	// Remove code block markers if present
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")

	// Remove common language identifiers
	for _, lang := range []string{"go", "golang", "js", "javascript", "python", "java"} {
		result = strings.TrimPrefix(result, lang)
	}

	// Trim again after removing markers
	result = strings.TrimSpace(result)

	// Limit length to 50 characters
	if len([]rune(result)) > 50 {
		runes := []rune(result)
		result = string(runes[:50])
	}

	return result
}
