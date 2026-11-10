#!/usr/bin/env bash
# PreToolUse guard — block adding a NEW local wire/DB construct under
# server/models/** without first defining it in github.com/meshery/schemas.
#
# meshery-cloud is schema-driven: wire/DB constructs (structs + json-tagged
# fields) are owned by meshery/schemas and regenerated here. Defining or
# extending one directly under server/models is the failure mode that produced
# the PR #5434 local-`AuthLinks` violation. This guard catches it at the moment
# of edit.
#
# This is the EARLY, in-session catch (it only governs Claude Code tool calls).
# The authoritative, environment-independent enforcement is CI:
# .github/workflows/schema-origin-guard.yml. Keep both.
#
# Contract: reads the PreToolUse JSON payload on stdin; exits 2 (deny, message
# shown to the agent) when the edited path is server/models/**.go (non-test)
# and the edit introduces a NET-NEW `type … struct` or `json:"…"` tag. A
# comment/logic edit that adds neither passes.
set -uo pipefail

command -v jq >/dev/null 2>&1 || exit 0 # no jq → fail open, CI still enforces

payload="$(cat)"
tool="$(jq -r '.tool_name // empty' <<<"$payload")"
path="$(jq -r '.tool_input.file_path // empty' <<<"$payload")"

case "$tool" in
  Edit | Write | MultiEdit) ;;
  *) exit 0 ;;
esac
case "$path" in
  */server/models/*.go) ;;
  *) exit 0 ;;
esac
case "$path" in
  *_test.go) exit 0 ;;
esac

# Text being ADDED (Write.content, Edit.new_string, MultiEdit.edits[].new_string)
new="$(jq -r '[ .tool_input.content, .tool_input.new_string, (.tool_input.edits[]?.new_string) ]
              | map(select(. != null)) | join("\n")' <<<"$payload")"
# Text being REMOVED/replaced — so we only flag tags/types that are NET-NEW
# (a Write has no old_string, so everything counts, which is what we want for a
# freshly authored model file).
old="$(jq -r '[ .tool_input.old_string, (.tool_input.edits[]?.old_string) ]
              | map(select(. != null)) | join("\n")' <<<"$payload")"

extract() { printf '%s' "$1" | grep -oE "$2" | sort -u; }
net_new() { comm -23 <(extract "$new" "$1") <(extract "$old" "$1"); }

added_tags="$(net_new '(json|db):"[^"]*"')"
added_structs="$(net_new 'type[[:space:]]+[A-Za-z0-9_]+[[:space:]]+struct')"

if [ -n "$added_tags$added_structs" ]; then
  {
    echo "⛔ schema-first guard — server/models/"
    echo
    echo "This edit to ${path##*/} introduces a new wire/DB construct:"
    [ -n "$added_structs" ] && echo "  • new struct type(s): $(echo "$added_structs" | paste -sd', ' -)"
    [ -n "$added_tags" ] && echo "  • new json/db tag(s): $(echo "$added_tags" | paste -sd', ' -)"
    echo
    echo "meshery-cloud is SCHEMA-DRIVEN. Constructs are owned by"
    echo "github.com/meshery/schemas and regenerated here — do NOT define or"
    echo "extend them locally. Instead:"
    echo "  1. Add/modify the construct in ../schemas (OpenAPI) and regenerate."
    echo "  2. Bump the schemas dep here and consume the generated type."
    echo "  3. cd ../schemas && go run ./cmd/consumer-audit \\"
    echo "       --cloud-repo ../meshery-cloud --meshery-repo ../meshery"
    echo
    echo "If the change is genuinely local (no schemas equivalent, no wire/DB"
    echo "tag), make it without a new json: tag. CI (schema-origin-guard) also"
    echo "verifies this on the PR."
  } >&2
  exit 2
fi

exit 0
