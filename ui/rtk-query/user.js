import { ctxUrl } from '@/utils/multi-ctx';
import { api } from './index';
import { initiateQuery } from './utils';

const Tags = {
  USER_PREF: 'userPref',
};

export const userApi = api
  .enhanceEndpoints({
    addTagTypes: [Tags.USER_PREF],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getLoggedInUser: builder.query({
        query: () => `user`,
        method: 'GET',
      }),
      getUserById: builder.query({
        query: (id) => `user/profile/${id}`,
        method: 'GET',
        providesTags: [Tags.USER_PREF],
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
          body: queryArg,
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
      getSystemVersion: builder.query({
        query: () => 'system/version',
        method: 'GET',
      }),
      handleFeedbackFormSubmission: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/identity/users/notify/feedback`,
          method: 'POST',
          body: queryArg.userFeedbackRequestBody,
        }),
        invalidatesTags: ['users'],
      }),
      getUsersForOrg: builder.query({
        query: (queryArg) => ({
          url: `extensions/api/identity/orgs/${queryArg.orgId}/users`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            filter: queryArg.filter,
            teamID: queryArg.teamId,
          },
        }),
        invalidatesTags: ['users'],
      }),
      removeUserFromTeam: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/identity/orgs/${queryArg.orgId}/teams/${queryArg.teamId}/users/${queryArg.userId}`,
          method: 'DELETE',
        }),
        invalidatesTags: ['teams'],
      }),
      handleUserInvite: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/identity/orgs/${queryArg.orgId}/users/invite`,
          method: 'POST',
          body: queryArg.userInvite,
        }),
        invalidatesTags: ['users'],
      }),
      getTeams: builder.query({
        query: (queryArg) => ({
          url: `extensions/api/identity/orgs/${queryArg.orgId}/teams`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
        }),
        invalidatesTags: ['teams'],
        providesTags: ['teams'],
      }),
    }),
  });

export const {
  useHandleUserInviteMutation,
  useGetLoggedInUserQuery,
  useGetUserByIdQuery,
  useLazyGetTokenQuery,
  useGetUserPrefQuery,
  useUpdateUserPrefMutation,
  useGetUserPrefWithContextQuery,
  useUpdateUserPrefWithContextMutation,
  useGetProviderCapabilitiesQuery,
  useHandleFeedbackFormSubmissionMutation,
  useGetUsersForOrgQuery,
  useRemoveUserFromTeamMutation,
  useGetTeamsQuery,
  useLazyGetTeamsQuery,
  useGetSystemVersionQuery,
} = userApi;

export const getProviderCapabilities = async () => {
  const res = await initiateQuery(userApi.endpoints.getProviderCapabilities);
  return res;
};

export const getSystemVersion = async () => {
  const res = await initiateQuery(userApi.endpoints.getSystemVersion);
  return res;
};
