import { renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

let catalogVisibility = true;
const dispatch = vi.fn((action) => {
  catalogVisibility = action.catalogVisibility;
});
let preferences: any;
let mutation: any;

vi.mock('react-redux', () => ({ useDispatch: () => dispatch }));
vi.mock('@/rtk-query/user', () => ({
  useGetUserPrefQuery: () => preferences,
  useUpdateUserPrefMutation: () => [vi.fn(), mutation],
}));
vi.mock('@/store/slices/mesheryUi', () => ({
  toggleCatalogContent: (value) => value,
}));

import { useSyncCatalogContentPreference } from './useSyncCatalogContentPreference';

describe('application-shell catalog preference sync', () => {
  beforeEach(() => {
    dispatch.mockClear();
    catalogVisibility = true;
    preferences = { data: { usersExtensionPreferences: { catalogContent: true } } };
    mutation = { isLoading: false };
  });

  it('applies a saved preference that arrives after shell startup', () => {
    preferences = { data: undefined };
    const { rerender } = renderHook(() => useSyncCatalogContentPreference());
    expect(dispatch).not.toHaveBeenCalled();

    preferences = { data: { usersExtensionPreferences: { catalogContent: false } } };
    rerender();
    expect(catalogVisibility).toBe(false);
  });

  it('does not replay the cached catalog value when unrelated preferences change', () => {
    const { rerender } = renderHook(() => useSyncCatalogContentPreference());
    catalogVisibility = false; // The preferences page applied an optimistic toggle.
    dispatch.mockClear();

    preferences = {
      data: { usersExtensionPreferences: { catalogContent: true, otherPreference: 'changed' } },
    };
    rerender();
    expect(dispatch).not.toHaveBeenCalled();
    expect(catalogVisibility).toBe(false);
  });

  it('keeps a newer choice when an older refetch arrives during its save', () => {
    const { rerender } = renderHook(() => useSyncCatalogContentPreference());
    dispatch.mockClear();
    mutation = { isLoading: true };
    preferences = { data: { usersExtensionPreferences: { catalogContent: false } } };
    rerender();
    expect(catalogVisibility).toBe(true);
    expect(dispatch).not.toHaveBeenCalled();

    mutation = { isLoading: false };
    rerender();
    expect(catalogVisibility).toBe(true);
    expect(dispatch).not.toHaveBeenCalled();

    preferences = { data: { usersExtensionPreferences: { catalogContent: true } } };
    rerender();
    expect(dispatch).toHaveBeenCalledWith({ catalogVisibility: true });
  });
});
