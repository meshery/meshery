---
name: cut-release
description: >-
  Cut (publish) a release of meshery-cloud by publishing the current Release
  Drafter draft with the gh CLI. Use this whenever the user asks to "cut a
  release", "publish the release", "ship a release", "release meshery-cloud",
  "do a release", "release the latest merge", or wants the recently merged PR to
  go out to production. Handles waiting for the Release Drafter workflow to fold
  the just-merged PR into the draft notes, then flipping the draft to published
  so the deploy + docs + helm workflows fire. The version tag is already set by
  Release Drafter and auto-increments after each release, so this skill never
  creates or bumps a tag.
---

# Cut a meshery-cloud release

Publishing a release here is intentionally a one-action step: **flip the existing
Release Drafter draft from draft to published.** Everything downstream is automated
by GitHub Actions.

## How releasing works in this repo

- `.github/workflows/release-drafter.yaml` runs on **every push to `master`**. As PRs
  merge, it maintains a single **draft** GitHub Release, titled `Layer5 Cloud v$NEXT_PATCH_VERSION`
  with tag `v$NEXT_PATCH_VERSION` (the patch version auto-increments). There is always
  exactly one draft waiting to be published.
- Publishing that draft (the release `released` event) triggers
  `.github/workflows/notify-release-and-deploy.yaml`, which deploys to **production** and
  releases the Helm charts, and `.github/workflows/release-notes.yaml`, which updates the docs.

So the tag is already chosen and the notes are already drafted. Your job is only to make
sure the just-merged PR is reflected in the draft, then publish it. **Do not create a tag,
do not bump a version, do not write release notes by hand** - Release Drafter owns all of that.

## Steps

### 1. Confirm the Release Drafter run for the latest master commit has finished

The draft only includes a PR after the drafter workflow run for that merge commit completes.
If you publish too early, the just-merged PR is missing from the notes.

```bash
# Latest master commit that should be in the release
git fetch origin master --quiet && git rev-parse origin/master

# Most recent Release Drafter runs (headSha should match origin/master, status "completed")
gh run list --workflow="release-drafter.yaml" --branch master --limit 3 \
  --json databaseId,status,conclusion,headSha,createdAt
```

If the run for the current `origin/master` SHA is still `in_progress` (or absent because the
push just landed), wait for it before publishing:

```bash
gh run watch <databaseId> --exit-status
```

A `conclusion` of `success` for the run whose `headSha` matches `origin/master` means the draft
is up to date. If the latest run **failed**, stop and surface that - do not publish stale notes.

### 2. Identify and inspect the draft

```bash
# Tag of the current draft (there is normally exactly one)
gh release list --limit 25 --json tagName,isDraft --jq '.[] | select(.isDraft) | .tagName'

# Read the draft notes and confirm the recently merged PR appears in them
gh release view <draftTag> --json tagName,name,isDraft,body --jq \
  '{tag: .tagName, name: .name, isDraft: .isDraft, body: .body}'
```

Verify the merged PR the user is releasing shows up under one of the category headings. If it
does not, the drafter run from step 1 has not yet indexed it - re-check step 1 rather than
publishing without it.

If more than one draft is ever returned, stop and ask the user which to publish rather than
guessing - publishing the wrong one ships an unintended version to production.

### 3. Publish the draft

Flipping `--draft=false` publishes it and fires the `released` event that drives the production
deploy, helm release, and docs update. Mark it `--latest` so it carries the "Latest" badge.

```bash
gh release edit <draftTag> --draft=false --latest
```

### 4. Confirm it published

```bash
gh release view <draftTag> --json tagName,isDraft,isLatest,publishedAt --jq \
  '{tag: .tagName, isDraft: .isDraft, isLatest: .isLatest, publishedAt: .publishedAt}'
```

`isDraft: false` and a non-null `publishedAt` confirm the release is out. Report the published
version to the user and note that the production deploy + docs + helm-chart workflows have been
triggered by the release event.

## What to watch for

- **Don't publish ahead of the drafter run.** The most common failure is publishing before the
  Release Drafter workflow has folded the merged PR into the notes, shipping a release whose notes
  omit the very change being released. Step 1 exists to prevent exactly this.
- **Never hand-author the tag or notes.** If you find yourself computing a version number or
  writing changelog entries, something is wrong - Release Drafter already did both.
- **One draft only.** If `gh release list` shows zero drafts, no drafter run has occurred since the
  last release; if it shows more than one, ask the user which to publish.
