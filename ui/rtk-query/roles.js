import { api } from './index';

const TAGS = {
  ROLES: 'roles',
};

const userRolesApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.ROLES],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getUserKeys: builder.query({
        query: (queryArgs) => ({
          url: `extensions/api/identity/orgs/${queryArgs.orgId}/roles`,
          params: {},
          method: 'GET',
          credentials: 'include',
        }),
      }),
    }),
  });

export const { useLazyGetUserKeysQuery } = userRolesApi;
