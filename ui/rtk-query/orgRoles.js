import { api } from './index';

const userOrgRolesApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getUserOrgRoles: builder.query({
      query: (queryArgs) => ({
        url: `extensions/api/identity/orgs/${queryArgs.orgId}/roles`,
        method: 'GET',
        credentials: 'include',
      }),
    }),
  }),
});

export const { useGetUserOrgRolesQuery } = userOrgRolesApi;
