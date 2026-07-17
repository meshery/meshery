import { beforeEach, describe, expect, it, vi } from 'vitest';

// CAN reads the global CASL ability; mock it so each case can drive the
// permission branch independently of any real ability state.
const canMock = vi.fn();
vi.mock('@/utils/can', () => ({ default: (...args: unknown[]) => canMock(...args) }));

import { canEditDesign } from './design-permissions';

describe('canEditDesign', () => {
  beforeEach(() => {
    canMock.mockReset();
  });

  it('allows any holder of the CatalogManagementEditDesign permission, even a non-owner', () => {
    canMock.mockReturnValue(true);
    expect(canEditDesign({ id: 'user-1' }, { userId: 'someone-else' })).toBe(true);
  });

  it('allows the owner even without the CatalogManagementEditDesign permission', () => {
    canMock.mockReturnValue(false);
    expect(canEditDesign({ id: 'owner-1' }, { userId: 'owner-1' })).toBe(true);
  });

  it('denies a non-owner who lacks the CatalogManagementEditDesign permission', () => {
    canMock.mockReturnValue(false);
    expect(canEditDesign({ id: 'user-1' }, { userId: 'owner-2' })).toBe(false);
  });

  it('denies an unauthenticated user on an owner-less design (no undefined === undefined match)', () => {
    canMock.mockReturnValue(false);
    expect(canEditDesign(undefined, undefined)).toBe(false);
    expect(canEditDesign({}, {})).toBe(false);
    expect(canEditDesign({ id: undefined }, { userId: undefined })).toBe(false);
  });
});
