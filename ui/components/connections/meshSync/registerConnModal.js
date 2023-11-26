import React from 'react';
import { DialogContent, Dialog } from '@material-ui/core';

import ConnectionWizard from '../../Connect/index.js';
import theme from '../../../themes/app.js';

const RegisterConnectionModal = (props) => {
  const { registerConnectionModalOpen, handleRegisterConnectionModalClose } = props;

  return (
    <div style={{ marginBottom: '1rem' }}>
      <Dialog
        open={registerConnectionModalOpen}
        onClose={handleRegisterConnectionModalClose}
        aria-labelledby="form-dialog-title"
        maxWidth="md"
        style={{ zIndex: 9999 }}
        PaperProps={{
          style: { borderRadius: 30 },
        }}
      >
        <DialogContent
          style={{
            padding: '0 1.5rem 1.5rem',
            borderRadius: '28px',
            border: `6px solid ${theme.palette.secondary.success}`,
          }}
        >
          <ConnectionWizard />
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default RegisterConnectionModal;
