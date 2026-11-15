import { describe, expect, it } from 'vitest';
import normalizeURI from '../normalizeURI';

describe('normalizeURI', () => {
  it.each([undefined, null, '', 0, 42, false, {}, [], { startsWith: () => true }])(
    'returns an empty string for invalid input %j',
    (input) => {
      expect(normalizeURI(input)).toBe('');
    },
  );

  it('prepends a slash to a URI without one', () => {
    expect(normalizeURI('users/me')).toBe('/users/me');
  });

  it('leaves a URI that already starts with a slash unchanged', () => {
    expect(normalizeURI('/users/me')).toBe('/users/me');
  });

  it('adds a slash for a single-character path', () => {
    expect(normalizeURI('a')).toBe('/a');
  });

  it('keeps deeply nested paths intact', () => {
    expect(normalizeURI('/api/v1/resource/42')).toBe('/api/v1/resource/42');
    expect(normalizeURI('api/v1/resource/42')).toBe('/api/v1/resource/42');
  });

  it('does not add another slash when the path is just "/"', () => {
    expect(normalizeURI('/')).toBe('/');
  });

  it('handles URIs containing query strings', () => {
    expect(normalizeURI('api?x=1&y=2')).toBe('/api?x=1&y=2');
    expect(normalizeURI('/api?x=1&y=2')).toBe('/api?x=1&y=2');
  });
});
