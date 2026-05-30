import { api } from './index';

const userProviderRolesApi = api.injectEndpoints({
  overrideExisting: module.hot?.status() === 'apply',
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
