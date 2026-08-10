/**
 * Create-model modal.
 *
 * Thin wrapper around the shared `Modal` primitive that mounts the URL-based
 * model creation stepper (`UrlStepper`). The stepper renders its own
 * `ModalBody` and `ModalFooter`, so we pass `disableBodyWrap` to keep the
 * shared `Modal` from nesting an extra `ModalBody`.
 *
 * Migrated to the shared modal primitives as part of Phase 5.b.6 (#18754).
 */
import { FC } from 'react';
import { Modal } from '@/components/shared/Modal';
import UrlStepper from './Stepper/UrlStepper';

export interface CreateModelModalProps {
  isCreateModalOpen: boolean;
  setIsCreateModalOpen: (_open: boolean) => void;
}

const CreateModelModal: FC<CreateModelModalProps> = ({
  isCreateModalOpen,
  setIsCreateModalOpen,
}) => {
  const handleClose = () => setIsCreateModalOpen(false);

  return (
    <Modal
      isOpen={isCreateModalOpen}
      onClose={handleClose}
      title="Create Model"
      size="sm"
      disableBodyWrap
      sx={{ zIndex: 1500 }}
    >
      <UrlStepper handleClose={handleClose} />
    </Modal>
  );
};

export default CreateModelModal;
