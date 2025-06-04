import React from 'react';
import { Modal } from '@sistent/sistent';
import UrlStepper from './Stepper/UrlStepper';

const CreateModelModal = ({ isCreateModalOpen, setIsCreateModalOpen }) => {
  return (
    <Modal
      maxWidth="sm"
      open={isCreateModalOpen}
      closeModal={() => setIsCreateModalOpen(false)}
      title="Create Model"
      style={{
        zIndex: 1500,
      }}
    >
      <UrlStepper handleClose={() => setIsCreateModalOpen(false)} />
    </Modal>
  );
};

export default CreateModelModal;
