import { ctxUrl } from '../utils/multi-ctx';
import {
  mesheryApi,
  useGetTeamsQuery as useSchemasGetTeamsQuery,
  useGetUserProfileByIdQuery,
} from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';
import { initiateQuery } from './utils';
import { useGetOrgsQuery } from './organization';
import { useGetWorkspacesQuery } from './workspace';
import { normalizeLoadTestPrefs } from '../lib/load-test-prefs';

const Tags = {
  USER_PREF: 'userPref',
  USER: 'User_users',
  LOAD_TEST_PREF: 'loadTestPref',
  PROVIDER_CAP: 'provider_capabilities',
};

const normalizeUserPreference = (preference = {}) => ({
  ...preference,
  selectedOrganizationId: preference.selectedOrganizationId ?? preference.selectedOrganizationID,
  grafana: preference.grafana
    ? {
        ...preference.grafana,
        grafanaUrl: preference.grafana.grafanaUrl ?? preference.grafana.grafanaURL,
        grafanaApiKey: preference.grafana.grafanaApiKey ?? preference.grafana.grafanaAPIKey,
      }
    : preference.grafana,
  prometheus: preference.prometheus
    ? {
        ...preference.prometheus,
        prometheusUrl: preference.prometheus.prometheusUrl ?? preference.prometheus.prometheusURL,
      }
    : preference.prometheus,
});

const canonicalizeUserPreferencePayload = (preference = {}) => {
  const payload = { ...preference };
  if (payload.selectedOrganizationID && !payload.selectedOrganizationId) {
    payload.selectedOrganizationId = payload.selectedOrganizationID;
  }
  delete payload.selectedOrganizationID;

  if (payload.grafana) {
    payload.grafana = {
      ...payload.grafana,
      grafanaUrl: payload.grafana.grafanaUrl ?? payload.grafana.grafanaURL,
      grafanaApiKey: payload.grafana.grafanaApiKey ?? payload.grafana.grafanaAPIKey,
    };
    delete payload.grafana.grafanaURL;
    delete payload.grafana.grafanaAPIKey;
  }

  if (payload.prometheus) {
    payload.prometheus = {
      ...payload.prometheus,
      prometheusUrl: payload.prometheus.prometheusUrl ?? payload.prometheus.prometheusURL,
    };
    delete payload.prometheus.prometheusURL;
  }

  return payload;
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
        transformResponse: (response) => normalizeLoadTestPrefs(response?.loadTestPrefs),
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
      getToken: builder.query({
        query: () => ({ url: `/api/token`, method: 'GET', credentials: 'include' }),
        method: 'GET',
      }),
      getUserPref: builder.query({
        query: () => '/api/identity/users/preferences',
        method: 'GET',
        providesTags: [Tags.USER_PREF],
        transformResponse: normalizeUserPreference,
      }),
      updateUserPref: builder.mutation({
        query: (queryArg) => ({
          url: '/api/identity/users/preferences',
          method: 'PUT',
          body: canonicalizeUserPreferencePayload(queryArg),
          credentials: 'include',
        }),
        invalidatesTags: [Tags.USER_PREF],
      }),
      getUserPrefWithContext: builder.query({
        query: (selectedK8sContexts) => ({
          url: ctxUrl('/api/user/prefs', selectedK8sContexts),
          method: 'GET',
          credentials: 'same-origin',
        }),
        providesTags: [Tags.USER_PREF],
      }),
      updateUserPrefWithContext: builder.mutation({
        query: (queryArg) => ({
          url: ctxUrl('/api/user/prefs', queryArg.selectedK8sContexts),
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
              Object.assign(draft, normalizeUserPreference(queryArg.body));
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
        query: () => '/api/provider/capabilities',
        method: 'GET',
      }),
      getExtensionsByType: builder.query({
        query: () => ({
          url: '/api/provider/capabilities',
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
          url: '/api/provider/capabilities',
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
        query: () => '/api/system/version',
        method: 'GET',
      }),
      handleFeedbackFormSubmission: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`extensions/api/identity/users/notify/feedback`),
          method: 'POST',
          body: queryArg.userFeedbackRequestBody,
        }),
        invalidatesTags: [Tags.USER],
      }),
      removeUserFromTeam: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `extensions/api/identity/orgs/${queryArg.orgId}/teams/${queryArg.teamId}/users/${queryArg.userId}`,
          ),
          method: 'DELETE',
        }),
        invalidatesTags: ['teams'],
      }),
      handleUserInvite: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`extensions/api/identity/orgs/${queryArg.orgId}/users/invite`),
          method: 'POST',
          body: queryArg.userInvite,
        }),
        invalidatesTags: [Tags.USER],
      }),
      getAccessToken: builder.query({
        query: () => ({
          url: `/api/token`,
        }),
        transformResponse: (response) => {
          return response?.token;
        },
      }),
    }),
    overrideExisting: true,
  });

export const {
  useGetExtensionsByTypeQuery,
  useLazyGetExtensionsByTypeQuery,
  useGetFullPageExtensionsQuery,
  useLazyGetFullPageExtensionsQuery,
  useGetLoadTestPrefsQuery,
  useUpdateLoadTestPrefsMutation,
  useHandleUserInviteMutation,
  useLazyGetTokenQuery,
  useGetUserPrefQuery,
  useUpdateUserPrefMutation,
  useGetUserPrefWithContextQuery,
  useUpdateUserPrefWithContextMutation,
  useGetProviderCapabilitiesQuery,
  useHandleFeedbackFormSubmissionMutation,
  useRemoveUserFromTeamMutation,
  useGetSystemVersionQuery,
} = userApi;

