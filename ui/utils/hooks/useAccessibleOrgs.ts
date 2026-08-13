import { useAccessibleOrgs as useSistentAccessibleOrgs } from '@sistent/sistent';
import type { PermissionKeySpec } from '@sistent/sistent';
import { useGetSelectedOrganization } from '@/rtk-query/user';

/**
 * Thin wrapper around sistent's `useAccessibleOrgs` that wires in
 * meshery's org data sources.
 *
 * Returns only organizations where the user holds the permission(s) described
 * by `permissionKey`, excluding the current org.
 */
export const useAccessibleOrgs = (permissionKey?: PermissionKeySpec) => {
  const {
    allOrganizations,
    selectedOrganization,
    isLoading: orgsLoading,
  } = useGetSelectedOrganization();

  return useSistentAccessibleOrgs({
    allOrgs: allOrganizations,
    currentOrgId: selectedOrganization?.id,
    orgsLoaded: !orgsLoading,
    permissionKey,
  });
};
