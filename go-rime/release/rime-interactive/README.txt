RIME Interactive Test Console
===============================

This is a standalone interactive terminal for testing the RIME input method engine.
Configured for Simplified Chinese (简体中文) by default.

## Quick Start

1. Run the launcher script:
   ./rime-interactive.sh

2. Or run directly:
   export DYLD_LIBRARY_PATH="$PWD/lib:$DYLD_LIBRARY_PATH"
   ./bin/rime-interactive

## Usage

Commands in interactive mode:
- <text>      Type pinyin (e.g., 'nihao')
- :<n>        Select candidate by number (e.g., ':1')
- /clear      Clear current composition
- /status     Show RIME status
- /committed  Show all committed text
- /help       Show help message
- /quit       Exit the program

## Example

> nihao
  Preedit: ni hao|
  
  Candidates:
    ▶ 1. 你好
      2. 妳好
      3. 逆号
      
> :1
  ✓ Committed: 你好

## Directory Structure

├── bin/              Executable binary
├── lib/              Dynamic libraries (librime)
├── data/             RIME data files (schemas, dictionaries)
│   └── default.custom.yaml  Simplified Chinese configuration
└── README.txt        This file

## Configuration

The input method is configured for Simplified Chinese by default.
Configuration file: data/default.custom.yaml

Default schema: Luna Pinyin Simplified (朙月拼音·简化字)

## Notes

- User data will be stored in: ~/.rime-interactive
- Shared data is in: ./data
- Library dependencies are in: ./lib

## Dependencies Included

- librime 1.x (RIME Input Method Engine)
- Luna Pinyin Simplified schema (朙月拼音·简化字)
- Essential dictionaries and configuration files