export {
  useGetUserProfileByIdQuery as useGetUserProfileSummaryByIdQuery,
  useGetUserQuery as useGetLoggedInUserQuery,
  useGetUsersQuery as useGetAllUsersQuery,
  useGetUsersQuery as useGetUsersForOrgQuery,
} from '@meshery/schemas/mesheryApi';

export const useGetUserByIdQuery = (id, options = {}) =>
  useGetUserProfileByIdQuery(
    {
      id,
    },
    // Falsy id must always skip — a caller's explicit skip can tighten but
    // not loosen that invariant. Merging options first and forcing skip last
    // prevents `{skip: false}` in options from re-enabling the query with an
    // empty/invalid UUID and reintroducing the 400/404 loop this wrapper
    // exists to prevent.
    { ...options, skip: !id || options?.skip },
  );

export const useGetTeamsQuery = (queryArg, options) =>
  useSchemasGetTeamsQuery(
    {
      orgId: queryArg?.orgId,
      search: queryArg?.search,
      order: queryArg?.order,
      page: queryArg?.page?.toString(),
      pagesize: queryArg?.pagesize?.toString(),
    },
    options,
  );

export const useLazyGetTeamsQuery = () => {
  const [trigger, result, lastPromiseInfo] = mesheryApi.endpoints.getTeams.useLazyQuery();

  const wrappedTrigger = (queryArg, preferCacheValue) =>
    trigger(
      {
        orgId: queryArg?.orgId,
        search: queryArg?.search,
        order: queryArg?.order,
        page: queryArg?.page?.toString(),
        pagesize: queryArg?.pagesize?.toString(),
      },
      preferCacheValue,
    );

  return [wrappedTrigger, result, lastPromiseInfo] as const;
};

export const getProviderCapabilities = async () => {
  const res = await initiateQuery(userApi.endpoints.getProviderCapabilities);
  return res;
};

export const getUserAccessToken = async () => {
  const accessToken = await initiateQuery(userApi.endpoints.getAccessToken, {}, {});
  return accessToken;
};

export const getUserProfile = async () => {
  const userProfile = await initiateQuery(mesheryApi.endpoints.getUser);
  return userProfile;
};

export const getSystemVersion = async () => {
  const res = await initiateQuery(userApi.endpoints.getSystemVersion);
  return res;
};

export const getAllUsers = async ({ page, pageSize, pagesize, search }) => {
  const users = await initiateQuery(mesheryApi.endpoints.getUsers, {
    page: page?.toString(),
    pageSize: (pageSize ?? pagesize)?.toString(),
    search,
  });
  return users;
};

export const useGetSelectedOrganization = () => {
  const {
    data: userPrefs,
    isLoading: isLoadingUserPrefs,
    error: errorLoadingUserPrefs,
  } = useGetUserPrefQuery();
  const {
    data: allOrgs,
    isLoading: isLoadingAllOrgs,
    error: errorLoadingAllOrgs,
  } = useGetOrgsQuery();

  const existingSelectedOrganization = allOrgs?.organizations?.find(
    (org) => org.id === userPrefs?.selectedOrganizationId,
  );

  const selectedOrganization = existingSelectedOrganization ?? allOrgs?.organizations?.[0];

  return {
    selectedOrganization,
    allOrganizations: allOrgs?.organizations || [],
    didFallback: !existingSelectedOrganization,
    isLoading: isLoadingUserPrefs || isLoadingAllOrgs,
    isError: errorLoadingUserPrefs || errorLoadingAllOrgs,
    error: errorLoadingUserPrefs || errorLoadingAllOrgs,
  };
};

export const useGetSelectedWorkspace = () => {
  const {
    selectedOrganization,
    isLoading: isLoadingOrganizations,
    error: errorGetSelectedOrg,
  } = useGetSelectedOrganization();
  const {
    data: workspacesData,
    isError: isWorkspacesError,
    isLoading: isLoadingingWorkspaces,
    error: errorGetWorkspaces,
  } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      order: 'updated_at desc',
      orgId: selectedOrganization?.id,
    },
    {
      skip: !selectedOrganization?.id,
    },
  );
  // const [updateSelectedWorkspace] = useUpdateSelectedWorkspaceMutation();
  const { data: userPrefs, isLoading: isLoadingPrefs } = useGetUserPrefQuery();
  const selectedWorkspaceID =
    userPrefs?.selectedWorkspaceForOrganizations?.[selectedOrganization?.id];

  const existingSelectedWorkspace = (workspacesData?.workspaces ?? []).find(
    (workspace) => workspace.id === selectedWorkspaceID,
  );

  const selectedWorkspace = existingSelectedWorkspace ?? workspacesData?.workspaces?.[0];

  const didFallback = !existingSelectedWorkspace;

  return {
    selectedWorkspace,
    didFallback,
    allWorkspaces: workspacesData?.workspaces,
    isLoading: isLoadingOrganizations || isLoadingingWorkspaces || isLoadingPrefs,
    isError: isWorkspacesError || errorGetSelectedOrg,
    error: errorGetWorkspaces || errorGetSelectedOrg,
  };
};

export const useUpdateSelectedOrganizationMutation = () => {
  const [updateUserPref, response] = useUpdateUserPrefMutation();

  const updateSelectedOrganization = async (orgId) => {
    await updateUserPref({ selectedOrganizationId: orgId });
  };

  return [updateSelectedOrganization, response];
};

export const useUpdateSelectedWorkspaceMutation = () => {
  const { data: userPrefs } = useGetUserPrefQuery();
  const [updateUserPref, response] = useUpdateUserPrefMutation();

  const updateSelectedWorkspace = async (orgId, workspaceId) => {
    await updateUserPref({
      selectedWorkspaceForOrganizations: {
        ...(userPrefs.selectedWorkspaceForOrganizations || {}),
        [orgId]: workspaceId,
      },
    });
  };

  return [updateSelectedWorkspace, response];
};
