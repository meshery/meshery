---
name: cut-release
description: >-
  Cut (publish) a release of meshery/meshery by publishing the current Release
  Drafter draft. Use this whenever the user asks to "cut a release", "publish the
  release", "ship a release", "release meshery", "do a release", "release the
  latest merge", or wants recently merged PRs to go out to users. Handles waiting
  for the Release Drafter workflow to fold the merged PRs into the draft notes,
  then flipping the draft to published so the stable build-and-release workflow
  family fires. The version tag is already set by Release Drafter and
  auto-increments after each release, so this skill never creates or bumps a tag.
---

# Cut a meshery/meshery release

Publishing a release here is intentionally a one-action step: **flip the existing
Release Drafter draft from draft to published.** Everything downstream is automated
by GitHub Actions.

Meshery has **no automatic release cadence** - nothing publishes the draft on a
schedule. The deliberate publish in this skill *is* the release.

## How releasing works in this repo

- `.github/workflows/release-drafter.yml` runs on **every push to `master`** (except
  pushes touching only `.github/**`, which it ignores). As PRs merge it maintains a
  single **draft** GitHub Release, titled `Meshery v$NEXT_PATCH_VERSION` with tag
  `v$NEXT_PATCH_VERSION` (patch auto-increments). There is always exactly one draft
  waiting to be published. Config: `.github/release-drafter.yml`.
- Publishing the draft creates the `v*` tag, which fires
  `.github/workflows/build-and-release-stable.yml` (tag push). That single workflow
  fans out to `ctlrelease` (mesheryctl via goreleaser), `call-dde-release-workflow`,
  `call-helm-chart-releaser`, and `email-meshery-release-notes-workflow` as internal
  jobs - those workflows are `workflow_call`/`workflow_dispatch` only and are *not*
  separately triggered by the release event, so their absence from the run list is
  expected, not a failure.
- The `released`/`published` event also fires
  `.github/workflows/release-notes.yml` (Release Notes Publisher) and Notify Remote
  Providers.

So the tag is already chosen and the notes are already drafted. Your job is only to
make sure the merged PRs are reflected in the draft, then publish it. **Do not create
a tag, do not bump a version, do not write release notes by hand** - Release Drafter
owns all of that.

## Steps

### 1. Confirm the Release Drafter run for the latest master commit has finished

The draft only includes a PR after the drafter workflow run for that merge commit
completes. If you publish too early, the just-merged PR is missing from the notes.

```bash
# Latest master commit that should be in the release
git fetch origin master --quiet && git rev-parse origin/master

# Most recent Release Drafter runs
gh run list --workflow=release-drafter.yml --branch master --limit 5 \
  --json databaseId,status,conclusion,headSha,createdAt
```

The run whose `headSha` matches `origin/master` must be `completed`/`success`. If it
is still `in_progress`, wait (`gh run watch <databaseId> --exit-status`). If it
**failed**, stop and surface that - do not publish stale notes.

Caveat: because of the `paths-ignore: '.github/**'` filter, a master HEAD whose commit
touched only `.github/**` has **no** drafter run at all. Match against the newest
commit that did trigger one rather than concluding the drafter is broken.

### 2. Identify and inspect the draft

```bash
gh api repos/meshery/meshery/releases \
  --jq '.[] | select(.draft) | {id, tag: .tag_name, name, target: .target_commitish}'

gh api repos/meshery/meshery/releases/<id> --jq .body
```

Verify the PRs being released show up under a category heading. If they do not, the
drafter run from step 1 has not indexed them - re-check step 1 rather than publishing
without them.

Confirm the commit range too, not just the notes - a PR can be named in notes carried
over from an earlier draft:

```bash
gh api repos/meshery/meshery/releases/latest --jq .tag_name   # e.g. v1.0.64
git log <lastTag>..origin/master --oneline --grep='#<PR>'
```

If more than one draft is ever returned, stop and ask which to publish rather than
guessing.

### 3. Publish the draft

```bash
gh release edit <draftTag> --draft=false --latest
```

### 4. Confirm it published - by re-reading, never by exit status

**`gh`/`gh-axi release edit --draft=false` can report success and leave the release a
draft.** Observed on 2026-08-07 cutting v1.0.65: `edit: ok` was printed and the
release was still `"draft": true` with a `untagged-<hash>` URL. Always re-read:

```bash
gh api repos/meshery/meshery/releases/<id> \
  --jq '{tag: .tag_name, draft, published_at, html_url}'
```

`draft: false` and a non-null `published_at` are the only proof. If the flag silently
failed, publish through the REST API instead:

```bash
gh api -X PATCH repos/meshery/meshery/releases/<id> -F draft=false -f make_latest=true
```

Note `-F draft=false` (typed) - `-f draft=false` sends the *string* `"false"`.

Then confirm the tag was created and points where you expect:

```bash
gh api repos/meshery/meshery/git/ref/tags/<draftTag> --jq '{ref, sha: .object.sha}'
```

### 5. Confirm the release-triggered workflows actually dispatched

A published release does not prove downstream ran. A release whose workflows never
dispatched is not a cut release - report it if that happens.

```bash
gh api "repos/meshery/meshery/actions/runs?created=>$(date -u -v-5M +%Y-%m-%dT%H:%M:%SZ)" \
  --jq '.workflow_runs[] | {name, event, status, url: .html_url}'
```

Expect at minimum **Meshery Build and Releaser (stable)** (event `push`, branch
`v<version>`) and **Release Notes Publisher** (event `release`). Report their run URLs.

## What to watch for

- **Don't publish ahead of the drafter run.** The most common failure is publishing
  before Release Drafter has folded the merged PR into the notes, shipping a release
  whose notes omit the very change being released. Step 1 exists to prevent this.
- **Never hand-author the tag or notes.** If you find yourself computing a version
  number or writing changelog entries, something is wrong.
- **One draft only.** Zero drafts means no drafter run since the last release; more
  than one means ask which to publish.
- **Trust the re-read, not the command.** Steps 4 and 5 are the verification; a
  zero exit code from `release edit` is not evidence of anything.
