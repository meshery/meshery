#!/usr/bin/env bash
# PostToolUse hook: compile-check the server package after any Go file edit.
# Exits non-zero on compile failure so Claude sees the error immediately.
set -uo pipefail

INPUT=$(cat)
FP=$(printf '%s' "$INPUT" | python3 -c "
import json, sys
d = json.load(sys.stdin)
print(d.get('tool_input', {}).get('file_path', ''))
" 2>/dev/null || true)

# Only trigger for .go files under server/
if printf '%s' "$FP" | grep -qE "\.go$" && printf '%s' "$FP" | grep -q "/server/"; then
    cd /Users/l/code/meshery-cloud
    echo "[go-build-check] compiling server/..." >&2
    go build ./server/... 2>&1
    exit $?
fi
