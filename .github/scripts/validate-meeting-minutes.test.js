'use strict';

// Unit tests for the Auto-Merge Meeting Minutes validation logic.
// Run with:  node --test .github/scripts/
//
// Regression coverage for issue #22684: the validation must enumerate EVERY
// changed file (across pagination), so a PR with more than 100 files cannot
// hide a file outside docs/meetings/ beyond the first page.

const test = require('node:test');
const assert = require('node:assert/strict');
const path = require('node:path');

const {
  validateFiles,
  listAllChangedFiles,
  MAX_ENUMERABLE_FILES,
} = require(path.join(__dirname, 'validate-meeting-minutes.js'));

const meetingFile = (i) => ({ filename: `docs/meetings/2026/week-${i}.md` });

test('rejects a >100-file PR whose out-of-scope file is beyond the first page', () => {
  // 150 files: all under docs/meetings/ except one at index 120 (page 2).
  const files = [];
  for (let i = 0; i < 150; i++) {
    files.push(i === 120 ? { filename: 'server/handlers/evil.go' } : meetingFile(i));
  }

  const result = validateFiles(files);

  assert.equal(result.filesCount, 150);
  assert.equal(result.isValid, false);
  assert.deepEqual(result.invalidFiles, ['server/handlers/evil.go']);
});

test('accepts a large PR that is entirely under docs/meetings/', () => {
  const files = Array.from({ length: 250 }, (_, i) => meetingFile(i));

  const result = validateFiles(files);

  assert.equal(result.isValid, true);
  assert.equal(result.invalidFiles.length, 0);
});

test('rejects an empty PR (nothing to merge)', () => {
  assert.equal(validateFiles([]).isValid, false);
});

test('fails closed when the file count hits the enumeration ceiling', () => {
  const files = Array.from({ length: MAX_ENUMERABLE_FILES }, (_, i) => meetingFile(i));

  const result = validateFiles(files);

  assert.equal(result.enumerationCapped, true);
  assert.equal(result.isValid, false);
});

test('listAllChangedFiles follows pagination (not a single 100-file page)', async () => {
  // Simulate a 127-file PR served across pages; github.paginate flattens them.
  const allFiles = Array.from({ length: 127 }, (_, i) => meetingFile(i));
  let calledWith;

  const github = {
    rest: { pulls: { listFiles: Symbol('listFiles') } },
    paginate: async (fn, params) => {
      calledWith = { fn, params };
      return allFiles;
    },
  };

  const files = await listAllChangedFiles(github, 'meshery', 'meshery', 20161);

  assert.equal(files.length, 127); // all pages, not just the first 100
  assert.equal(calledWith.fn, github.rest.pulls.listFiles);
  assert.equal(calledWith.params.per_page, 100);
  assert.equal(calledWith.params.pull_number, 20161);
});
