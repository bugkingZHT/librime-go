module github.com/bugkingzht/go-rime

go 1.22

toolchain go1.24.3

// CGO Configuration:
// This module wraps librime C API using the new rime_get_api() style
//
// Dependencies:
// - librime: brew install librime
// - Squirrel (optional): for RIME data files

require github.com/openai/openai-go v1.12.0

require (
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)
