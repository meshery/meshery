/**
 * Shared ConfirmModal primitive.
 *
 * Thin wrapper around the base `Modal` for yes/no confirmation flows. Use
 * this for "Are you sure you want to delete X?" style prompts so the action
 * affordances and copy stay consistent.
 *
 * The primary action defaults to `Confirm`; pass `variant="danger"` for
 * destructive operations so the confirm button styles match Sistent's
 * destructive button treatment.
 */
import { FC, ReactNode } from 'react';
import {
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalButtonDanger,
  Typography,
} from '@sistent/sistent';
import { Modal, ModalSize } from './Modal';

export type ConfirmModalVariant = 'primary' | 'danger';

export interface ConfirmModalProps {
  /** Whether the confirmation modal is visible. */
  isOpen: boolean;
  /** Called when the user dismisses, cancels, or completes the action. */
  onClose: () => void;
  /** Invoked when the user clicks the confirm button. */
  onConfirm: () => void;
  /** Modal title; e.g. "Delete environment". */
  title: string;
  /** Body copy or rich content describing the consequences of the action. */
  message: ReactNode;
  /** Confirm button label. Defaults to `Confirm`. */
  confirmText?: string;
  /** Cancel button label. Defaults to `Cancel`. */
  cancelText?: string;
  /**
   * Visual tone of the confirm button. Use `danger` for destructive actions
   * such as delete or undeploy.
   */
  variant?: ConfirmModalVariant;
  /** Optional icon rendered in the modal header. */
  headerIcon?: ReactNode;
  /** Disable the confirm button (e.g. while the action is in flight). */
  isConfirmDisabled?: boolean;
  /** Size token; defaults to `sm` for confirmations. */
  size?: ModalSize;
}

export const ConfirmModal: FC<ConfirmModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'primary',
  headerIcon,
  isConfirmDisabled = false,
  size = 'sm',
}) => {
  const ConfirmButton = variant === 'danger' ? ModalButtonDanger : ModalButtonPrimary;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      headerIcon={headerIcon}
      size={size}
      actions={
        <>
          <ModalButtonSecondary onClick={onClose}>{cancelText}</ModalButtonSecondary>
          <ConfirmButton onClick={onConfirm} disabled={isConfirmDisabled}>
            {confirmText}
          </ConfirmButton>
        </>
      }
    >
      {typeof message === 'string' ? <Typography variant="body1">{message}</Typography> : message}
    </Modal>
  );
};

ConfirmModal.displayName = 'ConfirmModal';

export default ConfirmModal;
