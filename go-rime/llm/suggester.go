package llm

import (
	"context"
	"fmt"

	"github.com/bugkingzht/go-rime/rime"
)

// Suggester integrates Rime and LLM for auto-suggestion
type Suggester struct {
	engine *Engine
	config *Config
}

// NewSuggester creates a new suggester
func NewSuggester(config *Config) (*Suggester, error) {
	engine, err := NewEngine(config)
	if err != nil {
		return nil, err
	}

	return &Suggester{
		engine: engine,
		config: config,
	}, nil
}

// GetEnhancedCandidates enhances Rime candidates with LLM suggestions
// Returns candidates with LLM suggestion at index 1 (second position)
func (s *Suggester) GetEnhancedCandidates(ctx context.Context, rimeEngine *rime.Rime, session *rime.Session) ([]rime.Candidate, error) {
	fmt.Println("[LLM] GetEnhancedCandidates called")
	
	// Get original Rime context
	rimeCtx, err := rimeEngine.GetContext(session)
	if err != nil || rimeCtx == nil {
		fmt.Printf("[LLM] Failed to get Rime context: %v\n", err)
		// Return empty on error, caller will use original context
		return []rime.Candidate{}, nil
	}

	// Get original candidates
	candidates := rimeCtx.Menu.Candidates
	fmt.Printf("[LLM] Original candidates count: %d\n", len(candidates))
	if len(candidates) == 0 {
		fmt.Println("[LLM] No candidates, returning empty")
		return candidates, nil
	}

	// Build prefix from preedit + first candidate
	prefix := buildPrefix(rimeCtx)
	fmt.Printf("[LLM] Prefix built: '%s'\n", prefix)
	if prefix == "" {
		fmt.Println("[LLM] Empty prefix, returning original candidates")
		return candidates, nil
	}

	// Get LLM completion with error recovery
	completion := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[LLM] Panic recovered: %v\n", r)
			}
		}()
		
		fmt.Println("[LLM] Calling LLM API...")
		var err error
		completion, err = s.engine.Complete(ctx, prefix, s.config)
		if err != nil {
			// If LLM fails, just return original candidates
			fmt.Printf("[LLM] Completion failed: %v\n", err)
			return
		}
		fmt.Printf("[LLM] Completion received: '%s'\n", completion)
	}()

	if completion == "" {
		fmt.Println("[LLM] Empty completion, returning original candidates")
		return candidates, nil
	}

	// Create LLM candidate
	llmCandidate := rime.Candidate{
		Text:    completion,
		Comment: "🤖 AI",
	}

	// Insert LLM suggestion at second position (index 1)
	enhancedCandidates := make([]rime.Candidate, 0, len(candidates)+1)
	if len(candidates) > 0 {
		enhancedCandidates = append(enhancedCandidates, candidates[0])
	}
	enhancedCandidates = append(enhancedCandidates, llmCandidate)
	if len(candidates) > 1 {
		enhancedCandidates = append(enhancedCandidates, candidates[1:]...)
	}

	fmt.Printf("[LLM] Returning %d enhanced candidates\n", len(enhancedCandidates))
	return enhancedCandidates, nil
}

// buildPrefix constructs the completion prefix from Rime context
func buildPrefix(ctx *rime.Context) string {
	// Only use the first candidate (most likely word) as context
	// DO NOT include preedit (pinyin) to avoid confusing the LLM
	if len(ctx.Menu.Candidates) > 0 {
		return ctx.Menu.Candidates[0].Text
	}
	return ""
}
