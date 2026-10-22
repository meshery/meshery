/**
 * Relationship-creation modal.
 *
 * Thin wrapper around the shared `Modal` primitive that mounts the
 * `RelationshipFormStepper` flow. The stepper renders its own `ModalBody`
 * and `ModalFooter` (back/next button bar), so we pass `disableBodyWrap`
 * to keep the shared `Modal` from nesting an extra `ModalBody`.
 *
 * Migrated from `components/relationship-builder/CreateRelationshipModal.tsx`
 * as part of Phase 5.b.6 (#18754).
 */
import { FC } from 'react';
import { Modal } from '@/components/shared/Modal';
import RelationshipFormStepper from '../relationship-builder/RelationshipFormStepper';

export interface CreateRelationshipModalProps {
  isRelationshipModalOpen: boolean;
  setIsRelationshipModalOpen: (_open: boolean) => void;
}

const CreateRelationshipModal: FC<CreateRelationshipModalProps> = ({
  isRelationshipModalOpen,
  setIsRelationshipModalOpen,
}) => {
  const handleClose = () => setIsRelationshipModalOpen(false);

  return (
    <Modal
      isOpen={isRelationshipModalOpen}
      onClose={handleClose}
      title="Create Relationship"
      size="md"
      disableBodyWrap
    >
      <RelationshipFormStepper handleClose={handleClose} />
    </Modal>
  );
};

export default CreateRelationshipModal;
