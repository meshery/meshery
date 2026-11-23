import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// organization.ts wraps useGetOrgsQuery from @meshery/schemas/mesheryApi.
// The wrapper coerces page/pagesize to strings and forwards search/order/all.
// The generated getOrgs endpoint (schemas v1beta2) reads `pageSize` only, so
// the wrapper must send the page size under that key or it never reaches the
// server. We mock the underlying schemas hook and assert the wrapper's contract.
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

  it('forwards stringified page and the page size as pageSize to the underlying hook', () => {
    useGetOrgsQuery({ page: 1, pagesize: 25, search: 'team', order: 'asc', all: true });

    expect(mockedSchemasGetOrgs).toHaveBeenCalledWith(
      {
        page: '1',
        pageSize: '25',
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
      expect.objectContaining({ page: '0', pageSize: '10' }),
      options,
    );
  });

  it('handles undefined queryArgs gracefully', () => {
    useGetOrgsQuery(undefined, undefined);

    expect(mockedSchemasGetOrgs).toHaveBeenCalledWith(
      {
        page: undefined,
        pageSize: undefined,
        search: undefined,
        order: undefined,
        all: undefined,
      },
      undefined,
    );
  });

  it("forwards pagesize 'all' as pageSize so every organization is requested", () => {
    useGetOrgsQuery({ pagesize: 'all' });

    const [args] = mockedSchemasGetOrgs.mock.calls[0];
    expect(args.pageSize).toBe('all');
    expect(args).not.toHaveProperty('pagesize');
  });

  it('returns the result of the underlying hook', () => {
    const expected = { data: { orgs: [{ id: 'org-1' }] }, isLoading: false };
    mockedSchemasGetOrgs.mockReturnValue(expected);

    const result = useGetOrgsQuery({ page: 0, pagesize: 10 });
    expect(result).toBe(expected);
  });
});
