'use strict';

/**
 * PR Review Suggestions Heuristic Module
 *
 * Provides deterministic, metadata-based PR review recommendations for Meshery.
 * Operates in memory on public GitHub API metadata without AI/LLM dependencies.
 */

// Common low-value stopwords to ignore during text matching
const STOP_WORDS = new Set([
  'the', 'a', 'an', 'this', 'that', 'is', 'are', 'was', 'were',
  'and', 'or', 'to', 'of', 'for', 'in', 'on', 'with', 'from',
  'fix', 'add', 'update', 'change', 'changes', 'issue', 'pull',
  'request', 'pr', 'improve', 'new',
  'by', 'at', 'it', 'as', 'be', 'all', 'not', 'we', 'you', 'i',
  'our', 'my', 'do', 'does', 'did', 'but', 'so', 'if', 'then',
  'into', 'also', 'has', 'have', 'had', 'been'
]);

// Persistent sticky comment marker to identify workflow-owned comments
const MARKER = '<!-- meshery-review-suggestions -->';

/**
 * Normalizes label arrays into language and area/component categories.
 * Matching is case-insensitive.
 *
 * @param {Array<string|object>} labels
 * @returns {{ languageLabels: string[], areaComponentLabels: string[] }}
 */
function normalizeLabels(labels) {
  const languageLabels = [];
  const areaComponentLabels = [];

  for (const label of labels || []) {
    const name = typeof label === 'string' ? label : (label && label.name);
    if (!name || typeof name !== 'string') continue;
    const lower = name.toLowerCase().trim();
    if (lower.startsWith('language/')) {
      languageLabels.push(lower);
    } else if (lower.startsWith('area/') || lower.startsWith('component/')) {
      areaComponentLabels.push(lower);
    }
  }

  return { languageLabels, areaComponentLabels };
}

/**
 * Normalizes and tokenizes text into unique technical tokens.
 * Converts punctuation (including hyphens and slashes) into whitespace delimiters.
 * Filters generic stopwords and single-character tokens.
 *
 * @param {string} text
 * @returns {Set<string>}
 */
function tokenizeText(text) {
  if (!text || typeof text !== 'string') return new Set();
  const cleaned = text.toLowerCase().replace(/[^a-z0-9]/g, ' ');
  const tokens = cleaned.split(/\s+/).filter(token => {
    return token.length > 1 && !STOP_WORDS.has(token);
  });
  return new Set(tokens);
}

/**
 * Sanitizes untrusted candidate PR titles for safe Markdown comment rendering.
 *
 * Defenses applied:
 * 1. Length bounding: Truncates titles to 200 characters max to prevent payload exhaustion.
 * 2. Control character normalization: Strips ASCII controls, carriage returns, and newlines.
 * 3. Unicode Bidi protection: Strips bidirectional override characters (\u202A-\u202E, \u2066-\u2069).
 * 4. Markdown escaping: Escapes backslashes, backticks, asterisks, underscores, brackets, parentheses, HTML tags, pipes, tildes, hashes.
 * 5. Mention neutralization: Inserts zero-width joiner (@\u200d) to prevent unsolicited notifications.
 * 6. Bare URL neutralization: Inserts zero-width space (:\u200b//) to prevent autolinking.
 *
 * @param {string} text
 * @returns {string}
 */
