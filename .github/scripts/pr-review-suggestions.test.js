#!/usr/bin/env node

// Regression and security tests for PR Review Suggestions workflow.
// Run with: node --test .github/scripts/pr-review-suggestions.test.js
// Uses Node's built-in test runner + assert (Node >= 18). Zero external dependencies.

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

// Extract and load the script directly from .github/workflows/pr-review-suggestions.yml
async function loadWorkflowModule() {
  const workflowPath = path.resolve(__dirname, '../workflows/pr-review-suggestions.yml');
  const yamlContent = fs.readFileSync(workflowPath, 'utf8');

  const marker = 'script: |\n';
  const startIndex = yamlContent.indexOf(marker);
  if (startIndex === -1) {
    throw new Error('Could not find "script: |" in workflow YAML');
  }

  const scriptLines = yamlContent.slice(startIndex + marker.length).split('\n');
  const scriptContent = scriptLines
    .map(line => (line.startsWith('            ') ? line.slice(12) : line))
    .join('\n');

  const mod = { exports: {} };
  const AsyncFunction = Object.getPrototypeOf(async () => {}).constructor;
  const fn = new AsyncFunction(
    'require',
    '__original_require__',
    'module',
    'exports',
    'console',
    'process',
    'Buffer',
    'global',
    'io',
    'exec',
    'core',
    'github',
    'context',
    scriptContent
  );

  // Initialize module to populate exports
  await fn(
    require,
    require,
    mod,
    mod.exports,
    console,
    process,
    Buffer,
    global,
    {},
    {},
    { info: () => {}, warning: () => {} },
    {},
    {}
  );

  return mod.exports;
}

// Mock GitHub API factory
function createMockGitHub({ prDetails = {}, openPRs = [], comments = [] } = {}) {
  const createdComments = [];
  const updatedComments = [];

  const github = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          if (prDetails[pull_number]) {
            return { data: prDetails[pull_number] };
          }
          throw new Error(`PR #${pull_number} not found`);
        },
        list: async ({ state, base, head } = {}) => {
          let filtered = [...openPRs];
          if (state) {
            filtered = filtered.filter(p => p.state === state);
          }
          if (base) {
            filtered = filtered.filter(p => !p.base || p.base.ref === base);
          }
          if (head) {
            filtered = filtered.filter(p => {
              if (!p.head) return false;
              const owner = p.head.user
                ? p.head.user.login
                : p.head.repo && p.head.repo.owner && p.head.repo.owner.login;
              const branch = p.head.ref;
              return `${owner}:${branch}` === head;
            });
          }
          return { data: filtered };
        }
      },
      repos: {
        listPullRequestsAssociatedWithCommit: async ({ commit_sha }) => {
          const matched = openPRs.filter(p => p.head && p.head.sha === commit_sha);
          return { data: matched };
        }
      },
      issues: {
        listComments: async () => ({ data: comments }),
        createComment: async ({ issue_number, body }) => {
          const newComment = {
            id: 999000 + createdComments.length + 1,
            body,
            user: { login: 'github-actions[bot]', type: 'Bot' }
          };
          createdComments.push({ issue_number, body, id: newComment.id });
          comments.push(newComment);
          return { data: newComment };
        },
        updateComment: async ({ comment_id, body }) => {
          updatedComments.push({ comment_id, body });
          const target = comments.find(c => c.id === comment_id);
          if (target) target.body = body;
          return { data: { id: comment_id, body } };
        }
      }
    },
    paginate: async (method, params) => {
      const res = await method(params);
      return res.data || [];
    }
  };

  return { github, createdComments, updatedComments };
}

function createMockCore() {
  const logs = { info: [], warning: [] };
  return {
    logs,
    core: {
      info: msg => logs.info.push(msg),
      warning: msg => logs.warning.push(msg)
    }
  };
}

