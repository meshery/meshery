import { api, mesheryApiPath } from './index';

const userOrgRolesApi = api.injectEndpoints({
  overrideExisting: module.hot?.status() === 'apply',
  endpoints: (builder) => ({
    getUserOrgRoles: builder.query({
      query: (queryArgs) => ({
        url: mesheryApiPath(`extensions/api/identity/orgs/${queryArgs.orgId}/roles`),
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
