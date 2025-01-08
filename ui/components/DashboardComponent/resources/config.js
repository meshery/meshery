import { ConfigurationTableConfig } from './configuration/config';
import { NetWorkTableConfig } from './network/config';
import { SecurityTypesConfig } from './security/config';
import { StorageTableConfig } from './storage/config';
import { WorkloadTableConfig } from './workloads/config';
import { NamespaceTableConfig } from './namespace/config';
import { NodeTableConfig } from './nodes/config';

export const ResourcesConfig = {
  Node: {
    tableConfig: NodeTableConfig,
    submenu: false,
  },
  Namespace: {
    tableConfig: NamespaceTableConfig,
    submenu: false,
  },
  Workload: {
    tableConfig: WorkloadTableConfig,
    submenu: true,
  },
  Configuration: {
    tableConfig: ConfigurationTableConfig,
    submenu: true,
  },
  Network: {
    tableConfig: NetWorkTableConfig,
    submenu: true,
  },
  Security: {
    tableConfig: SecurityTypesConfig,
    submenu: true,
  },
  Storage: {
    tableConfig: StorageTableConfig,
    submenu: true,
  },
};

export const ALL_VIEW = 'all';
export const SINGLE_VIEW = 'single';

const ResourceMenuConfig = {
  Node: [],
  Namespace: ['Namespace'],
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
  Network: ['Service', 'Endpoints', 'Ingress', 'IngressClass', 'NetworkPoliciy'],
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
  }

  // Construct the URL based on the found category and resource
  return `?resourceCategory=${resourceCategory}&resource=${resource}`;
}
