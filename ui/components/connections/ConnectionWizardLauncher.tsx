import { Button, Typography, useHasPermission } from '@sistent/sistent';
import { styled } from '@/theme';
import { Keys } from '@meshery/schemas/permissions';
import AddIconCircleBorder from '@/assets/icons/AddIconCircleBorder';
import { useConnectionWizardModal } from '@/utils/context/ConnectionWizardContextProvider';

const LaunchButton = styled(Button)({
  width: '100%',
  borderRadius: 5,
  padding: '8px',
});

/**
 * Connections-toolbar entry for Create Connection. Opens the app-level wizard
 * (no kind preset) so selection starts at "Choose Connection".
 */
const ConnectionWizardLauncher = () => {
  const { openCreateConnection } = useConnectionWizardModal();

  const canAddCluster = useHasPermission(Keys.LifecycleManagementAddCluster);
  const canConnectMetrics = useHasPermission(Keys.MesherySystemConnectMetrics);
  const hasPermission = canAddCluster || canConnectMetrics;

  return (
    <LaunchButton
      type="button"
      variant="contained"
      onClick={() => openCreateConnection()}
      permissionKey={!hasPermission ? Keys.LifecycleManagementAddCluster : undefined}
      data-testid="connection-create-connection"
    >
      <AddIconCircleBorder style={{ width: '20px', height: '20px' }} />
      <Typography
        style={{
          paddingLeft: '4px',
          width: 'max-content',
          marginRight: '4px',
        }}
      >
        Create Connection
      </Typography>
    </LaunchButton>
  );
};

export default ConnectionWizardLauncher;
