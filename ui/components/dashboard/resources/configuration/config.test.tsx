import { describe, expect, it, vi } from 'vitest';

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => () => undefined,
}));

vi.mock('./configmap-columns', () => ({
  buildConfigMapColumns: () => ({ name: 'ConfigMap', columns: [] }),
}));
vi.mock('./secret-columns', () => ({
  buildSecretColumns: () => ({ name: 'Secret', columns: [] }),
}));
vi.mock('./resource-quota-columns', () => ({
  buildResourceQuotaColumns: () => ({ name: 'ResourceQuota', columns: [] }),
}));
vi.mock('./limit-range-columns', () => ({
  buildLimitRangeColumns: () => ({ name: 'LimitRange', columns: [] }),
}));
vi.mock('./horizontal-pod-autoscaler-columns', () => ({
  buildHorizontalPodAutoscalerColumns: () => ({ name: 'HorizontalPodAutoscaler', columns: [] }),
}));
vi.mock('./vertical-pod-autoscaler-columns', () => ({
  buildVerticalPodAutoscalerColumns: () => ({ name: 'VerticalPodAutoscaler', columns: [] }),
}));
vi.mock('./pod-disruption-budget-columns', () => ({
  buildPodDisruptionBudgetColumns: () => ({ name: 'PodDisruptionBudget', columns: [] }),
}));
vi.mock('./priority-class-columns', () => ({
  buildPriorityClassColumns: () => ({ name: 'PriorityClass', columns: [] }),
}));
vi.mock('./runtime-class-columns', () => ({
  buildRuntimeClassColumns: () => ({ name: 'RuntimeClass', columns: [] }),
}));
vi.mock('./leases-columns', () => ({
  buildLeasesColumns: () => ({ name: 'Leases', columns: [] }),
}));
vi.mock('./mutating-webhook-configuration-columns', () => ({
  buildMutatingWebhookConfigurationColumns: () => ({
    name: 'MutatingWebhookConfiguration',
    columns: [],
  }),
}));

import { useConfigurationTableConfig } from './config';

describe('useConfigurationTableConfig', () => {
  it('exposes the canonical configuration workload configs', () => {
    const result = useConfigurationTableConfig(vi.fn(), [], {}, {}, 'ConfigMap');
    expect(Object.keys(result)).toEqual([
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
    ]);
    expect(result.ConfigMap.name).toBe('ConfigMap');
  });
});
