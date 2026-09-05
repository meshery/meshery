import { describe, expect, it, vi } from 'vitest';

vi.mock('uuid', () => ({
  v4: () => 'mocked-uuid-value',
}));

import { generateTestName, generateUUID, isValidDuration } from './helper';

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

describe('isValidDuration (re-exported from helper)', () => {
  it('accepts valid performance durations', () => {
    expect(isValidDuration('30s')).toBe(true);
    expect(isValidDuration('15s')).toBe(true);
    expect(isValidDuration('1m')).toBe(true);
    expect(isValidDuration('2h')).toBe(true);
  });

  it('rejects invalid performance durations', () => {
    expect(isValidDuration('0s')).toBe(false);
    expect(isValidDuration('abcs')).toBe(false);
    expect(isValidDuration('12abcs')).toBe(false);
    expect(isValidDuration('')).toBe(false);
    expect(isValidDuration(null)).toBe(false);
  });
});
