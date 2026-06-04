import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const schemaMocks = vi.hoisted(() => {
  const result = { isLoading: false };

  return {
    result,
    upsertProfileTrigger: vi.fn(),
    deleteProfileTrigger: vi.fn(),
    useGetPerformanceProfilesQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceProfiles',
      queryArg,
      options,
    })),
    useUpsertPerformanceProfileMutation: vi.fn(() => [vi.fn(), result] as const),
    useGetPerformanceResultsQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceResults',
      queryArg,
      options,
    })),
    useGetPerformanceProfileQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceProfile',
      queryArg,
      options,
    })),
    useDeletePerformanceProfileMutation: vi.fn(() => [vi.fn(), result] as const),
    useGetPerformanceProfileResultsQuery: vi.fn((queryArg, options) => ({
      hook: 'getPerformanceProfileResults',
      queryArg,
      options,
    })),
  };
});

vi.mock('@meshery/schemas/mesheryApi', () => ({
  useDeletePerformanceProfileMutation: schemaMocks.useDeletePerformanceProfileMutation,
  useGetPerformanceProfileQuery: schemaMocks.useGetPerformanceProfileQuery,
  useGetPerformanceProfileResultsQuery: schemaMocks.useGetPerformanceProfileResultsQuery,
  useGetPerformanceProfilesQuery: schemaMocks.useGetPerformanceProfilesQuery,
  useGetPerformanceResultsQuery: schemaMocks.useGetPerformanceResultsQuery,
  useUpsertPerformanceProfileMutation: schemaMocks.useUpsertPerformanceProfileMutation,
}));

const loadModule = async () => import('../performance-profile');

describe('performance-profile endpoints', () => {
  beforeEach(() => {
    schemaMocks.useGetPerformanceProfilesQuery.mockClear();
    schemaMocks.useUpsertPerformanceProfileMutation.mockClear();
    schemaMocks.useGetPerformanceResultsQuery.mockClear();
    schemaMocks.useGetPerformanceProfileQuery.mockClear();
    schemaMocks.useDeletePerformanceProfileMutation.mockClear();
    schemaMocks.useGetPerformanceProfileResultsQuery.mockClear();
    schemaMocks.upsertProfileTrigger = vi.fn();
    schemaMocks.deleteProfileTrigger = vi.fn();
    schemaMocks.useUpsertPerformanceProfileMutation.mockReturnValue([
      schemaMocks.upsertProfileTrigger,
      schemaMocks.result,
    ] as const);
    schemaMocks.useDeletePerformanceProfileMutation.mockReturnValue([
      schemaMocks.deleteProfileTrigger,
      schemaMocks.result,
    ] as const);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('exports all expected hooks', async () => {
    const mod = await loadModule();
    expect(mod.useGetPerformanceProfilesQuery).toBeTypeOf('function');
    expect(mod.useSavePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsQuery).toBeTypeOf('function');
    expect(mod.useGetPerformanceProfileByIdQuery).toBeTypeOf('function');
    expect(mod.useDeletePerformanceProfileMutation).toBeTypeOf('function');
    expect(mod.useGetProfileResultsByIdQuery).toBeTypeOf('function');
  });

  it('getPerformanceProfiles delegates to the schema-generated hook with normalized params', async () => {
    const { useGetPerformanceProfilesQuery } = await loadModule();
    const result = useGetPerformanceProfilesQuery(
      {
        page: 2,
        pagesize: 25,
        search: 'foo',
        order: 'name desc',
      },
      { skip: true },
    );

    expect(result).toEqual({
      hook: 'getPerformanceProfiles',
      queryArg: {
        page: '2',
        pagesize: '25',
        search: 'foo',
        order: 'name desc',
      },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceProfilesQuery).toHaveBeenCalledTimes(1);
  });

  it('savePerformanceProfile delegates to schema-generated upsertPerformanceProfile', async () => {
    const { useSavePerformanceProfileMutation } = await loadModule();
    const body = { name: 'profile-a', endpoints: ['/x'] };
    const [trigger, result] = useSavePerformanceProfileMutation();

    trigger({ body });

    expect(result).toBe(schemaMocks.result);
    expect(schemaMocks.useUpsertPerformanceProfileMutation).toHaveBeenCalledTimes(1);
    expect(schemaMocks.upsertProfileTrigger).toHaveBeenCalledWith({ body });
  });

  it('getProfileResults delegates to schema-generated getPerformanceResults', async () => {
    const { useGetProfileResultsQuery } = await loadModule();
    const result = useGetProfileResultsQuery(
      {
        page: 0,
        pagesize: 10,
        search: '',
        order: '',
        from: '2024-01-01',
        to: '2024-02-01',
      },
      { skip: true },
    );

    expect(result).toEqual({
      hook: 'getPerformanceResults',
      queryArg: {
        page: '0',
        pagesize: '10',
        search: '',
        order: '',
        from: '2024-01-01',
        to: '2024-02-01',
      },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceResultsQuery).toHaveBeenCalledTimes(1);
  });

  it('getPerformanceProfileById delegates to schema-generated getPerformanceProfile', async () => {
    const { useGetPerformanceProfileByIdQuery } = await loadModule();
    const result = useGetPerformanceProfileByIdQuery({ id: 'abc' }, { skip: true });

    expect(result).toEqual({
      hook: 'getPerformanceProfile',
      queryArg: { performanceProfileId: 'abc' },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceProfileQuery).toHaveBeenCalledTimes(1);
  });

  it('deletePerformanceProfile delegates to schema-generated deletePerformanceProfile', async () => {
    const { useDeletePerformanceProfileMutation } = await loadModule();
    const [trigger, result] = useDeletePerformanceProfileMutation();

    trigger({ id: 'xyz' });

    expect(result).toBe(schemaMocks.result);
    expect(schemaMocks.useDeletePerformanceProfileMutation).toHaveBeenCalledTimes(1);
    expect(schemaMocks.deleteProfileTrigger).toHaveBeenCalledWith({ performanceProfileId: 'xyz' });
  });

  it('getProfileResultsById delegates to schema-generated getPerformanceProfileResults', async () => {
    const { useGetProfileResultsByIdQuery } = await loadModule();
    const result = useGetProfileResultsByIdQuery(
      {
        id: 'pp-2',
        page: 1,
        pagesize: 5,
        search: 'q',
        order: 'created_at',
      },
      { skip: true },
    );

    expect(result).toEqual({
      hook: 'getPerformanceProfileResults',
      queryArg: {
        performanceProfileId: 'pp-2',
        page: '1',
        pagesize: '5',
        search: 'q',
        order: 'created_at',
      },
      options: { skip: true },
    });
    expect(schemaMocks.useGetPerformanceProfileResultsQuery).toHaveBeenCalledTimes(1);
  });
});
