// Interactive test console for RIME Go wrapper
// Build: go build -o rime-interactive
// Run: ./rime-interactive

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	
	"github.com/bugkingzht/go-rime/llm"
	"github.com/bugkingzht/go-rime/rime"
)

var (
	rimeEngine     *rime.Rime
	currentSession *rime.Session
	committedText  []string
	llmSuggester   *llm.Suggester
	enableLLM      bool
)

func printStatus(session *rime.Session) {
	status, err := rimeEngine.GetStatus(session)
	if err != nil || status == nil {
		return
	}

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Printf("║ Schema: %-30s ║\n", status.SchemaName)
	fmt.Printf("║ Composing: %-27v ║\n", status.IsComposing)
	fmt.Printf("║ ASCII Mode: %-26v ║\n", status.IsAsciiMode)
	fmt.Println("╚════════════════════════════════════════╝")
}

func printContext(session *rime.Session) {
	fmt.Println("[DEBUG] printContext called")
	ctx, err := rimeEngine.GetContext(session)
	if err != nil || ctx == nil {
		fmt.Printf("[DEBUG] No context: err=%v, ctx=%v\n", err, ctx)
		fmt.Println("  (No composition)")
		return
	}

	preedit := ctx.Composition.Preedit
	cursorPos := ctx.Composition.CursorPos
	fmt.Printf("[DEBUG] Preedit: '%s', CursorPos: %d\n", preedit, cursorPos)
	fmt.Printf("[DEBUG] Original candidates count: %d\n", len(ctx.Menu.Candidates))

	// Print preedit with cursor
	fmt.Print("\n  Preedit: ")
	if cursorPos >= 0 && cursorPos <= len(preedit) {
		fmt.Printf("%s|%s\n", preedit[:cursorPos], preedit[cursorPos:])
	} else {
		fmt.Println(preedit)
	}

	// Get candidates (enhanced with LLM if enabled)
	var candidates []rime.Candidate
	fmt.Printf("[DEBUG] LLM enabled: %v, suggester: %v\n", enableLLM, llmSuggester != nil)
	if enableLLM && llmSuggester != nil {
		enhancedCandidates, err := llmSuggester.GetEnhancedCandidates(context.Background(), rimeEngine, session)
		if err == nil && len(enhancedCandidates) > 0 {
			fmt.Printf("[DEBUG] Using enhanced candidates: %d\n", len(enhancedCandidates))
			candidates = enhancedCandidates
		} else {
			fmt.Printf("[DEBUG] Enhanced failed or empty, using original. err=%v, count=%d\n", err, len(enhancedCandidates))
			// Fall back to original candidates
			candidates = ctx.Menu.Candidates
		}
	} else {
		fmt.Println("[DEBUG] LLM disabled, using original candidates")
		candidates = ctx.Menu.Candidates
	}
	highlightedIdx := ctx.Menu.HighlightedCandidateIndex

	if len(candidates) > 0 {
		fmt.Println("\n  Candidates:")
		for i, cand := range candidates {
			marker := " "
			if i == highlightedIdx {
				marker = "▶"
			}
			comment := ""
			if cand.Comment != "" {
				comment = fmt.Sprintf(" [%s]", cand.Comment)
			}
			fmt.Printf("    %s %d. %s%s\n", marker, i+1, cand.Text, comment)
		}
	}
	fmt.Println()
}

func printCommitted() {
	if len(committedText) > 0 {
		fmt.Print("\n  ✓ Committed: ")
		for _, text := range committedText {
			fmt.Print(text)
		}
		fmt.Println()
		committedText = []string{}
	}
}

func checkCommit() {
	if text, err := rimeEngine.GetCommit(currentSession); err == nil && text != "" {
		committedText = append(committedText, text)
	}
}

