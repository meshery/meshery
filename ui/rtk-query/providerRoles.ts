import { api } from './index';
import { shouldOverrideExisting } from './utils';

const userProviderRolesApi = api.injectEndpoints({
  overrideExisting: shouldOverrideExisting,
  endpoints: (builder) => ({
    getUserProviderRoles: builder.query({
      query: () => ({
        url: `/api/user`,
        method: 'GET',
        credentials: 'include',
      }),
    }),
  }),
});

export const { useGetUserProviderRolesQuery } = userProviderRolesApi;
