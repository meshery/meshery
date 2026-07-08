import { api, mesheryApiPath } from './index';

// Layered Meshery Operator / MeshSync / Broker configuration endpoints.
// Server-wide defaults live under /api/system/controllers/config; the
// per-connection override rides the connection's metadata and is exposed at
// /api/integrations/connections/{id}/controllers/config. The wire contract is
// the controllers_config construct in meshery/schemas (v1alpha1).
const TAGS = {
  CONTROLLERS_CONFIG: 'controllers_config',
} as const;

export type ControllersConfigDoc = {
  schemaVersion?: string;
  operator?: {
    deploymentMode?: 'operator' | 'embedded';
    version?: string;
  };
  meshsync?: {
    version?: string;
    replicas?: number;
    watchList?: {
      whitelist?: { resource: string; events?: string[] }[];
      blacklist?: string[];
    };
    outputNamespaces?: string[];
    outputResources?: string[];
    redactSecrets?: boolean;
    brokerContentDedup?: boolean;
    debugLogging?: boolean;
  };
  broker?: {
    version?: string;
    replicas?: number;
    service?: {
      type?: 'ClusterIP' | 'NodePort' | 'LoadBalancer';
      annotations?: Record<string, string>;
      loadBalancerClass?: string;
      loadBalancerSourceRanges?: string[];
      externalEndpointOverride?: string;
    };
  };
};

export type ConnectionControllersConfigDoc = {
  override?: ControllersConfigDoc | null;
  default?: ControllersConfigDoc | null;
  effective: ControllersConfigDoc;
};

const controllersConfigApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.CONTROLLERS_CONFIG],
  })
  .injectEndpoints({
    overrideExisting: true,
    endpoints: (builder) => ({
      getControllersDefaultConfig: builder.query<ControllersConfigDoc, void>({
        query: () => ({
          url: mesheryApiPath('system/controllers/config'),
          method: 'GET',
        }),
        providesTags: [TAGS.CONTROLLERS_CONFIG],
      }),
      updateControllersDefaultConfig: builder.mutation<
        ControllersConfigDoc,
        { body: ControllersConfigDoc }
      >({
        query: (queryArg) => ({
          url: mesheryApiPath('system/controllers/config'),
          method: 'PUT',
          body: queryArg.body,
        }),
        invalidatesTags: [TAGS.CONTROLLERS_CONFIG],
      }),
      getConnectionControllersConfig: builder.query<
        ConnectionControllersConfigDoc,
        { connectionId: string }
      >({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `integrations/connections/${queryArg.connectionId}/controllers/config`,
          ),
          method: 'GET',
        }),
        providesTags: (_result, _error, arg) => [
          { type: TAGS.CONTROLLERS_CONFIG, id: arg.connectionId },
          TAGS.CONTROLLERS_CONFIG,
        ],
      }),
      updateConnectionControllersConfig: builder.mutation<
        ConnectionControllersConfigDoc,
        { connectionId: string; body: ControllersConfigDoc }
      >({
        query: (queryArg) => ({
          url: mesheryApiPath(
            `integrations/connections/${queryArg.connectionId}/controllers/config`,
          ),
          method: 'PUT',
          body: queryArg.body,
        }),
        invalidatesTags: (_result, _error, arg) => [
          { type: TAGS.CONTROLLERS_CONFIG, id: arg.connectionId },
          TAGS.CONTROLLERS_CONFIG,
        ],
      }),
    }),
  });

export const {
  useGetControllersDefaultConfigQuery,
  useUpdateControllersDefaultConfigMutation,
  useGetConnectionControllersConfigQuery,
  useUpdateConnectionControllersConfigMutation,
} = controllersConfigApi;