// ---------------------------------------------------------------------------
// 1. Text Normalization and Tokenization Tests
// ---------------------------------------------------------------------------
test('tokenizeText: handles hyphens and slashes as delimiters', async () => {
  const { tokenizeText } = await loadWorkflowModule();

  const tokens = tokenizeText('server-side component/server GraphQL-handler');
  assert.ok(tokens.has('server'), 'server token should exist');
  assert.ok(tokens.has('side'), 'side token should exist');
  assert.ok(tokens.has('component'), 'component token should exist');
  assert.ok(tokens.has('graphql'), 'graphql token should exist');
  assert.ok(tokens.has('handler'), 'handler token should exist');
  assert.ok(!tokens.has('server-side'), 'hyphenated compound should be split');
  assert.ok(!tokens.has('component/server'), 'slash compound should be split');
});

test('tokenizeText: filters generic stopwords and short tokens', async () => {
  const { tokenizeText } = await loadWorkflowModule();

  const tokens = tokenizeText('The fix for a new pr issue with an update');
  assert.equal(tokens.size, 0, 'all stopwords should be filtered');

  const singleChar = tokenizeText('a b c x y z');
  assert.equal(singleChar.size, 0, 'single characters should be ignored');
});

test('tokenizeText: deduplicates repeated tokens (set semantics)', async () => {
  const { tokenizeText } = await loadWorkflowModule();

  const tokens = tokenizeText('kubernetes kubernetes kubernetes cluster cluster');
  assert.equal(tokens.size, 2);
  assert.ok(tokens.has('kubernetes'));
  assert.ok(tokens.has('cluster'));
});

test('tokenizeText: handles punctuation-only and empty inputs safely', async () => {
  const { tokenizeText } = await loadWorkflowModule();

  assert.equal(tokenizeText('').size, 0);
  assert.equal(tokenizeText(null).size, 0);
  assert.equal(tokenizeText(undefined).size, 0);
  assert.equal(tokenizeText('   --- /// ___ ... !!! ???   ').size, 0);
});

// ---------------------------------------------------------------------------
// 2. Label Normalization Tests
// ---------------------------------------------------------------------------
test('normalizeLabels: categorizes language and area/component labels case-insensitively', async () => {
  const { normalizeLabels } = await loadWorkflowModule();

  const input = [
    { name: 'Language/Go' },
    { name: 'AREA/CI' },
    { name: 'Component/Server' },
    { name: 'kind/bug' },
    { name: 'help wanted' },
    'language/typescript'
  ];

  const { languageLabels, areaComponentLabels } = normalizeLabels(input);
  assert.deepEqual(languageLabels.sort(), ['language/go', 'language/typescript'].sort());
  assert.deepEqual(areaComponentLabels.sort(), ['area/ci', 'component/server'].sort());
});

test('normalizeLabels: ignores unrelated label families', async () => {
  const { normalizeLabels } = await loadWorkflowModule();

  const input = ['kind/enhancement', 'priority/high', 'status/review', 'hacktoberfest'];
  const { languageLabels, areaComponentLabels } = normalizeLabels(input);
  assert.equal(languageLabels.length, 0);
  assert.equal(areaComponentLabels.length, 0);
});

// ---------------------------------------------------------------------------
// 3. Markdown Escaping & Security Tests
// ---------------------------------------------------------------------------
test('escapeMarkdown: escapes markdown link syntax without breaking text and neutralizes bare URLs', async () => {
  const { escapeMarkdown } = await loadWorkflowModule();

  const title = 'Fix [review checklist](https://evil.invalid) bug';
  const escaped = escapeMarkdown(title);
  assert.ok(!escaped.includes('[review checklist](https://evil.invalid)'), 'link should not be unescaped');
  assert.ok(escaped.includes(':\u200b//'), 'URL protocol should have zero-width space inserted');
  assert.equal(escaped, 'Fix \\[review checklist\\]\\(https:\u200b//evil.invalid\\) bug');
});

test('escapeMarkdown: neutralizes @ mentions with zero-width joiner', async () => {
  const { escapeMarkdown } = await loadWorkflowModule();

  const title = 'Bug report for @meshery/maintainers and @leecalcote';
  const escaped = escapeMarkdown(title);
  assert.ok(escaped.includes('@\u200dmeshery/maintainers'));
  assert.ok(escaped.includes('@\u200dleecalcote'));
});

