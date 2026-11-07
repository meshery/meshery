#!/usr/bin/env node
'use strict';

// PR Review Suggestions Heuristic Test Suite
// Run with: node --test .github/scripts/pr-review-suggestions.test.js
// Uses Node.js built-in test runner + strict assert (Node >= 18). Zero external dependencies.

const test = require('node:test');
const assert = require('node:assert/strict');

const {
  STOP_WORDS,
  MARKER,
  normalizeLabels,
  tokenizeText,
  escapeMarkdown,
  calculateScore,
  isBotStickyComment,
  run
} = require('./pr-review-suggestions.js');

// ---------------------------------------------------------------------------
// Mock Helpers for GitHub API & Core Logging
// ---------------------------------------------------------------------------

function createMockGitHub({ prDetails = {}, openPRs = [], comments = [] } = {}) {
  const createdComments = [];
  const updatedComments = [];
  const deletedComments = [];

  const github = {
    rest: {
      pulls: {
        get: async ({ pull_number }) => {
          if (prDetails[pull_number]) {
            return { data: prDetails[pull_number] };
          }
          throw new Error(`PR #${pull_number} not found in mock store`);
        },
        list: async ({ state, base } = {}) => {
          let filtered = [...openPRs];
          if (state) {
            filtered = filtered.filter(p => p.state === state);
          }
          if (base) {
            filtered = filtered.filter(p => !p.base || p.base.ref === base);
          }
          return { data: filtered };
        }
      },
      issues: {
        listComments: async () => ({ data: [...comments] }),
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
        },
        deleteComment: async ({ comment_id }) => {
          deletedComments.push({ comment_id });
          const index = comments.findIndex(c => c.id === comment_id);
          if (index !== -1) comments.splice(index, 1);
          return { data: {} };
        }
      }
    },
    paginate: async (method, params) => {
      const res = await method(params);
      return res.data || [];
    }
  };

  return { github, createdComments, updatedComments, deletedComments, comments };
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

test('tokenizeText: handles hyphens and slashes as delimiters', () => {
  const tokens = tokenizeText('Fix server-side authentication in component/server and api_v1');
  assert.ok(tokens.has('server'), 'must include server');
  assert.ok(tokens.has('side'), 'must include side');
  assert.ok(tokens.has('authentication'), 'must include authentication');
  assert.ok(tokens.has('component'), 'must include component');
  assert.ok(tokens.has('api'), 'must include api');
  assert.ok(tokens.has('v1'), 'must include v1');
  assert.ok(!tokens.has('server-side'), 'compound hyphenated form should not remain');
});

test('tokenizeText: filters generic stopwords and short tokens', () => {
  const tokens = tokenizeText('a the this is for fix update test a v');
  assert.ok(!tokens.has('the'), 'stopword "the" should be filtered');
  assert.ok(!tokens.has('fix'), 'stopword "fix" should be filtered');
  assert.ok(!tokens.has('a'), '1-char token "a" should be filtered');
  assert.ok(!tokens.has('v'), '1-char token "v" should be filtered');
  assert.ok(tokens.has('test'), 'non-stopword "test" must be retained');
});

test('tokenizeText: deduplicates repeated tokens (set semantics)', () => {
  const tokens = tokenizeText('kubernetes kubernetes cluster cluster KUBERNETES');
  assert.equal(tokens.size, 2);
  assert.ok(tokens.has('kubernetes'));
  assert.ok(tokens.has('cluster'));
});

test('tokenizeText: handles punctuation-only and empty inputs safely', () => {
  assert.equal(tokenizeText('').size, 0);
  assert.equal(tokenizeText(null).size, 0);
  assert.equal(tokenizeText(undefined).size, 0);
  assert.equal(tokenizeText('---///...???!!!').size, 0);
});

// ---------------------------------------------------------------------------
// 2. Label Normalization Tests
// ---------------------------------------------------------------------------

