import React from 'react';

import {
  CheckCircleIcon,
  DeleteIcon,
  Modal,
  ModalBody,
  ModalButtonPrimary,
  ModalFooter,
  Typography,
  useTheme,
} from '@sistent/sistent';

const UnsavedChangesModal = ({ open, onClose, onDiscard, onSave }) => {
  const theme = useTheme();

  return (
    <Modal
      open={open}
      closeModal={onClose}
      title="Unsaved changes"
      maxWidth="xs"
      sx={{
        '& .MuiDialogTitle-root': {
          padding: '8px 16px',
        },
      }}
    >
      <ModalBody>
        <Typography variant="body1">
          You have unsaved changes. Leaving now will discard your edits.
        </Typography>
      </ModalBody>
      <ModalFooter variant="filled">
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            gap: '1rem',
            width: '100%',
          }}
        >
          <ModalButtonPrimary
            type="button"
            onClick={onDiscard}
            startIcon={<DeleteIcon fill={theme.palette.common.white} />}
            style={{
              background: theme.palette.background.error.default,
            }}
          >
            Discard Changes
          </ModalButtonPrimary>
          <ModalButtonPrimary type="button" onClick={onSave} startIcon={<CheckCircleIcon />}>
            Save Changes
          </ModalButtonPrimary>
        </div>
      </ModalFooter>
    </Modal>
  );
};

export default UnsavedChangesModal;
