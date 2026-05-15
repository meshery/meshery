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
  TNormalizedItem = TItem,
>(
  response: PaginatedCollectionResponse<TCollectionKey, TItem> | undefined,
  collectionKey: TCollectionKey,
  itemMapper?: (item: TItem) => TNormalizedItem,
) => {
  if (!response || typeof response !== 'object') {
    return response;
  }

  const items = Array.isArray(response[collectionKey]) ? response[collectionKey] : [];

  return {
    ...response,
    pageSize: response.pageSize ?? response.page_size,
    totalCount: response.totalCount ?? response.total_count,
    [collectionKey]: itemMapper ? items.map(itemMapper) : items,
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
