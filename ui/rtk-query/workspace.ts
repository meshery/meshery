import { urlEncodeParams } from '@/utils/utils';
import { mesheryApi } from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';
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
        keepUnusedDataFor: 0,
        queryFn: async (queryArgs, { dispatch }, _extraOptions, baseQuery) => {
          const { expandInfo, ...otherArgs } = queryArgs;
          const params = urlEncodeParams(otherArgs);
          const workspaces = await baseQuery({
            url: mesheryApiPath(`workspaces?${params}`),
            method: 'GET',
          });

          if (expandInfo && workspaces.data && !workspaces.error) {
            const modifiedWorkspaces = await Promise.all(
              workspaces.data.workspaces.map(async (workspace) => {
                const [designs, environments, views, teams] = await Promise.all([
                  dispatch(
                    workspacesApi.endpoints.getDesignsOfWorkspace.initiate({
                      workspaceId: workspace.id,
                      expandUser: true,
                      infiniteScroll: false,
                      page: 0,
                      pagesize: 1,
                    }),
                  ),
                  dispatch(
                    workspacesApi.endpoints.getEnvironmentsOfWorkspace.initiate({
                      workspaceId: workspace.id,
                      page: 0,
                      pagesize: 1,
                    }),
                  ),
                  dispatch(
                    workspacesApi.endpoints.getViewsOfWorkspace.initiate({
                      workspaceId: workspace.id,
                      page: 0,
                      pagesize: 1,
                    }),
                  ),
                  dispatch(
                    workspacesApi.endpoints.getTeamsOfWorkspace.initiate({
                      workspaceId: workspace.id,
                      page: 0,
                      pagesize: 1,
                    }),
                  ),
                ]);

                return {
                  ...workspace,
                  designCount: designs.data?.total_count || 0,
                  environmentCount: environments.data?.total_count || 0,
                  viewCount: views.data?.total_count || 0,
                  teamCount: teams.data?.total_count || 0,
                };
              }),
            );

            return _.merge({}, workspaces, { data: { workspaces: modifiedWorkspaces } });
          }
          return workspaces;
        },
        providesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      getEnvironmentsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces/${queryArg.workspaceId}/environments`),
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
          url: mesheryApiPath(
            `workspaces/${queryArg.workspaceId}/environments/${queryArg.environmentId}`,
          ),
          method: 'POST',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENTS }],
      }),

      unassignEnvironmentFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `workspaces/${queryArg.workspaceId}/environments/${queryArg.environmentId}`,
          ),
          method: 'DELETE',
        }),

        invalidatesTags: () => [{ type: TAGS.ENVIRONMENTS }],
      }),

      getDesignsOfWorkspace: builder.query({
        queryFn: async (queryArgs, { dispatch }, _extraOptions, baseQuery) => {
          const { expandUser, infiniteScroll: _infiniteScroll, ...otherArgs } = queryArgs;
          const params = urlEncodeParams(otherArgs);
          const designs = await baseQuery({
            url: mesheryApiPath(`workspaces/${queryArgs.workspaceId}/designs?${params}`),
            method: 'GET',
          });
          if (expandUser && designs.data && !designs.error) {
            const withUsersPromises = designs.data.designs.map(async (design) => {
              const user = await dispatch(
                mesheryApi.endpoints.getUserProfileById.initiate({ id: design.user_id }),
              );
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
      }),
      assignDesignToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces/${queryArg.workspaceId}/designs/${queryArg.designId}`),
          method: 'POST',
        }),
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),

      unassignDesignFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces/${queryArg.workspaceId}/designs/${queryArg.designId}`),
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.DESIGNS }],
      }),
      getViewsOfWorkspace: builder.query({
        queryFn: async (queryArg, { dispatch }, _extraOptions, baseQuery) => {
          const { expandUser, infiniteScroll: _infiniteScroll, ...otherArgs } = queryArg;
          const params = urlEncodeParams(otherArgs);
          const views = await baseQuery({
            url: mesheryApiPath(
              `extensions/api/workspaces/${queryArg.workspaceId}/views?${params}`,
            ),
            method: 'GET',
          });
          if (expandUser && views.data && !views.error) {
            const withUsersPromises = views.data.views.map(async (view) => {
              const user = await dispatch(
                mesheryApi.endpoints.getUserProfileById.initiate({ id: view.user_id }),
              );
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
      }),
      assignViewToWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `extensions/api/workspaces/${queryArg.workspaceId}/views/${queryArg.viewId}`,
          ),
          method: 'POST',
        }),
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),

      unassignViewFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `extensions/api/workspaces/${queryArg.workspaceId}/views/${queryArg.viewId}`,
          ),
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.VIEWS }],
      }),

      getTeamsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(`extensions/api/workspaces/${queryArg.workspaceId}/teams`),
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
          url: mesheryApiPath(
            `extensions/api/workspaces/${queryArg.workspaceId}/teams/${queryArg.teamId}`,
          ),
          method: 'POST',
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),

      unassignTeamFromWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `extensions/api/workspaces/${queryArg.workspaceId}/teams/${queryArg.teamId}`,
          ),
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),

      getEventsOfWorkspace: builder.query({
        query: (queryArg) => ({
          url: mesheryApiPath(`extensions/api/workspaces/${queryArg.workspaceId}/events`),
          params: {
            page: queryArg.page,
            pagesize: queryArg.pagesize,
            search: queryArg.search,
            order: queryArg.order,
          },
        }),
        invalidatesTags: () => [{ type: TAGS.TEAMS }],
      }),

      createWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces`),
          method: 'POST',
          body: {
            name: queryArg.name,
            description: queryArg.description,
            organization_id: queryArg.organization_id,
          },
        }),
        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      updateWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces/${queryArg.id}`),
          method: 'PUT',
          body: {
            name: queryArg.name,
            description: queryArg.description,
            organization_id: queryArg.organization_id,
          },
        }),
        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      deleteWorkspace: builder.mutation({
        query: (queryArg) => ({
          url: mesheryApiPath(`workspaces/${queryArg.id}`),
          method: 'DELETE',
        }),
        invalidatesTags: () => [{ type: TAGS.WORKSPACES }],
      }),
    }),
  });

export const {
  useGetWorkspacesQuery,
  useLazyGetWorkspacesQuery,
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

export const useCreateWorkspaceMutation = () => {
  const [trigger, result] = workspacesApi.endpoints.createWorkspace.useMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      name: queryArg.workspacePayload?.name,
      description: queryArg.workspacePayload?.description,
      organization_id: queryArg.workspacePayload?.organization_id,
    });

  return [wrappedTrigger, result] as const;
};

export const useUpdateWorkspaceMutation = () => {
  const [trigger, result] = workspacesApi.endpoints.updateWorkspace.useMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      id: queryArg.workspaceId,
      name: queryArg.workspacePayload?.name,
      description: queryArg.workspacePayload?.description,
      organization_id: queryArg.workspacePayload?.organization_id,
    });

  return [wrappedTrigger, result] as const;
};

export const useDeleteWorkspaceMutation = () => {
  const [trigger, result] = workspacesApi.endpoints.deleteWorkspace.useMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      id: queryArg.workspaceId,
    });

  return [wrappedTrigger, result] as const;
};