test('normalizeLabels: categorizes language and area/component labels case-insensitively', () => {
  const input = [
    { name: 'LANGUAGE/GO' },
    { name: 'Component/Server' },
    { name: 'area/docs' },
    'language/typescript'
  ];
  const { languageLabels, areaComponentLabels } = normalizeLabels(input);

  assert.deepEqual(languageLabels.sort(), ['language/go', 'language/typescript'].sort());
  assert.deepEqual(areaComponentLabels.sort(), ['area/docs', 'component/server'].sort());
});

test('normalizeLabels: ignores unrelated label families and invalid objects', () => {
  const input = [
    { name: 'kind/bug' },
    { name: 'help wanted' },
    { name: 'good first issue' },
    null,
    undefined,
    { invalid: true }
  ];
  const { languageLabels, areaComponentLabels } = normalizeLabels(input);

  assert.equal(languageLabels.length, 0);
  assert.equal(areaComponentLabels.length, 0);
});

// ---------------------------------------------------------------------------
// 3. Markdown Escaping & Security Sanitization Tests
// ---------------------------------------------------------------------------

test('escapeMarkdown: escapes markdown link syntax without breaking text and neutralizes bare URLs', () => {
  const title = 'Fix [review checklist](https://evil.invalid) bug and visit http://phish.test';
  const escaped = escapeMarkdown(title);
  assert.ok(!escaped.includes('[review checklist](https://evil.invalid)'), 'link should not be unescaped');
  assert.ok(escaped.includes(':\u200b//'), 'URL protocol should have zero-width space inserted');
  assert.equal(escaped, 'Fix \\[review checklist\\]\\(https:\u200b//evil.invalid\\) bug and visit http:\u200b//phish.test');
});

test('escapeMarkdown: escapes reference-style markdown links', () => {
  const title = 'Check reference link [malicious][1] here';
  const escaped = escapeMarkdown(title);
  assert.equal(escaped, 'Check reference link \\[malicious\\]\\[1\\] here');
});

test('escapeMarkdown: neutralizes @ mentions with zero-width joiner', () => {
  const title = 'Bug report for @meshery/maintainers and @leecalcote';
  const escaped = escapeMarkdown(title);
  assert.ok(escaped.includes('@\u200dmeshery/maintainers'));
  assert.ok(escaped.includes('@\u200dleecalcote'));
  assert.ok(!escaped.includes('@meshery'));
});

test('escapeMarkdown: escapes backticks, asterisks, underscores, pipes, and tags', () => {
  const title = '<script>alert(1)</script> | `rm -rf` | *bold* _italic_ #header ~strike~';
  const escaped = escapeMarkdown(title);
  assert.ok(escaped.includes('\\<script\\>alert\\(1\\)\\</script\\>'));
  assert.ok(escaped.includes('\\|'));
  assert.ok(escaped.includes('\\`rm -rf\\`'));
  assert.ok(escaped.includes('\\*bold\\*'));
  assert.ok(escaped.includes('\\_italic\\_'));
  assert.ok(escaped.includes('\\#header'));
  assert.ok(escaped.includes('\\~strike\\~'));
});

test('escapeMarkdown: normalizes newlines, carriage returns, and control characters', () => {
  const title = 'Multi-line\r\ntitle\nwith\ttabs\x00and\x1fcontrols';
  const escaped = escapeMarkdown(title);
  assert.equal(escaped, 'Multi-line title with tabs and controls');
});

test('escapeMarkdown: strips Unicode Bidirectional override characters', () => {
  const title = 'Benign text \u202Ereversed extension.exe \u202Dnormal';
  const escaped = escapeMarkdown(title);
  assert.ok(!escaped.includes('\u202E'), 'Bidi RLO character should be stripped');
  assert.ok(!escaped.includes('\u202D'), 'Bidi LRO character should be stripped');
  assert.equal(escaped, 'Benign text reversed extension.exe normal');
});

