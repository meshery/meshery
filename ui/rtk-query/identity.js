import { api } from './index';

const TAGS = {
  ORGS: 'organizations',
  KEYS: 'keys',
  WORKSPACES: 'workspaces',
};

export const identityApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getOrganizations: builder.query({
      query: () => ({
        url: 'identity/orgs',
        method: 'GET',
      }),
      providesTags: [TAGS.ORGS],
    }),

    getUserKeys: builder.query({
      query: (orgId) => ({
        url: `identity/orgs/${orgId}/users/keys`,
        method: 'GET',
      }),
      providesTags: (result, error, orgId) => [{ type: TAGS.KEYS, id: orgId }],
    }),

    getWorkspaces: builder.query({
      query: ({ search = '', order = '', page = 0, pagesize = 10, orgID }) => ({
        url: 'workspaces',
        method: 'GET',
        params: {
          search,
          order,
          page,
          pagesize,
          orgID,
        },
      }),
      providesTags: (result, error, { orgID }) => [{ type: TAGS.WORKSPACES, id: orgID }],
    }),
  }),
});

export const {
  useGetOrganizationsQuery,
  useLazyGetOrganizationsQuery,
  useGetUserKeysQuery,
  useLazyGetUserKeysQuery,
  useGetWorkspacesQuery,
  useLazyGetWorkspacesQuery,
} = identityApi;
