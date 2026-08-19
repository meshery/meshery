import { execFileSync, spawnSync } from 'node:child_process';
import { mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';

/**
 * `.agents/hooks/require-signoff.sh` is the sign-off guard that survives a
 * repointed `core.hooksPath`. `ui/.husky/commit-msg` applies the same rule, but
 * git reaches it only through that config, so a harness that points it at its
 * own hooks directory disarms it without saying so - which is how PR #21515
 * came to carry a commit with no `Signed-off-by` trailer, and by then the only
 * remedy left was rewriting the branch.
 *
 * The guard reads a proposed command line instead of git's hook plumbing, so
 * these run the real script against real repositories and assert the verdict it
 * returns: exit 1 blocks the command, exit 0 lets it through.
 */
// Resolved from the repository root rather than from this file's own URL or
// from `process.cwd()`: the guard lives outside `ui/`, and both of those move
// with how the suite was invoked.
const REPO_ROOT = execFileSync('git', ['rev-parse', '--show-toplevel'], {
  encoding: 'utf8',
}).trim();
const GUARD = join(REPO_ROOT, '.agents', 'hooks', 'require-signoff.sh');

const AUTHOR_NAME = 'Jane Dev';
const AUTHOR_EMAIL = 'jane@example.com';
const AUTHOR = `${AUTHOR_NAME} <${AUTHOR_EMAIL}>`;

/**
 * A developer's own git config is pinned away from these repositories. A global
 * `prepare-commit-msg` hook that appends the trailer - a common way to satisfy
 * DCO by habit - would otherwise sign every commit written here and make the
 * unsigned cases unreachable, so the suite would pass while proving nothing.
 */
const GIT_ENV = {
  ...process.env,
  GIT_CONFIG_GLOBAL: '/dev/null',
  GIT_CONFIG_SYSTEM: '/dev/null',
};

const git = (cwd: string, ...args: string[]) =>
  execFileSync('git', args, { cwd, env: GIT_ENV, encoding: 'utf8' });

const makeRepo = () => {
  const cwd = mkdtempSync(join(tmpdir(), 'require-signoff-'));

  git(cwd, 'init', '--quiet', '--initial-branch=main', '.');
  git(cwd, 'config', 'user.name', AUTHOR_NAME);
  git(cwd, 'config', 'user.email', AUTHOR_EMAIL);
  // The guard must not depend on a hook path being installed, and neither may
  // the fixtures that build these commits.
  git(cwd, 'config', 'core.hooksPath', join(cwd, 'no-such-hooks'));

  return cwd;
};

const guard = (cwd: string, command: string) => {
  const { status, stderr } = spawnSync(GUARD, [command], { cwd, env: GIT_ENV, encoding: 'utf8' });

  return { blocked: status !== 0, stderr };
};

// The guard is a bash script wired into a POSIX pre-command hook; on Windows the
// hook it belongs to does not run.
describe.skipIf(process.platform === 'win32')('git commit sign-off guard', () => {
  let unsignedHead: string;
  let signedHead: string;
  let merging: string;

  beforeAll(() => {
    unsignedHead = makeRepo();
    git(unsignedHead, 'commit', '--quiet', '--allow-empty', '-m', 'unsigned base');

    signedHead = makeRepo();
    git(signedHead, 'commit', '--quiet', '--allow-empty', '-s', '-m', 'signed base');

    merging = makeRepo();
    git(merging, 'commit', '--quiet', '--allow-empty', '-s', '-m', 'base');
    writeFileSync(join(merging, '.git', 'MERGE_HEAD'), git(merging, 'rev-parse', 'HEAD'));
  });

  afterAll(() => {
    for (const repo of [unsignedHead, signedHead, merging]) {
      rmSync(repo, { recursive: true, force: true });
    }
  });

  it('blocks the unsigned commit the DCO check fails on, and names the remedy', () => {
    const { blocked, stderr } = guard(unsignedHead, 'git commit -m "no sign-off"');

    expect(blocked).toBe(true);
    expect(stderr).toContain('git commit -s');
  });

  it('lets a signed commit through', () => {
    expect(guard(unsignedHead, 'git commit -s -m "signed"').blocked).toBe(false);
    expect(guard(unsignedHead, 'git commit --signoff -m "signed"').blocked).toBe(false);
  });

  it('reads a bundled -s, and does not read one that -m has swallowed', () => {
    // `-sm "msg"` is `-s -m "msg"`, but `-ms "msg"` is `-m s` - the sign-off
    // letter is the message body, and the commit is not signed at all.
    expect(guard(unsignedHead, 'git commit -sm "signed"').blocked).toBe(false);
    expect(guard(unsignedHead, 'git commit -ms "signed"').blocked).toBe(true);
  });

  it('accepts a trailer written out in the message', () => {
    expect(guard(unsignedHead, `git commit -m "body\n\nSigned-off-by: ${AUTHOR}"`).blocked).toBe(
      false,
    );
  });

  it('refuses an explicit --no-signoff even beside -s', () => {
    expect(guard(unsignedHead, 'git commit --no-signoff -s -m "x"').blocked).toBe(true);
  });

  it('does not read flags or a subcommand out of a quoted message body', () => {
    // A message may say anything, including the name of a flag it does not pass.
    expect(guard(unsignedHead, 'git commit -m "restore the -s flag"').blocked).toBe(true);
    expect(guard(unsignedHead, 'git commit -m "drop --signoff"').blocked).toBe(true);
    expect(guard(unsignedHead, 'echo "remember to git commit"').blocked).toBe(false);
  });

  it('ignores commands that create no commit', () => {
    expect(guard(unsignedHead, 'git status').blocked).toBe(false);
    expect(guard(unsignedHead, 'git log commit').blocked).toBe(false);
    expect(guard(unsignedHead, 'git push origin HEAD').blocked).toBe(false);
  });

  it('finds the commit behind global options and shell separators', () => {
    // `-C` takes its value as a separate token, so the subcommand is three
    // tokens along rather than adjacent.
    expect(guard(unsignedHead, `git -C ${unsignedHead} commit -m "x"`).blocked).toBe(true);
    expect(guard(unsignedHead, 'git --no-pager commit -m "x"').blocked).toBe(true);
    expect(guard(unsignedHead, 'cd ui && git commit -m "x"').blocked).toBe(true);
  });

  it('judges an amend by the message it will reuse', () => {
    expect(guard(unsignedHead, 'git commit --amend --no-edit').blocked).toBe(true);
    expect(guard(signedHead, 'git commit --amend --no-edit').blocked).toBe(false);
  });

  it('reads the message file -F names', () => {
    const signed = join(unsignedHead, 'signed.txt');
    const unsigned = join(unsignedHead, 'unsigned.txt');

    writeFileSync(signed, `body\n\nSigned-off-by: ${AUTHOR}\n`);
    writeFileSync(unsigned, 'body only\n');

    expect(guard(unsignedHead, `git commit -F ${signed}`).blocked).toBe(false);
    expect(guard(unsignedHead, `git commit -F ${unsigned}`).blocked).toBe(true);
  });

  it('exempts a merge commit, as the DCO check does', () => {
    expect(guard(merging, 'git commit -m "Merge branch \'topic\'"').blocked).toBe(false);
  });
});

/**
 * The rules above are only worth as much as the wiring that runs them. The
 * guard is harness-agnostic by design - a command line in, a verdict out as an
 * exit status - and Claude Code speaks neither: it hands a tool call to a hook
 * as JSON on stdin, and reads a plain non-zero exit as a *non-blocking* error,
 * which lets the commit proceed. `.claude/adapters/require-signoff.sh` adapts
 * between the two, and `.claude/settings.json` declares it.
 *
 * These drive that declared wiring the way the harness does: the command is
 * read out of the settings file and run through a shell with
 * `CLAUDE_PROJECT_DIR` set, so its quoting and expansion are exercised rather
 * than reproduced here, and the decision is read back off stdout.
 */
type WiredHook = { type?: string; command?: string };
type HookMatcher = { matcher?: string; hooks?: WiredHook[] };

const bashPreToolUseHooks = (): string[] => {
  const settings = JSON.parse(readFileSync(join(REPO_ROOT, '.claude', 'settings.json'), 'utf8'));
  const matchers: HookMatcher[] = settings?.hooks?.PreToolUse ?? [];

  return matchers
    .filter(
      ({ matcher }) => matcher === undefined || matcher === '*' || new RegExp(matcher).test('Bash'),
    )
    .flatMap(({ hooks }) => hooks ?? [])
    .filter((hook): hook is Required<WiredHook> => hook.type === 'command' && !!hook.command)
    .map(({ command }) => command);
};

const runWired = (cwd: string, payload: Record<string, unknown>) => {
  const results = bashPreToolUseHooks().map((command) =>
    spawnSync('/bin/sh', ['-c', command], {
      cwd,
      env: { ...GIT_ENV, CLAUDE_PROJECT_DIR: REPO_ROOT },
      input: JSON.stringify(payload),
      encoding: 'utf8',
    }),
  );

  const denial = results
    .map(({ status, stdout }) => ({
      status,
      output: stdout.trim() ? JSON.parse(stdout.trim()) : undefined,
    }))
    .find(({ output }) => output?.hookSpecificOutput?.permissionDecision === 'deny');

  return {
    decision: denial?.output.hookSpecificOutput.permissionDecision,
    reason: denial?.output.hookSpecificOutput.permissionDecisionReason ?? '',
    status: denial?.status,
  };
};

const bashCall = (cwd: string, command: string) => ({
  hook_event_name: 'PreToolUse',
  tool_name: 'Bash',
  cwd,
  tool_input: { command },
});

describe.skipIf(process.platform === 'win32')('git commit sign-off guard, as wired', () => {
  let repo: string;

  beforeAll(() => {
    repo = makeRepo();
    git(repo, 'commit', '--quiet', '--allow-empty', '-m', 'unsigned base');
  });

  afterAll(() => rmSync(repo, { recursive: true, force: true }));

  it('is armed for Bash tool calls', () => {
    expect(bashPreToolUseHooks().length).toBeGreaterThan(0);
  });

  it('denies the unsigned commit, and names the remedy', () => {
    const { decision, reason } = runWired(repo, bashCall(repo, 'git commit -m "no sign-off"'));

    expect(decision).toBe('deny');
    expect(reason).toContain('git commit -s');
  });

  it('reports the denial as a decision rather than as a hook failure', () => {
    // A non-zero exit is a *non-blocking* error here: the harness surfaces it
    // and runs the command anyway. Blocking has to come from the decision, so
    // the hook that denied has to have succeeded at denying.
    expect(runWired(repo, bashCall(repo, 'git commit -m "no sign-off"')).status).toBe(0);
  });

  it('lets a signed commit through', () => {
    expect(runWired(repo, bashCall(repo, 'git commit -s -m "signed"')).decision).toBeUndefined();
  });

  it('leaves a command that creates no commit alone', () => {
    expect(runWired(repo, bashCall(repo, 'git status')).decision).toBeUndefined();
  });

  it('leaves a tool call that is not Bash alone', () => {
    const write = {
      hook_event_name: 'PreToolUse',
      tool_name: 'Write',
      cwd: repo,
      tool_input: { file_path: join(repo, 'notes.md'), content: 'git commit -m "no sign-off"' },
    };

    expect(runWired(repo, write).decision).toBeUndefined();
  });
});
