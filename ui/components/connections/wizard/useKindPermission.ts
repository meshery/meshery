import { useHasPermission } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import type { ConnectionWizardKindConfig } from '../ConnectionWizard.helpers';

/**
 * Single authorization rule for "may this user create a connection of this
 * kind?". Kept in its own module (rather than alongside the step UI) so the
 * wizard hook and the kind chooser share one rule without either pulling in the
 * other's component tree.
 *
 * Returns a predicate rather than a boolean because callers evaluate it against
 * a kind that is only known at call time (a preset, or the row the user clicks).
 */
export const useKindPermission = () => {
  const canAddCluster = useHasPermission(Keys.LifecycleManagementAddCluster);
  const canConnectMetrics = useHasPermission(Keys.MesherySystemConnectMetrics);

  return (config?: ConnectionWizardKindConfig | null) => {
    if (!config) return false;
    return config.flow === 'kubernetes' ? canAddCluster : canConnectMetrics;
  };
};