test('escapeMarkdown: truncates long candidate titles to 200 characters with ellipsis', () => {
  const longTitle = 'A'.repeat(250);
  const escaped = escapeMarkdown(longTitle);
  assert.equal(escaped.length, 200);
  assert.ok(escaped.endsWith('...'));
  assert.equal(escaped, 'A'.repeat(197) + '...');
});

// ---------------------------------------------------------------------------
// 4. Scoring Heuristic Tests (V1 Model: Labels + Title Tokens, Zero Body Scoring)
// ---------------------------------------------------------------------------

test('calculateScore: candidate with no label overlap is strictly ineligible (score 0)', () => {
  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: new Set(['kubernetes', 'meshmodel', 'database'])
  };
  const candidate = {
    languageLabels: ['language/typescript'],
    areaComponentLabels: ['component/ui'],
    titleTokens: new Set(['kubernetes', 'meshmodel', 'database']) // identical titles
  };

  const score = calculateScore(current, candidate);
  assert.equal(score, 0, 'zero label overlap must yield score 0 despite title token overlap');
});

test('calculateScore: applies documented point values accurately (+40 lang, +40 area, +10 title token)', () => {
  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: new Set(['database', 'connection', 'pooling'])
  };
  const candidate = {
    languageLabels: ['language/go'],
    areaComponentLabels: ['component/server'],
    titleTokens: new Set(['database', 'migration']) // 1 shared token ('database')
  };

  // +40 (language) + 40 (area/component) + 10 (1 shared title token) = 90
  const score = calculateScore(current, candidate);
  assert.equal(score, 90);
});

test('calculateScore: PR template contamination produces 0 score (body scoring removed in V1)', () => {
  // Two unrelated PRs sharing the 48 tokens from .github/PULL_REQUEST_TEMPLATE.md
  // but with 1 shared label and unrelated titles:
  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: [],
    titleTokens: new Set(['server', 'database'])
  };
  const candidate = {
    languageLabels: ['language/go'],
    areaComponentLabels: [],
    titleTokens: new Set(['mesheryctl', 'context'])
  };

  // Only language match (+40), zero title token overlap.
  // Template body boilerplate is completely excluded.
  const score = calculateScore(current, candidate);
  assert.equal(score, 40, 'must score exactly +40 from language; template body cannot inflate score');
});

test('calculateScore: duplicate tokens in titles do not inflate score', () => {
  const current = {
    languageLabels: ['language/go'],
    areaComponentLabels: [],
    titleTokens: new Set(['adapter', 'provisioning'])
  };
  const candidate = {
    languageLabels: ['language/go'],
    areaComponentLabels: [],
    titleTokens: new Set(['adapter']) // only 1 shared token
  };

  const score = calculateScore(current, candidate);
  assert.equal(score, 50); // 40 (lang) + 10 (1 token)
});

// ---------------------------------------------------------------------------
// 5. Sticky Comment Matcher Tests
// ---------------------------------------------------------------------------

test('isBotStickyComment: only matches comments authored by github-actions[bot] containing marker', () => {
  const validBotComment = {
    user: { login: 'github-actions[bot]' },
    body: `${MARKER}\n\n### Suggestions`
  };
  assert.ok(isBotStickyComment(validBotComment, MARKER));

  const humanComment = {
    user: { login: 'contributor-user' },
    body: `${MARKER}\n\nSpoofed marker`
  };
  assert.ok(!isBotStickyComment(humanComment, MARKER), 'human comment must be rejected');

  const unrelatedBotComment = {
    user: { login: 'coderabbitai[bot]' },
    body: 'Review summary'
  };
  assert.ok(!isBotStickyComment(unrelatedBotComment, MARKER), 'unrelated bot must be rejected');

  const commentMissingMarker = {
    user: { login: 'github-actions[bot]' },
    body: 'Different workflow comment'
  };
  assert.ok(!isBotStickyComment(commentMissingMarker, MARKER), 'comment without marker must be rejected');
});

