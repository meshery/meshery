#!/usr/bin/env bash
# validate-schema-principles.sh — Enforce Meshery schema-driven development principles
#
# Usage:
#   ./validate-schema-principles.sh <file_path>
#
# Exits non-zero (blocks the edit) if the file violates Meshery's core principles:
#   1. No hand-rolled JSON request bodies — use @meshery/schemas generated types
#   2. No manual RTK Query endpoints for APIs covered by meshery/schemas
#   3. No snake_case field names on the wire (JSON tags / TS properties must be camelCase)
#
# RTK Query opt-out: If a file must keep manual endpoints because @meshery/schemas
# does not yet cover them, add this comment anywhere in the file:
#   // schema-rtk-exempt: <reason>
# This will skip the builder.query/builder.mutation check for that file only.
#
# Can be wired into any coding agent's post-edit hook system.
# Reference: https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md

set -euo pipefail

FILE="${1:-${CLAUDE_FILE_PATH:-}}"

if [[ -z "$FILE" ]]; then
  exit 0
fi

# Only check TypeScript, TSX, and Go files
if [[ ! "$FILE" =~ \.(ts|tsx|go)$ ]]; then
  exit 0
fi

# Fail closed: exit non-zero if the file cannot be read
if [[ ! -r "$FILE" ]]; then
  echo "[validate-schema-principles] ERROR: Cannot read file: $FILE" >&2
  exit 1
fi

BASENAME="$(basename "$FILE")"
VIOLATIONS=()

# ─── TypeScript / TSX checks ─────────────────────────────────────────────────
if [[ "$FILE" =~ \.(ts|tsx)$ ]]; then

  # 1. Hand-rolled JSON body with snake_case keys
  #    Use perl for syntax-aware multiline matching of JSON.stringify({ ... snake_case: ... })
  if perl -0777 -ne 'exit 0 if /JSON\.stringify\s*\(\s*\{[^}]*?['\''"]?[a-z0-9]+_[a-z0-9_]+['\''"]?\s*:/s; exit 1' -- "$FILE"; then
    VIOLATIONS+=("Hand-rolled JSON.stringify() body with snake_case key detected in $BASENAME. Use @meshery/schemas generated types for request payloads (wire format is camelCase).")
  fi

  # 2. Manual RTK Query endpoint definitions (builder.query / builder.mutation)
  #    in ui/rtk-query/** — these should come from @meshery/schemas/mesheryApi
  #    Files may opt out with:  // schema-rtk-exempt: <reason>
  if [[ "$FILE" == ui/rtk-query/* || "$FILE" == */ui/rtk-query/* ]]; then
    if ! grep -qE 'schema-rtk-exempt' -- "$FILE"; then
      if grep -qE 'builder\.(query|mutation)' -- "$FILE"; then
        VIOLATIONS+=("Manual RTK Query endpoint (builder.query/builder.mutation) in $BASENAME. Use generated hooks from @meshery/schemas/mesheryApi instead of hand-rolling endpoints. If this endpoint is not yet in schemas, add '// schema-rtk-exempt: <reason>' to the file.")
      fi
    fi
  fi

  # 3. Importing deprecated v1beta1 in new code (use v1beta3 or v1beta2)
  #    Use [[:space:]] — POSIX ERE compatible, works on both macOS BSD grep and GNU grep
  if grep -qE "from[[:space:]]+['\"]@meshery/schemas/models/v1beta1" -- "$FILE"; then
    VIOLATIONS+=("Deprecated v1beta1 import in $BASENAME. Use v1beta3 (or v1beta2 where v1beta3 is absent).")
  fi

fi

# ─── Go checks ───────────────────────────────────────────────────────────────
if [[ "$FILE" =~ \.go$ ]]; then

  # 4. Go struct json tags with snake_case values
  #    Pattern: json:"file_name" or json:"file_name,omitempty" — wire must be camelCase
  if grep -qE 'json:"[^"]*[a-z]+_[a-z0-9_]*[^"]*"' -- "$FILE"; then
    VIOLATIONS+=("Go struct in $BASENAME has a json tag with snake_case value. Wire is camelCase — json tags must be camelCase (e.g., json:\"fileName\" not json:\"file_name\").")
  fi

  # 5. Importing deprecated v1beta1 in new Go code
  if grep -qE '"github\.com/meshery/schemas/models/v1beta1"' -- "$FILE"; then
    VIOLATIONS+=("Deprecated v1beta1 import in $BASENAME. Use v1beta3 (or v1beta2 where v1beta3 is absent) in new Go code.")
  fi

fi

# ─── Report ──────────────────────────────────────────────────────────────────
if [[ ${#VIOLATIONS[@]} -gt 0 ]]; then
  echo ""
  echo "╔══════════════════════════════════════════════════════════════════════╗"
  echo "║   MESHERY SCHEMA-DRIVEN PRINCIPLES VIOLATION                        ║"
  echo "╚══════════════════════════════════════════════════════════════════════╝"
  echo ""
  echo "Meshery is a schema-driven project. All API types, wire field names,"
  echo "and RTK Query hooks must come from meshery/schemas — not be hand-rolled."
  echo "See: https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md"
  echo ""
  for v in "${VIOLATIONS[@]}"; do
    echo "  ✗  $v"
  done
  echo ""
  echo "Fix the violations above before proceeding."
  exit 1
fi

exit 0
