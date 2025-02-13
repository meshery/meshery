import React from 'react';

import { Modal, ModalBody } from '@layer5/sistent';
import ConnectionIcon from '@/assets/icons/Connection';
import ConnectionTable from '../connections/ConnectionTable';

const ConnectionModal = ({
  isOpenModal,
  setIsOpenModal,
  meshsyncControllerState,
  connectionMetadataState,
}) => {
  if (!isOpenModal) {
    return null;
  }
  return (
    <Modal
      open={isOpenModal}
      closeModal={() => setIsOpenModal(false)}
      headerIcon={<ConnectionIcon height={24} width={24} />}
      title="Connections"
      maxWidth="xl"
    >
      <ModalBody>
        <div style={{ marginBlock: '2rem', maxHeight: '65vh' }}>
          <ConnectionTable
            meshsyncControllerState={meshsyncControllerState}
            connectionMetadataState={connectionMetadataState}
            selectedFilter={'kubernetes'}
          />
        </div>
      </ModalBody>
    </Modal>
  );
};

export default ConnectionModal;
