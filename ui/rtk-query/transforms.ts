type PaginatedCollectionResponse<TCollectionKey extends string, TItem> = Partial<
  Record<TCollectionKey, TItem[]>
> & {
  page?: number;
  pageSize?: number;
  page_size?: number;
  totalCount?: number;
  total_count?: number;
};

type ProviderCapabilitiesResponse = {
  providerName?: string;
  provider_name?: string;
  providerType?: string;
  provider_type?: string;
  providerUrl?: string;
  provider_url?: string;
  providerDescription?: string[];
  provider_description?: string[];
  capabilities?: unknown[];
  extensions?: Record<string, unknown>;
  restrictedAccess?: Record<string, unknown>;
  [key: string]: unknown;
};

type KubernetesContextResponse = {
  totalCount?: number;
  total_count?: number;
  contexts?: Array<{
    createdBy?: string;
    created_by?: string;
    mesheryInstanceId?: string;
    meshery_instance_id?: string;
    kubernetesServerId?: string;
    kubernetes_server_id?: string;
    deploymentType?: string;
    deployment_type?: string;
    updatedAt?: string;
    updated_at?: string;
    createdAt?: string;
    created_at?: string;
    connectionId?: string;
    connection_id?: string;
    isCurrentContext?: boolean;
    is_current_context?: boolean;
    [key: string]: unknown;
  }>;
  [key: string]: unknown;
};

export const normalizePaginatedCollectionResponse = <
  TCollectionKey extends string,
  TItem = unknown,
>(
  response: PaginatedCollectionResponse<TCollectionKey, TItem> | undefined,
  collectionKey: TCollectionKey,
) => {
  if (!response || typeof response !== 'object') {
    return response;
  }

  return {
    ...response,
    pageSize: response.pageSize ?? response.page_size,
    totalCount: response.totalCount ?? response.total_count,
    [collectionKey]: Array.isArray(response[collectionKey]) ? response[collectionKey] : [],
  };
};

export const normalizeProviderCapabilities = (response?: ProviderCapabilitiesResponse) => {
  if (!response || typeof response !== 'object') {
    return undefined;
  }

  return {
    ...response,
    providerName: response.providerName ?? response.provider_name,
    providerType: response.providerType ?? response.provider_type,
    providerUrl: response.providerUrl ?? response.provider_url,
    providerDescription: response.providerDescription ?? response.provider_description,
    capabilities: Array.isArray(response.capabilities) ? response.capabilities : [],
  };
};

export const normalizeKubernetesContextsResponse = (response?: KubernetesContextResponse) => {
  if (!response || typeof response !== 'object') {
    return undefined;
  }

  return {
    ...response,
    totalCount: response.totalCount ?? response.total_count,
    contexts: Array.isArray(response.contexts)
      ? response.contexts.map((context) => ({
          ...context,
          createdBy: context.createdBy ?? context.created_by,
          mesheryInstanceId: context.mesheryInstanceId ?? context.meshery_instance_id,
          kubernetesServerId: context.kubernetesServerId ?? context.kubernetes_server_id,
          deploymentType: context.deploymentType ?? context.deployment_type,
          updatedAt: context.updatedAt ?? context.updated_at,
          createdAt: context.createdAt ?? context.created_at,
          connectionId: context.connectionId ?? context.connection_id,
          isCurrentContext: context.isCurrentContext ?? context.is_current_context,
        }))
      : [],
  };
};

type KubernetesConnection = {
  id?: string;
  status?: string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  metadata?: Record<string, unknown> | null;
  [key: string]: unknown;
};

/**
 * Maps kubernetes connections (from the connections REST API, kind=kubernetes)
 * into the K8sContext-list shape the app consumes. This mirrors the backend
 * `K8sContextFromConnection`: the connection's metadata carries the context
 * fields (id, name, server, kubernetesServerId, ...), and the connection id
 * becomes `connectionId`.
 *
 * Everything is connection-driven now — this replaces the deprecated
 * subscribeK8sContext GraphQL subscription and the k8s-contexts API as the
 * source of the k8sConfig list.
 */
export const connectionsToK8sContexts = (connections?: KubernetesConnection[]) => {
  const contexts = (Array.isArray(connections) ? connections : [])
    .filter((connection) => connection && connection.metadata)
    .map((connection) => ({
      ...(connection.metadata as Record<string, unknown>),
      connectionId: connection.id,
      // The connection's lifecycle status (connected/registered/discovered/...),
      // surfaced so consumers can gate on it (e.g. only poll controller status
      // for connected connections).
      connectionStatus: connection.status,
      createdAt: connection.created_at ?? connection.createdAt,
      updatedAt: connection.updated_at ?? connection.updatedAt,
    }));

  return normalizeKubernetesContextsResponse({ contexts, totalCount: contexts.length });
};

type LoggedInUserResponse = {
  id?: string;
  userId?: string;
  [key: string]: unknown;
};

// normalizeLoggedInUser adapts the current-user response for the v1beta3 Cloud
// account-consolidation cutover: schemas v1beta3 dropped the `userId` field (the
// canonical identifier is the `id` UUID). Backfill `userId` from `id` so
// ownership checks that compare `user.userId` against a resource's owner UUID
// keep working across every consumer of getLoggedInUser. Non-destructive: all
// other fields pass through unchanged.
export const normalizeLoggedInUser = (response?: LoggedInUserResponse) => {
  if (!response || typeof response !== 'object') {
    return undefined;
  }

  return {
    ...response,
    userId: response.userId ?? response.id,
  };
};
