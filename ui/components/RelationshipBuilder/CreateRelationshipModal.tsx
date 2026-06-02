import React from 'react';
import { Modal, useTheme } from '@sistent/sistent';
import RelationshipFormStepper from './RelationshipFormStepper';

const CreateRelationshipModal = ({ isRelationshipModalOpen, setIsRelationshipModalOpen }) => {
  const theme = useTheme();
  return (
    <Modal
      maxWidth="md"
      open={isRelationshipModalOpen}
      closeModal={() => setIsRelationshipModalOpen(false)}
      title="Create Relationship"
      style={{
        zIndex: theme.zIndex.modal + 100,
      }}
      disableBackdropClick={false}
      disableEscapeKeyDown={false}
    >
      <RelationshipFormStepper handleClose={() => setIsRelationshipModalOpen(false)} />
    </Modal>
  );
};

export default CreateRelationshipModal;
