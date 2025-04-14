import React from 'react';
import { Modal } from '@layer5/sistent';
import UrlStepper from './Stepper/UrlStepper';

const CreateModelModal = ({ isCreateModalOpen, setIsCreateModalOpen }) => {
  return (
    <Modal
      maxWidth="sm"
      open={isCreateModalOpen}
      closeModal={() => setIsCreateModalOpen(false)}
      title="Create Model"
    >
      <UrlStepper handleClose={() => setIsCreateModalOpen(false)} />
    </Modal>
  );
};

export default CreateModelModal;
