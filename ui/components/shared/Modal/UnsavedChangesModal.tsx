import React from 'react';
import Modal from './Modal';
import { CheckCircleIcon, DeleteIcon, Typography, useTheme, Button } from '@sistent/sistent';

interface UnsavedChangesModalProps {
  open: boolean;
  onClose: () => void;
  onDiscard: () => void;
  onSave: () => void | Promise<void>;
}

const UnsavedChangesModal = ({ open, onClose, onDiscard, onSave }: UnsavedChangesModalProps) => {
  const theme = useTheme();

  return (
    <Modal
      open={open}
      closeModal={onClose}
      title="Unsaved dashboard layout changes"
      maxWidth="xs"
      footer={
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            gap: '1rem',
            width: '100%',
          }}
        >
          <Button
            type="button"
            variant="contained"
            onClick={onDiscard}
            startIcon={<DeleteIcon fill={theme.palette.common.white} />}
            sx={{
              background: theme.palette.background.error.default,
              '&:hover': {
                background: theme.palette.background.error.dark,
              },
            }}
          >
            Discard Changes
          </Button>
          <Button
            type="button"
            variant="contained"
            color="primary"
            onClick={onSave}
            startIcon={<CheckCircleIcon />}
          >
            Save Changes
          </Button>
        </div>
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
