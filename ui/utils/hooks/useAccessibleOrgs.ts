import { useLazyGetUserKeysQuery } from '@meshery/schemas/mesheryApi';
import { useAccessibleOrgs as useSistentAccessibleOrgs } from '@sistent/sistent';
import type { PermissionKeySpec } from '@sistent/sistent';
import { useGetSelectedOrganization } from '@/rtk-query/user';

/**
 * Thin wrapper around sistent's `useAccessibleOrgs` that wires in
 * meshery's org data sources and mesheryApi trigger.
 */
export const useAccessibleOrgs = (permissionKey?: PermissionKeySpec) => {
  const {
    allOrganizations,
    selectedOrganization,
    isLoading: orgsLoading,
  } = useGetSelectedOrganization();
  const [triggerGetKeys] = useLazyGetUserKeysQuery();

  return useSistentAccessibleOrgs({
    allOrgs: allOrganizations,
    currentOrgId: selectedOrganization?.id,
    orgsLoaded: !orgsLoading,
    permissionKey,
    triggerGetKeys,
  });
};
