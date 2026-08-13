import { Button, Typography, useHasPermission } from '@sistent/sistent';
import { styled } from '@/theme';
import { Keys } from '@meshery/schemas/permissions';
import AddIconCircleBorder from '@/assets/icons/AddIconCircleBorder';
import { useConnectionWizardModal } from '@/utils/context/ConnectionWizardContextProvider';

const LaunchButton = styled(Button)(({ theme }) => ({
  width: '100%',
  borderRadius: 5,
  padding: '8px',
  [theme.breakpoints.down('sm')]: {
    padding: '4px 6px',
  },
}));

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
        sx={{
          paddingLeft: '4px',
          width: 'max-content',
          marginRight: '4px',
          fontSize: { xs: '0.75rem', sm: '1rem' },
        }}
      >
        Create Connection
      </Typography>
    </LaunchButton>
  );
};

export default ConnectionWizardLauncher;