func printHelp() {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          RIME Interactive Test - Commands                 ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║  <text>      - Type pinyin (e.g., 'nihao')                ║")
	fmt.Println("║  :<n>        - Select candidate by number (e.g., ':1')    ║")
	fmt.Println("║  /clear      - Clear current composition                  ║")
	fmt.Println("║  /status     - Show RIME status                           ║")
	fmt.Println("║  /committed  - Show all committed text                    ║")
	fmt.Println("║  /llm        - Toggle LLM suggestions (on/off)            ║")
	fmt.Println("║  /help       - Show this help message                     ║")
	fmt.Println("║  /quit       - Exit the program                           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}

func main() {
	// Parse command line flags
	var sharedDataDir string
	var userDataDir string
	var appName string
	
	flag.StringVar(&sharedDataDir, "shared-data", "", "Path to shared RIME data directory")
	flag.StringVar(&userDataDir, "user-data", "", "Path to user RIME data directory")
	flag.StringVar(&appName, "app-name", "rime-interactive", "Application name")
	flag.Parse()
	
	// Auto-detect data directories if not provided
	if sharedDataDir == "" {
		// Try executable directory first
		execPath, _ := os.Executable()
		execDir := filepath.Dir(execPath)
		dataPath := filepath.Join(filepath.Dir(execDir), "data")
		
		if _, err := os.Stat(dataPath); err == nil {
			sharedDataDir = dataPath
		} else {
			// Fallback to Squirrel
			sharedDataDir = "/Library/Input Methods/Squirrel.app/Contents/SharedSupport"
		}
	}
	
	if userDataDir == "" {
		homeDir, _ := os.UserHomeDir()
		userDataDir = filepath.Join(homeDir, ".rime-interactive")
	}
	
	// Ensure user data directory exists
	os.MkdirAll(userDataDir, 0755)
	
	// Initialize LLM suggester if API key is available
	llmConfig := llm.DefaultConfig()
	if llmConfig.APIKey != "" {
		var err error
		llmSuggester, err = llm.NewSuggester(llmConfig)
		if err != nil {
			fmt.Printf("⚠️  LLM initialization failed: %v\n", err)
			fmt.Println("    Continuing without LLM suggestions...")
		} else {
			enableLLM = true
			fmt.Printf("✓ LLM engine initialized: %s (%s)\n", llmConfig.Provider, llmConfig.Model)
			fmt.Println("  Use /llm to toggle suggestions")
		}
	} else {
		fmt.Println("ℹ️  No API key found (set DASHSCOPE_API_KEY or OPENAI_API_KEY)")
		fmt.Println("    LLM suggestions disabled")
	}
	
	// Copy default.custom.yaml to user data directory if not exists
	customConfigSrc := filepath.Join(sharedDataDir, "default.custom.yaml")
	customConfigDst := filepath.Join(userDataDir, "default.custom.yaml")
	if _, err := os.Stat(customConfigDst); os.IsNotExist(err) {
		if data, err := os.ReadFile(customConfigSrc); err == nil {
			os.WriteFile(customConfigDst, data, 0644)
		}
	}
	
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        RIME Go Wrapper - Interactive Test Console         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Shared Data: %s\n", sharedDataDir)
	fmt.Printf("User Data:   %s\n", userDataDir)
	fmt.Println()

	// Initialize RIME
	fmt.Print("Initializing RIME... ")
	traits := rime.Traits{
		SharedDataDir: sharedDataDir,
		UserDataDir:   userDataDir,
		AppName:       appName,
	}
	
	rimeEngine = rime.New(traits)
	if rimeEngine == nil {
		fmt.Println("✗ Failed")
		return
	}
	fmt.Println("✓ Success")
	defer rimeEngine.Shutdown()

	// Create session
	fmt.Print("Creating session... ")
	currentSession = rimeEngine.CreateSession()
	if currentSession == nil {
		fmt.Println("✗ Failed")
		return
	}
	fmt.Println("✓ Success")
	defer rimeEngine.DestroySession(currentSession)

	printStatus(currentSession)
	printHelp()

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(input, "/") {
			switch input {
			case "/quit", "/exit", "/q":
				fmt.Println("Goodbye!")
				return
			case "/help", "/h":
				printHelp()
			case "/status":
				printStatus(currentSession)
			case "/clear":
				rimeEngine.ClearComposition(currentSession)
				fmt.Println("  ✓ Composition cleared")
			case "/committed":
				if len(committedText) > 0 {
					fmt.Println("\n  All committed text:")
					for i, text := range committedText {
						fmt.Printf("    %d. %s\n", i+1, text)
					}
				} else {
					fmt.Println("  (No committed text)")
				}
			case "/llm":
				if llmSuggester != nil {
					enableLLM = !enableLLM
					if enableLLM {
						fmt.Println("  ✓ LLM suggestions enabled")
					} else {
						fmt.Println("  ✓ LLM suggestions disabled")
					}
				} else {
					fmt.Println("  ✗ LLM not available (set DASHSCOPE_API_KEY or OPENAI_API_KEY)")
				}
			default:
				fmt.Printf("  ✗ Unknown command: %s (type /help for commands)\n", input)
			}
			continue
		}

		// Handle candidate selection
		if strings.HasPrefix(input, ":") {
			indexStr := strings.TrimPrefix(input, ":")
			if index, err := strconv.Atoi(indexStr); err == nil && index > 0 {
				if rimeEngine.SelectCandidate(currentSession, index-1) {
					checkCommit()
					printCommitted()
					printContext(currentSession)
				} else {
					fmt.Println("  ✗ Invalid candidate index")
				}
			} else {
				fmt.Println("  ✗ Invalid format. Use ':1', ':2', etc.")
			}
			continue
		}

		// Type input
		if rimeEngine.SimulateKeySequence(currentSession, input) {
			checkCommit()
			if len(committedText) > 0 {
				printCommitted()
			}
			printContext(currentSession)
		} else {
			fmt.Println("  ✗ Failed to process input")
		}
	}
}
