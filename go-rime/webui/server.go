package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bugkingzht/go-rime/llm"
	"github.com/bugkingzht/go-rime/rime"
)

// Server manages the web UI for Rime testing
type Server struct {
	rimeEngine   *rime.Rime
	llmSuggester *llm.Suggester
	session      *rime.Session
	mu           sync.Mutex
	enableLLM    bool
}

// NewServer creates a new web UI server
func NewServer(rimeEngine *rime.Rime, llmSuggester *llm.Suggester) *Server {
	session := rimeEngine.CreateSession()
	return &Server{
		rimeEngine:   rimeEngine,
		llmSuggester: llmSuggester,
		session:      session,
		enableLLM:    llmSuggester != nil,
	}
}

// InputRequest represents input from user
type InputRequest struct {
	Text string `json:"text"`
}

// CandidateResponse represents candidates to display
type CandidateResponse struct {
	Preedit    string             `json:"preedit"`
	Candidates []CandidateDisplay `json:"candidates"`
	Success    bool               `json:"success"`
	Error      string             `json:"error,omitempty"`
}

// CandidateDisplay represents a single candidate
type CandidateDisplay struct {
	Text    string `json:"text"`
	Comment string `json:"comment"`
	IsAI    bool   `json:"isAI"`
}

// CommitResponse represents commit result
type CommitResponse struct {
	Text    string `json:"text"`
	Success bool   `json:"success"`
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/input", s.handleInput)
	http.HandleFunc("/api/select", s.handleSelect)
	http.HandleFunc("/api/clear", s.handleClear)
	http.HandleFunc("/api/toggle-llm", s.handleToggleLLM)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting web UI server at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlTemplate))
}

func (s *Server) handleInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var req InputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request")
		return
	}

	log.Printf("[WebUI] Input: '%s'\n", req.Text)

	// Send input to Rime
	if !s.rimeEngine.SimulateKeySequence(s.session, req.Text) {
		s.sendError(w, "Failed to process input")
		return
	}

	// Get context
	ctx, err := s.rimeEngine.GetContext(s.session)
	if err != nil || ctx == nil {
		s.sendError(w, "Failed to get context")
		return
	}

	// Get candidates
	candidates := ctx.Menu.Candidates
	preedit := ctx.Composition.Preedit

	log.Printf("[WebUI] Preedit: '%s', Candidates: %d\n", preedit, len(candidates))

	// Enhance with LLM if enabled
	if s.enableLLM && s.llmSuggester != nil && preedit != "" {
		enhanced := s.llmSuggester.GetEnhancedCandidates(
			context.Background(),
			preedit,
			candidates,
		)
		candidates = enhanced
		log.Printf("[WebUI] Enhanced candidates: %d\n", len(candidates))
	}

	// Build response
	resp := CandidateResponse{
		Preedit:    preedit,
		Candidates: make([]CandidateDisplay, 0, len(candidates)),
		Success:    true,
	}

	for _, cand := range candidates {
		resp.Candidates = append(resp.Candidates, CandidateDisplay{
			Text:    cand.Text,
			Comment: cand.Comment,
			IsAI:    cand.Comment == "🤖 AI",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleSelect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var req struct {
		Index int `json:"index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request")
		return
	}

	log.Printf("[WebUI] Select candidate: %d\n", req.Index)

	// Select candidate
	if !s.rimeEngine.SelectCandidate(s.session, req.Index) {
		s.sendError(w, "Failed to select candidate")
		return
	}

	// Check for commit
	commitText, _ := s.rimeEngine.GetCommit(s.session)
	
	// Record commit in LLM history
	if commitText != "" && s.llmSuggester != nil {
		s.llmSuggester.RecordCommit(commitText)
		log.Printf("[WebUI] Committed: '%s'\n", commitText)
	}

	resp := CommitResponse{
		Text:    commitText,
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.rimeEngine.ClearComposition(s.session)
	log.Println("[WebUI] Cleared composition")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *Server) handleToggleLLM(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.enableLLM = !s.enableLLM
	log.Printf("[WebUI] LLM toggled: %v\n", s.enableLLM)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"enabled": s.enableLLM,
	})
}

func (s *Server) sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(CandidateResponse{
		Success: false,
		Error:   message,
	})
}
