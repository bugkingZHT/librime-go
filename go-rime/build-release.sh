#!/bin/bash
# Release packaging script for rime-interactive
# This script creates a standalone release package with all dependencies

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_NAME="rime-interactive"
RELEASE_DIR="$SCRIPT_DIR/release"
PKG_DIR="$RELEASE_DIR/$APP_NAME"

echo "==================================="
echo "  RIME Interactive Release Builder"
echo "==================================="
echo

# Clean previous build
echo "1. Cleaning previous build..."
rm -rf "$RELEASE_DIR"
mkdir -p "$PKG_DIR/bin"
mkdir -p "$PKG_DIR/lib"
mkdir -p "$PKG_DIR/data"

# Build the executable
echo "2. Building executable..."
cd "$SCRIPT_DIR"
go build -o "$PKG_DIR/bin/$APP_NAME" main.go

# Copy librime dynamic library
echo "3. Copying librime library..."
LIBRIME_PATH="/opt/homebrew/opt/librime/lib"
if [ -d "$LIBRIME_PATH" ]; then
    cp "$LIBRIME_PATH/librime.1.dylib" "$PKG_DIR/lib/" || \
    cp "$LIBRIME_PATH/librime.dylib" "$PKG_DIR/lib/librime.1.dylib"
    
    # Update library path in executable
    install_name_tool -change \
        "$LIBRIME_PATH/librime.1.dylib" \
        "@executable_path/../lib/librime.1.dylib" \
        "$PKG_DIR/bin/$APP_NAME"
else
    echo "   Warning: librime not found at $LIBRIME_PATH"
fi

# Copy RIME data files from Squirrel
echo "4. Copying RIME data files..."
SQUIRREL_DATA="/Library/Input Methods/Squirrel.app/Contents/SharedSupport"
if [ -d "$SQUIRREL_DATA" ]; then
    # Copy essential data files
    for file in default.yaml symbols.yaml; do
        if [ -f "$SQUIRREL_DATA/$file" ]; then
            cp "$SQUIRREL_DATA/$file" "$PKG_DIR/data/"
        fi
    done
    
    # Copy schema files (拼音输入方案 - 重点复制简体方案)
    for schema in luna_pinyin_simp luna_pinyin terra_pinyin; do
        cp "$SQUIRREL_DATA/${schema}"*.yaml "$PKG_DIR/data/" 2>/dev/null || true
    done
    
    # Copy dictionary files (优先复制简体相关的预编译文件)
    # 通用字典
    cp "$SQUIRREL_DATA/"*.txt "$PKG_DIR/data/" 2>/dev/null || true
    cp "$SQUIRREL_DATA/"*.gram "$PKG_DIR/data/" 2>/dev/null || true
    
    # 预编译文件 - 简体拼音优先
    for schema in luna_pinyin_simp luna_pinyin; do
        cp "$SQUIRREL_DATA/${schema}.table.bin" "$PKG_DIR/data/" 2>/dev/null || true
        cp "$SQUIRREL_DATA/${schema}.prism.bin" "$PKG_DIR/data/" 2>/dev/null || true
        cp "$SQUIRREL_DATA/${schema}.reverse.bin" "$PKG_DIR/data/" 2>/dev/null || true
    done
    
    # 其他预编译文件
    cp "$SQUIRREL_DATA/"*.table.bin "$PKG_DIR/data/" 2>/dev/null || true
    cp "$SQUIRREL_DATA/"*.prism.bin "$PKG_DIR/data/" 2>/dev/null || true
    cp "$SQUIRREL_DATA/"*.reverse.bin "$PKG_DIR/data/" 2>/dev/null || true
else
    echo "   Warning: Squirrel data not found at $SQUIRREL_DATA"
    echo "   Please install Squirrel or download RIME data manually"
fi

# Copy custom configuration for simplified Chinese
echo "   Copying custom configuration (simplified Chinese)..."
if [ -f "$SCRIPT_DIR/data/default.custom.yaml" ]; then
    cp "$SCRIPT_DIR/data/default.custom.yaml" "$PKG_DIR/data/"
    echo "   ✓ Using simplified Chinese configuration"
fi

# Copy user database from Squirrel if available
echo "5. Copying user database template..."
USER_RIME="$HOME/Library/Rime"
if [ -d "$USER_RIME/luna_pinyin.userdb" ]; then
    mkdir -p "$PKG_DIR/data/userdb_template"
    cp -r "$USER_RIME/luna_pinyin.userdb" "$PKG_DIR/data/userdb_template/"
    echo "   ✓ User database template copied from Squirrel"
else
    echo "   Warning: No user database found in Squirrel"
    echo "   Will initialize fresh database on first run"
fi

# Note: We don't copy the build directory template because RIME needs to
# generate the proper .table.bin and .reverse.bin files during deployment.
# The launcher script will trigger deployment on first run.

# Create wrapper script
echo "6. Creating launcher script..."
cat > "$PKG_DIR/rime-interactive.sh" << 'EOF'
#!/bin/bash
# Launcher script for rime-interactive

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set library path
export DYLD_LIBRARY_PATH="$SCRIPT_DIR/lib:$DYLD_LIBRARY_PATH"

# Create user data directory if not exists
USER_DATA="$HOME/.rime-interactive"
mkdir -p "$USER_DATA"

# Initialize user database from template if needed
if [ -d "$SCRIPT_DIR/data/userdb_template/luna_pinyin.userdb" ] && [ ! -d "$USER_DATA/luna_pinyin.userdb" ]; then
    echo "Initializing user database from template..."
    cp -r "$SCRIPT_DIR/data/userdb_template/luna_pinyin.userdb" "$USER_DATA/"
    echo "✓ User database initialized"
fi

# Note: Build directory will be generated automatically by RIME deployment on first run

# Run the application
"$SCRIPT_DIR/bin/rime-interactive" \
    --shared-data="$SCRIPT_DIR/data" \
    --user-data="$USER_DATA" \
    "$@"
EOF

chmod +x "$PKG_DIR/rime-interactive.sh"

# Create README
echo "7. Creating README..."
cat > "$PKG_DIR/README.txt" << 'EOF'
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

EOF

# Create package info
echo "8. Creating package info..."
cat > "$PKG_DIR/VERSION.txt" << EOF
Version: 1.0.0
Build Date: $(date '+%Y-%m-%d %H:%M:%S')
Platform: $(uname -s) $(uname -m)
Go Version: $(go version)
EOF

# Create archive
echo "9. Creating archive..."
cd "$RELEASE_DIR"
tar -czf "${APP_NAME}-$(uname -s)-$(uname -m).tar.gz" "$APP_NAME"

echo
echo "==================================="
echo "  Build Complete!"
echo "==================================="
echo
echo "Package location: $RELEASE_DIR"
echo "Archive: ${APP_NAME}-$(uname -s)-$(uname -m).tar.gz"
echo
echo "To test locally:"
echo "  cd $PKG_DIR"
echo "  ./rime-interactive.sh"
echo
