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
