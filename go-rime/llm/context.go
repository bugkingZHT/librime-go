package llm

import (
	"fmt"
	"sync"
)

// CommitHistory stores the history of committed text for context
type CommitHistory struct {
	mu      sync.RWMutex
	entries []string
	maxSize int
}

// NewCommitHistory creates a new commit history cache
func NewCommitHistory(maxSize int) *CommitHistory {
	if maxSize <= 0 {
		maxSize = 20 // Default: keep last 20 commits
	}
	return &CommitHistory{
		entries: make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds a new committed text to history
func (ch *CommitHistory) Add(text string) {
	if text == "" {
		return
	}

	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.entries = append(ch.entries, text)
	
	// Keep only the last maxSize entries
	if len(ch.entries) > ch.maxSize {
		ch.entries = ch.entries[len(ch.entries)-ch.maxSize:]
	}
	
	fmt.Printf("[CommitHistory] Added: '%s', Total: %d\n", text, len(ch.entries))
}

// GetRecent returns the last N committed texts
func (ch *CommitHistory) GetRecent(n int) []string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if n <= 0 || len(ch.entries) == 0 {
		return []string{}
	}

	start := len(ch.entries) - n
	if start < 0 {
		start = 0
	}

	result := make([]string, len(ch.entries)-start)
	copy(result, ch.entries[start:])
	return result
}

// GetAll returns all committed texts
func (ch *CommitHistory) GetAll() []string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	result := make([]string, len(ch.entries))
	copy(result, ch.entries)
	return result
}

// Clear clears all history
func (ch *CommitHistory) Clear() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.entries = ch.entries[:0]
	fmt.Println("[CommitHistory] Cleared")
}

// Size returns the current number of entries
func (ch *CommitHistory) Size() int {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.entries)
}
