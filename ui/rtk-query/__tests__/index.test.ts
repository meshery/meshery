import { describe, expect, it } from 'vitest';
import { api, mesheryApiPath } from '../index';

// ---------------------------------------------------------------------------
// Unit tests for the mesheryApiPath helper exported from rtk-query/index.ts.
// The helper is used by every endpoint to produce a fully-qualified URL,
// so its behavior is exhaustively documented here.
// ---------------------------------------------------------------------------

describe('api re-export', () => {
  it('re-exports mesheryApi from the schemas package', () => {
    expect(api).toBeDefined();
    expect(typeof api.injectEndpoints).toBe('function');
    expect(typeof api.enhanceEndpoints).toBe('function');
    expect(api.reducerPath).toBe('mesheryRtkSchemasApi');
  });
});

describe('mesheryApiPath – default branch', () => {
  it('returns the API prefix when called with no argument', () => {
    expect(mesheryApiPath()).toBe('/api');
  });

  it('returns the API prefix when called with empty string', () => {
    expect(mesheryApiPath('')).toBe('/api');
  });
});

describe('mesheryApiPath – absolute URLs', () => {
  it('returns http URLs untouched', () => {
    expect(mesheryApiPath('http://example.com/foo')).toBe('http://example.com/foo');
  });

  it('returns https URLs untouched', () => {
    expect(mesheryApiPath('https://example.com/foo')).toBe('https://example.com/foo');
  });

  it('returns ws URLs untouched', () => {
    expect(mesheryApiPath('ws://localhost/foo')).toBe('ws://localhost/foo');
  });

  it('returns protocol-relative URLs untouched', () => {
    expect(mesheryApiPath('//cdn.example.com/foo')).toBe('//cdn.example.com/foo');
  });
});

describe('mesheryApiPath – extension routes', () => {
  it('returns the extension prefix when called with just the prefix', () => {
    expect(mesheryApiPath('/api/extensions')).toBe('/api/extensions');
  });

  it('preserves /api/extensions/* paths verbatim', () => {
    expect(mesheryApiPath('/api/extensions/foo')).toBe('/api/extensions/foo');
  });

  it('prefixes /extensions/* with /api', () => {
    expect(mesheryApiPath('/extensions/api/content/patterns')).toBe(
      '/api/extensions/api/content/patterns',
    );
  });

  it('prefixes extensions/* with /api/', () => {
    expect(mesheryApiPath('extensions/api/content/patterns')).toBe(
      '/api/extensions/api/content/patterns',
    );
  });
});

describe('mesheryApiPath – /api prefixed paths', () => {
  it('returns /api unchanged', () => {
    expect(mesheryApiPath('/api')).toBe('/api');
  });

  it('returns /api/x unchanged', () => {
    expect(mesheryApiPath('/api/pattern')).toBe('/api/pattern');
  });
});

describe('mesheryApiPath – /api/registry/relationships/evaluate', () => {
  it('passes the evaluate endpoint through unchanged as an /api/* path (meshery/schemas#916)', () => {
    expect(mesheryApiPath('/api/registry/relationships/evaluate')).toBe(
      '/api/registry/relationships/evaluate',
    );
  });
});

describe('mesheryApiPath – relative paths', () => {
  it('prefixes a bare path with /api/', () => {
    expect(mesheryApiPath('pattern')).toBe('/api/pattern');
  });

  it('prefixes a leading-slash path with /api/ once', () => {
    expect(mesheryApiPath('/pattern')).toBe('/api/pattern');
  });

  it('treats path starting with "//" as protocol-relative URL and returns it unchanged', () => {
    // '//pattern' matches startsWith('//') and is returned as-is.
    expect(mesheryApiPath('//pattern')).toBe('//pattern');
  });

  it('preserves query strings on bare paths', () => {
    expect(mesheryApiPath('pattern?page=0&pagesize=10')).toBe('/api/pattern?page=0&pagesize=10');
  });

  it('preserves path segments on bare paths', () => {
    expect(mesheryApiPath('integrations/credentials/abc-123')).toBe(
      '/api/integrations/credentials/abc-123',
    );
  });
});
