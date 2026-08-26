import { urlEncodeParams } from '@/utils/utils';
import {
  mesheryApi,
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useCreateWorkspaceMutation as useSchemasCreateWorkspaceMutation,
  useDeleteWorkspaceMutation as useSchemasDeleteWorkspaceMutation,
  useGetEnvironmentsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUpdateWorkspaceMutation as useSchemasUpdateWorkspaceMutation,
} from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';
import _ from 'lodash';
import { normalizePaginatedCollectionResponse } from './transforms';
import { normalizeUserProfileSummary } from './userProfile';
import { appendInvalidatesTags } from './utils';

const TAGS = {
  WORKSPACES: 'workspaces',
  DESIGNS: 'workspaces_designs',
  ENVIRONMENTS: 'workspaces_environments',
  VIEWS: 'workspaces_views',
  TEAMS: 'workspaces_teams',
};

// The workspace designs/views endpoints return the schemas v1beta1 workspace
// page contract, whose owner id is a `core.Id` ([16]byte) value. Go's
// `omitempty` cannot omit a zero-valued array, so an ownerless or published
// resource serializes its owner id as the all-zeros UUID rather than omitting
// it. That string is truthy, so a naive by-id profile lookup fires a
// guaranteed-miss request per resource. Guard against the nil UUID (and an
// absent id) before dispatching the lookup.
const NIL_UUID = '00000000-0000-0000-0000-000000000000';
const isResolvableOwnerId = (id?: string): id is string => Boolean(id) && id !== NIL_UUID;

// KNOWN DEAD CODE - meshery/meshery#21175.
//
// The endpoints injected below whose names @meshery/schemas also defines
// (getWorkspaces, getDesignsOfWorkspace, getViewsOfWorkspace,
// getEnvironmentsOfWorkspace, and the view and team assign/unassign pairs)
// never take effect: `injectEndpoints` without
// `overrideExisting: true` silently discards a colliding name, so the schemas
// definitions serve every call. The `expandInfo` counts, the `expandUser`
// owner resolution and the infinite-scroll merge below therefore do not run.
//
// Restoring them changes what the UI renders and needs verification against a
// running server - and the view/team declarations additionally disagree with
// schemas about the path (`/api/extensions/api/workspaces/...` vs
// `/api/workspaces/...`), so flipping the flag blind would trade a rendering
// bug for a routing one. Tracked in #21175 rather than fixed blind here.
const workspacesApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.WORKSPACES, TAGS.DESIGNS, TAGS.ENVIRONMENTS, TAGS.VIEWS, TAGS.TEAMS],
    endpoints: {
      // The four workspace assign/unassign mutations are the schemas-generated
      // endpoints - the requests are NOT re-declared here. They were, byte for
      // byte (same path, same method, same args), which forked the contract for
      // no gain. `appendInvalidatesTags` adds the local tags below on top of the
      // tags schemas declares, and cannot drop them.
      //
      // Neither local tag has a live provider today, so both currently
      // invalidate nothing: `workspaces_designs` is provided only by the
      // shadowed `getDesignsOfWorkspace` in the dead block below, and
      // `workspaces_environments` is provided by nothing at all - this module
      // no longer declares a `getEnvironmentsOfWorkspace` at all, having
      // dropped the shadowed local one and re-exported the generated query,
      // which provides the schemas-side tags instead.
      // Refetching still happens - the schemas definitions invalidate
      // `Workspace_workspaces`, which is what the effective queries provide.
      // The entries stay registered because they become load-bearing the moment
      // meshery/meshery#21175 restores the shadowed list queries; removing them
      // would just reintroduce the stale-list bug they guard against.
      assignDesignToWorkspace: appendInvalidatesTags('assignDesignToWorkspace', {
        type: TAGS.DESIGNS,
      }),
      unassignDesignFromWorkspace: appendInvalidatesTags('unassignDesignFromWorkspace', {
        type: TAGS.DESIGNS,
      }),
      assignEnvironmentToWorkspace: appendInvalidatesTags('assignEnvironmentToWorkspace', {
        type: TAGS.ENVIRONMENTS,
      }),
      unassignEnvironmentFromWorkspace: appendInvalidatesTags('unassignEnvironmentFromWorkspace', {
        type: TAGS.ENVIRONMENTS,
      }),
    },
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
                  designCount: designs.data?.totalCount ?? designs.data?.total_count ?? 0,
                  environmentCount:
                    environments.data?.totalCount ?? environments.data?.total_count ?? 0,
                  viewCount: views.data?.totalCount ?? views.data?.total_count ?? 0,
                  teamCount: teams.data?.totalCount ?? teams.data?.total_count ?? 0,
                };
              }),
            );

            return _.merge({}, workspaces, { data: { workspaces: modifiedWorkspaces } });
          }
          return workspaces;
        },
        providesTags: () => [{ type: TAGS.WORKSPACES }],
      }),

      getDesignsOfWorkspace: builder.query({
        queryFn: async (queryArgs, { dispatch }, _extraOptions, baseQuery) => {
          const { expandUser, infiniteScroll: _infiniteScroll, ...otherArgs } = queryArgs;
          const params = urlEncodeParams(otherArgs);
          const designs = await baseQuery({
            url: mesheryApiPath(`workspaces/${queryArgs.workspaceId}/designs?${params}`),
            method: 'GET',
          });
          const normalizedDesigns =
            designs.data && !designs.error
              ? { ...designs, data: normalizePaginatedCollectionResponse(designs.data, 'designs') }
              : designs;
          if (expandUser && normalizedDesigns.data && !normalizedDesigns.error) {
            // The `user_id` fallback is live, not legacy: the workspace designs
            // endpoint returns the schemas v1beta1 workspace.MesheryDesignPage,
            // whose design contract spells the owner `user_id` (unlike the
            // v1beta3 design contract's `userId`).
            const withUsersPromises = normalizedDesigns.data.designs.map(async (design) => {
              const ownerId = design.userId ?? design.user_id;
              const user = isResolvableOwnerId(ownerId)
                ? await dispatch(mesheryApi.endpoints.getUserProfileById.initiate({ id: ownerId }))
                : undefined;
              const normalizedUser = normalizeUserProfileSummary(user?.data);
              return {
                ...design,
                firstName: normalizedUser?.firstName || '[deleted]',
                lastName: normalizedUser?.lastName || '',
                avatarUrl: normalizedUser?.avatarUrl || '',
                userId: normalizedUser?.id || design.userId || design.user_id || '',
                email: normalizedUser?.email || '',
              };
            });

            const modifiedDesigns = await Promise.all(withUsersPromises);
            return _.merge({}, normalizedDesigns, { data: { designs: modifiedDesigns } });
          }

          return normalizedDesigns;
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
          const normalizedViews =
            views.data && !views.error
              ? { ...views, data: normalizePaginatedCollectionResponse(views.data, 'views') }
              : views;
          if (expandUser && normalizedViews.data && !normalizedViews.error) {
            const withUsersPromises = normalizedViews.data.views.map(async (view) => {
              const ownerId = view.userId ?? view.user_id;
              const user = isResolvableOwnerId(ownerId)
                ? await dispatch(mesheryApi.endpoints.getUserProfileById.initiate({ id: ownerId }))
                : undefined;
              const normalizedUser = normalizeUserProfileSummary(user?.data);
              return {
                ...view,
                firstName: normalizedUser?.firstName || '[deleted]',
                lastName: normalizedUser?.lastName || '',
                avatarUrl: normalizedUser?.avatarUrl || '',
                userId: normalizedUser?.id || view.userId || view.user_id || '',
                email: normalizedUser?.email || '',
              };
            });
            const modifiedViews = await Promise.all(withUsersPromises);
            return _.merge({}, normalizedViews, { data: { views: modifiedViews } });
          }

          return normalizedViews;
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
    }),
  });

