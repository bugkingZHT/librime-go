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

// Complete generates suggestions based on pinyin input and context
func (e *Engine) Complete(ctx context.Context, pinyin string, contextHistory []string, config *Config) ([]string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Build the prompt
	prompt := buildPromptFromPinyin(pinyin, contextHistory)
	fmt.Printf("[LLM] Prompt:\n%s\n", prompt)

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
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
	}

	result := completion.Choices[0].Message.Content
	return parseSuggestions(result), nil
}

// buildPromptFromPinyin constructs prompt from pinyin input and context
func buildPromptFromPinyin(pinyin string, contextHistory []string) string {
	contextStr := ""
	if len(contextHistory) > 0 {
		// Use last 5 commits as context
		start := len(contextHistory) - 5
		if start < 0 {
			start = 0
		}
		recentContext := contextHistory[start:]
		for _, text := range recentContext {
			contextStr += text
		}
	}

	return fmt.Sprintf(`You are a Chinese input method AI assistant. Based on pinyin input and context, generate appropriate Chinese word/phrase suggestions.

RULES:
1. Output 1-3 Chinese suggestions ONLY, one per line
2. Each suggestion should be ≤ 10 characters
3. Prioritize common words/phrases that match the pinyin
4. Consider the context to generate contextually relevant suggestions
5. DO NOT include pinyin, explanations, or any other text
6. DO NOT use markdown, numbers, or bullets

EXAMPLE:
<CONTEXT>今天天气很</CONTEXT>
<PINYIN>hao</PINYIN>
OUTPUT:
好
好的

<CONTEXT>%s</CONTEXT>
<PINYIN>%s</PINYIN>
OUTPUT:`, contextStr, pinyin)
}

// parseSuggestions extracts suggestions from LLM response
func parseSuggestions(response string) []string {
	lines := strings.Split(strings.TrimSpace(response), "\n")
	suggestions := make([]string, 0, 3)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Clean up the line
		line = cleanSuggestion(line)
		if line != "" && len([]rune(line)) <= 10 {
			suggestions = append(suggestions, line)
			if len(suggestions) >= 3 {
				break
			}
		}
	}
	
	fmt.Printf("[LLM] Parsed %d suggestions from response\n", len(suggestions))
	return suggestions
}

// cleanSuggestion cleans a single suggestion
func cleanSuggestion(s string) string {
	// Remove common prefixes/markers
	s = strings.TrimPrefix(s, "-")
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimPrefix(s, "•")
	
	// Remove numbers at start (e.g., "1. ", "1) ")
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '.' || r == ')' || r == ' ' {
			continue
		}
		s = s[i:]
		break
	}
	
	s = strings.TrimSpace(s)
	
	// Remove code block markers
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	
	return strings.TrimSpace(s)
}
