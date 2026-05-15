import { describe, expect, it, vi } from 'vitest';

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => () => undefined,
}));

vi.mock('./pod-columns', () => ({ buildPodColumns: () => ({ name: 'Pod', columns: [] }) }));
vi.mock('./deployment-columns', () => ({
  buildDeploymentColumns: () => ({ name: 'Deployment', columns: [] }),
}));
vi.mock('./daemonset-columns', () => ({
  buildDaemonSetColumns: () => ({ name: 'DaemonSet', columns: [] }),
}));
vi.mock('./statefulset-columns', () => ({
  buildStatefulSetColumns: () => ({ name: 'StatefulSet', columns: [] }),
}));
vi.mock('./replicaset-columns', () => ({
  buildReplicaSetColumns: () => ({ name: 'ReplicaSet', columns: [] }),
}));
vi.mock('./replication-controller-columns', () => ({
  buildReplicationControllerColumns: () => ({ name: 'ReplicationController', columns: [] }),
}));
vi.mock('./job-columns', () => ({ buildJobColumns: () => ({ name: 'Job', columns: [] }) }));
vi.mock('./cronjob-columns', () => ({
  buildCronJobColumns: () => ({ name: 'CronJob', columns: [] }),
}));

import { WorkloadTableConfig } from './config';

describe('WorkloadTableConfig', () => {
  it('aggregates all the workload column builders into a kind-keyed object', () => {
    const result = WorkloadTableConfig(vi.fn(), [], {}, {}, 'Pod');

    expect(Object.keys(result)).toEqual([
      'Pod',
      'Deployment',
      'DaemonSet',
      'StatefulSet',
      'ReplicaSet',
      'ReplicationController',
      'Job',
      'CronJob',
    ]);
    expect(result.Pod.name).toBe('Pod');
  });
});
