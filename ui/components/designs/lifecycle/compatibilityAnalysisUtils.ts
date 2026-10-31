/**
 * Utilities for normalizing and analyzing cross-cluster design validation responses
 * into a consolidated compatibility matrix.
 */

export type CompatibilityStatus = 'success' | 'warning' | 'failed' | 'incomplete';

export interface ClusterInfo {
  id: string;
  name: string;
  version: string;
  status: string;
}

export interface ResourceValidationDetail {
  status: CompatibilityStatus;
  errorMessage?: string;
  errorType?: string;
  fieldPath?: string;
  possibleFix?: string;
  dependencyImpact?: string;
}

export interface ResourceRow {
  resourceName: string;
  kind: string;
  apiVersion: string;
  component?: any;
  clusterResults: Record<string, ResourceValidationDetail>;
}

export interface ValidationMetrics {
  totalClusters: number;
  totalResources: number;
  succeededCount: number;
  warningCount: number;
  failedCount: number;
  incompleteCount?: number;
}

export interface ClusterSummaryItem {
  clusterId: string;
  clusterName: string;
  version: string;
  passed: number;
  warnings: number;
  failed: number;
  incomplete?: number;
  total: number;
}

export interface CompatibilityAnalysisData {
  clusters: ClusterInfo[];
  resources: ResourceRow[];
  metrics: ValidationMetrics;
  clusterSummaries: ClusterSummaryItem[];
  nextSteps: string[];
}

/**
 * Generates heuristic possible fix advice based on the error message and resource kind.
 */
export function generatePossibleFix(
  errorMessage: string = '',
  kind: string = 'Resource',
  clusterName: string = 'the target cluster',
): string {
  const lowerMsg = errorMessage.toLowerCase();

  if (
    lowerMsg.includes('no matches for kind') ||
    lowerMsg.includes('crd') ||
    lowerMsg.includes('the server could not find the requested resource')
  ) {
    return `The required Custom Resource Definition (CRD) for ${kind} is not installed on ${clusterName}. Install the CRD on this cluster or use a supported resource alternative.`;
  }

  if (
    lowerMsg.includes('apiversion') ||
    lowerMsg.includes('no matches for version') ||
    lowerMsg.includes('deprecated')
  ) {
    return `The specified apiVersion is unavailable or deprecated on Kubernetes ${clusterName}. Verify the cluster Kubernetes version and update the apiVersion in your design.`;
  }

  if (
    lowerMsg.includes('namespace') ||
    lowerMsg.includes('namespaces') ||
    lowerMsg.includes('forbidden')
  ) {
    return `The target namespace may be missing or unauthorized on ${clusterName}. Verify namespace existence and cluster connection permissions.`;
  }

  if (lowerMsg.includes('quota') || lowerMsg.includes('exceeded') || lowerMsg.includes('limit')) {
    return `Resource limits or quota exceeded on ${clusterName}. Reduce resource requests or increase cluster quota limits.`;
  }

  if (lowerMsg.includes('immutable') || lowerMsg.includes('forbidden: field is immutable')) {
    return `An immutable field change was attempted. Recreate the resource or retain the original immutable field values.`;
  }

  return `Verify that the configuration of ${kind} adheres to the target schema for ${clusterName} and that required cluster operators are healthy.`;
}

/**
 * Generates dependency impact explanation for dependent workloads.
 */
export function generateDependencyImpact(
  resourceName: string,
  kind: string,
  allResourceNames: string[] = [],
): string {
  const lowerKind = kind.toLowerCase();
  const lowerName = resourceName.toLowerCase();

  const isDatabaseOrStore =
    lowerKind.includes('statefulset') ||
    lowerKind.includes('postgres') ||
    lowerKind.includes('mysql') ||
    lowerKind.includes('redis') ||
    lowerKind.includes('mongo') ||
    lowerName.includes('db') ||
    lowerName.includes('postgres') ||
    lowerName.includes('redis');

  if (isDatabaseOrStore) {
    const dependentWorkloads = allResourceNames.filter(
      (name) => name !== resourceName && !name.toLowerCase().includes('ingress'),
    );
    const dependentsStr =
      dependentWorkloads.length > 0 ? dependentWorkloads.join(', ') : 'backend services';
    return `This failure may affect dependent resources (e.g., ${dependentsStr}) that require persistent data or database connectivity to initialize.`;
  }

  if (
    lowerKind.includes('ingress') ||
    lowerKind.includes('gateway') ||
    lowerKind.includes('route')
  ) {
    return `External ingress routing and SSL termination will fail; services exposed by this ingress will not be reachable outside the cluster.`;
  }

  if (lowerKind.includes('configmap') || lowerKind.includes('secret')) {
    return `Workloads referencing this ${kind} as environment variables or volume mounts will enter CrashLoopBackOff or fail to start.`;
  }

  if (lowerKind.includes('service') && !lowerKind.includes('mesh')) {
    return `Internal service discovery and inter-pod network communication targeting ${resourceName} will fail.`;
  }

  return `Components that depend on or communicate with ${resourceName} may encounter degraded performance or runtime errors.`;
}

