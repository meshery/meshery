import { ctxUrl } from '@/utils/multi-ctx';
import { api } from './index';

const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getLoggedInUser: builder.query({
      query: () => `user`,
      method: 'GET',
    }),
    getUserById: builder.query({
      query: (id) => `user/profile/${id}`,
      method: 'GET',
    }),
    getToken: builder.query({
      query: () => `token`,
      method: 'GET',
    }),
    getUserPref: builder.query({
      query: () => 'user/prefs',
      method: 'GET',
    }),
    updateUserPref: builder.mutation({
      query: (queryArg) => ({
        url: 'user/prefs',
        method: 'POST',
        body: queryArg.usersExtensionPreferences,
        credentials: 'include',
      }),
    }),
    getUserPrefWithContext: builder.query({
      query: (selectedK8sContexts) => ({
        url: ctxUrl('user/prefs', selectedK8sContexts),
        method: 'GET',
        credentials: 'same-origin',
      }),
    }),
    updateUserPrefWithContext: builder.mutation({
      query: (queryArg) => ({
        url: ctxUrl('/user/prefs', queryArg.selectedK8sContexts),
        method: 'POST',
        headers: {
          'Content-Type': 'application/json;charset=UTF-8',
        },
        body: queryArg.body,
      }),
    }),
    getProviderCapabilities: builder.query({
      query: () => 'provider/capabilities',
      method: 'GET',
    }),
  }),
});

export const {
  useGetLoggedInUserQuery,
  useGetUserByIdQuery,
  useLazyGetTokenQuery,
  useGetUserPrefQuery,
  useUpdateUserPrefMutation,
  useGetUserPrefWithContextQuery,
  useUpdateUserPrefWithContextMutation,
  useGetProviderCapabilitiesQuery,
} = userApi;
