import React from 'react';
import Environments from '@/components/environments';
import { EnvironmentIcon, useTheme, Box } from '@sistent/sistent';
import Modal from '@/components/shared/Modal/Modal';

const EnvironmentFormModal = ({ isOpenModal, setIsOpenModal }) => {
  const theme = useTheme();

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
      <Box sx={{ p: 1, maxHeight: '65vh' }}>
        <Environments />
      </Box>
    </Modal>
  );
};

export default EnvironmentFormModal;
