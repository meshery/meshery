import React from 'react';
import { Modal, ModalBody } from '@layer5/sistent';
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
          message: 'Connection registered!',
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
      <Modal
        open={openRegistrationModal}
        closeModal={handleClose}
        aria-labelledby="form-dialog-title"
        maxWidth="md"
      >
        <ModalBody>
          <CustomizedSteppers
            formConnectionIdRef
            onClose={handleClose}
            connectionData={connectionData}
            sharedData={sharedData}
            setSharedData={setSharedData}
            handleRegistrationComplete={handleRegistrationComplete}
          />
        </ModalBody>
      </Modal>
    </div>
  );
};

export default RegisterConnectionModal;
