import React from 'react';
import { iconMedium } from '../../css/icons.styles';
import { AddCircleIcon as AddIcon, Button, useTheme } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import { CoreConnectionKinds } from '@/utils/Enum';
import { useConnectionWizardModal } from '@/utils/context/ConnectionWizardContextProvider';

/**
 * Empty-state / chart "Connect Clusters" entry point.
 *
 * Opens the shared Create Connection wizard in place with Kubernetes
 * pre-selected at Import Kubeconfig — same host used by Add Cluster,
 * telemetry, and the Connections toolbar.
 */
function ConnectClustersBtn() {
  const theme = useTheme();
  const { openCreateConnection } = useConnectionWizardModal();

  const handleClick = () => {
    openCreateConnection({
      kind: CoreConnectionKinds.kubernetes,
      skipKindSelection: true,
    });
  };

  return (
    <Button
      type="button"
      variant="contained"
      color="primary"
      size="large"
      onClick={handleClick}
      permissionKey={Keys.LifecycleManagementAddCluster}
      style={{ margin: '0.5rem 0.5rem', whiteSpace: 'nowrap' }}
    >
      <AddIcon
        style={{
          width: theme.spacing(2.5),
          paddingRight: theme.spacing(0.5),
          ...iconMedium,
        }}
      />
      Connect Clusters
    </Button>
  );
}

export default ConnectClustersBtn;
