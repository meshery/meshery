#!/usr/bin/env bash
set -euo pipefail

site_dir=$1
pages_dir=$2

cd "$pages_dir"
git config user.name 'github-actions[bot]'
git config user.email 'github-actions[bot]@users.noreply.github.com'

for attempt in 1 2 3 4 5; do
  # Start every attempt from the latest branch so previews added or removed
  # while the site was building (or while a push was in flight) are preserved.
  git fetch --quiet origin gh-pages
  git reset --quiet --hard FETCH_HEAD
  rsync -ac --delete --exclude='/.git/' --exclude='/pr-preview/' "$site_dir/" ./
  git add --all

  if git diff --cached --quiet; then
    echo 'GitHub Pages is already up to date.'
    exit 0
  fi

  git commit -m "deploy: ${GITHUB_SHA}"
  if git push origin HEAD:gh-pages; then
    exit 0
  fi

  echo "GitHub Pages push failed (attempt ${attempt}/5); refreshing and retrying."
  if [ "$attempt" -lt 5 ]; then
    sleep "$((attempt * 5))"
  fi
done

echo '::error::Could not deploy to GitHub Pages after 5 attempts.'
exit 1
