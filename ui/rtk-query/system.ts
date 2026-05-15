import { api, mesheryApiPath } from './index';
import { normalizeKubernetesContextsResponse } from './transforms';

const TAGS = {
  SYSTEM: 'system',
  ADAPTERS: 'adapters',
  SYNC: 'sync',
};

const systemApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getDatabaseSummary: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath(`system/database`),
        params: {
          page: queryArg.page,
          pagesize: queryArg.pagesize,
          search: queryArg.search,
          order: queryArg.order,
        },
        method: 'GET',
      }),
      providesTags: () => [{ type: TAGS.SYSTEM }],
    }),
    getAdapters: builder.query({
      query: () => ({
        url: mesheryApiPath('system/adapters'),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    getAvailableAdapters: builder.query({
      query: () => ({
        url: mesheryApiPath('system/availableAdapters'),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.ADAPTERS],
    }),

    pingAdapter: builder.query({
      query: (adapterLoc) => ({
        url: mesheryApiPath(`system/adapters`),
        params: { adapter: adapterLoc },
        credentials: 'include',
      }),
      providesTags: (result, error, adapterLoc) => [{ type: TAGS.ADAPTERS, id: adapterLoc }],
    }),
    getSystemSync: builder.query({
      query: () => ({
        url: mesheryApiPath('system/sync'),
        method: 'GET',
        credentials: 'include',
      }),
      providesTags: [TAGS.SYNC],
    }),

    getKubernetesContexts: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath('system/kubernetes/contexts'),
        params: {
          pagesize: queryArg?.pagesize || 10,
          search: queryArg?.search || '',
        },
        method: 'GET',
      }),
      transformResponse: normalizeKubernetesContextsResponse,
      providesTags: [TAGS.SYSTEM],
    }),

    adapterOperation: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(queryArg.url || 'system/adapter/operation'),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: queryArg.body,
      }),
    }),

    getSmiResults: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath('smi/results'),
        params: {
          page: queryArg?.page,
          pagesize: queryArg?.pagesize,
          search: queryArg?.search,
          order: queryArg?.order,
        },
        method: 'GET',
        credentials: 'include',
      }),
    }),

    manageAdapter: builder.mutation({
      query: (queryArg) => {
        if (queryArg.method === 'DELETE') {
          return {
            url: mesheryApiPath(`system/adapter/manage`),
            method: 'DELETE',
            credentials: 'include',
            params: { adapter: queryArg.adapter },
          };
        }

        return {
          url: mesheryApiPath('system/adapter/manage'),
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
          body: `meshLocationURL=${encodeURIComponent(queryArg.meshLocationURL)}`,
        };
      },
      invalidatesTags: [TAGS.ADAPTERS],
    }),

    getMeshAddons: builder.query<
      { addonsState: Array<{ name: string; owner: string; endpoint?: string }> },
      { meshType?: string; clusterId?: string | string[] }
    >({
      query: (queryArg) => ({
        url: mesheryApiPath('system/meshes/addons'),
        params: {
          meshType: queryArg?.meshType,
          clusterId: queryArg?.clusterId,
        },
        method: 'GET',
        credentials: 'include',
      }),
      // The legacy GraphQL aliased getAvailableAddons as `addonsState`.
      // Preserve that consumer contract regardless of whether the REST
      // endpoint returns the array directly or wraps it.
      transformResponse: (
        response:
          | { addonsState?: Array<{ name: string; owner: string; endpoint?: string }> }
          | { addons?: Array<{ name: string; owner: string; endpoint?: string }> }
          | Array<{ name: string; owner: string; endpoint?: string }>
          | undefined,
      ) => {
        if (Array.isArray(response)) {
          return { addonsState: response };
        }
        if (response && Array.isArray((response as { addonsState?: unknown[] }).addonsState)) {
          return {
            addonsState: (response as { addonsState: typeof response.addonsState }).addonsState,
          };
        }
        if (response && Array.isArray((response as { addons?: unknown[] }).addons)) {
          return { addonsState: (response as { addons: typeof response.addons }).addons };
        }
        return { addonsState: [] };
      },
      providesTags: [TAGS.SYSTEM],
    }),

    getControlPlanes: builder.query<
      {
        controlPlanesState: Array<{
          name: string;
          members: Array<{ name: string; version: string; component: string; namespace: string }>;
        }>;
      },
      { type?: string; clusterId?: string | string[] }
    >({
      query: (queryArg) => ({
        url: mesheryApiPath('system/meshes/control-planes'),
        params: {
          type: queryArg?.type,
          clusterId: queryArg?.clusterId,
        },
        method: 'GET',
        credentials: 'include',
      }),
      // The legacy GraphQL aliased getControlPlanes as `controlPlanesState`.
      // Preserve that consumer contract.
      transformResponse: (
        response:
          | {
              controlPlanesState?: Array<{
                name: string;
                members: Array<{
                  name: string;
                  version: string;
                  component: string;
                  namespace: string;
                }>;
              }>;
            }
          | {
              control_planes?: Array<{
                name: string;
                members: Array<{
                  name: string;
                  version: string;
                  component: string;
                  namespace: string;
                }>;
              }>;
            }
          | Array<{
              name: string;
              members: Array<{
                name: string;
                version: string;
                component: string;
                namespace: string;
              }>;
            }>
          | undefined,
      ) => {
        if (Array.isArray(response)) {
          return { controlPlanesState: response };
        }
        if (
          response &&
          Array.isArray((response as { controlPlanesState?: unknown[] }).controlPlanesState)
        ) {
          return {
            controlPlanesState: (
              response as { controlPlanesState: typeof response.controlPlanesState }
            ).controlPlanesState,
          };
        }
        if (
          response &&
          Array.isArray((response as { control_planes?: unknown[] }).control_planes)
        ) {
          return {
            controlPlanesState: (response as { control_planes: typeof response.control_planes })
              .control_planes,
          };
        }
        return { controlPlanesState: [] };
      },
      providesTags: [TAGS.SYSTEM],
    }),

    resetDatabase: builder.mutation<
      { message?: string },
      { k8scontextID?: string; clearDB?: string; ReSync?: string; hardReset?: string }
    >({
      query: (queryArg) => ({
        url: mesheryApiPath('system/database/reset'),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: {
          k8scontextID: queryArg?.k8scontextID || '',
          clearDB: queryArg?.clearDB,
          ReSync: queryArg?.ReSync,
          hardReset: queryArg?.hardReset,
        },
      }),
      invalidatesTags: [TAGS.SYSTEM],
    }),

    changeOperatorStatus: builder.mutation<
      unknown,
      { contextID: string; targetStatus: 'ENABLED' | 'DISABLED' }
    >({
      query: (queryArg) => ({
        url: mesheryApiPath(
          `system/kubernetes/contexts/${encodeURIComponent(queryArg.contextID)}/operator`,
        ),
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: { targetStatus: queryArg.targetStatus },
      }),
      invalidatesTags: [TAGS.SYSTEM],
    }),
  }),
});

export const {
  useGetDatabaseSummaryQuery,
  useGetAdaptersQuery,
  useGetAvailableAdaptersQuery,
  useLazyPingAdapterQuery,
  useManageAdapterMutation,
  useGetSystemSyncQuery,
  useLazyGetSystemSyncQuery,
  useGetKubernetesContextsQuery,
  useLazyGetKubernetesContextsQuery,
  useAdapterOperationMutation,
  useLazyGetSmiResultsQuery,
  useGetMeshAddonsQuery,
  useLazyGetMeshAddonsQuery,
  useGetControlPlanesQuery,
  useLazyGetControlPlanesQuery,
  useResetDatabaseMutation,
  useChangeOperatorStatusMutation,
} = systemApi;
