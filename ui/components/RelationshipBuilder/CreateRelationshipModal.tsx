import React from 'react';
import { Modal } from '@sistent/sistent';
import RelationshipFormStepper from './RelationshipFormStepper';

const CreateRelationshipModal = ({
  isRelationshipModalOpen,
  setIsRelationshipModalOpen,
}: {
  isRelationshipModalOpen: boolean;
  setIsRelationshipModalOpen: (_value: boolean) => void;
}) => {
  const modalProps = {
    maxWidth: 'md' as const,
    open: isRelationshipModalOpen,
    closeModal: () => setIsRelationshipModalOpen(false),
    title: 'Create Relationship',
    style: {
      zIndex: 1500,
    },
    disableBackdropClick: false,
    disableEscapeKeyDown: false,
  };

  return (
    <Modal {...(modalProps as any)}>
      <RelationshipFormStepper handleClose={() => setIsRelationshipModalOpen(false)} />
    </Modal>
  );
};

export default CreateRelationshipModal;
