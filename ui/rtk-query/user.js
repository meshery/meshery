import { ctxUrl } from '@/utils/multi-ctx';
import { api } from './index';
import { initiateQuery } from './utils';

const Tags = {
  USER_PREF: 'userPref',
  LOAD_TEST_PREF: 'loadTestPref',
  PROVIDER_CAP: 'provider_capabilities',
};

export const userApi = api
  .enhanceEndpoints({
    addTagTypes: [Tags.USER_PREF, Tags.LOAD_TEST_PREF, Tags.PROVIDER_CAP],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getLoadTestPrefs: builder.query({
        query: (selectedK8sContexts) => ({
          url: ctxUrl('/api/user/prefs', selectedK8sContexts),
          method: 'GET',
          credentials: 'include',
        }),
        providesTags: [Tags.LOAD_TEST_PREF],
        // Transform response to directly get the loadTestPrefs
        transformResponse: (response) => response?.loadTestPrefs || {},
      }),

      updateLoadTestPrefs: builder.mutation({
        query: (queryArg) => ({
          url: ctxUrl('/api/user/prefs', queryArg.selectedK8sContexts),
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json;charset=UTF-8' },
          body: JSON.stringify({ loadTestPrefs: queryArg.loadTestPrefs }),
        }),
        invalidatesTags: [Tags.LOAD_TEST_PREF],
      }),
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
      getUserProfileSummaryById: builder.query({
        query: (queryArg) => ({
          url: `/user/profile/${queryArg.id}`,
        }),
        transformResponse: (response) => {
          // Modify the response data to keep only necessary fields
          return {
            id: response.id,
            email: response?.email,
            user_id: response?.user_id,
            avatar_url: response?.avatar_url,
            first_name: response?.first_name,
            last_name: response?.last_name,
          };
        },
      }),
      getExtensionsByType: builder.query({
        query: () => ({
          url: 'provider/capabilities',
          method: 'GET',
          credentials: 'include',
        }),
        transformResponse: (response, _, type) => {
          if (!response?.extensions || !response?.extensions[type]) {
            return [];
          }

          try {
            const ExtensionPointSchemaValidator =
              require('../utils/ExtensionPointSchemaValidator').default;
            return ExtensionPointSchemaValidator(type)(response?.extensions[type]);
          } catch (error) {
            console.group('extension error');
            console.error(error);
            console.groupEnd();
            return [];
          }
        },
        providesTags: [Tags.PROVIDER_CAP],
      }),
      getFullPageExtensions: builder.query({
        query: () => ({
          url: 'provider/capabilities',
          method: 'GET',
          credentials: 'include',
        }),
        transformResponse: (response) => {
          if (!response?.extensions) {
            return [];
          }

          let extNames = [];
          for (var key of Object.keys(response.extensions)) {
            if (Array.isArray(response.extensions[key])) {
              response.extensions[key].forEach((comp) => {
                if (comp?.type === 'full_page') {
                  let ext = {
                    name: key,
                    uri: comp?.href?.uri,
                  };
                  extNames.push(ext);
                }
              });
            }
          }

          return extNames;
        },
        // Make sure we have proper tag
        providesTags: [Tags.PROVIDER_CAP],
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
      getAllUsers: builder.query({
        query: (queryArg) => ({
          url: `identity/users`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
            filter: queryArg.filter,
          },
        }),
        providesTags: ['users'],
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
  useGetUserProfileSummaryByIdQuery,
  useGetExtensionsByTypeQuery,
  useLazyGetExtensionsByTypeQuery,
  useGetFullPageExtensionsQuery,
  useLazyGetFullPageExtensionsQuery,
  useGetLoadTestPrefsQuery,
  useUpdateLoadTestPrefsMutation,
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
  useGetAllUsersQuery,
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

export const getAllUsers = async ({ page, pagesize, search }) => {
  const users = await initiateQuery(
    userApi.endpoints.getAllUsers,
    { page, pagesize, search },
    { skip: !search },
  );
  return users;
};
