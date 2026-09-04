import { describe, expect, it } from 'vitest';
import {
  isArrayEmpty,
  isFieldEmpty,
  isValidDuration,
  isValidJSON,
  matchesSearch,
  normalizeSearchTerm,
} from '../validators';

describe('isFieldEmpty', () => {
  it('treats null/undefined as empty', () => {
    expect(isFieldEmpty(null)).toBe(true);
    expect(isFieldEmpty(undefined)).toBe(true);
  });

  it('treats whitespace-only strings as empty', () => {
    expect(isFieldEmpty('')).toBe(true);
    expect(isFieldEmpty('   ')).toBe(true);
    expect(isFieldEmpty('\t\n')).toBe(true);
  });

  it('treats populated strings as non-empty', () => {
    expect(isFieldEmpty('a')).toBe(false);
    expect(isFieldEmpty('  x  ')).toBe(false);
  });

  it('handles empty arrays', () => {
    expect(isFieldEmpty([])).toBe(true);
    expect(isFieldEmpty([0])).toBe(false);
  });
});

describe('isArrayEmpty', () => {
  it('returns true for null/undefined/empty', () => {
    expect(isArrayEmpty(null)).toBe(true);
    expect(isArrayEmpty(undefined)).toBe(true);
    expect(isArrayEmpty([])).toBe(true);
  });

  it('returns false for non-empty arrays', () => {
    expect(isArrayEmpty([1])).toBe(false);
    expect(isArrayEmpty(['', null])).toBe(false);
  });
});

describe('isValidJSON', () => {
  it('returns true for valid JSON strings', () => {
    expect(isValidJSON('{}')).toBe(true);
    expect(isValidJSON('{"a":1}')).toBe(true);
    expect(isValidJSON('[1,2,3]')).toBe(true);
    expect(isValidJSON('"hello"')).toBe(true);
    expect(isValidJSON('true')).toBe(true);
  });

  it('returns false for invalid input', () => {
    expect(isValidJSON('not json')).toBe(false);
    expect(isValidJSON('{a:1}')).toBe(false);
    expect(isValidJSON('')).toBe(false);
    expect(isValidJSON('   ')).toBe(false);
    expect(isValidJSON(null)).toBe(false);
    expect(isValidJSON(undefined)).toBe(false);
    expect(isValidJSON(123)).toBe(false);
  });
});

describe('normalizeSearchTerm / matchesSearch', () => {
  it('lowercases and trims', () => {
    expect(normalizeSearchTerm('  HelloWorld  ')).toBe('helloworld');
  });

  it('returns true for empty search term', () => {
    expect(matchesSearch('anything', '')).toBe(true);
    expect(matchesSearch('anything', '   ')).toBe(true);
  });

  it('matches case-insensitively', () => {
    expect(matchesSearch('Production cluster', 'CLUSTER')).toBe(true);
    expect(matchesSearch('Production cluster', '  prod  ')).toBe(true);
  });

  it('returns false on no match', () => {
    expect(matchesSearch('alpha', 'beta')).toBe(false);
  });
});

describe('isValidDuration', () => {
  it('returns true for valid durations with standard units (s, m, h)', () => {
    expect(isValidDuration('30s')).toBe(true);
    expect(isValidDuration('15s')).toBe(true);
    expect(isValidDuration('1m')).toBe(true);
    expect(isValidDuration('5m')).toBe(true);
    expect(isValidDuration('1h')).toBe(true);
    expect(isValidDuration('24h')).toBe(true);
  });

  it('handles case-insensitivity for unit suffix', () => {
    expect(isValidDuration('30S')).toBe(true);
    expect(isValidDuration('5M')).toBe(true);
    expect(isValidDuration('2H')).toBe(true);
  });

  it('trims leading and trailing whitespace', () => {
    expect(isValidDuration('  30s  ')).toBe(true);
    expect(isValidDuration('\t10m\n')).toBe(true);
  });

  it('rejects zero or negative values', () => {
    expect(isValidDuration('0s')).toBe(false);
    expect(isValidDuration('00s')).toBe(false);
    expect(isValidDuration('0m')).toBe(false);
    expect(isValidDuration('0h')).toBe(false);
    expect(isValidDuration('-5s')).toBe(false);
  });

  it('rejects non-numeric characters before unit', () => {
    expect(isValidDuration('abcs')).toBe(false);
    expect(isValidDuration('tests')).toBe(false);
    expect(isValidDuration('testm')).toBe(false);
    expect(isValidDuration('12abcs')).toBe(false);
    expect(isValidDuration('12 34s')).toBe(false);
    expect(isValidDuration('1.5s')).toBe(false);
  });

  it('rejects missing numeric part or unit suffix', () => {
    expect(isValidDuration('s')).toBe(false);
    expect(isValidDuration('m')).toBe(false);
    expect(isValidDuration('h')).toBe(false);
    expect(isValidDuration('30')).toBe(false);
    expect(isValidDuration('100')).toBe(false);
  });

  it('rejects unsupported units by default', () => {
    expect(isValidDuration('10x')).toBe(false);
    expect(isValidDuration('1d')).toBe(false);
    expect(isValidDuration('500ms')).toBe(false);
  });

  it('supports custom allowed units when provided', () => {
    expect(isValidDuration('1d', ['d', 'h', 'm', 's'])).toBe(true);
    expect(isValidDuration('500ms', ['ms', 's'])).toBe(false); // suffix length > 1
  });

  it('rejects empty strings and non-string types', () => {
    expect(isValidDuration('')).toBe(false);
    expect(isValidDuration('   ')).toBe(false);
    expect(isValidDuration(null)).toBe(false);
    expect(isValidDuration(undefined)).toBe(false);
    expect(isValidDuration(30)).toBe(false);
    expect(isValidDuration({})).toBe(false);
    expect(isValidDuration([])).toBe(false);
  });
});
