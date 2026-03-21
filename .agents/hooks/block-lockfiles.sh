#!/usr/bin/env bash
# block-lockfiles.sh — Prevent direct edits to generated lock files
#
# Usage:
#   ./block-lockfiles.sh <file_path>
#
# Exits non-zero (blocks the edit) if the file is a lock file.
# Exits 0 (allows the edit) otherwise.
#
# Can be wired into any coding agent's pre-edit hook system.

set -euo pipefail

FILE="${1:-${CLAUDE_FILE_PATH:-}}"

if [[ -z "$FILE" ]]; then
  exit 0
fi

BASENAME="$(basename "$FILE")"

case "$BASENAME" in
  go.sum|package-lock.json|yarn.lock|pnpm-lock.yaml|Gemfile.lock|composer.lock)
    echo "BLOCKED: Do not edit $BASENAME directly. Use the appropriate package manager instead."
    exit 1
    ;;
esac

exit 0
