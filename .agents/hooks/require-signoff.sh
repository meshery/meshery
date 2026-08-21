#!/usr/bin/env bash
# require-signoff.sh — Block a `git commit` that would produce a commit with no
# Signed-off-by trailer.
#
# Usage:
#   ./require-signoff.sh '<shell command line>'
#   CLAUDE_TOOL_INPUT='<shell command line>' ./require-signoff.sh
#
# Exits non-zero (blocks the command) if the command line runs `git commit` and
# nothing in it would put a Signed-off-by trailer on the resulting message.
# Exits 0 (allows the command) otherwise.
#
# Can be wired into any coding agent's pre-command hook system.
#
# Why this exists beside ui/.husky/commit-msg, which already applies the same
# rule: git reaches that hook only through `core.hooksPath`, so it is silently
# absent both on a checkout where `make ui-setup` never ran and under any
# harness that repoints `core.hooksPath` at its own hooks directory. This guard
# inspects the command instead of git's hook plumbing, so repointing
# `core.hooksPath` does not disarm it. An absent guard reads exactly like a
# guard that passed, and by the time the DCO check reports the omission the only
# remedy left is rewriting the branch.

set -uo pipefail

COMMAND="${1:-${CLAUDE_TOOL_INPUT:-}}"

if [[ -z "$COMMAND" ]]; then
  exit 0
fi

# Flags are read from the command with quoted spans blanked out, so a message
# body cannot be mistaken for one: `git commit -m "drop the -s flag"` carries no
# sign-off, and `echo "git commit"` runs no commit at all. The unmodified
# command is kept for reading flag *values*, which are often quoted.
FLAGS=$(printf '%s' "$COMMAND" | sed -E 's/"[^"]*"/ /g; s/'"'"'[^'"'"']*'"'"'/ /g')

