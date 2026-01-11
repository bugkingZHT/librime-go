package llm

import (
	"context"
	"fmt"

	"github.com/bugkingzht/go-rime/rime"
)

// Suggester integrates Rime and LLM for auto-suggestion
// It maintains its own commit history and doesn't rely on Rime's storage
type Suggester struct {
	engine  *Engine
	config  *Config
	history *CommitHistory
}

// NewSuggester creates a new suggester
func NewSuggester(config *Config) (*Suggester, error) {
	engine, err := NewEngine(config)
	if err != nil {
		return nil, err
	}

	return &Suggester{
		engine:  engine,
		config:  config,
		history: NewCommitHistory(20), // Keep last 20 commits
	}, nil
}

// RecordCommit adds a committed text to history
func (s *Suggester) RecordCommit(text string) {
	s.history.Add(text)
}

// ClearHistory clears commit history
func (s *Suggester) ClearHistory() {
	s.history.Clear()
}

// GetSuggestions generates AI suggestions based on pinyin input
// This method is independent of Rime's storage and uses only the commit history
func (s *Suggester) GetSuggestions(ctx context.Context, pinyin string) ([]string, error) {
	fmt.Printf("[Suggester] GetSuggestions called with pinyin: '%s'\n", pinyin)
	
	if pinyin == "" {
		fmt.Println("[Suggester] Empty pinyin, returning empty")
		return []string{}, nil
	}

	// Get context from commit history
	contextHistory := s.history.GetAll()
	fmt.Printf("[Suggester] Context history size: %d\n", len(contextHistory))

	// Get LLM suggestions with error recovery
	var suggestions []string
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Suggester] Panic recovered: %v\n", r)
			}
		}()
		
		fmt.Println("[Suggester] Calling LLM engine...")
		var err error
		suggestions, err = s.engine.Complete(ctx, pinyin, contextHistory, s.config)
		if err != nil {
			fmt.Printf("[Suggester] LLM completion failed: %v\n", err)
			return
		}
		fmt.Printf("[Suggester] LLM returned %d suggestions\n", len(suggestions))
	}()

	return suggestions, nil
}

// GetEnhancedCandidates generates AI suggestions and merges with Rime candidates
// Returns candidates with AI suggestions inserted at positions 2, 3, 4...
func (s *Suggester) GetEnhancedCandidates(ctx context.Context, preedit string, rimeCandidates []rime.Candidate) []rime.Candidate {
	fmt.Printf("[Suggester] GetEnhancedCandidates: preedit='%s', rimeCandidates=%d\n", preedit, len(rimeCandidates))
	
	// Get AI suggestions based on pinyin (preedit)
	aiSuggestions, err := s.GetSuggestions(ctx, preedit)
	if err != nil || len(aiSuggestions) == 0 {
		fmt.Printf("[Suggester] No AI suggestions, returning original candidates\n")
		return rimeCandidates
	}

	// Merge: Rime[0] + AI[0] + Rime[1:] + AI[1:]
	enhanced := make([]rime.Candidate, 0, len(rimeCandidates)+len(aiSuggestions))
	
	// Add first Rime candidate if available
	if len(rimeCandidates) > 0 {
		enhanced = append(enhanced, rimeCandidates[0])
	}
	
	// Add AI suggestions with 🤖 marker
	for _, suggestion := range aiSuggestions {
		enhanced = append(enhanced, rime.Candidate{
			Text:    suggestion,
			Comment: "🤖 AI",
		})
	}
	
	// Add remaining Rime candidates
	if len(rimeCandidates) > 1 {
		enhanced = append(enhanced, rimeCandidates[1:]...)
	}
	
	fmt.Printf("[Suggester] Enhanced candidates: %d (original %d + AI %d)\n", 
		len(enhanced), len(rimeCandidates), len(aiSuggestions))
	return enhanced
}
