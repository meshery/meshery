"""Exercise a GitHub Pages deployment racing a preview push using local repos."""

import os
import pathlib
import shutil
import subprocess
import tempfile


def run(*args, cwd=None, env=None):
    return subprocess.run(
        args, cwd=cwd, env=env, check=True, text=True,
        stdout=subprocess.PIPE, stderr=subprocess.PIPE,
    )


with tempfile.TemporaryDirectory(prefix='meshery-pages-race-') as temp_dir:
    root = pathlib.Path(temp_dir)
    remote = root / 'remote.git'
    seed, pages, rival = (root / name for name in ('seed', 'pages', 'rival'))
    site, tools, verify = (root / name for name in ('site', 'bin', 'verify'))

    run('git', 'init', '--bare', str(remote))
    run('git', 'clone', str(remote), str(seed))
    run('git', 'config', 'user.name', 'Test', cwd=seed)
    run('git', 'config', 'user.email', 'test@example.com', cwd=seed)
    run('git', 'checkout', '-b', 'gh-pages', cwd=seed)
    (seed / 'pr-preview' / 'pr-1').mkdir(parents=True)
    (seed / 'pr-preview' / 'pr-1' / 'index.html').write_text('old preview')
    (seed / 'old.html').write_text('stale')
    (seed / 'index.html').write_text('old site')
    run('git', 'add', '.', cwd=seed)
    run('git', 'commit', '-m', 'seed', cwd=seed)
    run('git', 'push', 'origin', 'gh-pages', cwd=seed)
    run('git', 'clone', '-b', 'gh-pages', str(remote), str(pages))
    run('git', 'clone', '-b', 'gh-pages', str(remote), str(rival))
    run('git', 'config', 'user.name', 'Rival', cwd=rival)
    run('git', 'config', 'user.email', 'rival@example.com', cwd=rival)
    site.mkdir()
    (site / 'index.html').write_text('new site')
    (site / '.nojekyll').touch()
    tools.mkdir()

    # Inject a preview push immediately before the first deployment push.
    wrapper = tools / 'git'
    wrapper.write_text('''#!/bin/bash
if [ "$1" = push ] && [ ! -e "$RACE_MARKER" ]; then
  touch "$RACE_MARKER"
  mkdir -p "$RIVAL/pr-preview/pr-2"
  printf 'new preview' > "$RIVAL/pr-preview/pr-2/index.html"
  "$REAL_GIT" -C "$RIVAL" add pr-preview
  "$REAL_GIT" -C "$RIVAL" commit -qm 'concurrent preview'
  "$REAL_GIT" -C "$RIVAL" push -q origin HEAD:gh-pages
fi
exec "$REAL_GIT" "$@"
''')
    wrapper.chmod(0o755)
    env = os.environ.copy()
    env.update(
        PATH=str(tools) + os.pathsep + env['PATH'],
        RACE_MARKER=str(root / 'race'), RIVAL=str(rival),
        REAL_GIT=shutil.which('git'), GITHUB_SHA='test-sha',
    )
    script = pathlib.Path(__file__).resolve().parent.parent / 'deploy-docs-pages.sh'
    result = run('bash', str(script), str(site), str(pages), env=env)
    run('git', 'clone', '-b', 'gh-pages', str(remote), str(verify))
    assert (root / 'race').exists() and 'refreshing and retrying' in result.stdout
    assert (verify / 'index.html').read_text() == 'new site'
    assert (verify / 'pr-preview' / 'pr-1' / 'index.html').read_text() == 'old preview'
    assert (verify / 'pr-preview' / 'pr-2' / 'index.html').read_text() == 'new preview'
    assert (verify / '.nojekyll').exists()
    assert not (verify / 'old.html').exists()
    result = run('bash', str(script), str(site), str(pages), env=env)
    assert 'already up to date' in result.stdout
    print('PASS: retry preserves both previews and deploys the updated site')
