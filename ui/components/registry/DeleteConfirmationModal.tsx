import React from 'react';
import {
  Box,
  Modal,
  ModalBody,
  ModalFooter,
  ModalButtonDanger,
  ModalButtonSecondary,
  Typography,
} from '@sistent/sistent';

type DeleteConfirmationModalProps = {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  registrantName: string;
  modelCount: number;
  isDeleting: boolean;
};

const DeleteConfirmationModal = ({
  open,
  onClose,
  onConfirm,
  registrantName,
  modelCount,
  isDeleting,
}: DeleteConfirmationModalProps) => {
  return (
    <Modal
      open={open}
      closeModal={onClose}
      title={`Delete all models registered by "${registrantName}"?`}
    >
      <ModalBody>
        <Typography variant="body2" sx={{ my: 1 }}>
          Are you sure you want to delete{' '}
          {modelCount === 1 ? 'this model' : `all ${modelCount} models`} for this registrant? This
          will delete all associated components, relationships, and policies from the database.
        </Typography>
      </ModalBody>
      <ModalFooter variant="filled">
        <Box sx={{ width: '100%', display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
          <ModalButtonSecondary onClick={onClose} disabled={isDeleting}>
            Cancel
          </ModalButtonSecondary>
          <ModalButtonDanger onClick={onConfirm} disabled={isDeleting}>
            {isDeleting ? 'Deleting...' : 'Delete All Models'}
          </ModalButtonDanger>
        </Box>
      </ModalFooter>
    </Modal>
  );
};

export default DeleteConfirmationModal;
