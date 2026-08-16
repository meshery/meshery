---
title: Build Environment Gotchas
description: Local build and dependency failures in Meshery that report a cause other than their own.
categories: [contributing]
---

Each entry here is a failure whose message points somewhere other than the actual problem. Read the one that matches your symptom; none of them need to be read up front.

## Go: "does not match go tool version" on dozens of dependencies

A build that fails with `compile: version "goX.Y.Z" does not match go tool version` across many dependencies is an environment mismatch, not a code problem. A stale `GOROOT` points at one Go installation while the `go` on `PATH` is a different one.

Any toolchain at or above the version in `go.mod` works - `go 1.26.4` there does not mean 1.26.5 is wrong - so the fix is to make the two agree, not to pin an exact patch release. Drop the stale `GOROOT` and use one installation's own binary:

{{< code code=`env -u GOROOT PATH="$HOME/.gvm/gos/go1.26.4/bin:$PATH" go test ./...` >}}

## `ui/tsconfig.tsbuildinfo` gets committed by accident

`ui/tsconfig.tsbuildinfo` is a tracked build artifact and is not gitignored, so any local `tsc --noEmit` leaves it modified. It has repeatedly been committed by accident.

Stage explicit paths rather than `git add -A`, and run `git checkout -- ui/tsconfig.tsbuildinfo` before committing.

## A locally built sibling package makes a green run meaningless

When a change depends on an unreleased `@sistent/sistent` (or any sibling-repo package), `ui/node_modules/@sistent/sistent` is often overwritten in place with a locally built dist. A local test run is then green against code that is not published, and CI fails on the same commit.

Re-verify with `npm ci` after any local sibling build before trusting a green run or declaring a dependency bump done. A local build usually keeps the published version string, so **matching versions are not evidence that the installed contents are the published ones**. A version *mismatch* against `ui/package.json` is a useful tell that this has happened, but a match proves nothing.

## A `@sistent/sistent` bump must cover three manifests, and its peers are not optional

Three `package.json` files consume `@sistent/sistent` - `ui/`, `provider-ui/` and `install/docker-extension/ui/` - and a bump must cover all three.

Sistent's peers (`@mui/x-date-pickers`, `date-fns`, the `@rjsf/*` set, `xstate`/`@xstate/react`) are not optional. A consumer that omits one still installs cleanly and fails only at bundle time with `Module not found` pointing *inside* `@sistent/sistent/dist`, which reads as a sistent bug rather than a missing peer.

`install/docker-extension/ui` installs with `--legacy-peer-deps` because `@docker/docker-mui-theme` pins MUI <= 6 against sistent's MUI 9 - see its Dockerfile.

## `@meshery/schemas` keeps its `latest` tag below its highest version

`@meshery/schemas` deliberately keeps its `latest` dist-tag *below* its highest semver (1.4.0 is stale). That is safe because npm prefers the `latest`-tagged version whenever it satisfies the range, so the `^1.3.x` carets do not jump to 1.4.0 - but verify the resolved version in every regenerated lockfile rather than assuming it.

## `make docs-mesheryctl` rewrites ~100 pages with your home directory

`make docs-mesheryctl` (that is, `cd mesheryctl/doc && go run doc.go`) bakes the machine's `$HOME` into every generated page's "Options inherited from parent commands" block, via the `--config` default path. Running it locally rewrites all ~100 pages under `docs/content/en/reference/references/mesheryctl/` with your local home directory, even when only one command changed.

Committed docs use `/home/runner/...`, the GitHub Actions runner home. After regenerating:

1. `git diff --stat` the docs directory.
2. `git checkout --` every file whose only change is that path.
3. Manually restore `/home/runner/...` in the pages you actually intended to change.
