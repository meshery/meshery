'use strict';

// Validation logic for the "Auto-Merge Meeting Minutes" workflow
// (.github/workflows/auto-merge-meeting-minutes.yml).
//
// That workflow auto-approves and auto-merges a PR to master, using a token
// that bypasses branch protection, if every changed file is under
// docs/meetings/. The check therefore has to be exhaustive: a single file
// outside docs/meetings/ must block the auto-merge.
//
// Extracted here so the enumeration + validation can be unit-tested without
// running the workflow. See issue #22684.

// Only files under this prefix may be auto-merged by this workflow.
const ALLOWED_PATH_PREFIX = 'docs/meetings/';

// GitHub's pulls.listFiles enumerates at most 3000 files (30 pages of 100).
// If a PR reaches that ceiling we cannot prove every file is in scope, so we
// fail closed rather than auto-merging a PR we could not fully inspect.
const MAX_ENUMERABLE_FILES = 3000;

// listAllChangedFiles returns EVERY file changed in the PR by following
// pagination. The previous single-page listFiles({ per_page: 100 }) call
// silently ignored files beyond the first 100, which let a PR with >100 files
// smuggle files outside docs/meetings/ past this gate and get auto-merged to
// master (bypassing branch protection). See issue #22684.
async function listAllChangedFiles(github, owner, repo, pullNumber) {
  return github.paginate(github.rest.pulls.listFiles, {
    owner,
    repo,
    pull_number: pullNumber,
    per_page: 100,
  });
}

// validateFiles decides whether a PR is eligible for auto-merge. It must change
// at least one file, every file must be under ALLOWED_PATH_PREFIX, and the file
// list must not have hit the enumeration ceiling (fail closed).
function validateFiles(files) {
  const validFiles = [];
  const invalidFiles = [];

  for (const file of files) {
    if (file.filename.startsWith(ALLOWED_PATH_PREFIX)) {
      validFiles.push(file.filename);
    } else {
      invalidFiles.push(file.filename);
    }
  }

  const enumerationCapped = files.length >= MAX_ENUMERABLE_FILES;
  const isValid =
    files.length > 0 && invalidFiles.length === 0 && !enumerationCapped;

  return {
    filesCount: files.length,
    validFiles,
    invalidFiles,
    enumerationCapped,
    isValid,
  };
}

module.exports = {
  ALLOWED_PATH_PREFIX,
  MAX_ENUMERABLE_FILES,
  listAllChangedFiles,
  validateFiles,
};
