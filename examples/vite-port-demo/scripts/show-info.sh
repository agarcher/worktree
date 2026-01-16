#!/bin/bash
# Info hook: Displays the dev server URL for this worktree
#
# This hook is run by `wt info` and `wt list -v` to show custom
# worktree-specific information. It calculates the port the same
# way as setup-ports.sh and outputs a clickable URL.

if [ -z "$WT_INDEX" ]; then
    exit 0
fi

# ANSI codes for bold
BOLD='\033[1m'
RESET='\033[0m'

# Check if a port is listening
is_port_running() {
    lsof -i ":$1" >/dev/null 2>&1
}

# Format status with optional bold
format_status() {
    local port=$1
    if is_port_running "$port"; then
        echo -e "(${BOLD}running${RESET})"
    else
        echo "(stopped)"
    fi
}

PORT=$((5173 + WT_INDEX * 10))
echo -e "URL: http://localhost:$PORT $(format_status "$PORT")"
