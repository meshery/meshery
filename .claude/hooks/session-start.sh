#!/bin/bash
set -euo pipefail

if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

if command -v gh &>/dev/null; then
  exit 0
fi

apt-get install -y -qq gh
