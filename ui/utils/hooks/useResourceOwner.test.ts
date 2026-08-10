import { beforeEach, describe, expect, it, vi } from 'vitest';

const mockGetUserById = vi.fn();
const mockGetProviderCapabilities = vi.fn();

vi.mock('@/rtk-query/user', () => ({
  useGetUserByIdQuery: (id?: string, options?: unknown) => mockGetUserById(id, options),
  useGetProviderCapabilitiesQuery: () => mockGetProviderCapabilities(),
}));

import { useResourceOwner } from './useResourceOwner';

const fetchedOwner = { id: 'u1', firstName: 'Bob', lastName: 'Jones', avatarUrl: '/b.png' };
const embeddedOwner = { id: 'u9', firstName: 'Ada', lastName: 'Lovelace', avatarUrl: '/a.png' };

describe('useResourceOwner', () => {
  beforeEach(() => {
    mockGetUserById.mockReset();
    mockGetUserById.mockReturnValue({ data: fetchedOwner });
    mockGetProviderCapabilities.mockReset();
    mockGetProviderCapabilities.mockReturnValue({ data: { providerType: 'remote' } });
  });

  it('prefers the embedded owner profile and skips the by-id lookup', () => {
    const { owner } = useResourceOwner('u9', embeddedOwner);

    expect(owner).toBe(embeddedOwner);
    // The guard must actually suppress the request, not merely lose the `??`
    // race: the built-in single-user provider resolves no other id, so an
    // unskipped lookup is a guaranteed-miss request per rendered resource.
    expect(mockGetUserById).toHaveBeenCalledWith('u9', expect.objectContaining({ skip: true }));
  });

  it('falls back to the by-id lookup when no profile is embedded', () => {
    const { owner } = useResourceOwner('u1');

    expect(owner).toBe(fetchedOwner);
    expect(mockGetUserById).toHaveBeenCalledWith('u1', expect.objectContaining({ skip: false }));
  });

  it('does not skip the lookup for a null or empty embedded profile', () => {
    useResourceOwner('u1', null);

    expect(mockGetUserById).toHaveBeenCalledWith('u1', expect.objectContaining({ skip: false }));
  });

  it('reports a cloud profile on a remote provider', () => {
    expect(useResourceOwner('u1').hasCloudProfile).toBe(true);
  });

  it('reports no cloud profile on the built-in local provider', () => {
    mockGetProviderCapabilities.mockReturnValue({ data: { providerType: 'local' } });

    expect(useResourceOwner('u1').hasCloudProfile).toBe(false);
  });

  it('reports no cloud profile while the provider type is not yet known', () => {
    // On cold load the provider-capabilities query is still in flight, so the
    // type is unresolved. Reporting a cloud profile here would flash a dead,
    // 404-bound link before the provider resolves, so an unresolved provider
    // must fail closed to no cloud profile.
    mockGetProviderCapabilities.mockReturnValue({ data: undefined });

    expect(useResourceOwner('u1').hasCloudProfile).toBe(false);
  });
});
