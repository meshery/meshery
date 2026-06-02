import React from 'react';
import { Modal, useTheme } from '@sistent/sistent';
import UrlStepper from './Stepper/UrlStepper';

type CreateModelModalProps = {
  isCreateModalOpen: boolean;
  setIsCreateModalOpen: (_open: boolean) => void;
};

const CreateModelModal = ({ isCreateModalOpen, setIsCreateModalOpen }: CreateModelModalProps) => {
  const theme = useTheme();
  return (
    <Modal
      maxWidth="sm"
      open={isCreateModalOpen}
      closeModal={() => setIsCreateModalOpen(false)}
      title="Create Model"
      style={{
        zIndex: theme.zIndex.modal + 100,
      }}
    >
      <UrlStepper handleClose={() => setIsCreateModalOpen(false)} />
    </Modal>
  );
};

export default CreateModelModal;
