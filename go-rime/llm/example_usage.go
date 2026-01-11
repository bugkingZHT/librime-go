package llm

// Example usage documentation for the LLM-enhanced Rime input method
//
// The LLM engine supports two providers:
//
// 1. DashScope (Alibaba Cloud Qwen) - Priority provider
//    - Set: export DASHSCOPE_API_KEY="your-api-key"
//    - Default model: qwen-turbo
//    - Endpoint: https://dashscope.aliyuncs.com/compatible-mode/v1
//
// 2. OpenAI (Fallback)
//    - Set: export OPENAI_API_KEY="your-api-key"
//    - Optional: export OPENAI_BASE_URL="custom-endpoint"
//    - Default model: gpt-4o-mini
//
// Usage Example:
//
//   // Set up your API key
//   export DASHSCOPE_API_KEY="sk-xxxxxxxxxxxxx"
//
//   // Build and run
//   go build -o rime-interactive
//   ./rime-interactive
//
//   // In the interactive console:
//   > yongGoFengZhuangRimeheLLMYinQingWanQuanKe  // Type pinyin
//   // The system will show:
//   // 1. 用Go封装Rime和LLM引擎完全可 (from Rime)
//   // 2. 行 🤖 AI                        (from LLM)
//   // 3. 能 (other Rime candidates)
//
//   // Commands:
//   > :2         // Select the LLM suggestion
//   > /llm       // Toggle LLM on/off
//   > /help      // Show all commands
//
// Architecture:
//
//   ┌─────────────┐
//   │ User Input  │
//   └──────┬──────┘
//          │
//   ┌──────▼──────────────────────────┐
//   │ main.go (Interactive Console)   │
//   └──────┬──────────────────────────┘
//          │
//          ├─────────────────────┬──────────────────┐
//          │                     │                  │
//   ┌──────▼──────┐      ┌──────▼───────┐   ┌─────▼─────┐
//   │ rime.Rime   │      │ llm.Suggester│   │ llm.Engine│
//   │ (CGO)       │      └──────┬───────┘   └─────┬─────┘
//   └──────┬──────┘             │                 │
//          │                    │                 │
//   ┌──────▼──────┐      ┌──────▼─────────────────▼─────┐
//   │ librime.so  │      │ OpenAI-compatible API        │
//   │ (C Library) │      │ (DashScope/OpenAI)           │
//   └─────────────┘      └──────────────────────────────┘
//
// Flow:
//   1. User types pinyin → Rime generates candidates
//   2. Suggester takes first candidate + preedit as context
//   3. LLM completes the text (≤50 chars)
//   4. LLM suggestion inserted at position 2 in candidate list
