import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// organization.ts wraps useGetOrgsQuery from @meshery/schemas/mesheryApi.
// The wrapper coerces page/pagesize to strings and forwards search/order/all.
// We mock the underlying schemas hook and assert the wrapper's contract.
// ---------------------------------------------------------------------------

const { mockedSchemasGetOrgs } = vi.hoisted(() => ({ mockedSchemasGetOrgs: vi.fn() }));
vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetOrgsQuery: (...args: unknown[]) => mockedSchemasGetOrgs(...args),
}));

import { useGetOrgsQuery } from '../organization';

describe('useGetOrgsQuery', () => {
  beforeEach(() => {
    mockedSchemasGetOrgs.mockReset();
    mockedSchemasGetOrgs.mockReturnValue({ data: undefined, isLoading: true });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('forwards stringified page and pagesize to the underlying hook', () => {
    useGetOrgsQuery({ page: 1, pagesize: 25, search: 'team', order: 'asc', all: true });

    expect(mockedSchemasGetOrgs).toHaveBeenCalledWith(
      {
        page: '1',
        pagesize: '25',
        search: 'team',
        order: 'asc',
        all: true,
      },
      undefined,
    );
  });

  it('passes options through as the second argument', () => {
    const options = { skip: true };
    useGetOrgsQuery({ page: 0, pagesize: 10 }, options);

    expect(mockedSchemasGetOrgs).toHaveBeenCalledWith(
      expect.objectContaining({ page: '0', pagesize: '10' }),
      options,
    );
  });

  it('handles undefined queryArgs gracefully', () => {
    useGetOrgsQuery(undefined, undefined);

    expect(mockedSchemasGetOrgs).toHaveBeenCalledWith(
      {
        page: undefined,
        pagesize: undefined,
        search: undefined,
        order: undefined,
        all: undefined,
      },
      undefined,
    );
  });

  it('returns the result of the underlying hook', () => {
    const expected = { data: { orgs: [{ id: 'org-1' }] }, isLoading: false };
    mockedSchemasGetOrgs.mockReturnValue(expected);

    const result = useGetOrgsQuery({ page: 0, pagesize: 10 });
    expect(result).toBe(expected);
  });
});