function escapeMarkdown(text) {
  if (!text || typeof text !== 'string') return '';

  // 1. Truncate long titles (defense-in-depth against GitHub comment limit exhaustion)
  const bounded = text.length > 200 ? text.slice(0, 197).trim() + '...' : text;

  // 2. Normalize whitespace, remove ASCII controls, and strip Unicode Bidi overrides
  const sanitized = bounded
    .replace(/[\r\n\t\x00-\x1f\x7f]/g, ' ')
    .replace(/[\u202A-\u202E\u2066-\u2069]/g, '')
    .replace(/\s+/g, ' ')
    .trim();

  // 3. Escape Markdown formatting syntax, neutralize @ mentions, and break bare URLs
  return sanitized
    .replace(/\\/g, '\\\\')
    .replace(/([`*_\[\]()<>|~#])/g, '\\$1')
    .replace(/@/g, '@\u200d')
    .replace(/:\/\//g, ':\u200b//');
}

/**
 * Calculates deterministic relevance score between current PR and candidate PR.
 *
 * Scoring Model (V1):
 * - Shared language label: +40
 * - Shared area or component label: +40
 * - Shared unique title keyword: +10 per word
 *
 * Mandatory Eligibility Gate: Candidate must share at least one language,
 * area, or component label with current PR. If not, returns 0.
 *
 * @param {{ languageLabels: string[], areaComponentLabels: string[], titleTokens: Set<string> }} current
 * @param {{ languageLabels: string[], areaComponentLabels: string[], titleTokens: Set<string> }} candidate
 * @returns {number}
 */
function calculateScore(current, candidate) {
  const hasLanguageMatch = current.languageLabels.some(l => candidate.languageLabels.includes(l));
  const hasAreaComponentMatch = current.areaComponentLabels.some(l => candidate.areaComponentLabels.includes(l));

  if (!hasLanguageMatch && !hasAreaComponentMatch) {
    return 0; // Ineligible: candidate must share at least one language/area/component label
  }

  let score = 0;
  if (hasLanguageMatch) score += 40;
  if (hasAreaComponentMatch) score += 40;

  let sharedTitleCount = 0;
  for (const token of current.titleTokens) {
    if (candidate.titleTokens.has(token)) {
      sharedTitleCount++;
    }
  }
  score += sharedTitleCount * 10;

  return score;
}

/**
 * Asserts that a comment is authored by the official GitHub Actions bot
 * and contains the workflow marker. Prevents participant comment hijacking.
 *
 * @param {object} comment
 * @param {string} marker
 * @returns {boolean}
 */
function isBotStickyComment(comment, marker) {
  if (!comment || !comment.body || !comment.user) return false;
  const isGithubActionsBot = comment.user.login === 'github-actions[bot]';
  return isGithubActionsBot && comment.body.includes(marker);
}

/**
 * Main orchestration entrypoint for PR review recommendations.
 *
 * @param {{ github: object, context: object, core: object }} param0
 */
async function run({ github, context, core }) {
  try {
    const pr = context.payload && context.payload.pull_request;
    if (!pr || !pr.number) {
      core.info('No pull request found in event payload; skipping review suggestions.');
      return;
    }

    const prNumber = pr.number;

    // Fetch up-to-date PR details from the API to guarantee current labels and state
    let currentPRData;
    try {
      const { data } = await github.rest.pulls.get({
        owner: context.repo.owner,
        repo: context.repo.repo,
        pull_number: prNumber
      });
      currentPRData = data;
    } catch (err) {
      core.warning(`Could not fetch current PR details for #${prNumber}: ${err.message}`);
      return;
    }

    // Helper: Safely retire any existing workflow-owned comment when PR becomes ineligible
    async function cleanupExistingCommentIfAny(reason) {
      try {
        const comments = await github.paginate(github.rest.issues.listComments, {
          owner: context.repo.owner,
          repo: context.repo.repo,
          issue_number: prNumber,
          per_page: 100
        });
        const botComments = (comments || []).filter(c => isBotStickyComment(c, MARKER));
        for (const comment of botComments) {
          core.info(`Retiring existing suggestion comment (id: ${comment.id}) due to: ${reason}`);
          await github.rest.issues.deleteComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            comment_id: comment.id
          });
        }
      } catch (err) {
        core.warning(`Failed to clean up stale comments on PR #${prNumber}: ${err.message}`);
      }
    }

    // 1. Verify PR state: Preserve historical suggestion comments after a PR is closed or merged
    if (currentPRData.state !== 'open') {
      core.info(`PR #${prNumber} is ${currentPRData.state}; preserving historical suggestion comment and skipping updates.`);
      return;
    }

    // 2. Verify target branch (master only)
    if (currentPRData.base && currentPRData.base.ref !== 'master') {
      core.info(`PR #${prNumber} targets ${currentPRData.base.ref}, not master; retiring any existing suggestion comment.`);
      await cleanupExistingCommentIfAny(`PR targets ${currentPRData.base.ref}`);
      return;
    }

    // 3. Verify draft status
    if (currentPRData.draft) {
      core.info(`PR #${prNumber} is a draft; retiring any existing suggestion comment.`);
      await cleanupExistingCommentIfAny('PR is converted to draft');
      return;
    }

    // 4. Verify PR author is not a bot (feature encourages human reciprocal peer review)
    const isBotAuthor = currentPRData.user && (
      currentPRData.user.type === 'Bot' ||
      (typeof currentPRData.user.login === 'string' && currentPRData.user.login.endsWith('[bot]'))
    );
    if (isBotAuthor) {
      core.info(`PR #${prNumber} is authored by bot (${currentPRData.user.login}); skipping review suggestions.`);
      return;
    }

    core.info(`Processing review suggestions for PR #${prNumber}: "${currentPRData.title}"`);

    const currentLabels = normalizeLabels(currentPRData.labels);
    const currentItem = {
      number: prNumber,
      title: currentPRData.title || '',
      languageLabels: currentLabels.languageLabels,
      areaComponentLabels: currentLabels.areaComponentLabels,
      titleTokens: tokenizeText(currentPRData.title || '')
    };

    // If current PR has no qualifying labels, it cannot match any candidates
    if (currentItem.languageLabels.length === 0 && currentItem.areaComponentLabels.length === 0) {
      core.info(`PR #${prNumber} has no language/*, area/*, or component/* labels; no candidates can be eligible.`);
      await cleanupExistingCommentIfAny('No qualifying labels on triggering PR');
      return;
    }

    // Fetch candidate open PRs (up to 100 most recently updated targeting master)
    let openPRs = [];
    try {
      const response = await github.rest.pulls.list({
        owner: context.repo.owner,
        repo: context.repo.repo,
        state: 'open',
        base: 'master',
        per_page: 100,
        sort: 'updated',
        direction: 'desc'
      });
      openPRs = response.data || [];
    } catch (err) {
      core.warning(`Failed to fetch open PRs: ${err.message}`);
      return;
    }

    const prAuthor = currentPRData.user ? currentPRData.user.login : '';
    const candidatePRs = [];
    const seenNumbers = new Set();

    for (const candidate of openPRs) {
      if (!candidate || typeof candidate.number !== 'number') continue;
      // Exclude triggering PR
      if (candidate.number === prNumber) continue;
      // Exclude PRs by same author
      if (candidate.user && candidate.user.login === prAuthor) continue;
      // Exclude bot-authored PRs (Dependabot, model automation, etc.)
      if (candidate.user && (
        candidate.user.type === 'Bot' ||
        (typeof candidate.user.login === 'string' && candidate.user.login.endsWith('[bot]'))
      )) continue;
      // Exclude drafts and non-open PRs
      if (candidate.draft) continue;
      if (candidate.state !== 'open') continue;
      // Exclude PRs targeting non-master branches
      if (candidate.base && candidate.base.ref !== 'master') continue;
      // Exclude duplicates
      if (seenNumbers.has(candidate.number)) continue;
      seenNumbers.add(candidate.number);

      candidatePRs.push(candidate);
    }

    // Score eligible candidates
    const scoredCandidates = [];

    for (const candidate of candidatePRs) {
      const candidateLabels = normalizeLabels(candidate.labels);
      const candidateItem = {
        number: candidate.number,
        title: candidate.title || '',
        html_url: candidate.html_url,
        languageLabels: candidateLabels.languageLabels,
        areaComponentLabels: candidateLabels.areaComponentLabels,
        titleTokens: tokenizeText(candidate.title || '')
      };

      const score = calculateScore(currentItem, candidateItem);
      if (score > 0) {
        scoredCandidates.push({
          number: candidate.number,
          title: candidate.title || '',
          html_url: candidate.html_url,
          score
        });
      }
    }

    // Rank candidates: descending score, tie-break by PR number descending
    scoredCandidates.sort((a, b) => {
      if (b.score !== a.score) {
        return b.score - a.score;
      }
      return b.number - a.number;
    });

    const topCandidates = scoredCandidates.slice(0, 3);
    core.info(`Scored candidates count: ${scoredCandidates.length}; selected top ${topCandidates.length}`);

    // Retrieve existing PR comments
    let comments = [];
    try {
      comments = await github.paginate(github.rest.issues.listComments, {
        owner: context.repo.owner,
        repo: context.repo.repo,
        issue_number: prNumber,
        per_page: 100
      });
    } catch (err) {
      core.warning(`Failed to list comments on PR #${prNumber}: ${err.message}`);
      return;
    }

    const botComments = (comments || []).filter(c => isBotStickyComment(c, MARKER));

    // Handle zero eligible candidates: Suppress comment creation, retire stale comment if present
    if (topCandidates.length === 0) {
      core.info(`No eligible review suggestions found for PR #${prNumber}.`);
      if (botComments.length > 0) {
        for (const comment of botComments) {
          core.info(`Retiring stale suggestion comment (id: ${comment.id}).`);
          try {
            await github.rest.issues.deleteComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              comment_id: comment.id
            });
          } catch (err) {
            core.warning(`Failed to delete stale comment ${comment.id}: ${err.message}`);
          }
        }
      } else {
        core.info('No existing bot comment present; skipping comment creation (zero noise).');
      }
      return;
    }

    // Build comment content for eligible candidates
    const prList = topCandidates
      .map(c => `- [#${c.number}](${c.html_url}) — ${escapeMarkdown(c.title)}`)
      .join('\n');

    const commentSections = [
      MARKER,
      '### 🔎 PRs you may want to review',
      'Based on the language, area, and context of this PR, here are a few related open PRs that may be useful to review:\n\n' + prList,
      'Reviewing related work can help contributors learn the codebase and contribute review capacity back to the community.',
      '_This is only a suggestion and does not block your PR._'
    ];
    const commentBody = commentSections.join('\n\n');

    // Deterministic duplicate cleanup: if more than one bot comment exists, keep the first and delete the rest
    if (botComments.length > 1) {
      for (let i = 1; i < botComments.length; i++) {
        core.info(`Removing duplicate bot sticky comment (id: ${botComments[i].id}).`);
        try {
          await github.rest.issues.deleteComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            comment_id: botComments[i].id
          });
        } catch (err) {
          core.warning(`Failed to remove duplicate comment ${botComments[i].id}: ${err.message}`);
        }
      }
    }

    const primaryComment = botComments[0];

    if (primaryComment) {
      if (primaryComment.body.trim() === commentBody.trim()) {
        core.info('Existing suggestion comment is already up-to-date; skipping API write.');
      } else {
        core.info(`Updating existing suggestion comment (id: ${primaryComment.id}).`);
        try {
          await github.rest.issues.updateComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            comment_id: primaryComment.id,
            body: commentBody
          });
        } catch (err) {
          core.warning(`Failed to update suggestion comment ${primaryComment.id}: ${err.message}`);
        }
      }
    } else {
      core.info('Creating new suggestion comment.');
      try {
        await github.rest.issues.createComment({
          owner: context.repo.owner,
          repo: context.repo.repo,
          issue_number: prNumber,
          body: commentBody
        });
      } catch (err) {
        core.warning(`Failed to create suggestion comment on PR #${prNumber}: ${err.message}`);
      }
    }
  } catch (error) {
    core.warning(`PR review suggestions failed safely: ${error.message}`);
  }
}

module.exports = {
  STOP_WORDS,
  MARKER,
  normalizeLabels,
  tokenizeText,
  escapeMarkdown,
  calculateScore,
  isBotStickyComment,
  run
};
