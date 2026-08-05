import { describe, expect, it } from 'vitest';
import { isSafeInternalNavPath } from '../hooks/useNotification';

describe('isSafeInternalNavPath', () => {
  it('allows app-relative paths used by View connections', () => {
    expect(isSafeInternalNavPath('/management/connections')).toBe(true);
    expect(isSafeInternalNavPath('/management/connections?tab=connections')).toBe(true);
    expect(isSafeInternalNavPath('/dashboard')).toBe(true);
  });

  it('rejects open-redirect and non-path values', () => {
    expect(isSafeInternalNavPath('https://evil.example')).toBe(false);
    expect(isSafeInternalNavPath('//evil.example/phish')).toBe(false);
    expect(isSafeInternalNavPath('javascript:alert(1)')).toBe(false);
    expect(isSafeInternalNavPath('/\\evil')).toBe(false);
    expect(isSafeInternalNavPath('')).toBe(false);
    expect(isSafeInternalNavPath(null)).toBe(false);
    expect(isSafeInternalNavPath(undefined)).toBe(false);
  });
});
