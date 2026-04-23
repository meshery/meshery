import React from 'react';
import Modal, { ModalProps } from './Modal';
import {
  ModalFooter,
  ModalButtonPrimary,
  ModalButtonSecondary,
  Typography,
  ModalBody,
} from '@sistent/sistent';

interface ConfirmModalProps extends Omit<ModalProps, 'children'> {
  message: string | React.ReactNode;
  onConfirm: () => void;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: 'primary' | 'danger';
}

/**
 * Standardized ConfirmModal component.
 */
const ConfirmModal: React.FC<ConfirmModalProps> = ({
  message,
  onConfirm,
  confirmLabel = 'Confirm',
  cancelLabel = 'Cancel',
  closeModal,
  variant = 'primary',
  ...props
}) => {
  return (
    <Modal {...props} closeModal={closeModal}>
      <ModalBody>
        <Typography variant="body1">{message}</Typography>
      </ModalBody>
      <ModalFooter>
        <ModalButtonSecondary onClick={closeModal}>{cancelLabel}</ModalButtonSecondary>
        <ModalButtonPrimary
          onClick={() => {
            onConfirm();
            closeModal();
          }}
          sx={{
            backgroundColor: variant === 'danger' ? 'red' : undefined,
            '&:hover': {
              backgroundColor: variant === 'danger' ? 'darkred' : undefined,
            },
          }}
        >
          {confirmLabel}
        </ModalButtonPrimary>
      </ModalFooter>
    </Modal>
  );
};

export default ConfirmModal;
