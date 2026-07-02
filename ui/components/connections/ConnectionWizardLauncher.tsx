import { useState } from 'react';
import { Button, Typography } from '@sistent/sistent';
import { styled } from '@/theme';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import AddIconCircleBorder from '@/assets/icons/AddIconCircleBorder';
import ConnectionWizardModal from './ConnectionWizardModal';
const LaunchButton = styled(Button)({
  width: '100%',
  borderRadius: 5,
  padding: '8px',
});

const canOpenConnectionWizard = () =>
  CAN(keys.ADD_CLUSTER.action, keys.ADD_CLUSTER.subject) ||
  CAN(keys.CONNECT_METRICS.action, keys.CONNECT_METRICS.subject);

const ConnectionWizardLauncher = () => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <LaunchButton
        type="button"
        variant="contained"
        onClick={() => setIsOpen(true)}
        disabled={!canOpenConnectionWizard()}
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
      {isOpen && <ConnectionWizardModal isOpen={isOpen} onClose={() => setIsOpen(false)} />}
    </>
  );
};

export default ConnectionWizardLauncher;