// ---------------------------------------------------------------------------
// 6. Workflow Orchestration and Candidate Filtering Tests
// ---------------------------------------------------------------------------

test('workflow run: creates new sticky comment for eligible candidates', async () => {
  const pr100 = {
    number: 100,
    title: '[Server] Add database connection pooling',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }, { name: 'component/server' }]
  };

  const openPRs = [
    pr100,
    {
      number: 101,
      title: '[Server] Database migration runner',
      user: { login: 'author-bob', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/101',
      labels: [{ name: 'language/go' }, { name: 'component/server' }]
    },
    {
      number: 102,
      title: '[Docs] Update server docs',
      user: { login: 'author-charlie', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/102',
      labels: [{ name: 'area/docs' }] // Ineligible: zero label overlap
    },
    {
      number: 103,
      title: '[Draft] Draft server work',
      user: { login: 'author-dave', type: 'User' },
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
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 100 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1, 'should create exactly one sticky comment');
  assert.ok(createdComments[0].body.includes(MARKER));
  assert.ok(createdComments[0].body.includes('[#101](https://github.com/meshery/meshery/pull/101)'));
  assert.ok(!createdComments[0].body.includes('#102'), 'ineligible PR #102 must not appear');
  assert.ok(!createdComments[0].body.includes('#103'), 'draft PR #103 must not appear');
});

test('workflow run: updates existing bot sticky comment and preserves human comments', async () => {
  const pr100 = {
    number: 100,
    title: '[Server] Add database connection pooling',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }]
  };

  const humanComment = {
    id: 111,
    body: `${MARKER} Human wrote this marker`,
    user: { login: 'contributor-joe', type: 'User' }
  };

  const botComment = {
    id: 222,
    body: `${MARKER}\n\nOld suggestions`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const openPRs = [
    pr100,
    {
      number: 105,
      title: '[mesheryctl] Go CLI improvement',
      user: { login: 'author-bob', type: 'User' },
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
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 100 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0, 'should not create a new comment when bot comment exists');
  assert.equal(updatedComments.length, 1, 'should update existing bot comment');
  assert.equal(updatedComments[0].comment_id, 222, 'must update bot comment #222, NOT human comment #111');
});

test('workflow run: idempotency skips update when suggestion comment body is identical', async () => {
  const pr100 = {
    number: 100,
    title: '[Server] Database pooling',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  const candidate = {
    number: 105,
    title: '[Server] Database migration',
    user: { login: 'author-bob', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/105',
    labels: [{ name: 'component/server' }]
  };

  const expectedBody = [
    MARKER,
    '### 🔎 PRs you may want to review',
    'Based on the language, area, and context of this PR, here are a few related open PRs that may be useful to review:\n\n- [#105](https://github.com/meshery/meshery/pull/105) — \\[Server\\] Database migration',
    'Reviewing related work can help contributors learn the codebase and contribute review capacity back to the community.',
    '_This is only a suggestion and does not block your PR._'
  ].join('\n\n');

  const botComment = {
    id: 555,
    body: expectedBody,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github, updatedComments, createdComments } = createMockGitHub({
    prDetails: { 100: pr100 },
    openPRs: [pr100, candidate],
    comments: [botComment]
  });
  const { core, logs } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 100 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0);
  assert.equal(updatedComments.length, 0, 'idempotent run should not issue an update request');
  assert.ok(logs.info.some(msg => msg.includes('already up-to-date')));
});

test('workflow run: does NOT create comment when zero candidates are eligible (zero noise policy)', async () => {
  const pr200 = {
    number: 200,
    title: '[Docs] Update contributing guide',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'area/docs' }]
  };

  // No other PRs share area/docs
  const openPRs = [
    pr200,
    {
      number: 201,
      title: '[Server] Core engine',
      user: { login: 'author-bob', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/201',
      labels: [{ name: 'component/server' }]
    }
  ];

  const { github, createdComments } = createMockGitHub({
    prDetails: { 200: pr200 },
    openPRs,
    comments: []
  });
  const { core, logs } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 200 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0, 'must NOT post a comment when 0 candidates exist');
  assert.ok(logs.info.some(msg => msg.includes('zero noise') || msg.includes('skipping comment creation')));
});