export const {
  useGetWorkspacesQuery,
  useLazyGetWorkspacesQuery,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useAssignViewToWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
  useGetTeamsOfWorkspaceQuery,
  useAssignTeamToWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useGetEventsOfWorkspaceQuery,
} = workspacesApi;

// Re-exported straight from the schemas client: identical path, method and
// args, so there is nothing for this module to add beyond the cache tags
// attached via `enhanceEndpoints` above. Callers keep importing from here so
// the module stays the single place workspace hooks are sourced.
export {
  useGetEnvironmentsOfWorkspaceQuery,
  useAssignEnvironmentToWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useAssignDesignToWorkspaceMutation,
  useUnassignDesignFromWorkspaceMutation,
};

// Workspace CRUD is the schemas-generated create/update/deleteWorkspace. These
// wrappers exist only to translate the callers' argument shape into the
// generated `{ workspaceId, body }` args.
//
// They used to translate into the shape of this module's own local
// createWorkspace/updateWorkspace/deleteWorkspace declarations - but those
// never ran. `injectEndpoints` without `overrideExisting: true` silently
// ignores a name @meshery/schemas already defines, so the generated endpoints
// served every call and read `body` and `workspaceId` as `undefined`: create
// posted an empty body, and update and delete addressed
// `/api/workspaces/undefined`. Deleting the dead declarations and mapping onto
// the generated args is the fix; see __tests__/workspace-mutation-wrappers.
//
// The accepted caller shapes (`workspacePayload`, `body`, or a bare payload;
// `workspaceId` or `id`) are carried over from #21184, which widened them for
// callers that do not use the `workspacePayload` spelling. That PR applied the
// widening to the dead local declarations, so it never reached the wire - it is
// applied here, on the wrappers that feed the endpoints actually serving.
const toWorkspacePayload = (queryArg) => queryArg?.workspacePayload || queryArg?.body || queryArg;

const toWorkspaceId = (queryArg) => queryArg?.workspaceId || queryArg?.id;

const toWorkspaceBody = (queryArg) => {
  const payload = toWorkspacePayload(queryArg);
  return {
    name: payload?.name,
    description: payload?.description,
    organizationId: payload?.organizationId || payload?.organization_id,
  };
};

export const useCreateWorkspaceMutation = () => {
  const [trigger, result] = useSchemasCreateWorkspaceMutation();

  const wrappedTrigger = (queryArg) => trigger({ body: toWorkspaceBody(queryArg) });

  return [wrappedTrigger, result] as const;
};

export const useUpdateWorkspaceMutation = () => {
  const [trigger, result] = useSchemasUpdateWorkspaceMutation();

  const wrappedTrigger = (queryArg) =>
    trigger({
      workspaceId: toWorkspaceId(queryArg),
      body: toWorkspaceBody(queryArg),
    });

  return [wrappedTrigger, result] as const;
};

export const useDeleteWorkspaceMutation = () => {
  const [trigger, result] = useSchemasDeleteWorkspaceMutation();

  const wrappedTrigger = (queryArg) => trigger({ workspaceId: toWorkspaceId(queryArg) });

  return [wrappedTrigger, result] as const;
};
