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
