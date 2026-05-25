import { useMemo } from 'react';
import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import { buildServiceColumns } from './service-columns';
import { buildEndpointsColumns } from './endpoints-columns';
import { buildEndpointSliceColumns } from './endpoint-slice-columns';
import { buildIngressColumns } from './ingress-columns';
import { buildIngressClassColumns } from './ingress-class-columns';
import { buildNetworkPolicyColumns } from './network-policy-columns';

export const NetWorkTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
) => {
  const ping = useKubernetesHook();
  return useMemo(() => {
    const args = {
      switchView,
      meshSyncResources,
      k8sConfig,
      connectionMetadataState,
      workloadType,
      ping,
    };
    return {
      Service: buildServiceColumns(args),
      Endpoints: buildEndpointsColumns(args),
      EndpointSlice: buildEndpointSliceColumns(args),
      Ingress: buildIngressColumns(args),
      IngressClass: buildIngressClassColumns(args),
      NetworkPolicy: buildNetworkPolicyColumns(args),
    };
  }, [switchView, meshSyncResources, k8sConfig, connectionMetadataState, workloadType, ping]);
};
