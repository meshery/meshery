#!/usr/bin/env bash
# PreToolUse hook: block edits to staging and production environment files.
# These files contain live credentials and auth tokens.
set -euo pipefail

INPUT=$(cat)
FP=$(printf '%s' "$INPUT" | python3 -c "
import json, sys
d = json.load(sys.stdin)
print(d.get('tool_input', {}).get('file_path', ''))
" 2>/dev/null || true)

if printf '%s' "$FP" | grep -qE "install/environment/(staging|production)"; then
    echo "BLOCKED: '$(basename "$FP")' is a live-environment credential file." >&2
    echo "Edit install/environment/PRIVATE.md for guidance, or explicitly confirm this is intentional." >&2
    exit 1
fi
