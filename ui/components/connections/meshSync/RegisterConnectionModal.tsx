import React from 'react';
import { Modal, ModalBody } from '@sistent/sistent';
import CustomizedSteppers from './Stepper';

import { useCancelConnectionRegisterMutation } from '@/rtk-query/connection';
import { useDeleteMeshsyncResourceMutation } from '@/rtk-query/meshsync';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

type ConnectionRegistrationData = {
  connection?: {
    id?: string;
  };
  resourceID?: string;
} & Record<string, unknown>;

type RegisterConnectionModalProps = {
  openRegistrationModal: boolean;
  connectionData: ConnectionRegistrationData;
  handleRegistrationModalClose: () => void;
};

const RegisterConnectionModal = ({
  openRegistrationModal,
  connectionData,
  handleRegistrationModalClose,
}: RegisterConnectionModalProps) => {
  const [sharedData, setSharedData] = React.useState<ConnectionRegistrationData | null>(null);
  const { notify } = useNotification();
  const [cancelConnection] = useCancelConnectionRegisterMutation();
  const [deleteMeshsyncResource] = useDeleteMeshsyncResourceMutation();

  const cancelConnectionRegister = (id?: string) => {
    if (!id) {
      return;
    }

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

  const handleRegistrationComplete = (resourceId?: string) => {
    if (!resourceId) {
      return;
    }

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
