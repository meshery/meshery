import { ctxUrl } from '@/utils/multi-ctx';
import { api } from './index';

const Tags = {
  USER_PREF: 'userPref',
};

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
      providesTags: [Tags.USER_PREF],
    }),
    updateUserPref: builder.mutation({
      query: (queryArg) => ({
        url: 'user/prefs',
        method: 'POST',
        body: queryArg.usersExtensionPreferences,
        credentials: 'include',
      }),
      invalidatesTags: [Tags.USER_PREF],
    }),
    getUserPrefWithContext: builder.query({
      query: (selectedK8sContexts) => ({
        url: ctxUrl('user/prefs', selectedK8sContexts),
        method: 'GET',
        credentials: 'same-origin',
      }),
      providesTags: [Tags.USER_PREF],
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
      invalidatesTags: [Tags.USER_PREF],
      // Perform optimistic update
      onQueryStarted: async (queryArg, { dispatch, queryFulfilled }) => {
        // Optimistically update the cache with the new preferences
        const patchResult = dispatch(
          api.util.updateQueryData('getUserPref', queryArg.selectedK8sContexts, (draft) => {
            Object.assign(draft, queryArg.body);
          }),
        );
        try {
          // Wait for the mutation to complete
          await queryFulfilled;
        } catch {
          // If the mutation fails, revert the optimistic update
          patchResult.undo();
        }
      },
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
