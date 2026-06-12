import { useConfigurationTableConfig } from './configuration/config';
import { useNetWorkTableConfig } from './network/config';
import { useSecurityTypesConfig } from './security/config';
import { useStorageTableConfig } from './storage/config';
import { useWorkloadTableConfig } from './workloads/config';
import { useNamespaceTableConfig } from './namespace/config';
import { useNodeTableConfig } from './nodes/config';
import { useCustomResourceConfig } from './crds/config';
import _ from 'lodash';

export const ResourcesConfig = {
  Node: {
    useTableConfig: useNodeTableConfig,
    submenu: false,
  },
  Namespace: {
    useTableConfig: useNamespaceTableConfig,
    submenu: false,
  },
  Workload: {
    useTableConfig: useWorkloadTableConfig,
    submenu: true,
  },
  Configuration: {
    useTableConfig: useConfigurationTableConfig,
    submenu: true,
  },
  Network: {
    useTableConfig: useNetWorkTableConfig,
    submenu: true,
  },
  Security: {
    useTableConfig: useSecurityTypesConfig,
    submenu: true,
  },
  Storage: {
    useTableConfig: useStorageTableConfig,
    submenu: true,
  },
  CRDS: {
    useTableConfig: useCustomResourceConfig,
    submenu: true,
  },
};

export const ALL_VIEW = 'all';
export const SINGLE_VIEW = 'single';

export const ResourceMenuConfig = {
  Node: [],
  Namespace: [],
  Workload: [
    'Pod',
    'Deployment',
    'StatefulSet',
    'DaemonSet',
    'Job',
    'CronJob',
    'ReplicaSet',
    'ReplicationController',
  ],
  Configuration: [
    'ConfigMap',
    'Secret',
    'ResourceQuota',
    'LimitRange',
    'HorizontalPodAutoscaler',
    'VerticalPodAutoscaler',
    'PodDisruptionBudget',
    'PriorityClass',
    'RuntimeClass',
    'Leases',
    'MutatingWebhookConfiguration',
  ],
  Network: ['Service', 'Endpoints', 'EndpointSlice', 'Ingress', 'IngressClass', 'NetworkPoliciy'],
  Security: ['ServiceAccount', 'ClusterRole', 'Role', 'ClusterRoleBinding', 'RoleBinding'],
  Storage: ['PersistentVolume', 'PersistentVolumeClaim', 'StorageClass'],
};

export function generateDynamicURL(kind) {
  let resourceCategory = '';
  let resource = '';

  // Iterate through ResourceMenuConfig to find the matching category and resource
  for (const [category, kinds] of Object.entries(ResourceMenuConfig)) {
    if (category === kind) {
      resourceCategory = category;
      resource = ''; // No specific resource for Node and Namespace
      break;
    }

    if (kinds.includes(kind)) {
      resourceCategory = category;
      resource = kind;
      break;
    }
    resourceCategory = 'CRDS';
    resource = kind;
  }

  // Construct the URL based on the found category and resource
  return `?resourceCategory=${resourceCategory}&resource=${resource}`;
}

export const getAllCustomResourceDefinitionsKinds = (kinds) => {
  const resourceMenuArray = ['Node', 'Namespace', ..._.flatten(Object.values(ResourceMenuConfig))];
  const customResources = kinds?.filter((kind) => {
    return !resourceMenuArray.includes(kind.Kind);
  });
  return customResources ?? [];
};
