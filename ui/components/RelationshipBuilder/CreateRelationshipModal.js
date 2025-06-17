import React from 'react';
import { Modal } from '@sistent/sistent';
import RelationshipFormStepper from './RelationshipFormStepper';

const CreateRelationshipModal = ({ isRelationshipModalOpen, setIsRelationshipModalOpen }) => {
  return (
    <Modal
      maxWidth="sm"
      open={isRelationshipModalOpen}
      closeModal={() => setIsRelationshipModalOpen(false)}
      title="Create Relationship"
      style={{
        zIndex: 1500,
      }}
      disableBackdropClick={false}
      disableEscapeKeyDown={false}
    >
      <RelationshipFormStepper handleClose={() => setIsRelationshipModalOpen(false)} />
    </Modal>
  );
};

export default CreateRelationshipModal;
