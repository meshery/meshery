import { CONNECTION_STATES, MESHSYNC_DEPLOYMENT_TYPE } from './Enum';

// sessionStorage key for the user's include/exclude selection of Kubernetes
// contexts. Session-scoped on purpose: the selection is a view preference, not
// durable configuration, so it survives navigation and reloads within a
// browser session without any server round-trip.
const SELECTED_K8S_CONTEXTS_STORAGE_KEY = 'selectedK8sContexts';

/**
 * Reads the persisted context selection for this browser session.
 *
 * @returns {Array.<string>|null} The persisted selection (may be `[]` when the
 *   user deselected every context, or `['all']`), or `null` when nothing valid
 *   is persisted and the caller should fall back to the default selection.
 */
export function loadSelectedK8sContexts(): string[] | null {
  if (typeof window === 'undefined') return null;
  try {
    const raw = window.sessionStorage.getItem(SELECTED_K8S_CONTEXTS_STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) && parsed.every((id) => typeof id === 'string') ? parsed : null;
  } catch {
    return null;
  }
}

/**
 * Persists the context selection for this browser session.
 *
 * @param {Array.<string>} selectedK8sContexts Context ids, or `['all']`.
 */
export function persistSelectedK8sContexts(selectedK8sContexts: string[]): void {
  if (typeof window === 'undefined' || !Array.isArray(selectedK8sContexts)) return;
  try {
    window.sessionStorage.setItem(
      SELECTED_K8S_CONTEXTS_STORAGE_KEY,
      JSON.stringify(selectedK8sContexts),
    );
  } catch {
    // Quota exceeded or private browsing — degrade silently; the selection
    // simply won't survive navigation.
  }
}

/**
 * A function to be used by the requests sent for the
 * operations based on multi-context support
 *
 * @param {string} url The request URL
 * @param {Array.<string>} ctx The context Array
 * @returns {string} The final query-parametrised URL
 */
export function ctxUrl(url, ctx) {
  if (ctx?.length) {
    const contextQuery = ctx.map((context) => `contexts=${context}`).join('&');
    return `${url}?${contextQuery}`;
  }
  return url;
}

/**
 * The function takes in all the context and returns
 *  their respective cluster IDs associated to them
 *
 * @param {Array.<string>} selectedContexts
 * @param {Array.<string>} k8sconfig
 * @returns
 */
export const getK8sClusterIdsFromCtxId = (selectedContexts, k8sconfig) => {
  if (!selectedContexts || !k8sconfig || selectedContexts.length === 0) {
    return [];
  }

  if (selectedContexts.includes('all')) {
    // Drop configs without a resolved kubernetesServerId (e.g. a connection
    // registered while its cluster was unreachable, so no server ID was assigned
    // yet). Sending an undefined/empty cluster id contributes a junk `clusterId`
    // query param that matches no MeshSync rows, and matches the truthy filter the
    // explicit-selection branch below already applies.
    return k8sconfig.map((cfg) => cfg?.kubernetesServerId).filter(Boolean);
  }
  const clusterIds = [];
  selectedContexts.forEach((context) => {
    const clusterId = k8sconfig.find((cfg) => cfg.id === context)?.kubernetesServerId;
    if (clusterId) {
      clusterIds.push(clusterId);
    }
  });

  return clusterIds;
};

/**
 *
 * @param {Array.<string>} selectedK8sContexts
 * @param {Array.<string>} k8sConfig
 * @returns {string} The context ID
 */
export function getFirstCtxIdFromSelectedCtxIds(selectedK8sContexts, k8sConfig) {
  if (!selectedK8sContexts?.length) {
    return '';
  }

  if (selectedK8sContexts?.includes('all')) {
    return k8sConfig[0]?.id;
  }

  return selectedK8sContexts[0];
}

/**
 * Get the k8sConfigIds of K8sconfig
 * @param {Array.<Object>} k8sConfig
 * @returns
 */
export function getK8sConfigIdsFromK8sConfig(k8sConfig) {
  if (!k8sConfig || !k8sConfig.length) {
    return [];
  }

  return k8sConfig.map((cfg) => cfg.id);
}

export const getK8sClusterNamesFromCtxId = (selectedContexts, k8sconfig) => {
  if (selectedContexts.length === 0) {
    return [];
  }

  if (selectedContexts.includes('all')) {
    return ['all'];
  }

  const clusterNames = [];

  selectedContexts.forEach((context) => {
    const name = k8sconfig.find((cfg) => cfg.id === context)?.name;
    if (name) {
      clusterNames.push(name);
    }
  });

  return clusterNames;
};

