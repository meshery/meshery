import { api, mesheryApiPath } from './index';

const TAGS = {
  KUBERNETES: 'kubernetes',
};

interface Controller {
  name: string;
  version: string;
  status: string;
}

interface OperatorStatus {
  status: string;
  controller?: string;
  connectionID?: string;
  version?: string;
  error?: unknown;
  controllers?: Controller[];
}

const kubernetesApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.KUBERNETES],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getKubernetesNamespaces: builder.query<
        { namespaces: Array<{ namespace: string }> },
        { clusterIds?: string | string[] }
      >({
        query: (queryArg) => ({
          url: mesheryApiPath('system/kubernetes/namespaces'),
          params: {
            clusterIds: queryArg?.clusterIds,
          },
          method: 'GET',
          credentials: 'include',
        }),
        // Legacy GraphQL aliased getAvailableNamespaces as `namespaces`.
        // Coerce whatever the REST endpoint returns into that shape.
        transformResponse: (
          response:
            | { namespaces?: Array<{ namespace: string }> }
            | Array<{ namespace: string } | string>
            | undefined,
        ) => {
          if (Array.isArray(response)) {
            const normalized = response.map((entry) =>
              typeof entry === 'string' ? { namespace: entry } : entry,
            );
            return { namespaces: normalized };
          }
          if (response && Array.isArray(response.namespaces)) {
            return { namespaces: response.namespaces };
          }
          return { namespaces: [] };
        },
        providesTags: [TAGS.KUBERNETES],
      }),

      getMesheryOperatorStatus: builder.query<
        { operator: OperatorStatus | null },
        { connectionID: string }
      >({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `system/kubernetes/connections/${encodeURIComponent(queryArg.connectionID)}/operator/status`,
          ),
          method: 'GET',
          credentials: 'include',
        }),
        // Legacy GraphQL aliased getOperatorStatus as `operator`. The REST
        // endpoint may return the status object directly or already wrapped.
        transformResponse: (
          response: { operator?: OperatorStatus } | OperatorStatus | undefined,
        ) => {
          if (!response) return { operator: null };
          if (
            typeof response === 'object' &&
            'operator' in response &&
            (response as { operator?: OperatorStatus }).operator !== undefined
          ) {
            return { operator: (response as { operator: OperatorStatus | null }).operator };
          }
          return { operator: response as OperatorStatus };
        },
        providesTags: (_result, _error, arg) => [
          { type: TAGS.KUBERNETES, id: `operator-${arg?.connectionID}` },
        ],
      }),

      getMeshsyncStatus: builder.query<{ controller: Controller }, { connectionID: string }>({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `system/kubernetes/connections/${encodeURIComponent(queryArg.connectionID)}/meshsync/status`,
          ),
          method: 'GET',
          credentials: 'include',
        }),
        // Legacy GraphQL aliased getMeshsyncStatus as `controller`.
        transformResponse: (response: { controller?: Controller } | Controller | undefined) => {
          if (!response) return { controller: { name: '', version: '', status: '' } };
          if (
            typeof response === 'object' &&
            'controller' in response &&
            (response as { controller?: Controller }).controller !== undefined
          ) {
            return { controller: (response as { controller: Controller }).controller };
          }
          return { controller: response as Controller };
        },
        providesTags: (_result, _error, arg) => [
          { type: TAGS.KUBERNETES, id: `meshsync-${arg?.connectionID}` },
        ],
      }),

      resyncCluster: builder.mutation<
        { status?: string },
        {
          contextID: string;
          reSync?: boolean;
          clearDb?: boolean;
          hardReset?: boolean;
        }
      >({
        // Replaces the legacy resyncCluster GraphQL operation. Hits the
        // dedicated REST endpoint added in #19397 (server/rest-and-sse-
        // endpoints) instead of /api/system/database/reset — those are
        // separate operations: this one resyncs MeshSync for a single
        // cluster, the reset endpoint drops the entire database.
        query: ({ contextID, ...body }) => ({
          url: mesheryApiPath(`system/kubernetes/contexts/${encodeURIComponent(contextID)}/resync`),
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: {
            reSync: body.reSync ?? false,
            clearDb: body.clearDb ?? false,
            hardReset: body.hardReset ?? false,
          },
        }),
        invalidatesTags: [TAGS.KUBERNETES],
      }),

      getNatsStatus: builder.query<{ controller: Controller }, { connectionID: string }>({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `system/kubernetes/connections/${encodeURIComponent(queryArg.connectionID)}/nats/status`,
          ),
          method: 'GET',
          credentials: 'include',
        }),
        // Legacy GraphQL aliased getNatsStatus as `controller`.
        transformResponse: (response: { controller?: Controller } | Controller | undefined) => {
          if (!response) return { controller: { name: '', version: '', status: '' } };
          if (
            typeof response === 'object' &&
            'controller' in response &&
            (response as { controller?: Controller }).controller !== undefined
          ) {
            return { controller: (response as { controller: Controller }).controller };
          }
          return { controller: response as Controller };
        },
        providesTags: (_result, _error, arg) => [
          { type: TAGS.KUBERNETES, id: `nats-${arg?.connectionID}` },
        ],
      }),
    }),
  });

export const {
  useGetKubernetesNamespacesQuery,
  useLazyGetKubernetesNamespacesQuery,
  useGetMesheryOperatorStatusQuery,
  useLazyGetMesheryOperatorStatusQuery,
  useGetMeshsyncStatusQuery,
  useLazyGetMeshsyncStatusQuery,
  useGetNatsStatusQuery,
  useLazyGetNatsStatusQuery,
  useResyncClusterMutation,
} = kubernetesApi;
