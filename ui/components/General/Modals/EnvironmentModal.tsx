import React from 'react';
import { EnvironmentComponent } from '../../Lifecycle';

import { EnvironmentIcon, Modal, ModalBody, useTheme, Box } from '@sistent/sistent';

type EnvironmentModalProps = {
  isOpenModal: boolean;
  setIsOpenModal: (isOpenModal: boolean) => void;
};

const EnvironmentModal = ({ isOpenModal, setIsOpenModal }: EnvironmentModalProps) => {
  const theme = useTheme();
  if (!isOpenModal) {
    return null;
  }
  return (
    <Modal
      open={isOpenModal}
      closeModal={() => setIsOpenModal(false)}
      headerIcon={
        <EnvironmentIcon height={24} width={24} fill={theme.palette.background.constant.white} />
      }
      title="Environment"
      maxWidth="xl"
    >
      <ModalBody>
        <Box maxHeight={'65vh'}>
          <EnvironmentComponent />
        </Box>
      </ModalBody>
    </Modal>
  );
};

export default EnvironmentModal;
