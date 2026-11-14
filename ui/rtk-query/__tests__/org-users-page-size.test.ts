import { describe, expect, it, vi } from 'vitest';

vi.mock('@/utils/utils', () => ({ urlEncodeArrayParam: () => '', urlEncodeParams: () => '' }));
vi.mock('../utils', () => ({ appendInvalidatesTags: () => () => {}, initiateQuery: vi.fn() }));

const { getUsersForOrg } = vi.hoisted(() => ({ getUsersForOrg: vi.fn() }));
vi.mock('@meshery/schemas/mesheryApi', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@meshery/schemas/mesheryApi')>()),
  useGetUsersForOrgQuery: getUsersForOrg,
}));

import { useGetUsersForOrgQuery } from '../user';

describe('useGetUsersForOrgQuery', () => {
  it.each([10, 20, 25])('forwards page size %i to the generated endpoint', (pagesize) => {
    const options = { skip: true };
    useGetUsersForOrgQuery({ orgId: 'org-1', page: 0, pagesize }, options);
    expect(getUsersForOrg).toHaveBeenLastCalledWith(
      expect.objectContaining({ orgId: 'org-1', page: '0', pageSize: String(pagesize) }),
      options,
    );
    expect(getUsersForOrg.mock.lastCall?.[0]).not.toHaveProperty('pagesize');
  });
});
