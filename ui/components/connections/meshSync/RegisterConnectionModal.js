import React from 'react';
import { DialogContent, Dialog } from '@material-ui/core';

import theme from '../../../themes/app.js';
import CustomizedSteppers from './Stepper/index.js';

import { useCancelConnectionRegisterMutation } from '@/rtk-query/connection.js';
import { useDeleteMeshsyncResourceMutation } from '@/rtk-query/meshsync.js';
import { useNotification } from '@/utils/hooks/useNotification.js';
import { EVENT_TYPES } from 'lib/event-types.js';

const RegisterConnectionModal = ({
  openRegistrationModal,
  connectionData,
  handleRegistrationModalClose,
}) => {
  const [sharedData, setSharedData] = React.useState(null);
  const { notify } = useNotification();
  const [cancelConnection] = useCancelConnectionRegisterMutation();
  const [deleteMeshsyncResource] = useDeleteMeshsyncResourceMutation();

  const cancelConnectionRegister = (id) => {
    cancelConnection({ body: JSON.stringify({ id }) })
      .unwrap()
      .then(() => {
        notify({
          message: 'Connection registration cancelled.',
          event_type: EVENT_TYPES.INFO,
        });
      });
  };
  const handleClose = () => {
    handleRegistrationModalClose();
    cancelConnectionRegister(sharedData?.connection?.id);
  };

  const handleRegistrationComplete = (resourceId) => {
    deleteMeshsyncResource({ resourceId: resourceId })
      .unwrap()
      .then(() => {
        notify({
          message: 'Connection registered successfully!',
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((error) => {
        notify({
          message: `Failed to register connection: ${error}`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  return (
    <div style={{ marginBottom: '1rem' }}>
      <Dialog
        open={openRegistrationModal}
        onClose={handleClose}
        aria-labelledby="form-dialog-title"
        maxWidth="md"
        style={{ zIndex: 9999 }}
        PaperProps={{
          style: { borderRadius: 30 },
        }}
      >
        <DialogContent
          style={{
            padding: '0 1.5rem 1.5rem',
            borderRadius: '28px',
            border: `6px solid ${theme.palette.secondary.success}`,
          }}
        >
          <CustomizedSteppers
            formConnectionIdRef
            onClose={handleClose}
            connectionData={connectionData}
            sharedData={sharedData}
            setSharedData={setSharedData}
            handleRegistrationComplete={handleRegistrationComplete}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default RegisterConnectionModal;
