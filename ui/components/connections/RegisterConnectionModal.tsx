/**
 * RegisterConnectionModal
 *
 * Wraps the connection-registration stepper in the shared `Modal` primitive
 * so the modal chrome matches the rest of the UI surface. Side-effects
 * around cancelling an in-flight registration and clearing the underlying
 * MeshSync resource on completion are preserved from the legacy module.
 */
import { FC, useState } from 'react';
import { Modal } from '@/components/shared/Modal';
import { useCancelConnectionRegisterMutation } from '@/rtk-query/connection';
import { useDeleteMeshsyncResourceMutation } from '@/rtk-query/meshsync';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import CustomizedSteppers from './meshSync/Stepper';

interface RegisterConnectionModalProps {
  openRegistrationModal: boolean;
  handleRegistrationModalClose: () => void;
  connectionData: unknown;
}

interface SharedData {
  connection?: { id?: string };
  [key: string]: unknown;
}

const RegisterConnectionModal: FC<RegisterConnectionModalProps> = ({
  openRegistrationModal,
  connectionData,
  handleRegistrationModalClose,
}) => {
  const [sharedData, setSharedData] = useState<SharedData | null>(null);
  const { notify } = useNotification();
  const [cancelConnection] = useCancelConnectionRegisterMutation();
  const [deleteMeshsyncResource] = useDeleteMeshsyncResourceMutation();

  const cancelConnectionRegister = (id?: string) => {
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
    deleteMeshsyncResource({ resourceId })
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
    <Modal
      isOpen={openRegistrationModal}
      onClose={handleClose}
      title="Register Connection"
      size="md"
      aria-labelledby="form-dialog-title"
    >
      <CustomizedSteppers
        formConnectionIdRef
        onClose={handleClose}
        connectionData={connectionData}
        sharedData={sharedData}
        setSharedData={setSharedData}
        handleRegistrationComplete={handleRegistrationComplete}
      />
    </Modal>
  );
};

export default RegisterConnectionModal;
