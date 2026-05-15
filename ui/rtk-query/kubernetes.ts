// kubernetes.ts — RTK Query hooks for the per-context Kubernetes endpoints
// introduced in the GraphQL→SSE migration.
//
// TODO(schemas-canonical): These hooks are defined locally as a known
// tech-debt measure.  Once meshery/schemas publishes OpenAPI specs for these
// five paths and generates @meshery/schemas@>=1.3.0, this file must be
// deleted and consumers must import the generated hooks from
// @meshery/schemas/mesheryApi instead.
//
// Tracked in: meshery/meshery#19424

import { api } from './index';

// ──────────────────────────────────────────────────────────────────────────────
// Response shapes
//
// TODO(schemas-canonical): Once ControllerStatus and OperatorControllerStatus
// are exported from @meshery/schemas/constructs, replace these local types with
// imports and remove this block.
// ──────────────────────────────────────────────────────────────────────────────

/** Mirrors server/internal/graphql/model.OperatorControllerStatus */
export interface OperatorControllerStatus {
  connectionID: string;
  name: string;
  version: string;
  status: string;
  error?: { title: string; description: string } | null;
}

/** Mirrors server/internal/graphql/model.MesheryControllersStatusListItem */
export interface ControllerStatus {
  connectionID: string;
  controller: 'OPERATOR' | 'MESHSYNC' | 'BROKER';
  status: string;
  version: string;
}

// ──────────────────────────────────────────────────────────────────────────────
// Endpoints
// ──────────────────────────────────────────────────────────────────────────────

const kubernetesApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getKubernetesNamespaces: builder.query<{ namespaces: string[] }, { contextID: string }>({
      query: ({ contextID }) => ({
        url: `system/kubernetes/contexts/${contextID}/namespaces`,
        method: 'GET',
      }),
    }),

    getMesheryOperatorStatus: builder.query<ControllerStatus, { contextID: string }>({
      query: ({ contextID }) => ({
        url: `system/kubernetes/contexts/${contextID}/operator/status`,
        method: 'GET',
      }),
    }),

    getMeshsyncStatus: builder.query<OperatorControllerStatus, { contextID: string }>({
      query: ({ contextID }) => ({
        url: `system/kubernetes/contexts/${contextID}/meshsync/status`,
        method: 'GET',
      }),
    }),

    getNatsStatus: builder.query<OperatorControllerStatus, { contextID: string }>({
      query: ({ contextID }) => ({
        url: `system/kubernetes/contexts/${contextID}/nats/status`,
        method: 'GET',
      }),
    }),

    resyncCluster: builder.mutation<
      { status: string },
      { contextID: string; clearDb?: string; hardReset?: string; reSync?: string }
    >({
      query: ({ contextID, clearDb = 'false', hardReset = 'false', reSync = 'false' }) => ({
        url: `system/kubernetes/contexts/${contextID}/resync`,
        method: 'POST',
        body: { clearDb, hardReset, reSync },
      }),
    }),
  }),
});

export const {
  useGetKubernetesNamespacesQuery,
  useGetMesheryOperatorStatusQuery,
  useGetMeshsyncStatusQuery,
  useGetNatsStatusQuery,
  useResyncClusterMutation,
} = kubernetesApi;
