import { describe, expect, it, vi } from 'vitest';

vi.mock('uuid', () => ({
  v4: () => 'mocked-uuid-value',
}));

import { generateTestName, generateUUID } from './helper';

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
