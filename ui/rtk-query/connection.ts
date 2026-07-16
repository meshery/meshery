import {
  mesheryApi,
  useGetConnectionsQuery as useSchemasGetConnectionsQuery,
  useGetControllerDiagnosticsQuery as useSchemasGetControllerDiagnosticsQuery,
  useGetUserCredentialsQuery as useSchemasGetUserCredentialsQuery,
  useUpdateConnectionMutation as useSchemasUpdateConnectionMutation,
} from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';

// These must match the tag types declared on the shared `mesheryApi`
// (see @meshery/schemas/mesheryApi) — the connections list query
// (`getConnections`) provides `Connection_API_Connections`, so mutations have to
// invalidate that exact tag to make the table refetch. A bare 'connections'
// string isn't a registered tag type and silently invalidates nothing.
const TAGS = {
  CONNECTIONS: 'Connection_API_Connections',
};

const connectionsApi = api.injectEndpoints({
  overrideExisting: true,
  endpoints: (builder) => ({
    verifyAndRegisterConnection: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath('integrations/connections/register'),
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: [TAGS.CONNECTIONS],
    }),

    connectToConnection: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath('integrations/connections/register'),
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: [TAGS.CONNECTIONS],
    }),
    getConnectionDetails: builder.query({
      query: (queryArg) => ({
        url: mesheryApiPath(`integrations/connections/${queryArg.connectionKind}/details`),
        params: { id: queryArg.repoURL },
      }),
    }),
    verifyConnectionURL: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`integrations/connections/${queryArg.connectionKind}/verify`),
        method: 'POST',
        params: { id: queryArg.repoURL },
      }),
    }),
    connectionMetaData: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`integrations/connections/${queryArg.connectionKind}/metadata`),
        method: 'POST',
        body: queryArg.body,
      }),
    }),
    configureConnection: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`integrations/connections/${queryArg.connectionKind}/configure`),
        method: 'POST',
        body: queryArg.body,
      }),
    }),
    cancelConnectionRegister: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`integrations/connections/register`),
        method: 'DELETE',
        body: queryArg.body,
      }),
    }),
    pingKubernetes: builder.query({
      query: (connectionId) => ({
        url: mesheryApiPath(`system/kubernetes/ping`),
        params: { connectionId: connectionId },
        credentials: 'include',
      }),
    }),
    updateConnectionStatus: builder.mutation({
      query: ({ kind, body }) => ({
        url: mesheryApiPath(`integrations/connections/${kind}/status`),
        method: 'PUT',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      }),
      invalidatesTags: () => [{ type: TAGS.CONNECTIONS }],
    }),
    addKubernetesConfig: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`system/kubernetes`),
        method: 'POST',
        body: queryArg.body,
      }),
      invalidatesTags: () => [{ type: TAGS.CONNECTIONS }],
    }),
    // Parses a kubeconfig and returns its contexts (including unreachable ones,
    // flagged) WITHOUT persisting them, so the wizard can let the user pick
    // which to import before any connection is created.
    discoverKubernetesContexts: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`system/kubernetes/contexts`),
        method: 'POST',
        body: queryArg.body,
      }),
    }),
  }),
});

export const {
  useVerifyAndRegisterConnectionMutation,
  useConnectToConnectionMutation,
  useLazyGetConnectionDetailsQuery,
  useVerifyConnectionURLMutation,
  useConnectionMetaDataMutation,
  useConfigureConnectionMutation,
  useCancelConnectionRegisterMutation,
  useAddKubernetesConfigMutation,
  useDiscoverKubernetesContextsMutation,
  useLazyPingKubernetesQuery,
  useUpdateConnectionStatusMutation,
} = connectionsApi;

// Backed by the schemas-generated `getUserCredentials` (GET
// /api/integrations/credentials) rather than a local re-declaration of the same
// endpoint. Callers pass no list args, so every schemas param stays undefined
// and the request is the bare GET this module used to build. queryArg is
// forwarded as-is rather than defaulted to `{}`, because RTK derives the cache
// key from it: `{}` would key separately from a plain schemas call.
export const useGetCredentialsQuery = (queryArg?: undefined, options?: object) =>
  useSchemasGetUserCredentialsQuery(queryArg, options);

