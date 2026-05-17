import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import { buildPodColumns } from './pod-columns';
import { buildDeploymentColumns } from './deployment-columns';
import { buildDaemonSetColumns } from './daemonset-columns';
import { buildStatefulSetColumns } from './statefulset-columns';
import { buildReplicaSetColumns } from './replicaset-columns';
import { buildReplicationControllerColumns } from './replication-controller-columns';
import { buildJobColumns } from './job-columns';
import { buildCronJobColumns } from './cronjob-columns';

export const WorkloadTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
) => {
  const ping = useKubernetesHook();
  const args = {
    switchView,
    meshSyncResources,
    k8sConfig,
    connectionMetadataState,
    workloadType,
    ping,
  };
  return {
    Pod: buildPodColumns(args),
    Deployment: buildDeploymentColumns(args),
    DaemonSet: buildDaemonSetColumns(args),
    StatefulSet: buildStatefulSetColumns(args),
    ReplicaSet: buildReplicaSetColumns(args),
    ReplicationController: buildReplicationControllerColumns(args),
    Job: buildJobColumns(args),
    CronJob: buildCronJobColumns(args),
  };
};
