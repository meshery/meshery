#!/usr/bin/env bash
# format-frontend.sh — Auto-format JS/TS/JSX/TSX files in ui/ with Prettier
#
# Usage:
#   ./format-frontend.sh <file_path>
#
# This script checks if the given file is a frontend source file inside ui/
# and runs Prettier on it. Exit 0 always (formatting is best-effort).
#
# Can be wired into any coding agent's post-edit hook system,
# or run manually / from a pre-commit hook.

set -euo pipefail

FILE="${1:-${CLAUDE_FILE_PATH:-}}"

if [[ -z "$FILE" ]]; then
  exit 0
fi

# Only format frontend files in ui/
if [[ "$FILE" != ui/* && "$FILE" != */ui/* ]]; then
  exit 0
fi

# Only format JS/TS files
if [[ ! "$FILE" =~ \.(js|jsx|ts|tsx)$ ]]; then
  exit 0
fi

# Find the ui directory relative to the file
UI_DIR="${FILE%ui/*}ui"

if [[ -d "$UI_DIR" ]] && [[ -f "$UI_DIR/node_modules/.bin/prettier" ]]; then
  "$UI_DIR/node_modules/.bin/prettier" --write "$FILE" 2>/dev/null || true
elif command -v npx &>/dev/null; then
  cd "$UI_DIR" && npx prettier --write "$FILE" 2>/dev/null || true
fi

exit 0
