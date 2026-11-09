import { describe, expect, it, vi } from 'vitest';

vi.mock('uuid', () => ({
  v4: () => 'mocked-uuid-value',
}));

import {
  buildPerformanceRunUrl,
  generateTestName,
  generateUUID,
  getStoredAdditionalOptions,
} from './helper';

describe('generateTestName', () => {
  it('returns the original name when one is supplied', () => {
    expect(generateTestName('my-test', 'istio')).toBe('my-test');
  });

  it('falls back to a "No mesh"-prefixed timestamped name for an empty mesh', () => {
    const name = generateTestName('', '');
    expect(name.startsWith('No mesh_')).toBe(true);
  });

  it('treats "None" the same as an empty mesh', () => {
    const name = generateTestName('', 'None');
    expect(name.startsWith('No mesh_')).toBe(true);
  });

  it('uses the supplied mesh name when generating a fallback', () => {
    const name = generateTestName('', 'linkerd');
    expect(name.startsWith('linkerd_')).toBe(true);
  });

  it('treats whitespace-only names as empty', () => {
    const name = generateTestName('   ', 'consul');
    expect(name.startsWith('consul_')).toBe(true);
  });
});

describe('generateUUID', () => {
  it('proxies to the uuid v4 implementation', () => {
    expect(generateUUID()).toBe('mocked-uuid-value');
  });
});

describe('getStoredAdditionalOptions', () => {
  it('reads the first entry of options saved from the UI', () => {
    expect(getStoredAdditionalOptions(['{"keepalive":true}'])).toBe('{"keepalive":true}');
  });

  it('keeps options stored as a plain string by mesheryctl', () => {
    expect(getStoredAdditionalOptions('{"keepalive":true}')).toBe('{"keepalive":true}');
  });

  it('returns an empty string when nothing is stored', () => {
    expect(getStoredAdditionalOptions(undefined)).toBe('');
    expect(getStoredAdditionalOptions([])).toBe('');
  });
});

describe('buildPerformanceRunUrl', () => {
  it('constructs valid URL when no Kubernetes contexts are selected', () => {
    const url = buildPerformanceRunUrl('test-profile-id', [], 'name=load-test&qps=50');
    expect(url).toBe(
      '/api/user/performance/profiles/test-profile-id/run?cert=true&name=load-test&qps=50',
    );
  });

  it('constructs valid URL when Kubernetes contexts are selected', () => {
    const url = buildPerformanceRunUrl(
      'test-profile-id',
      ['ctx-1', 'ctx-2'],
      'name=load-test&qps=50',
    );
    expect(url).toBe(
      '/api/user/performance/profiles/test-profile-id/run?contexts=ctx-1&contexts=ctx-2&cert=true&name=load-test&qps=50',
    );
  });

  it('handles null or undefined selectedK8sContexts', () => {
    const url = buildPerformanceRunUrl('test-profile-id', null, 'name=load-test');
    expect(url).toBe('/api/user/performance/profiles/test-profile-id/run?cert=true&name=load-test');
  });

  it('handles empty params string', () => {
    const url = buildPerformanceRunUrl('test-profile-id', []);
    expect(url).toBe('/api/user/performance/profiles/test-profile-id/run?cert=true');
  });
});