test('workflow run: deletes stale bot comment when candidates become zero after label removal', async () => {
  const pr200 = {
    number: 200,
    title: '[Docs] Small typo fix',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'kind/documentation' }] // No qualifying language/area/component label
  };

  const botComment = {
    id: 333,
    body: `${MARKER}\n\n- [#101](https://github.com/meshery/meshery/pull/101) — Stale candidate`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github, deletedComments, createdComments } = createMockGitHub({
    prDetails: { 200: pr200 },
    openPRs: [pr200],
    comments: [botComment]
  });
  const { core } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 200 }, action: 'unlabeled' }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0);
  assert.equal(deletedComments.length, 1, 'must delete stale bot sticky comment');
  assert.equal(deletedComments[0].comment_id, 333);
});

test('workflow run: deletes existing bot comment when PR is converted to draft', async () => {
  const draftPR = {
    number: 250,
    title: '[UI] Work in progress',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: true, // Converted to draft
    labels: [{ name: 'component/ui' }]
  };

  const botComment = {
    id: 444,
    body: `${MARKER}\n\nSuggestions`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github, deletedComments, createdComments } = createMockGitHub({
    prDetails: { 250: draftPR },
    openPRs: [draftPR],
    comments: [botComment]
  });
  const { core } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 250 }, action: 'converted_to_draft' }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0);
  assert.equal(deletedComments.length, 1, 'must delete bot comment when PR becomes draft');
  assert.equal(deletedComments[0].comment_id, 444);
});

test('workflow run: ignores bot-authored triggering PRs (no self-spam for automated PRs)', async () => {
  const botPR = {
    number: 300,
    title: 'Bump golangci-lint from 1.60 to 1.61',
    user: { login: 'dependabot[bot]', type: 'Bot' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 300: botPR },
    openPRs: [botPR],
    comments: []
  });
  const { core, logs } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 300 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 0, 'must not create review suggestions for bot-authored PRs');
  assert.ok(logs.info.some(msg => msg.includes('authored by bot')));
});

test('workflow run: excludes candidate PRs authored by the same author', async () => {
  const pr310 = {
    number: 310,
    title: '[Server] Add endpoint A',
    user: { login: 'developer-dave', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  const sameAuthorCandidate = {
    number: 311,
    title: '[Server] Add endpoint B (also by Dave)',
    user: { login: 'developer-dave', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/311',
    labels: [{ name: 'component/server' }]
  };

  const otherAuthorCandidate = {
    number: 312,
    title: '[Server] Fix database pooling',
    user: { login: 'developer-elena', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/312',
    labels: [{ name: 'component/server' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 310: pr310 },
    openPRs: [pr310, sameAuthorCandidate, otherAuthorCandidate],
    comments: []
  });
  const { core } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 310 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1);
  const body = createdComments[0].body;
  assert.ok(body.includes('#312'), 'candidate from other author must appear');
  assert.ok(!body.includes('#311'), 'candidate from the same author must be strictly excluded');
});

test('workflow run: excludes candidate PRs authored by bot accounts', async () => {
  const pr320 = {
    number: 320,
    title: '[Server] Update dependencies',
    user: { login: 'contributor-frank', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'language/go' }]
  };

  const dependabotCandidate = {
    number: 321,
    title: 'Bump github.com/stretchr/testify',
    user: { login: 'dependabot[bot]', type: 'Bot' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/321',
    labels: [{ name: 'language/go' }]
  };

  const modelBotCandidate = {
    number: 322,
    title: 'Update generated models',
    user: { login: 'meshery[bot]', type: 'User' }, // ends with [bot]
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/322',
    labels: [{ name: 'language/go' }]
  };

  const humanCandidate = {
    number: 323,
    title: '[Server] Add health check ping',
    user: { login: 'contributor-gina', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/323',
    labels: [{ name: 'language/go' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 320: pr320 },
    openPRs: [pr320, dependabotCandidate, modelBotCandidate, humanCandidate],
    comments: []
  });
  const { core } = createMockCore();

  const context = {
    repo: { owner: 'meshery', repo: 'meshery' },
    payload: { pull_request: { number: 320 } }
  };

  await run({ github, context, core });

  assert.equal(createdComments.length, 1);
  const body = createdComments[0].body;
  assert.ok(body.includes('#323'), 'human candidate must be suggested');
  assert.ok(!body.includes('#321'), 'dependabot candidate must be excluded');
  assert.ok(!body.includes('#322'), 'meshery[bot] candidate must be excluded');
});

test('workflow run: ignores non-master base branch PRs safely', async () => {
  const nonMasterPR = {
    number: 401,
    title: 'PR targeting release branch',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'release-0.8' },
    labels: [{ name: 'language/go' }]
  };

  const { github, createdComments } = createMockGitHub({
    prDetails: { 401: nonMasterPR },
    openPRs: [],
    comments: []
  });
  const { core } = createMockCore();

  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 401 } }
    },
    core
  });
  assert.equal(createdComments.length, 0, 'non-master PR should not receive comments');
});

