import { urlEncodeParams } from '@/utils/utils';
import { api } from './index';
import { userApi } from './user';
import _ from 'lodash';

const TAGS = {
  WORKSPACES: 'workspaces',
  DESIGNS: 'workspaces_designs',
  ENVIRONMENTS: 'workspaces_environments',
  VIEWS: 'workspaces_views',
  TEAMS: 'workspaces_teams',
};
const workspacesApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.WORKSPACES, TAGS.DESIGNS, TAGS.ENVIRONMENTS, TAGS.VIEWS, TAGS.TEAMS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getWorkspaces: builder.query({
        query: (queryArg) => ({
          url: `workspaces`,
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page || 0,
            pagesize: queryArg.pagesize || 'all',
            orgID: queryArg.orgId,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      createWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces`,
          method: 'POST',
          body: queryArg.workspacePayload,
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      updateWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}`,
          method: 'PUT',
          body: queryArg.workspacePayload,
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      deleteWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}`,
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      getEnvironmentsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}/environments`,
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            filter: queryArg.filter,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.ENVIRONMENTS }],
      }),

      assignEnvironmentToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}/environments/${queryArg.environmentId}`,
          method: 'POST',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENTS }],
      }),

      unassignEnvironmentFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}/environments/${queryArg.environmentId}`,
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENTS }],
      }),

      getDesignsOfWorkspace: builder.query({
        queryFn: async (queryArgs, { dispatch }, _extraOptions, baseQuery) => {
          // eslint-disable-next-line no-unused-vars
          const { expandUser, infiniteScroll, ...otherArgs } = queryArgs;
          const params = urlEncodeParams(otherArgs);
          const designs = await baseQuery({
            url: `workspaces/${queryArgs.workspaceId}/designs?${params}`,
            method: 'GET',
          });
          if (expandUser && designs.data && !designs.error) {
            const withUsersPromises = designs.data.designs.map(async (design) => {
              const user = await dispatch(userApi.endpoints.getUserById.initiate(design.user_id));
              return {
                ...design,
                first_name: user.data?.first_name || '[deleted]',
                last_name: user.data?.last_name || '',
                avatar_url: user.data?.avatar_url || '',
                user_id: user.data?.id || '',
                email: user.data?.email || '',
              };
            });

            const modifiedDesigns = await Promise.all(withUsersPromises);
            return _.merge({}, designs, { data: { designs: modifiedDesigns } });
          }

          return designs;
        },
        serializeQueryArgs: ({ endpointName, queryArgs }) => {
          if (queryArgs?.infiniteScroll) {
            return endpointName;
          }
          return `${endpointName}-${JSON.stringify(queryArgs)}`;
        },
        merge: (currentCache, newItems, { arg }) => {
          if (!arg.infiniteScroll) {
            return newItems;
          }

          if (arg.page === 0) {
            return newItems;
          }
          return {
            ...(currentCache || {}),
            ...(newItems || {}),
            designs: [...(currentCache?.designs || []), ...(newItems?.designs || [])],
          };
        },
        forceRefetch({ currentArg, previousArg }) {
          if (!currentArg.infiniteScroll) {
            return true;
          }
          return !_.eq(currentArg, previousArg);
        },
        providesTags: () => [{ type: TAGS.DESIGNS }],
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      assignDesignToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}/designs/${queryArg.designId}`,
          method: 'POST',
        }),

        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),

      unassignDesignFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `workspaces/${queryArg.workspaceId}/designs/${queryArg.designId}`,
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      getViewsOfWorkspace: builder.query({
        queryFn: async (queryArg, { dispatch }, _extraOptions, baseQuery) => {
          // eslint-disable-next-line no-unused-vars
          const { expandUser, infiniteScroll, ...otherArgs } = queryArg;
          const params = urlEncodeParams(otherArgs);
          const views = await baseQuery({
            url: `extensions/api/workspaces/${queryArg.workspaceId}/views?${params}`,
            method: 'GET',
          });
          if (expandUser && views.data && !views.error) {
            const withUsersPromises = views.data.views.map(async (view) => {
              const user = await dispatch(userApi.endpoints.getUserById.initiate(view.user_id));
              return {
                ...view,
                first_name: user.data?.first_name || '[deleted]',
                last_name: user.data?.last_name || '',
                avatar_url: user.data?.avatar_url || '',
                user_id: user.data?.id || '',
                email: user.data?.email || '',
              };
            });
            const modifiedViews = await Promise.all(withUsersPromises);
            return _.merge({}, views, { data: { views: modifiedViews } });
          }

          return views;
        },
        serializeQueryArgs: ({ endpointName, queryArgs }) => {
          if (queryArgs?.infiniteScroll) {
            return endpointName;
          }
          return `${endpointName}-${JSON.stringify(queryArgs)}`;
        },
        merge: (currentCache, newItems, { arg }) => {
          if (!arg.infiniteScroll) {
            return newItems;
          }

          if (arg.page === 0) {
            return newItems;
          }
          return {
            ...(currentCache || {}),
            ...(newItems || {}),
            views: [...(currentCache?.views || []), ...(newItems?.views || [])],
          };
        },
        forceRefetch({ currentArg, previousArg }) {
          if (!currentArg.infiniteScroll) {
            return true;
          }
          return !_.eq(currentArg, previousArg);
        },
        providesTags: () => [{ type: TAGS.VIEWS }],
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),
      assignViewToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/views/${queryArg.viewId}`,
          method: 'POST',
        }),
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),

      unassignViewFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/views/${queryArg.viewId}`,
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),

      getTeamsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/teams`,
          params: {
            search: queryArg.search,
            order: queryArg.order,
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            filter: queryArg.filter,
          },
          method: 'GET',
        }),
        providesTags: () => [{ type: TAGS.TEAMS }],
      }),

      assignTeamToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/teams/${queryArg.teamId}`,
          method: 'POST',
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),

      unassignTeamFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/teams/${queryArg.teamId}`,
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),

      getEventsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: `extensions/api/workspaces/${queryArg.workspaceId}/events`,
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),
    }),
  });

export const {
  useGetWorkspacesQuery,
  useCreateWorkspaceMutation,
  useUpdateWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetEnvironmentsOfWorkspaceQuery,
  useAssignEnvironmentToWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useAssignDesignToWorkspaceMutation,
  useUnassignDesignFromWorkspaceMutation,
  useGetViewsOfWorkspaceQuery,
  useAssignViewToWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
  useGetTeamsOfWorkspaceQuery,
  useAssignTeamToWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useGetEventsOfWorkspaceQuery,
} = workspacesApi;
