import { spawnSync } from 'node:child_process';
import { existsSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { afterEach, describe, expect, it } from 'vitest';

// Behaviour tests for `ui/.husky/commit-msg`, the sign-off guard husky installs
// with the UI dependencies. They drive the hook the way git does - through
// `git commit` with `core.hooksPath` pointed at it, and by invoking it directly
// with a message file - and assert on what is observable: whether a commit was
// created, and what trailer it carries.
//
// This lives outside `tests/e2e` for the same reason as `testPlanLink.test.ts`:
// vitest excludes `tests/e2e`, and Playwright would otherwise try to run it.

// vitest rewrites `import.meta.url` away from a file URL under the jsdom
// environment, so locate the hook from the working directory instead - which is
// `ui/` for `npm test`, and the repository root for a run started from there.
const HOOK = [
  path.resolve(process.cwd(), '.husky/commit-msg'),
  path.resolve(process.cwd(), 'ui/.husky/commit-msg'),
].find((candidate) => existsSync(candidate)) as string;

// The hook is a POSIX shell script; git only ever runs it through a shell.
const posix = process.platform !== 'win32';

const repos: string[] = [];

/** A git repository isolated from the developer's global and system config. */
function initRepo(): string {
  const repo = mkdtempSync(path.join(tmpdir(), 'meshery-commit-msg-'));
  repos.push(repo);

  git(repo, ['init', '--quiet', '--initial-branch=master']);
  git(repo, ['config', 'user.name', 'Jane Smith']);
  git(repo, ['config', 'user.email', 'jane@example.com']);
  git(repo, ['config', 'commit.gpgsign', 'false']);

  // Install the hook under test the way husky does - its `_/h` runner executes
  // `sh -e <hook> "$@"`, and errexit changes how the hook has to be written - but
  // in a directory of its own, so the sibling `pre-commit` hook (which shells
  // into `ui/`) stays out of the way.
  const hooks = path.join(repo, 'hooks');
  mkdirSync(hooks);
  writeFileSync(path.join(hooks, 'commit-msg'), `#!/usr/bin/env sh\nsh -e "${HOOK}" "$@"\n`, {
    mode: 0o755,
  });
  git(repo, ['config', 'core.hooksPath', hooks]);

  return repo;
}

function git(repo: string, args: string[], env: NodeJS.ProcessEnv = {}) {
  return spawnSync('git', args, {
    cwd: repo,
    encoding: 'utf8',
    env: {
      ...process.env,
      GIT_CONFIG_GLOBAL: '/dev/null',
      GIT_CONFIG_SYSTEM: '/dev/null',
      ...env,
    },
  });
}

/** Stage a new file so there is always something to commit. */
function stageChange(repo: string, name: string) {
  writeFileSync(path.join(repo, name), `${name}\n`);
  git(repo, ['add', name]);
}

function commitCount(repo: string): number {
  const listed = git(repo, ['rev-list', '--count', 'HEAD']);
  return listed.status === 0 ? Number(listed.stdout.trim()) : 0;
}

/** Invoke the hook with the path to a message file, as husky's runner does. */
function runHook(repo: string, message: string) {
  const file = path.join(repo, 'COMMIT_MSG_UNDER_TEST');
  writeFileSync(file, message);
  return spawnSync('sh', ['-e', HOOK, file], { cwd: repo, encoding: 'utf8' });
}

afterEach(() => {
  while (repos.length) {
    rmSync(repos.pop() as string, { recursive: true, force: true });
  }
});

describe.skipIf(!posix)('ui/.husky/commit-msg', () => {
  it('refuses a commit whose message carries no sign-off', () => {
    const repo = initRepo();
    stageChange(repo, 'a.txt');

    const committed = git(repo, ['commit', '-m', '[Docs] unsigned']);

    expect(committed.status).not.toBe(0);
    expect(commitCount(repo)).toBe(0);
    expect(committed.stderr).toContain('not signed off');
    // The message is not lost - the hook tells the author how to recover it.
    expect(committed.stderr).toContain('git commit -s');
  });

  it('accepts a commit signed off with -s and records the trailer', () => {
    const repo = initRepo();
    stageChange(repo, 'a.txt');

    const committed = git(repo, ['commit', '-s', '-m', '[Docs] signed']);

    expect(committed.status).toBe(0);
    expect(commitCount(repo)).toBe(1);
    expect(git(repo, ['log', '-1', '--format=%(trailers:only=true,unfold=true)']).stdout).toContain(
      'Signed-off-by: Jane Smith <jane@example.com>',
    );
  });

  // The DCO check accepts any spacing after the colon, so the hook has to as
  // well: rejecting a trailer CI accepts sends the author rewriting a commit
  // that was never wrong.
  it('accepts a sign-off written with more than one space after the colon', () => {
    const repo = initRepo();
    stageChange(repo, 'a.txt');

    const committed = git(repo, [
      'commit',
      '-m',
      '[Docs] spaced trailer\n\nSigned-off-by:  Jane Smith <jane@example.com>',
    ]);

    expect(committed.status).toBe(0);
    expect(commitCount(repo)).toBe(1);
  });

  it('accepts a sign-off naming the committer when the author differs', () => {
    const repo = initRepo();
    stageChange(repo, 'a.txt');

    const committed = git(
      repo,
      ['commit', '-m', '[Docs] carried\n\nSigned-off-by: Jane Smith <jane@example.com>'],
      { GIT_AUTHOR_NAME: 'Other Person', GIT_AUTHOR_EMAIL: 'other@example.com' },
    );

    expect(committed.status).toBe(0);
    expect(commitCount(repo)).toBe(1);
  });

  it('refuses a sign-off that names neither the author nor the committer', () => {
    const repo = initRepo();
    stageChange(repo, 'a.txt');

    const committed = git(repo, [
      'commit',
      '-m',
      '[Docs] wrong signer\n\nSigned-off-by: John Doe <john@example.com>',
    ]);

    expect(committed.status).not.toBe(0);
    expect(commitCount(repo)).toBe(0);
    expect(committed.stderr).toContain('Signed-off-by: John Doe <john@example.com>');
    expect(committed.stderr).toContain('Expected: Signed-off-by: Jane Smith <jane@example.com>');
  });

  it('exempts a merge commit, as the DCO check does', () => {
    const repo = initRepo();
    stageChange(repo, 'f.txt');
    git(repo, ['commit', '-s', '-m', '[Docs] base']);
    git(repo, ['checkout', '--quiet', '-b', 'topic']);
    writeFileSync(path.join(repo, 'f.txt'), 'topic\n');
    git(repo, ['commit', '-s', '-a', '-m', '[Docs] topic']);
    git(repo, ['checkout', '--quiet', 'master']);
    writeFileSync(path.join(repo, 'f.txt'), 'master\n');
    git(repo, ['commit', '-s', '-a', '-m', '[Docs] master']);

    // Conflict, so the merge commit goes through `git commit` - the path the
    // hook actually sees - with the merge message and no sign-off.
    expect(git(repo, ['merge', 'topic']).status).not.toBe(0);
    writeFileSync(path.join(repo, 'f.txt'), 'resolved\n');
    git(repo, ['add', 'f.txt']);
    const merged = git(repo, ['commit', '--no-edit']);

    expect(merged.status).toBe(0);
    expect(git(repo, ['log', '-1', '--format=%p']).stdout.trim().split(' ')).toHaveLength(2);
  });

  it('does not count a sign-off that git will strip as a comment', () => {
    const repo = initRepo();

    const checked = runHook(
      repo,
      '[Docs] commented out\n\n# Signed-off-by: Jane Smith <jane@example.com>\n',
    );

    expect(checked.status).not.toBe(0);
    expect(checked.stderr).toContain('carries no Signed-off-by trailer');
  });

  it('does not count a sign-off from the diff `git commit --verbose` appends', () => {
    const repo = initRepo();

    const checked = runHook(
      repo,
      [
        '[Docs] verbose',
        '# ------------------------ >8 ------------------------',
        'diff --git a/CONTRIBUTING.md b/CONTRIBUTING.md',
        '+Signed-off-by: Jane Smith <jane@example.com>',
        '',
      ].join('\n'),
    );

    expect(checked.status).not.toBe(0);
    expect(checked.stderr).toContain('carries no Signed-off-by trailer');
  });
});
