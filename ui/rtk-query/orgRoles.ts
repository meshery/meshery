import { api } from './index';

const userOrgRolesApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getUserOrgRoles: builder.query({
      query: (queryArgs) => ({
        url: `extensions/api/identity/orgs/${queryArgs.orgId}/roles`,
        method: 'GET',
        credentials: 'include',
        params: {
          page: queryArgs.page,
          pagesize: queryArgs.pagesize,
          search: queryArgs.search,
          order: queryArgs.order,
          all: queryArgs.all,
          selector: queryArgs.selector,
        },
      }),
    }),
  }),
});

export const { useGetUserOrgRolesQuery } = userOrgRolesApi;
