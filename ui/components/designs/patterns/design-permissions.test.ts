import { describe, expect, it } from 'vitest';
import { canEditDesign } from './design-permissions';

describe('canEditDesign', () => {
  it('allows any holder of the CatalogManagementEditDesign permission, even a non-owner', () => {
    expect(canEditDesign({ id: 'user-1' }, { userId: 'someone-else' }, true)).toBe(true);
  });

  it('allows the owner even without the CatalogManagementEditDesign permission', () => {
    expect(canEditDesign({ id: 'owner-1' }, { userId: 'owner-1' }, false)).toBe(true);
  });

  it('denies a non-owner who lacks the CatalogManagementEditDesign permission', () => {
    expect(canEditDesign({ id: 'user-1' }, { userId: 'owner-2' }, false)).toBe(false);
  });

  it('denies an unauthenticated user on an owner-less design (no undefined === undefined match)', () => {
    expect(canEditDesign(undefined, undefined, false)).toBe(false);
    expect(canEditDesign({}, {}, false)).toBe(false);
    expect(canEditDesign({ id: undefined }, { userId: undefined }, false)).toBe(false);
  });
});