/**
 * Extracts and normalizes dry-run response data into a unified cross-cluster matrix.
 */
export function buildCompatibilityAnalysis({
  dryRunResponse,
  design,
  selectedK8sConnections = [],
  selectedK8sContexts = [],
}: {
  dryRunResponse?: any;
  design?: any;
  selectedK8sConnections?: any[];
  selectedK8sContexts?: string[];
}): CompatibilityAnalysisData {
  // 1. Resolve cluster list
  const clusterMap: Record<string, ClusterInfo> = {};

  // Build clusters from selected connections first (skip 'all' which is a context selector, not a cluster ID)
  selectedK8sConnections.forEach((conn) => {
    const id = conn.metadata?.id || conn.id;
    if (id && id !== 'all') {
      const name = conn.metadata?.name || conn.name || `cluster-${id.slice(0, 6)}`;
      const version =
        conn.metadata?.server_version ||
        conn.metadata?.version ||
        conn.metadata?.k8s_version ||
        'Unknown';
      const status = conn.status || 'unknown';
      clusterMap[id] = { id, name, version, status };
    }
  });

  // Also include any concrete contexts in selectedK8sContexts (excluding 'all')
  const effectiveContexts = selectedK8sContexts.filter((ctxId) => ctxId && ctxId !== 'all');
  effectiveContexts.forEach((ctxId) => {
    if (!clusterMap[ctxId]) {
      clusterMap[ctxId] = {
        id: ctxId,
        name: `cluster-${ctxId.slice(0, 6)}`,
        version: 'Unknown',
        status: 'unknown',
      };
    }
  });

  // Also discover any cluster IDs in the dryRunResponse (excluding 'all')
  if (dryRunResponse && typeof dryRunResponse === 'object') {
    Object.values(dryRunResponse).forEach((contextMap: any) => {
      if (contextMap && typeof contextMap === 'object') {
        Object.keys(contextMap).forEach((ctxId) => {
          if (ctxId && ctxId !== 'all' && !clusterMap[ctxId]) {
            clusterMap[ctxId] = {
              id: ctxId,
              name: `cluster-${ctxId.slice(0, 6)}`,
              version: 'Unknown',
              status: 'unknown',
            };
          }
        });
      }
    });
  }

  const clusters = Object.values(clusterMap);
  if (clusters.length === 0) {
    clusters.push({
      id: 'default-cluster',
      name: 'kubernetes-cluster',
      version: 'Unknown',
      status: 'unknown',
    });
  }

  // 2. Resolve resources from design or dryRunResponse
  const resourcesMap: Record<string, ResourceRow> = {};

  // If design components are present, seed the resources
  const designComponents = design?.components || design?.pattern_file?.components || [];
  const allResourceNames: string[] = [];

  designComponents.forEach((comp: any) => {
    const name = comp.displayName || comp.name || comp.id;
    if (name) {
      allResourceNames.push(name);
      resourcesMap[name] = {
        resourceName: name,
        kind: comp.component?.kind || comp.kind || 'Resource',
        apiVersion: comp.component?.version || comp.apiVersion || 'v1',
        component: comp,
        clusterResults: {},
      };
    }
  });

  // Parse dryRunResponse entries
  if (dryRunResponse && typeof dryRunResponse === 'object') {
    Object.entries(dryRunResponse).forEach(([compName, contextMap]: [string, any]) => {
      if (!resourcesMap[compName]) {
        allResourceNames.push(compName);
        resourcesMap[compName] = {
          resourceName: compName,
          kind: 'Component',
          apiVersion: 'v1',
          clusterResults: {},
        };
      }

      if (contextMap && typeof contextMap === 'object') {
        Object.entries(contextMap).forEach(([clusterId, result]: [string, any]) => {
          const clusterName = clusterMap[clusterId]?.name || clusterId;
          const isSuccess = Boolean(result?.success);

          if (result?.component) {
            if (result.component.kind) {
              resourcesMap[compName].kind = result.component.kind;
            }
            if (result.component.version) {
              resourcesMap[compName].apiVersion = result.component.version;
            }
          }

          let status: CompatibilityStatus = isSuccess ? 'success' : 'failed';
          let errorMessage = '';
          let errorType = '';
          let fieldPath = '';

          if (!isSuccess && result?.error) {
            const err = result.error;
            if (Array.isArray(err.Causes) && err.Causes.length > 0) {
              errorMessage = err.Causes.map((c: any) => c.Message || c).join('; ');
              errorType = err.Causes[0]?.Type || '';
              fieldPath = err.Causes[0]?.FieldPath || '';
            } else if (err.Status) {
              errorMessage = String(err.Status);
            } else if (typeof err === 'string') {
              errorMessage = err;
            } else {
              errorMessage = JSON.stringify(err);
            }

            // Detect warnings (e.g. non-blocking annotations or deprecations)
            if (
              errorMessage.toLowerCase().includes('warning') ||
              errorMessage.toLowerCase().includes('deprecated')
            ) {
              status = 'warning';
            }
          }

          const possibleFix = !isSuccess
            ? generatePossibleFix(errorMessage, resourcesMap[compName].kind, clusterName)
            : undefined;

          const dependencyImpact = !isSuccess
            ? generateDependencyImpact(compName, resourcesMap[compName].kind, allResourceNames)
            : undefined;

          resourcesMap[compName].clusterResults[clusterId] = {
            status,
            errorMessage: errorMessage || undefined,
            errorType: errorType || undefined,
            fieldPath: fieldPath || undefined,
            possibleFix,
            dependencyImpact,
          };
        });
      }
    });
  }

  // Ensure every cluster has a result status for each resource (mark 'incomplete' if absent)
  const resources = Object.values(resourcesMap);
  resources.forEach((resource) => {
    clusters.forEach((cluster) => {
      if (!resource.clusterResults[cluster.id]) {
        // If there's no result for this cluster, mark as incomplete
        resource.clusterResults[cluster.id] = {
          status: 'incomplete',
          errorMessage: `No dry-run result was returned for ${cluster.name}.`,
        };
      }
    });
  });

  // 3. Compute Metrics
  let succeededCount = 0;
  let warningCount = 0;
  let failedCount = 0;
  let incompleteCount = 0;

  // Aggregate resource-level overall status
  resources.forEach((resource) => {
    const statuses = Object.values(resource.clusterResults).map((r) => r.status);
    if (statuses.includes('failed')) {
      failedCount++;
    } else if (statuses.includes('incomplete')) {
      incompleteCount++;
    } else if (statuses.includes('warning')) {
      warningCount++;
    } else {
      succeededCount++;
    }
  });

  const metrics: ValidationMetrics = {
    totalClusters: clusters.length,
    totalResources: resources.length,
    succeededCount,
    warningCount,
    failedCount,
    incompleteCount,
  };

  // 4. Compute per-cluster summaries
  const clusterSummaries: ClusterSummaryItem[] = clusters.map((cluster) => {
    let passed = 0;
    let warnings = 0;
    let failed = 0;
    let incomplete = 0;

    resources.forEach((res) => {
      const resStatus = res.clusterResults[cluster.id]?.status || 'incomplete';
      if (resStatus === 'success') {
        passed++;
      } else if (resStatus === 'warning') {
        warnings++;
      } else if (resStatus === 'failed') {
        failed++;
      } else if (resStatus === 'incomplete') {
        incomplete++;
      }
    });

    return {
      clusterId: cluster.id,
      clusterName: cluster.name,
      version: cluster.version,
      passed,
      warnings,
      failed,
      incomplete,
      total: resources.length,
    };
  });

  // 5. Generate Next Steps
  const nextSteps: string[] = [];
  const incompleteClusters = new Set<string>();
  resources.forEach((res) => {
    clusters.forEach((cluster) => {
      const r = res.clusterResults[cluster.id];
      if (r?.status === 'failed') {
        nextSteps.push(`Resolve the failed resource (${res.resourceName}) on ${cluster.name}.`);
      } else if (r?.status === 'warning') {
        nextSteps.push(`Review the warning for ${res.resourceName} on ${cluster.name}.`);
      } else if (r?.status === 'incomplete') {
        incompleteClusters.add(cluster.name);
      }
    });
  });

  incompleteClusters.forEach((clusterName) => {
    nextSteps.push(
      `Obtain validation results for unavailable or incomplete cluster ${clusterName}.`,
    );
  });

  if (nextSteps.length === 0) {
    nextSteps.push('All resources successfully validated across all targeted clusters.');
    nextSteps.push('Design is ready for deployment.');
  } else {
    nextSteps.push('Re-run the analysis after making the necessary changes.');
  }

  return {
    clusters,
    resources,
    metrics,
    clusterSummaries,
    nextSteps,
  };
}