/**
 *
 * @param {string} clusterId Kubernetes Cluster ID
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {string} Kubernetes context
 */
export function getK8sContextFromClusterId(clusterId, k8sConfig) {
  const cluster = k8sConfig.find((cfg) => cfg.kubernetesServerId === clusterId);
  if (!cluster) {
    return {};
  }
  return cluster;
}

/**
 *
 * @param {string} clusterId Kubernetes Cluster ID
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {string} Kubernetes cluster name
 */
export function getClusterNameFromClusterId(clusterId, k8sConfig) {
  const cluster = k8sConfig.find((cfg) => cfg.kubernetesServerId === clusterId);
  if (!cluster) {
    return '';
  }
  return cluster.name;
}

/**
 *
 * @param {string} connectionId Kubernetes Connection ID
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {string} Kubernetes cluster name
 */
export function getClusterNameFromConnectionId(connId, k8sConfig) {
  const cluster = k8sConfig.find((cfg) => cfg.connectionId === connId);
  if (!cluster) {
    return '';
  }
  return cluster.name;
}

/**
 *
 * @param {string} clusterId Kubernetes Cluster ID
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {string} Kubernetes connection ID
 */
export function getConnectionIdFromClusterId(clusterId, k8sConfig) {
  const cluster = k8sConfig.find((cfg) => cfg.kubernetesServerId === clusterId);
  if (!cluster) {
    return '';
  }
  return cluster.connectionId;
}

/**
 *
 * @param {string} ctxId Kubernetes context ID
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {string} Kubernetes cluster name
 */
export function getClusterNameFromCtxId(ctxId, k8sConfig) {
  const cluster = k8sConfig.find((cfg) => cfg.id === ctxId);
  if (!cluster) {
    return '';
  }
  return cluster.name;
}

/**
 *
 * @param {Array<Object>} contextIDs Kubernetes context ids
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {Array<string>} array of connection ID for given kubernetes contexts
 */
export function getConnectionIDsFromContextIds(contexts, k8sConfig) {
  const filteredK8sConnfigs = k8sConfig.filter((config) =>
    contexts.some((context) => context == config.id),
  );
  return filteredK8sConnfigs.map((config) => config.connectionId);
}

// Reads the MeshSync deployment mode from a k8sConfig/context entry, tolerating
// both camelCase and snake_case metadata keys.
function getMeshsyncDeploymentMode(config) {
  return config?.meshsyncDeploymentMode ?? config?.meshsync_deployment_mode;
}

// Reads the connection lifecycle status from a k8sConfig/context entry,
// tolerating the connectionStatus (mapped from the connection) or a status key.
function getConnectionStatus(config) {
  return config?.connectionStatus ?? config?.status;
}

/**
 * Returns the connection IDs whose controller status is worth streaming, i.e.
 * connections that are BOTH:
 *   - in `operator` MeshSync deployment mode — the operator/broker/meshsync
 *     controllers only exist in-cluster in operator mode; embedded connections
 *     run MeshSync in-process and have no such resources, so polling their
 *     controller status just 404s every tick; and
 *   - in the `connected` state — a discovered/registered/disconnected
 *     connection has no live controllers running, so there is nothing to poll.
 *
 * Scoping the controller-status stream this way avoids pointless per-tick work
 * and log noise. A missing/embedded mode is treated as non-operator (matching
 * the embedded default), and a missing status is treated as not-connected.
 *
 * @param {Array<Object>} contexts Kubernetes context ids
 * @param {Array<Object>} k8sConfig Kubernetes config
 * @returns {Array<string>} connection IDs eligible for controller-status polling
 */
export function getControllerPollConnectionIDsFromContextIds(contexts, k8sConfig) {
  if (!Array.isArray(contexts) || !Array.isArray(k8sConfig)) {
    return [];
  }
  return k8sConfig
    .filter((config) => contexts.some((context) => context == config.id))
    .filter((config) => getMeshsyncDeploymentMode(config) === MESHSYNC_DEPLOYMENT_TYPE.OPERATOR)
    .filter((config) => getConnectionStatus(config) === CONNECTION_STATES.CONNECTED)
    .map((config) => config.connectionId);
}
