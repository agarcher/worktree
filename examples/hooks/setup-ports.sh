#!/bin/bash
# Example post-create hook: Configure ports based on WT_INDEX
#
# This script uses the WT_INDEX environment variable to calculate
# unique port offsets for each worktree, avoiding port conflicts
# when running multiple development environments simultaneously.
#
# Usage in .wt.yaml:
#   hooks:
#     post_create:
#       - script: ./examples/hooks/setup-ports.sh

set -e

# Check if WT_INDEX is set
if [ -z "$WT_INDEX" ]; then
    echo "Warning: WT_INDEX not set, using default ports"
    exit 0
fi

# Calculate port offset based on WT_INDEX
# Each worktree gets a block of 10 ports
PORT_OFFSET=$((WT_INDEX * 10))

# Write port configuration to a file that can be sourced
cat > "$WT_PATH/.wt-ports.env" << EOF
# Auto-generated port configuration for worktree: $WT_NAME
# Index: $WT_INDEX, Offset: $PORT_OFFSET
PORT_OFFSET=$PORT_OFFSET
VITE_PORT=$((5173 + PORT_OFFSET))
API_PORT=$((3000 + PORT_OFFSET))
DB_PORT=$((5432 + PORT_OFFSET))
REDIS_PORT=$((6379 + PORT_OFFSET))
EOF

echo "Configured ports with offset $PORT_OFFSET (index #$WT_INDEX)"
echo "  VITE_PORT=$((5173 + PORT_OFFSET))"
echo "  API_PORT=$((3000 + PORT_OFFSET))"
echo "Port config written to .wt-ports.env"
