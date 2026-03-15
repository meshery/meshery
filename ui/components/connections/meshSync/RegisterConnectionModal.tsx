import React from 'react';
import { Modal, ModalBody } from '@sistent/sistent';
import CustomizedSteppers from './Stepper';

import { useCancelConnectionRegisterMutation } from '@/rtk-query/connection';
import { useDeleteMeshsyncResourceMutation } from '@/rtk-query/meshsync';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

type RegisterConnectionModalProps = {
  openRegistrationModal: boolean;
  connectionData: any;
  handleRegistrationModalClose: () => void;
};

const RegisterConnectionModal = ({
  openRegistrationModal,
  connectionData,
  handleRegistrationModalClose,
}: RegisterConnectionModalProps) => {
  const [sharedData, setSharedData] = React.useState<any>(null);
  const { notify } = useNotification();
  const [cancelConnection] = useCancelConnectionRegisterMutation();
  const [deleteMeshsyncResource] = useDeleteMeshsyncResourceMutation();

  const cancelConnectionRegister = (id: string) => {
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
    if (sharedData?.connection?.id) {
      cancelConnectionRegister(sharedData.connection.id);
    }
  };

  const handleRegistrationComplete = (resourceId: string) => {
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
        title="Register Connection"
        aria-labelledby="form-dialog-title"
        maxWidth="md"
      >
        <ModalBody>
          <CustomizedSteppers
            onClose={handleClose}
            handleRegistrationComplete={handleRegistrationComplete}
            connectionData={connectionData}
            sharedData={sharedData}
            setSharedData={setSharedData}
          />
        </ModalBody>
      </Modal>
    </div>
  );
};

export default RegisterConnectionModal;
