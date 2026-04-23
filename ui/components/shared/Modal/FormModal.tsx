import React from 'react';
import Modal, { ModalProps } from './Modal';
import { ModalFooter, ModalButtonPrimary, ModalButtonSecondary } from '@sistent/sistent';

interface FormModalProps extends ModalProps {
  onCancel?: () => void;
  onSubmit: (e?: React.FormEvent) => void;
  submitLabel?: string;
  cancelLabel?: string;
  submitting?: boolean;
}

/**
 * Standardized FormModal component.
 */
const FormModal: React.FC<FormModalProps> = ({
  children,
  onCancel,
  onSubmit,
  submitLabel = 'Submit',
  cancelLabel = 'Cancel',
  closeModal,
  submitting = false,
  ...props
}) => {
  return (
    <Modal {...props} closeModal={closeModal}>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          onSubmit(e);
        }}
      >
        {children}
        <ModalFooter>
          <ModalButtonSecondary onClick={onCancel || closeModal} disabled={submitting}>
            {cancelLabel}
          </ModalButtonSecondary>
          <ModalButtonPrimary type="submit" disabled={submitting}>
            {submitting ? 'Submitting...' : submitLabel}
          </ModalButtonPrimary>
        </ModalFooter>
      </form>
    </Modal>
  );
};

export default FormModal;
