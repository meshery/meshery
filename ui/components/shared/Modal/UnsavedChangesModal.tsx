/**
 * UnsavedChangesModal — cross-cutting confirmation for the dashboard's
 * "you have unsaved layout changes" flow.
 *
 * Phase 5.b.2 migration target: composes the shared `Modal` primitive instead
 * of using Sistent's `Modal` / `ModalFooter` directly. The two action buttons
 * (Discard / Save) and their iconography are preserved so the rendered UX is
 * identical to the legacy implementation.
 *
 * The shape of the underlying choice — discard or save before navigation —
 * doesn't fit the `confirm/cancel` shape of `ConfirmModal`, so this wraps the
 * base `Modal` and threads its own button row into `actions`. Callers should
 * keep using the existing prop names (`open`, `onClose`, `onDiscard`,
 * `onSave`) which remain unchanged.
 */
import { FC } from 'react';
import { CheckCircleIcon, DeleteIcon, ModalButtonPrimary, Typography } from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { Modal } from './Modal';

const ActionsRow = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  gap: '1rem',
  width: '100%',
});

const DiscardButton = styled(ModalButtonPrimary)(({ theme }) => ({
  '&&': {
    backgroundColor: theme.palette.background.error.default,
    '&:hover': {
      backgroundColor: theme.palette.background.error.hover,
    },
  },
}));

export interface UnsavedChangesModalProps {
  /** Whether the modal is visible. */
  open: boolean;
  /** Called when the user dismisses the modal without choosing an action. */
  onClose: () => void;
  /** Invoked when the user discards unsaved changes. */
  onDiscard: () => void;
  /** Invoked when the user saves and continues. */
  onSave: () => void | Promise<void>;
}

const UnsavedChangesModal: FC<UnsavedChangesModalProps> = ({
  open,
  onClose,
  onDiscard,
  onSave,
}) => {
  const theme = useTheme();

  return (
    <Modal
      isOpen={open}
      onClose={onClose}
      title="Unsaved dashboard layout changes"
      size="sm"
      actions={
        <ActionsRow>
          <DiscardButton
            type="button"
            onClick={onDiscard}
            startIcon={<DeleteIcon fill={theme.palette.common.white} />}
          >
            Discard Changes
          </DiscardButton>
          <ModalButtonPrimary type="button" onClick={onSave} startIcon={<CheckCircleIcon />}>
            Save Changes
          </ModalButtonPrimary>
        </ActionsRow>
      }
    >
      <Typography variant="body1">
        You have unsaved changes to your dashboard layout. If you leave now, your widget arrangement
        and edits will be lost.
      </Typography>
    </Modal>
  );
};

export default UnsavedChangesModal;