test('escapeMarkdown: escapes backticks, asterisks, underscores, pipes, and tags', async () => {
  const { escapeMarkdown } = await loadWorkflowModule();

  const title = '`rm -rf` **bold** _italic_ <script>alert(1)</script> | pipe #123';
  const escaped = escapeMarkdown(title);
  assert.ok(escaped.includes('\\`rm -rf\\`'));
  assert.ok(escaped.includes('\\*\\*bold\\*\\*'));
  assert.ok(escaped.includes('\\_italic\\_'));
  assert.ok(escaped.includes('\\<script\\>alert\\(1\\)\\</script\\>'));
  assert.ok(escaped.includes('\\|'));
  assert.ok(escaped.includes('\\#123'));
});

test('escapeMarkdown: normalizes newlines and carriage returns', async () => {
  const { escapeMarkdown } = await loadWorkflowModule();

  const title = 'Title with\r\nnewlines and \n tabs\tand spaces';
  const escaped = escapeMarkdown(title);
  assert.equal(escaped, 'Title with newlines and tabs and spaces');
});

// ---------------------------------------------------------------------------
// 4. Scoring, Eligibility, and Ranking Tests
// ---------------------------------------------------------------------------
test('calculateScore: candidate with no label overlap is strictly ineligible (score 0)', async () => {
  const { calculateScore, tokenizeText } = await loadWorkflowModule();

  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: tokenizeText('fix environment pagination bug'),
    bodyTokens: tokenizeText('environment pagination controller')
  };

  const candidateNoLabelMatch = {
    languageLabels: ['language/typescript'],
    areaComponentLabels: ['area/docs'],
    titleTokens: tokenizeText('fix environment pagination bug'),
    bodyTokens: tokenizeText('environment pagination controller')
  };

  const score = calculateScore(current, candidateNoLabelMatch);
  assert.equal(score, 0, 'text overlap without label overlap must score 0');
});

test('calculateScore: applies documented point values accurately', async () => {
  const { calculateScore, tokenizeText } = await loadWorkflowModule();

  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: tokenizeText('database migration handler'),
    bodyTokens: tokenizeText('postgresql connection pool')
  };

  // Case 1: Shared language only (+40) + 1 title token (+10) = 50
  const candidateLangOnly = {
    languageLabels: ['language/go'],
    areaComponentLabels: [],
    titleTokens: tokenizeText('database schema generator'),
    bodyTokens: tokenizeText('unrelated description')
  };
  assert.equal(calculateScore(current, candidateLangOnly), 50);

  // Case 2: Shared area/component only (+40) + 1 body token (+5) = 45
  const candidateAreaOnly = {
    languageLabels: [],
    areaComponentLabels: ['component/server'],
    titleTokens: tokenizeText('different title'),
    bodyTokens: tokenizeText('postgresql adapter')
  };
  assert.equal(calculateScore(current, candidateAreaOnly), 45);

  // Case 3: Both language (+40) and component (+40) + 2 title (+20) + 2 body (+10) = 110
  const candidateBoth = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: tokenizeText('database migration scripts'),
    bodyTokens: tokenizeText('postgresql connection metrics')
  };
  assert.equal(calculateScore(current, candidateBoth), 110);
});

// ---------------------------------------------------------------------------
// 5. Sticky Comment Ownership Tests
// ---------------------------------------------------------------------------
test('isBotStickyComment: only matches comments authored by github-actions[bot] containing marker', async () => {
  const { isBotStickyComment } = await loadWorkflowModule();
  const MARKER = '<!-- meshery-review-suggestions -->';

  // Valid bot comment
  assert.ok(
    isBotStickyComment(
      { body: `${MARKER}\nContent`, user: { login: 'github-actions[bot]', type: 'Bot' } },
      MARKER
    )
  );

  // Human comment with marker: must NOT match
  assert.ok(
    !isBotStickyComment(
      { body: `${MARKER}\nHuman mimic comment`, user: { login: 'random-contributor', type: 'User' } },
      MARKER
    )
  );

  // Unrelated bot comment with marker: must NOT match
  assert.ok(
    !isBotStickyComment(
      { body: `${MARKER}\nCodeRabbit review`, user: { login: 'coderabbitai[bot]', type: 'Bot' } },
      MARKER
    )
  );

  // Bot comment without marker: must NOT match
  assert.ok(
    !isBotStickyComment(
      { body: 'Build succeeded', user: { login: 'github-actions[bot]', type: 'Bot' } },
      MARKER
    )
  );
});

