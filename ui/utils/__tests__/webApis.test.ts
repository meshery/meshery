import { describe, expect, it, beforeEach, afterEach } from 'vitest';
import {
  getProtocol,
  getWebAdress,
  getQueryParam,
  getRawUrlFromCssUrlString,
  getNavigator,
} from '../webApis';

// jsdom exposes `window.location` but it is a read-only `Location` object.
// Re-defining the descriptor lets each test simulate different URLs.

const setLocation = (url: string) => {
  const u = new URL(url);
  Object.defineProperty(window, 'location', {
    configurable: true,
    value: {
      ...window.location,
      href: u.href,
      origin: u.origin,
      protocol: u.protocol,
      host: u.host,
      hostname: u.hostname,
      port: u.port,
      pathname: u.pathname,
      search: u.search,
      hash: u.hash,
    },
  });
};

describe('webApis', () => {
  let originalLocation: Location;

  beforeEach(() => {
    originalLocation = window.location;
  });
  afterEach(() => {
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation,
    });
  });

  describe('getProtocol / getWebAdress', () => {
    it('returns the current protocol', () => {
      setLocation('http://meshery.test/path');
      expect(getProtocol()).toBe('http:');
    });

    it('builds the full web address with protocol + host', () => {
      setLocation('https://meshery.test:9081/x');
      expect(getWebAdress()).toBe('https://meshery.test:9081');
    });
  });

  describe('getQueryParam', () => {
    it('extracts a single query value', () => {
      setLocation('http://localhost:9081/x?foo=bar');
      expect(getQueryParam('foo')).toBe('bar');
    });

    it('returns the first occurrence when a key is repeated', () => {
      setLocation('http://localhost:9081/x?a=1&a=2');
      expect(getQueryParam('a')).toBe('1');
    });

    it('returns empty string when the key is absent', () => {
      setLocation('http://localhost:9081/x?foo=bar');
      expect(getQueryParam('missing')).toBe('');
    });

    it('returns empty string when there is no query string', () => {
      setLocation('http://localhost:9081/x');
      expect(getQueryParam('foo')).toBe('');
    });

    it('does not match a prefix; only exact key equality', () => {
      setLocation('http://localhost:9081/x?abc=1&ab=2');
      expect(getQueryParam('ab')).toBe('2');
      expect(getQueryParam('abc')).toBe('1');
    });
  });

  describe('getRawUrlFromCssUrlString', () => {
    it('returns undefined for nullish input', () => {
      expect(getRawUrlFromCssUrlString(undefined)).toBeUndefined();
      expect(getRawUrlFromCssUrlString(null)).toBeUndefined();
      expect(getRawUrlFromCssUrlString('')).toBeUndefined();
    });

    it('strips a url(...) wrapper', () => {
      // url(http://localhost:9081/path/to/svg) -> /path/to/svg (after http stripping)
      expect(getRawUrlFromCssUrlString('url(http://localhost:9081/path/to/svg)')).toBe(
        'path/to/svg',
      );
    });

    it('returns the host-relative path for absolute http URLs', () => {
      expect(getRawUrlFromCssUrlString('http://localhost:9081/path/to/svg')).toBe('path/to/svg');
      expect(getRawUrlFromCssUrlString('https://meshery.io/icon.svg')).toBe('icon.svg');
    });

    it('returns relative paths verbatim when they are neither url(...) nor http', () => {
      expect(getRawUrlFromCssUrlString('/local/path')).toBe('/local/path');
      expect(getRawUrlFromCssUrlString('static/icon.png')).toBe('static/icon.png');
    });
  });

  describe('SSR / null safety guards', () => {
    it('handles undefined window.location gracefully', () => {
      Object.defineProperty(window, 'location', {
        configurable: true,
        value: undefined,
      });
      expect(getProtocol()).toBe('');
      expect(getWebAdress()).toBe('');
      expect(getQueryParam('foo')).toBe('');
    });

    it('handles invalid or non-string queryKey in getQueryParam', () => {
      setLocation('http://localhost:9081/x?foo=bar');
      expect(getQueryParam(undefined)).toBe('');
      expect(getQueryParam(null as any)).toBe('');
      expect(getQueryParam(123 as any)).toBe('');
    });

    it('handles non-string inputs safely in getRawUrlFromCssUrlString', () => {
      expect(getRawUrlFromCssUrlString(123 as any)).toBeUndefined();
      expect(getRawUrlFromCssUrlString({} as any)).toBeUndefined();
      expect(getRawUrlFromCssUrlString([] as any)).toBeUndefined();
    });

    it('returns navigator when present in browser environment', () => {
      expect(getNavigator()).toBe(window.navigator);
    });
  });
});
