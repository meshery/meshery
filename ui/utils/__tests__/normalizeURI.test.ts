import { describe, expect, it } from 'vitest';
import normalizeURI from '../normalizeURI';

describe('normalizeURI', () => {
  it('leaves already normalized URIs starting with a slash unchanged', () => {
    expect(normalizeURI('/api/v1/test')).toBe('/api/v1/test');
    expect(normalizeURI('/users/me')).toBe('/users/me');
    expect(normalizeURI('/api/v1/resource/42')).toBe('/api/v1/resource/42');
  });

  it('prepends a slash to unnormalized URIs without one', () => {
    expect(normalizeURI('api/v1/test')).toBe('/api/v1/test');
    expect(normalizeURI('users/me')).toBe('/users/me');
    expect(normalizeURI('api/v1/resource/42')).toBe('/api/v1/resource/42');
    expect(normalizeURI('a')).toBe('/a');
  });

  it('returns an empty string when given an empty string', () => {
    expect(normalizeURI('')).toBe('');
  });

  it('handles undefined and null inputs safely without throwing TypeError', () => {
    expect(() => normalizeURI(undefined)).not.toThrow();
    expect(() => normalizeURI(null)).not.toThrow();
    expect(() => normalizeURI()).not.toThrow();
    expect(normalizeURI(undefined)).toBe('');
    expect(normalizeURI(null)).toBe('');
    expect(normalizeURI()).toBe('');
  });

  it('handles non-string inputs safely without throwing TypeError', () => {
    expect(() => normalizeURI(123)).not.toThrow();
    expect(() => normalizeURI({})).not.toThrow();
    expect(() => normalizeURI([])).not.toThrow();
    expect(normalizeURI(123)).toBe('');
    expect(normalizeURI(0)).toBe('');
    expect(normalizeURI(NaN)).toBe('');
    expect(normalizeURI(true)).toBe('');
    expect(normalizeURI(false)).toBe('');
    expect(normalizeURI({})).toBe('');
    expect(normalizeURI({ path: '/test' })).toBe('');
    expect(normalizeURI([])).toBe('');
    expect(normalizeURI([1, 2, 3])).toBe('');
    expect(normalizeURI(() => '/test')).toBe('');
    expect(normalizeURI(Symbol('uri'))).toBe('');
  });

  it('does not add another slash when the path is just "/"', () => {
    expect(normalizeURI('/')).toBe('/');
  });

  it('handles URIs containing query strings', () => {
    expect(normalizeURI('api?x=1&y=2')).toBe('/api?x=1&y=2');
    expect(normalizeURI('/api?x=1&y=2')).toBe('/api?x=1&y=2');
  });
});