// ---------------------------------------------------------------------------
// 6. Full End-to-End Workflow Execution Scenarios
// ---------------------------------------------------------------------------
test('workflow run: creates new sticky comment for eligible candidates', async () => {
  const { run } = await loadWorkflowModule();

  const pr100 = {
    number: 100,
    title: '[Server] Add database connection pooling',
    body: 'Implements connection pool for PostgreSQL database',
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }, { name: 'component/server' }]
  };

  const openPRs = [
    pr100,
    {
      number: 101,
      title: '[Server] Database migration runner',
      body: 'Runs database migrations during startup',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/101',
      labels: [{ name: 'language/go' }, { name: 'component/server' }]
    },
    {
      number: 102,
      title: '[Docs] Update server docs',
      body: 'Documentation for server configuration',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/102',
      labels: [{ name: 'area/docs' }] // Ineligible: zero label overlap
    },
    {
      number: 103,
      title: '[Draft] Draft server work',
      body: 'Unfinished draft PR',
      state: 'open',
      base: { ref: 'master' },
      draft: true, // Ineligible: draft
      html_url: 'https://github.com/meshery/meshery/pull/103',
      labels: [{ name: 'language/go' }]
    }
  ];

  const { github, createdComments } = createMockGitHub({
    prDetails: { 100: pr100 },
    openPRs,
    comments: []
  });
  const { core } = createMockCore();

  const context = {
    eventName: 'pull_request_target',
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 100 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1, 'should create exactly one sticky comment');
  assert.ok(createdComments[0].body.includes('<!-- meshery-review-suggestions -->'));
  assert.ok(createdComments[0].body.includes('[#101](https://github.com/meshery/meshery/pull/101)'));
  assert.ok(!createdComments[0].body.includes('#102'), 'ineligible PR #102 must not appear');
  assert.ok(!createdComments[0].body.includes('#103'), 'draft PR #103 must not appear');
});

test('workflow run: updates existing bot sticky comment and ignores human comments', async () => {
  const { run } = await loadWorkflowModule();

  const pr100 = {
    number: 100,
    title: '[Server] Add database connection pooling',
    body: 'Implements connection pool',
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }]
  };

  const humanComment = {
    id: 111,
    body: '<!-- meshery-review-suggestions --> Human wrote this marker',
    user: { login: 'contributor-joe', type: 'User' }
  };

  const botComment = {
    id: 222,
    body: '<!-- meshery-review-suggestions -->\n\nOld suggestions',
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const openPRs = [
    pr100,
    {
      number: 105,
      title: '[mesheryctl] Go CLI improvement',
      body: 'CLI change',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/105',
      labels: [{ name: 'language/go' }]
    }
  ];

  const { github, createdComments, updatedComments } = createMockGitHub({
    prDetails: { 100: pr100 },
    openPRs,
    comments: [humanComment, botComment]
  });
  const { core } = createMockCore();

  const context = {
    eventName: 'pull_request_target',
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 100 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0, 'should not create a new comment when bot comment exists');
  assert.equal(updatedComments.length, 1, 'should update existing bot comment');
  assert.equal(updatedComments[0].comment_id, 222, 'must update bot comment #222, NOT human comment #111');
});

test('workflow run: replaces stale recommendations with "no closely related open PRs" when labels are removed', async () => {
  const { run } = await loadWorkflowModule();

  // PR has no eligible labels remaining after an unlabeled event
  const pr200 = {
    number: 200,
    title: '[Docs] Small typo fix',
    body: 'Typo fix',
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'kind/documentation' }] // No language/area/component label
  };

  const botComment = {
    id: 333,
    body: '<!-- meshery-review-suggestions -->\n\n- [#101](https://github.com/meshery/meshery/pull/101) — Stale candidate',
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const openPRs = [
    pr200,
    {
      number: 101,
      title: '[Server] Some candidate',
      body: 'Some candidate',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/101',
      labels: [{ name: 'language/go' }]
    }
  ];

  const { github, updatedComments } = createMockGitHub({
    prDetails: { 200: pr200 },
    openPRs,
    comments: [botComment]
  });
  const { core } = createMockCore();

  const context = {
    eventName: 'pull_request_target',
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 200 }, action: 'unlabeled' }
  };

  await run({ github, context, core });

  assert.equal(updatedComments.length, 1);
  assert.ok(
    updatedComments[0].body.includes('no closely related open PRs were found'),
    'stale recommendation must be replaced with no-candidate notice'
  );
});

