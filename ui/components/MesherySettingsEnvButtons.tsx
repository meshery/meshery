import { Button, Typography } from '@sistent/sistent';
import AddIconCircleBorder from '../assets/icons/AddIconCircleBorder';
import { Keys } from '@meshery/schemas/permissions';
import useTestIDsGenerator from '@/utils/hooks/useTestIDs';
import CAN from '@/utils/can';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { useConnectionWizardModal } from '@/utils/context/ConnectionWizardContextProvider';

type MesherySettingsEnvButtonsProps = {
  /** Called after the wizard is opened (e.g. close the context-switcher menu). */
  onOpened?: () => void;
};

/**
 * Context-switcher "Add Cluster" entry point.
 *
 * Opens the shared Create Connection wizard in place (any page) with Kubernetes
 * pre-selected at Import Kubeconfig — same host used by telemetry and the
 * Connections toolbar.
 */
const MesherySettingsEnvButtons = ({ onOpened }: MesherySettingsEnvButtonsProps) => {
  const testIDs = useTestIDsGenerator('connection');
  const { openCreateConnection } = useConnectionWizardModal();

  const handleClick = () => {
    openCreateConnection({
      kind: CONNECTION_KINDS.KUBERNETES,
      skipKindSelection: true,
    });
    onOpened?.();
  };

  return (
    <div>
      <Button
        type="button"
        variant="contained"
        onClick={handleClick}
        style={{
          width: '100%',
          borderRadius: 5,
          padding: '8px',
        }}
        disabled={
          !CAN(Keys.LifecycleManagementAddCluster.id, Keys.LifecycleManagementAddCluster.function)
        }
        data-cy="btnAddCluster"
      >
        <AddIconCircleBorder style={{ width: '20px', height: '20px' }} />
        <Typography
          style={{
            paddingLeft: '4px',
            width: 'max-content',
            marginRight: '4px',
          }}
          data-testid={testIDs('addCluster')}
        >
          Add Cluster
        </Typography>
      </Button>
    </div>
  );
};

export default MesherySettingsEnvButtons;
