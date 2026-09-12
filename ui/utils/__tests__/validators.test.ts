import { describe, expect, it } from 'vitest';
import {
  isArrayEmpty,
  isFieldEmpty,
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