// Backed by the schemas-generated `updateConnection` (PUT
// /api/integrations/connections/{connectionId}). The body is still narrowed to
// status + metadata: those are the only fields the server honours here, and
// forwarding a caller's whole object would start sending the rest.
export const useUpdateConnectionByIdMutation = () => {
  const [trigger, result] = useSchemasUpdateConnectionMutation();

  const wrappedTrigger = (queryArg: {
    connectionId: string;
    body?: { status?: unknown; metadata?: unknown };
  }) =>
    trigger({
      connectionId: queryArg.connectionId,
      body: {
        status: queryArg.body?.status,
        metadata: queryArg.body?.metadata,
      },
    });

  return [wrappedTrigger, result] as const;
};

// One-shot controller status pings, backed by the schemas-generated endpoints
// (getOperatorControllerStatus / getMeshsyncControllerStatus /
// getBrokerControllerStatus). Live status is delivered via the SSE stream in
// lib/controllersStatusSubscription.ts. The schemas trigger takes
// `{ connectionId }`; these wrappers accept a bare id so callers stay simple.
const wrapControllerStatusLazyQuery = (endpoint: {
  useLazyQuery: () => readonly [
    (arg: { connectionId: string }, preferCacheValue?: boolean) => unknown,
    ...unknown[],
  ];
}) => {
  return () => {
    const [trigger, ...rest] = endpoint.useLazyQuery();
    const wrappedTrigger = (connectionId: string, preferCacheValue?: boolean) =>
      trigger({ connectionId }, preferCacheValue);
    return [wrappedTrigger, ...rest] as const;
  };
};

// The generated performConnectionAction mutation (POST /connections/{id}/actions)
// returns the updated connection; adding invalidatesTags makes the connections
// list refetch so the new MeshSync mode is reflected immediately.
const connectionActionsApi = api.enhanceEndpoints({
  endpoints: {
    performConnectionAction: {
      invalidatesTags: [TAGS.CONNECTIONS],
    },
  },
});
export const { usePerformConnectionActionMutation } = connectionActionsApi;

export const useLazyGetOperatorStatusQuery = wrapControllerStatusLazyQuery(
  mesheryApi.endpoints.getOperatorControllerStatus,
);
export const useLazyGetMeshsyncStatusQuery = wrapControllerStatusLazyQuery(
  mesheryApi.endpoints.getMeshsyncControllerStatus,
);
export const useLazyGetBrokerStatusQuery = wrapControllerStatusLazyQuery(
  mesheryApi.endpoints.getBrokerControllerStatus,
);

// Per-connection controller diagnostics + remediation, fetched on demand by the
// connection detail view. Skips when no connectionId is available.
export const useGetControllerDiagnosticsQuery = (connectionId, options = {}) =>
  useSchemasGetControllerDiagnosticsQuery({ connectionId }, { skip: !connectionId, ...options });

export const useGetConnectionsQuery = (queryArg, options) =>
  useSchemasGetConnectionsQuery(
    {
      page: queryArg?.page?.toString(),
      // Schemas uses camelCase `pageSize` on the wire; accept either spelling
      // from callers but always forward the canonical one so it reaches the
      // server (which reads `pageSize`).
      pageSize: (queryArg?.pageSize ?? queryArg?.pagesize)?.toString(),
      search: queryArg?.search,
      order: queryArg?.order,
      // Filters are repeated query params (kind=a&kind=b); pass the value(s)
      // straight through — no JSON encoding.
      status: queryArg?.status,
      kind: queryArg?.kind,
      type: queryArg?.type,
      name: queryArg?.name,
    },
    options,
  );

export const useLazyGetConnectionsQuery = () => {
  const [trigger, result, lastPromiseInfo] = mesheryApi.endpoints.getConnections.useLazyQuery();

  const wrappedTrigger = (queryArg, preferCacheValue) =>
    trigger(
      {
        page: queryArg?.page?.toString(),
        pagesize: queryArg?.pagesize?.toString(),
        search: queryArg?.search,
        order: queryArg?.order,
        status: queryArg?.status,
        kind: queryArg?.kind,
        type: queryArg?.type,
        name: queryArg?.name,
      },
      preferCacheValue,
    );

  return [wrappedTrigger, result, lastPromiseInfo] as const;
};