test('workflow run: handles workflow_run handoff after automated labeling completes', async () => {
  const { run } = await loadWorkflowModule();

  const pr300 = {
    number: 300,
    title: '[UI] Add dark mode toggle',
    body: 'Implements dark mode in UI',
    state: 'open',
    base: { ref: 'master' },
    head: { sha: 'abcdef1234567890' },
    labels: [{ name: 'component/ui' }, { name: 'language/typescript' }]
  };

  const openPRs = [
    pr300,
    {
      number: 301,
      title: '[UI] Fix theme switcher',
      body: 'Theme switcher bug fix',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/301',
      labels: [{ name: 'component/ui' }]
    }
  ];

  const { github, createdComments } = createMockGitHub({
    prDetails: { 300: pr300 },
    openPRs,
    comments: []
  });
  const { core } = createMockCore();

  const context = {
    eventName: 'workflow_run',
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: {
      workflow_run: {
        id: 777888,
        conclusion: 'success',
        pull_requests: [{ number: 300 }],
        head_sha: 'abcdef1234567890'
      }
    }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1);
  assert.ok(createdComments[0].body.includes('[#301](https://github.com/meshery/meshery/pull/301)'));
});

test('workflow run: handles fork workflow_run handoff with empty pull_requests array via headFilter', async () => {
  const { run } = await loadWorkflowModule();

  const forkPR = {
    number: 350,
    title: '[UI] Community fork feature',
    body: 'Fork contribution',
    state: 'open',
    base: { ref: 'master' },
    head: {
      sha: 'forksha123456789',
      ref: 'feature-branch',
      user: { login: 'contributor-user' }
    },
    labels: [{ name: 'component/ui' }]
  };

  const openPRs = [
    forkPR,
    {
      number: 351,
      title: '[UI] Component suggestions',
      body: 'Related suggestion',
      state: 'open',
      draft: false,
      base: { ref: 'master' },
      html_url: 'https://github.com/meshery/meshery/pull/351',
      labels: [{ name: 'component/ui' }]
    }
  ];

  const { github, createdComments } = createMockGitHub({
    prDetails: { 350: forkPR },
    openPRs,
    comments: []
  });
  const { core } = createMockCore();

  // Emulate GitHub Actions payload for fork PR: pull_requests is empty!
  const context = {
    eventName: 'workflow_run',
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: {
      workflow_run: {
        id: 998877,
        conclusion: 'success',
        pull_requests: [], // Empty on fork PRs!
        head_sha: 'forksha123456789',
        head_branch: 'feature-branch',
        head_repository: {
          owner: { login: 'contributor-user' }
        }
      }
    }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1, 'should locate fork PR and create sticky comment');
  assert.ok(createdComments[0].body.includes('[#351](https://github.com/meshery/meshery/pull/351)'));
});

test('workflow run: ignores non-master and closed PRs safely', async () => {
  const { run } = await loadWorkflowModule();

  const closedPR = {
    number: 400,
    title: 'Closed PR',
    state: 'closed',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }]
  };

  const nonMasterPR = {
    number: 401,
    title: 'PR targeting release branch',
    state: 'open',
    base: { ref: 'release-0.8' },
    labels: [{ name: 'language/go' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 400: closedPR, 401: nonMasterPR },
    openPRs: [],
    comments: []
  });
  const { core } = createMockCore();

  // Test closed PR
  await run({
    github,
    context: {
      eventName: 'pull_request_target',
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 400 } }
    },
    core
  });
  assert.equal(createdComments.length, 0, 'closed PR should not be commented on');

  // Test non-master PR
  await run({
    github,
    context: {
      eventName: 'pull_request_target',
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 401 } }
    },
    core
  });
  assert.equal(createdComments.length, 0, 'non-master PR should not be commented on');
});