# Only `git commit` creates the commits this guard is about. Walk the tokens
# rather than pattern-matching the line, so `git commit` is still recognised
# behind git's global options - including the ones that take a separate value,
# as in `git -C ui commit` - and after a shell separator, as in
# `cd ui && git commit -m ...`, while `git log commit` is not mistaken for one.
runs_git_commit() {
  local -a tokens
  local i j
  read -ra tokens <<<"$FLAGS"

  for ((i = 0; i < ${#tokens[@]}; i++)); do
    [[ "${tokens[i]}" == *git ]] || continue

    for ((j = i + 1; j < ${#tokens[@]}; j++)); do
      case "${tokens[j]}" in
        commit) return 0 ;;
        -C | -c | --git-dir | --work-tree | --namespace | --exec-path | --config-env)
          ((j++)) ;;
        -*) ;;
        *) break ;;
      esac
    done
  done

  return 1
}

if ! runs_git_commit; then
  exit 0
fi

allow() { exit 0; }

block() {
  echo "BLOCKED: this git commit would not be signed off, and the DCO check would fail on it." >&2
  echo >&2
  echo "  $1" >&2
  echo >&2
  echo "  Re-run the commit with -s so it carries a Signed-off-by trailer:" >&2
  echo >&2
  echo "    git commit -s ..." >&2
  echo >&2
  echo "  A sign-off cannot be added afterwards without rewriting the branch: the" >&2
  echo "  DCO check (https://github.com/dcoapp/app) reads its configuration from" >&2
  echo "  the default branch only, so no change on a branch can relax it." >&2
  exit 1
}

# A merge commit carries no sign-off requirement - the DCO check exempts it, and
# so does ui/.husky/commit-msg.
if [[ -f "$(git rev-parse --git-path MERGE_HEAD 2>/dev/null || echo /nonexistent)" ]]; then
  allow
fi

# `--no-signoff` is an explicit request to leave the trailer off, which is
# exactly what this guard is here to refuse.
if [[ "$FLAGS" =~ (^|[[:space:]])--no-signoff([[:space:]]|$) ]]; then
  block "The command passes --no-signoff."
fi

# An explicit trailer counts however it was written: `--trailer
# Signed-off-by=...`, or spelled out inside the message body.
if [[ "$COMMAND" == *Signed-off-by* ]]; then
  allow
fi

if [[ "$FLAGS" =~ (^|[[:space:]])--signoff([[:space:]]|=|$) ]]; then
  allow
fi

# Short options bundle, and a bundled option that takes a value swallows the
# rest of its cluster: `-sm "msg"` is `-s -m "msg"`, but `-ms "msg"` is `-m s`
# with no sign-off at all. So `s` only counts when it appears ahead of the first
# value-taking letter in its cluster.
while read -r cluster; do
  [[ -z "$cluster" ]] && continue
  if [[ "${cluster%%[mFcCt]*}" == *s* ]]; then
    allow
  fi
done < <(printf '%s\n' "$FLAGS" | grep -oE '(^|[[:space:]])-[a-zA-Z]+' | sed 's/^[[:space:]]*-//')

# Nothing on the command line signs off, so the message the commit ends up with
# decides. Where that message comes from existing text - an amend reusing HEAD's
# message, `-F` reading a file, `-c`/`-C` reusing another commit's - the trailer
# may already be in it.
signed_for() {
  # $1 message text, remaining args the identities a trailer may name.
  local message="$1" ident signoffs
  shift

  # Git drops comment lines before recording the message, and `git commit
  # --verbose` appends the staged diff below a scissors line. Neither can carry
  # a sign-off, so remove both before looking for one.
  message=$(printf '%s\n' "$message" |
    sed -e '/^#[[:space:]]*-\{2,\}[[:space:]]*>8[[:space:]]*-\{2,\}/,$d' -e '/^#/d')

  # Normalize the trailers to one canonical prefix before matching identities as
  # fixed strings: "Signed-off-by:  Jane <jane@example.com>" is a trailer the DCO
  # check accepts, and unnormalized it would match no identity here.
  signoffs=$(printf '%s\n' "$message" |
    grep -i '^Signed-off-by:[[:space:]]*..*<..*>' |
    sed -e 's/^[^:]*:[[:space:]]*/Signed-off-by: /')

  [[ -z "$signoffs" ]] && return 1

  for ident in "$@"; do
    [[ -z "$ident" ]] && continue
    printf '%s\n' "$signoffs" | grep -qiF "Signed-off-by: $ident" && return 0
  done

  return 1
}

# `git var` renders an identity as "Name <email> 1700000000 +0000"; drop the
# timestamp and the timezone to leave the "Name <email>" a trailer must carry.
identity() {
  local ident
  ident=$(git var "$1" 2>/dev/null) || return 0
  ident=${ident% *}
  printf '%s' "${ident% *}"
}

committer=$(identity GIT_COMMITTER_IDENT)

# A new message given on the command line is the whole message, so an unsigned
# one stays unsigned however HEAD reads.
if [[ "$FLAGS" =~ (^|[[:space:]])(-m|--message)([[:space:]]|=) ]]; then
  block "The message passed with -m carries no Signed-off-by trailer."
fi

flag_value() {
  # Echo the value of the first of the given flags present on the command, in
  # either `--flag value` or `--flag=value` spelling.
  local flag value
  for flag in "$@"; do
    value=$(printf '%s\n' "$COMMAND" |
      grep -oE "(^|[[:space:]])${flag}([[:space:]]+|=)[^[:space:]]+" |
      head -n 1 |
      sed -E "s/.*${flag}([[:space:]]+|=)//" |
      tr -d "\"'")
    if [[ -n "$value" ]]; then
      printf '%s' "$value"
      return 0
    fi
  done
  return 1
}

if file=$(flag_value --file -F); then
  if [[ -r "$file" ]] && signed_for "$(cat "$file")" "$(identity GIT_AUTHOR_IDENT)" "$committer"; then
    allow
  fi
  block "The message file '$file' carries no Signed-off-by trailer naming the author or the committer."
fi

if ref=$(flag_value --reuse-message --reedit-message -C -c); then
  if signed_for "$(git log -1 --format=%B "$ref" 2>/dev/null)" \
    "$(git log -1 --format='%an <%ae>' "$ref" 2>/dev/null)" "$committer"; then
    allow
  fi
  block "The message reused from '$ref' carries no Signed-off-by trailer naming the author or the committer."
fi

if [[ "$FLAGS" =~ (^|[[:space:]])--amend([[:space:]]|$) ]]; then
  # An amend with no new message source reuses HEAD's message, keeping HEAD's
  # author and re-stamping the committer as whoever runs it.
  if signed_for "$(git log -1 --format=%B 2>/dev/null)" \
    "$(git log -1 --format='%an <%ae>' 2>/dev/null)" "$committer"; then
    allow
  fi
  block "Amending HEAD reuses its message, which carries no Signed-off-by trailer naming its author or the committer."
fi

# No message source at all: the message is about to be typed into an editor, so
# whether it will be signed cannot be known from here. The commit-msg hook is
# what would normally decide, and this guard exists because that hook may be
# unreachable - so require the trailer to be asked for up front.
block "The commit carries no sign-off flag."
