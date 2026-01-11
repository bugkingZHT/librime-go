// Web-based interactive test for RIME Go wrapper with LLM integration
// Build: go build -o rime-web web-main.go
// Run: ./rime-web
// Then open: http://localhost:8080

package webui

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugkingzht/go-rime/llm"
	"github.com/bugkingzht/go-rime/rime"
)

func main() {
	// Parse command line flags
	var sharedDataDir string
	var userDataDir string
	var appName string
	var port int

	flag.StringVar(&sharedDataDir, "shared-data", "", "Path to shared RIME data directory")
	flag.StringVar(&userDataDir, "user-data", "", "Path to user RIME data directory")
	flag.StringVar(&appName, "app-name", "rime-web", "Application name")
	flag.IntVar(&port, "port", 8080, "HTTP server port")
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
		userDataDir = filepath.Join(homeDir, ".rime-web")
	}

	// Ensure user data directory exists
	os.MkdirAll(userDataDir, 0755)

	// Initialize LLM suggester if API key is available
	var llmSuggester *llm.Suggester
	llmConfig := llm.DefaultConfig()
	if llmConfig.APIKey != "" {
		var err error
		llmSuggester, err = llm.NewSuggester(llmConfig)
		if err != nil {
			fmt.Printf("⚠️  LLM initialization failed: %v\n", err)
			fmt.Println("    Continuing without LLM suggestions...")
		} else {
			fmt.Printf("✓ LLM engine initialized: %s (%s)\n", llmConfig.Provider, llmConfig.Model)
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
	fmt.Println("║        RIME Go Wrapper - Web UI Test Console              ║")
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

	rimeEngine := rime.New(traits)
	if rimeEngine == nil {
		fmt.Println("✗ Failed")
		return
	}
	fmt.Println("✓ Success")
	defer rimeEngine.Shutdown()

	// Create and start web server
	server := NewServer(rimeEngine, llmSuggester)
	fmt.Printf("\n🌐 Starting web server on http://localhost:%d\n", port)
	fmt.Println("   Press Ctrl+C to stop")
	fmt.Println()

	if err := server.Start(port); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
