#!/usr/bin/env bash
# PreToolUse guard — block edits that violate Meshery's schema-driven principles.
#
# Meshery is schema-driven: all API types, wire field names, and RTK Query hooks
# are owned by github.com/meshery/schemas. @sistent/sistent owns the UI design
# system. Hand-rolling any of these silently diverges the wire contract.
#
# Checks (NET-NEW only — migrating away from a violation is never blocked):
#   1. Hand-rolled JSON.stringify() with snake_case keys (wire is camelCase)
#   2. Manual builder.query/builder.mutation in ui/rtk-query/ (use schemas hooks)
#   3. Deprecated v1beta1 import declarations (use v1beta3 or v1beta2)
#   4. Go json: tags with snake_case values (wire is camelCase)
#
# RTK opt-out: files that intentionally keep manual endpoints because
# meshery/schemas does not yet cover them must contain:
#   // schema-rtk-exempt: <reason>
#
# Contract: reads the PreToolUse JSON payload on stdin; exits 2 (deny, message
# shown to the agent). Anything else passes (exit 0).
set -uo pipefail

command -v jq    >/dev/null 2>&1 || exit 0  # no jq → fail open
command -v perl  >/dev/null 2>&1 || exit 0  # no perl → fail open

payload="$(cat)"
tool="$(jq -r '.tool_name // empty' <<<"$payload")"
path="$(jq -r '.tool_input.file_path // empty' <<<"$payload")"

case "$tool" in
  Edit | Write | MultiEdit) ;;
  *) exit 0 ;;
esac

case "$path" in
  *.ts | *.tsx | *.go) ;;
  *) exit 0 ;;
esac

# Text being ADDED
new="$(jq -r '[ .tool_input.content, .tool_input.new_string, (.tool_input.edits[]?.new_string) ]
              | map(select(. != null)) | join("\n")' <<<"$payload")"

# Text being REMOVED.
# For Write (full-file overwrite), load the current on-disk content as the
# baseline so pre-existing violations in an unrelated edit are never blocked.
if [ "$tool" = "Write" ] && [ -f "$path" ]; then
  old="$(cat -- "$path")"
else
  old="$(jq -r '[ .tool_input.old_string, (.tool_input.edits[]?.old_string) ]
                | map(select(. != null)) | join("\n")' <<<"$payload")"
fi

[ -z "$new" ] && exit 0

# count <text> <ERE-pattern> — returns integer count (0 if none, never empty)
count() { printf '%s' "$1" | grep -oE "$2" 2>/dev/null | grep -c . || true; }

VIOLATIONS=()

# ─── TypeScript / TSX checks ─────────────────────────────────────────────────
case "$path" in
  *.ts | *.tsx)

    # 1. Hand-rolled JSON.stringify with snake_case keys (multiline-aware)
    new_s=$(count "$new" 'JSON\.stringify')
    old_s=$(count "$old" 'JSON\.stringify')
    if [ "$new_s" -gt "$old_s" ]; then
      if perl -0777 -ne 'exit 0 if /JSON\.stringify\s*\(\s*\{[^}]*?[a-z0-9]+_[a-z0-9_]+\s*:/s; exit 1' <<<"$new"; then
        VIOLATIONS+=("Hand-rolled JSON.stringify() with snake_case key. Use @meshery/schemas generated types — wire format is camelCase (e.g. fileName not file_name).")
      fi
    fi

    # 2. Manual RTK Query endpoints in ui/rtk-query/ (net-new)
    #    Match both absolute and repository-relative paths
    case "$path" in
      ui/rtk-query/* | */ui/rtk-query/*)
        # Check for schema-rtk-exempt in BOTH the new content AND the existing file
        existing_file_content=""
        [ -f "$path" ] && existing_file_content="$(cat -- "$path")"
        if ! printf '%s\n%s' "$new" "$existing_file_content" | grep -qE 'schema-rtk-exempt'; then
          new_b=$(count "$new" 'builder\.(query|mutation)')
          old_b=$(count "$old" 'builder\.(query|mutation)')
          if [ "$new_b" -gt "$old_b" ]; then
            VIOLATIONS+=("New manual RTK Query endpoint (builder.query/builder.mutation). Use generated hooks from @meshery/schemas/mesheryApi. If not yet in schemas, add '// schema-rtk-exempt: <reason>' to the file.")
          fi
        fi
        ;;
    esac

    # 3. Deprecated v1beta1 import declarations only (not comments or docs)
    new_v=$(count "$new" "from[[:space:]]+['\"]@meshery/schemas/models/v1beta1")
    old_v=$(count "$old" "from[[:space:]]+['\"]@meshery/schemas/models/v1beta1")
    if [ "$new_v" -gt "$old_v" ]; then
      VIOLATIONS+=("Deprecated v1beta1 import declaration. Use v1beta3 (or v1beta2 where v1beta3 is absent).")
    fi
    ;;
esac

# ─── Go checks ───────────────────────────────────────────────────────────────
case "$path" in
  *.go)

    # 4. Go json tags with snake_case values (net-new)
    new_snake=$(count "$new" 'json:"[^"]*[a-z]+_[a-z0-9_]*[^"]*"')
    old_snake=$(count "$old" 'json:"[^"]*[a-z]+_[a-z0-9_]*[^"]*"')
    if [ "$new_snake" -gt "$old_snake" ]; then
      tags="$(printf '%s' "$new" | grep -oE 'json:"[^"]*[a-z]+_[a-z0-9_]*[^"]*"' | paste -sd', ' -)"
      VIOLATIONS+=("Go json tag with snake_case value: $tags. Wire is camelCase — use json:\"fileName\" not json:\"file_name\".")
    fi

    # 5. Deprecated v1beta1 Go import paths only (not comments or docs)
    #    Match only inside import declarations: "github.com/meshery/schemas/models/v1beta1"
    new_gv=$(count "$new" '"github\.com/meshery/schemas/models/v1beta1[^"]*"')
    old_gv=$(count "$old" '"github\.com/meshery/schemas/models/v1beta1[^"]*"')
    if [ "$new_gv" -gt "$old_gv" ]; then
      VIOLATIONS+=("Deprecated v1beta1 Go import path. Use v1beta3 (or v1beta2 where v1beta3 is absent).")
    fi
    ;;
esac

# ─── Report ──────────────────────────────────────────────────────────────────
if [ ${#VIOLATIONS[@]} -gt 0 ]; then
  {
    echo "⛔ schema-principles guard — ${path##*/}"
    echo
    echo "Meshery is a SCHEMA-DRIVEN project. This edit introduces code that"
    echo "bypasses the schema-first workflow:"
    echo
    for v in "${VIOLATIONS[@]}"; do
      echo "  ✗  $v"
    done
    echo
    echo "All API types, wire field names, and RTK Query hooks are owned by"
    echo "meshery/schemas. UI components are owned by @sistent/sistent."
    echo "Import the generated types and hooks — do NOT hand-roll them."
    echo
    echo "Ref: https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md"
    echo
    echo "Note: only NET-NEW violations are blocked. Migrating existing code"
    echo "(removing as many violations as you add) passes cleanly."
  } >&2
  exit 2
fi

exit 0
