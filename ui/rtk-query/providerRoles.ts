import { api } from './index';

const userProviderRolesApi = api.injectEndpoints({
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