test('workflow run: excludes candidate PRs targeting non-master branches', async () => {
  const pr700 = {
    number: 700,
    title: '[Server] Core feature',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  const candidateMaster = {
    number: 701,
    title: '[Server] Master candidate',
    user: { login: 'author-bob', type: 'User' },
    state: 'open',
    draft: false,
    base: { ref: 'master' },
    html_url: 'https://github.com/meshery/meshery/pull/701',
    labels: [{ name: 'component/server' }]
  };

  const candidateRelease = {
    number: 702,
    title: '[Server] Release branch candidate',
    user: { login: 'author-charlie', type: 'User' },
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
  const pr500 = {
    number: 500,
    title: '[Server] Server handler update',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/server' }]
  };

  const openPRs = [
    pr500,
    {
      number: 501,
      title: '[Server] Candidate 501',
      user: { login: 'user-1', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/501',
      labels: [{ name: 'component/server' }] // score: 40
    },
    {
      number: 502,
      title: '[Server] Candidate 502',
      user: { login: 'user-2', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/502',
      labels: [{ name: 'component/server' }] // score: 40
    },
    {
      number: 503,
      title: '[Server] Server handler Candidate 503', // shared title tokens: 'server', 'handler' (+20) -> score: 60
      user: { login: 'user-3', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/503',
      labels: [{ name: 'component/server' }]
    },
    {
      number: 504,
      title: '[Server] Candidate 504',
      user: { login: 'user-4', type: 'User' },
      state: 'open',
      base: { ref: 'master' },
      draft: false,
      html_url: 'https://github.com/meshery/meshery/pull/504',
      labels: [{ name: 'component/server' }] // score: 40 (newer than 502)
    },
    {
      number: 505,
      title: '[Server] Candidate 505',
      user: { login: 'user-5', type: 'User' },
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

test('workflow run: cleans up duplicate bot comments if multiple exist', async () => {
  const pr600 = {
    number: 600,
    title: '[UI] Component update',
    user: { login: 'author-alice', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    labels: [{ name: 'component/ui' }]
  };

  const botComment1 = {
    id: 1001,
    body: `${MARKER}\n\nOld suggestions 1`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };
  const botComment2 = {
    id: 1002,
    body: `${MARKER}\n\nOld suggestions 2 (duplicate)`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const candidate = {
    number: 601,
    title: '[UI] Another UI PR',
    user: { login: 'author-bob', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/601',
    labels: [{ name: 'component/ui' }]
  };

  const { github, updatedComments, deletedComments } = createMockGitHub({
    prDetails: { 600: pr600 },
    openPRs: [pr600, candidate],
    comments: [botComment1, botComment2]
  });
  const { core } = createMockCore();

  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 600 } }
    },
    core
  });

  // botComment2 must be deleted as duplicate; botComment1 must be updated
  assert.equal(deletedComments.length, 1);
  assert.equal(deletedComments[0].comment_id, 1002);
  assert.equal(updatedComments.length, 1);
  assert.equal(updatedComments[0].comment_id, 1001);
});

test('workflow run: API errors fail safely with core.warning without crashing or throwing', async () => {
  const brokenGitHub = {
    rest: {
      pulls: {
        get: async () => {
          throw new Error('GitHub API rate limit exceeded (HTTP 403)');
        }
      }
    }
  };
  const { core, logs } = createMockCore();

  await run({
    github: brokenGitHub,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 999 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Could not fetch current PR details')),
    'API failure should be captured as core.warning'
  );
});

test('workflow run: preserves historical suggestion comments after PR is closed or merged', async () => {
  const closedPR = {
    number: 800,
    title: '[UI] Closed feature PR',
    user: { login: 'author-alice', type: 'User' },
    state: 'closed',
    base: { ref: 'master' },
    draft: false,
    labels: [{ name: 'component/ui' }]
  };

  const existingComment = {
    id: 5001,
    body: `${MARKER}\n\n### 🔎 PRs you may want to review\n\n- [#801](https://github.com/meshery/meshery/pull/801) — Other PR`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github, createdComments, updatedComments, deletedComments, comments } = createMockGitHub({
    prDetails: { 800: closedPR },
    openPRs: [],
    comments: [existingComment]
  });
  const { core, logs } = createMockCore();

  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 800 } }
    },
    core
  });

  // Policy verification: Historical comments are preserved upon close/merge
  assert.equal(deletedComments.length, 0, 'must not delete existing suggestion comment on closed PR');
  assert.equal(createdComments.length, 0, 'must not create comments on closed PR');
  assert.equal(updatedComments.length, 0, 'must not update comments on closed PR');
  assert.equal(comments.length, 1, 'comment must remain in repository discussion history');
  assert.equal(comments[0].id, 5001);
  assert.ok(
    logs.info.some(i => i.includes('preserving historical suggestion comment')),
    'must log info indicating historical comment preservation'
  );
});

test('workflow run: handles error when deleting stale comments safely', async () => {
  const pr810 = {
    number: 810,
    title: '[Server] Draft PR with existing comment',
    user: { login: 'author-bob', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: true, // Converted to draft -> triggers stale comment cleanup
    labels: [{ name: 'component/server' }]
  };

  const botComment = {
    id: 5010,
    body: `${MARKER}\n\nStale suggestion`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github } = createMockGitHub({
    prDetails: { 810: pr810 },
    openPRs: [],
    comments: [botComment]
  });
  const { core, logs } = createMockCore();

  // Simulate API failure on deleteComment
  github.rest.issues.deleteComment = async () => {
    throw new Error('GitHub API 500: Server error on deleteComment');
  };

  // Must not throw or crash
  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 810 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Failed to clean up stale comments')),
    'must capture delete error as core.warning'
  );
});

test('workflow run: handles error when deleting duplicate comments safely', async () => {
  const pr820 = {
    number: 820,
    title: '[UI] PR with duplicate comments',
    user: { login: 'author-charlie', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    labels: [{ name: 'component/ui' }]
  };

  const candidate = {
    number: 821,
    title: '[UI] Related UI work',
    user: { login: 'author-dave', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/821',
    labels: [{ name: 'component/ui' }]
  };

  const botComment1 = {
    id: 5021,
    body: `${MARKER}\n\nOld suggestions 1`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };
  const botComment2 = {
    id: 5022,
    body: `${MARKER}\n\nOld suggestions 2`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github, updatedComments } = createMockGitHub({
    prDetails: { 820: pr820 },
    openPRs: [pr820, candidate],
    comments: [botComment1, botComment2]
  });
  const { core, logs } = createMockCore();

  // Simulate API failure when deleting the duplicate comment
  github.rest.issues.deleteComment = async () => {
    throw new Error('GitHub API 403: Forbidden deleting duplicate');
  };

  // Must not crash; must log warning and still update the primary comment
  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 820 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Failed to remove duplicate comment 5022')),
    'must log warning when duplicate comment deletion fails'
  );
  assert.equal(updatedComments.length, 1, 'must still proceed to update primary comment');
  assert.equal(updatedComments[0].comment_id, 5021);
});

test('workflow run: handles error when updating suggestion comment safely', async () => {
  const pr830 = {
    number: 830,
    title: '[UI] PR needing update',
    user: { login: 'author-eve', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    labels: [{ name: 'component/ui' }]
  };

  const candidate = {
    number: 831,
    title: '[UI] Related UI work',
    user: { login: 'author-frank', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/831',
    labels: [{ name: 'component/ui' }]
  };

  const botComment = {
    id: 5031,
    body: `${MARKER}\n\nOutdated body`,
    user: { login: 'github-actions[bot]', type: 'Bot' }
  };

  const { github } = createMockGitHub({
    prDetails: { 830: pr830 },
    openPRs: [pr830, candidate],
    comments: [botComment]
  });
  const { core, logs } = createMockCore();

  // Simulate API failure on updateComment
  github.rest.issues.updateComment = async () => {
    throw new Error('Secondary rate limit on updateComment');
  };

  // Must not throw or crash
  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 830 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Failed to update suggestion comment 5031')),
    'must log warning when updateComment fails'
  );
});

test('workflow run: handles error when creating suggestion comment safely', async () => {
  const pr840 = {
    number: 840,
    title: '[Server] Fresh PR with suggestions',
    user: { login: 'author-grace', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    labels: [{ name: 'component/server' }]
  };

  const candidate = {
    number: 841,
    title: '[Server] Another server PR',
    user: { login: 'author-heidi', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    html_url: 'https://github.com/meshery/meshery/pull/841',
    labels: [{ name: 'component/server' }]
  };

  const { github } = createMockGitHub({
    prDetails: { 840: pr840 },
    openPRs: [pr840, candidate],
    comments: []
  });
  const { core, logs } = createMockCore();

  // Simulate API failure on createComment
  github.rest.issues.createComment = async () => {
    throw new Error('Resource locked on createComment');
  };

  // Must not throw or crash
  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 840 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Failed to create suggestion comment on PR #840')),
    'must log warning when createComment fails'
  );
});

test('workflow run: handles error when fetching open candidate PRs safely', async () => {
  const pr850 = {
    number: 850,
    title: '[Server] Valid PR triggering candidate fetch',
    user: { login: 'author-ivan', type: 'User' },
    state: 'open',
    base: { ref: 'master' },
    draft: false,
    labels: [{ name: 'component/server' }]
  };

  const { github, createdComments, updatedComments, deletedComments } = createMockGitHub({
    prDetails: { 850: pr850 },
    openPRs: [],
    comments: []
  });
  const { core, logs } = createMockCore();

  // Simulate API failure on pulls.list
  github.rest.pulls.list = async () => {
    throw new Error('API gateway timeout 504 on pulls.list');
  };

  // Must not throw or crash, must log warning and exit early
  await run({
    github,
    context: {
      repo: { owner: 'meshery', repo: 'meshery' },
      payload: { pull_request: { number: 850 } }
    },
    core
  });

  assert.ok(
    logs.warning.some(w => w.includes('Failed to fetch open PRs: API gateway timeout 504')),
    'must log warning when candidate PR listing fails'
  );
  assert.equal(createdComments.length, 0, 'must not create comments on fetch failure');
  assert.equal(updatedComments.length, 0, 'must not update comments on fetch failure');
  assert.equal(deletedComments.length, 0, 'must not delete comments on fetch failure');
});
