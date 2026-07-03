import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import { buildConfigMapColumns } from './configmap-columns';
import { buildSecretColumns } from './secret-columns';
import { buildResourceQuotaColumns } from './resource-quota-columns';
import { buildLimitRangeColumns } from './limit-range-columns';
import { buildHorizontalPodAutoscalerColumns } from './horizontal-pod-autoscaler-columns';
import { buildVerticalPodAutoscalerColumns } from './vertical-pod-autoscaler-columns';
import { buildPodDisruptionBudgetColumns } from './pod-disruption-budget-columns';
import { buildPriorityClassColumns } from './priority-class-columns';
import { buildRuntimeClassColumns } from './runtime-class-columns';
import { buildLeasesColumns } from './leases-columns';
import { buildMutatingWebhookConfigurationColumns } from './mutating-webhook-configuration-columns';

export const useConfigurationTableConfig = (
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
    ConfigMap: buildConfigMapColumns(args),
    Secret: buildSecretColumns(args),
    ResourceQuota: buildResourceQuotaColumns(args),
    LimitRange: buildLimitRangeColumns(args),
    HorizontalPodAutoscaler: buildHorizontalPodAutoscalerColumns(args),
    VerticalPodAutoscaler: buildVerticalPodAutoscalerColumns(args),
    PodDisruptionBudget: buildPodDisruptionBudgetColumns(args),
    PriorityClass: buildPriorityClassColumns(args),
    RuntimeClass: buildRuntimeClassColumns(args),
    Leases: buildLeasesColumns(args),
    MutatingWebhookConfiguration: buildMutatingWebhookConfigurationColumns(args),
  };
};