test('workflow run: excludes candidate PRs targeting non-master branches', async () => {
  const { run } = await loadWorkflowModule();

  const pr700 = {
    number: 700,
    title: '[Server] Core feature',
    body: 'Core feature',
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  const candidateMaster = {
    number: 701,
    title: '[Server] Master candidate',
    body: 'Master candidate',
    state: 'open',
    draft: false,
    base: { ref: 'master' },
    html_url: 'https://github.com/meshery/meshery/pull/701',
    labels: [{ name: 'component/server' }]
  };

  const candidateRelease = {
    number: 702,
    title: '[Server] Release branch candidate',
    body: 'Release candidate',
    state: 'open',
    draft: false,
    base: { ref: 'release-0.8' }, // Non-master candidate
    html_url: 'https://github.com/meshery/meshery/pull/702',
    labels: [{ name: 'component/server' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 700: pr700 },
    openPRs: [pr700, candidateMaster, candidateRelease],
    comments: []
  });
  const { core } = createMockCore();

  await run({
    github,
    context: {
      eventName: 'pull_request_target',
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 700 } }
    },
    core
  });

  assert.equal(createdComments.length, 1);
  const comment = createdComments[0].body;
  assert.ok(comment.includes('#701'), 'master-targeting candidate must appear');
  assert.ok(!comment.includes('#702'), 'release-targeting candidate must be excluded');
});

test('workflow run: deterministic tie-breaking and max-3 candidate selection', async () => {
  const { run } = await loadWorkflowModule();

  const pr500 = {
    number: 500,
    title: '[Server] Server handler update',
    body: 'Server update',
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  // 5 eligible candidates with varying scores and duplicate scores
  const openPRs = [
    pr500,
    {
      number: 501,
      title: '[Server] Candidate 501',
      body: 'Candidate body',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/501',
      labels: [{ name: 'component/server' }] // score: 40
    },
    {
      number: 502,
      title: '[Server] Candidate 502',
      body: 'Candidate body',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/502',
      labels: [{ name: 'component/server' }] // score: 40 (newer than 501)
    },
    {
      number: 503,
      title: '[Server] Server handler Candidate 503', // shared title: 'server', 'handler' (+20) -> score: 60
      body: 'Candidate body',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/503',
      labels: [{ name: 'component/server' }]
    },
    {
      number: 504,
      title: '[Server] Candidate 504',
      body: 'Candidate body',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/504',
      labels: [{ name: 'component/server' }] // score: 40 (newer than 502)
    },
    {
      number: 505,
      title: '[Server] Candidate 505',
      body: 'Candidate body',
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/505',
      labels: [{ name: 'component/server' }] // score: 40 (newer than 504)
    }
  ];

  const { github, createdComments } = createMockGitHub({
    prDetails: { 500: pr500 },
    openPRs,
    comments: []
  });
  const { core } = createMockCore();

  await run({
    github,
    context: {
      eventName: 'pull_request_target',
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 500 } }
    },
    core
  });

  assert.equal(createdComments.length, 1);
  const comment = createdComments[0].body;

  // Highest score #503 (score 60) must be first
  // Then tie-break by PR number descending among score 40: #505, then #504
  // #502 and #501 must be sliced off (only top 3)
  assert.ok(comment.includes('#503'), 'top score #503 must appear');
  assert.ok(comment.includes('#505'), 'tie-break #505 must appear');
  assert.ok(comment.includes('#504'), 'tie-break #504 must appear');
  assert.ok(!comment.includes('#502'), 'candidate #502 must be excluded by top-3 cap');
  assert.ok(!comment.includes('#501'), 'candidate #501 must be excluded by top-3 cap');
});

test('workflow run: API errors fail safely without crashing or throwing', async () => {
  const { run } = await loadWorkflowModule();

  const brokenGitHub = {
    rest: {
      pulls: {
        get: async () => {
          throw new Error('GitHub API rate limit exceeded');
        }
      }
    }
  };
  const { core, logs } = createMockCore();

  await run({
    github: brokenGitHub,
    context: {
      eventName: 'pull_request_target',
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 600 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Could not fetch current PR details from API')),
    'API failure should be captured as core.warning'
  );
});
