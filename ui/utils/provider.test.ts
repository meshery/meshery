import { describe, expect, it } from 'vitest';
import { isLocalProvider, isRemoteProvider } from './provider';

describe('isLocalProvider', () => {
  it('returns true when providerType is "local"', () => {
    expect(isLocalProvider({ providerType: 'local' })).toBe(true);
  });

  it('returns false when providerType is "remote"', () => {
    expect(isLocalProvider({ providerType: 'remote' })).toBe(false);
  });

  it('returns false when capabilities are undefined or null', () => {
    expect(isLocalProvider(undefined)).toBe(false);
    expect(isLocalProvider(null)).toBe(false);
  });

  it('returns false when providerType is absent', () => {
    expect(isLocalProvider({})).toBe(false);
  });
});

describe('isRemoteProvider', () => {
  it('returns true when providerType is "remote"', () => {
    expect(isRemoteProvider({ providerType: 'remote' })).toBe(true);
  });

  it('returns false when providerType is "local"', () => {
    expect(isRemoteProvider({ providerType: 'local' })).toBe(false);
  });

  it('returns false when capabilities are undefined or null', () => {
    expect(isRemoteProvider(undefined)).toBe(false);
    expect(isRemoteProvider(null)).toBe(false);
  });
});
