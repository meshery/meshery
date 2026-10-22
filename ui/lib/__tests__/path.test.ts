import { describe, expect, it, beforeEach, afterEach } from 'vitest';
import { getPath } from '../path';

describe('getPath', () => {
  let originalLocation: Location;

  beforeEach(() => {
    originalLocation = window.location;
  });

  afterEach(() => {
    // Restore window.location
    Object.defineProperty(window, 'location', {
      configurable: true,
      writable: true,
      value: originalLocation,
    });
  });

  const setLocation = (pathname: string, search = '') => {
    Object.defineProperty(window, 'location', {
      configurable: true,
      writable: true,
      value: { ...originalLocation, pathname, search },
    });
  };

  it('strips the trailing segment after the last slash', () => {
    setLocation('/foo/bar');
    expect(getPath()).toBe('/foo');
  });

  it('strips the last segment but preserves search params', () => {
    setLocation('/foo/bar', '?baz=qux');
    expect(getPath()).toBe('/foo?baz=qux');
  });

  it('returns the same path when there is only one leading slash and no parent', () => {
    setLocation('/foo');
    // lastIndexOf('/') is 0, which is not greater than 0 — no stripping happens
    expect(getPath()).toBe('/foo');
  });

  it('returns search-only when pathname is just "/"', () => {
    setLocation('/', '?a=1');
    // lastIndexOf('/') is 0, not greater than 0 — pathname stays as "/"
    expect(getPath()).toBe('/?a=1');
  });

  it('handles a deeply nested path', () => {
    setLocation('/a/b/c/d');
    expect(getPath()).toBe('/a/b/c');
  });

  it('handles an empty pathname', () => {
    setLocation('');
    expect(getPath()).toBe('');
  });

  it('handles a path with a trailing slash by stripping the empty segment', () => {
    setLocation('/foo/');
    expect(getPath()).toBe('/foo');
  });

  it('returns just the search when path is empty and search is present', () => {
    setLocation('', '?x=1');
    expect(getPath()).toBe('?x=1');
  });

  it('preserves complex search params', () => {
    setLocation('/foo/bar', '?a=1&b=2&c=hello%20world');
    expect(getPath()).toBe('/foo?a=1&b=2&c=hello%20world');
  });
});
