import React from 'react';
import { DialogContent, Dialog } from '@material-ui/core';

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
      >
        <DialogContent style={{ padding: '0 1.5rem 1.5rem' }}>Register Connection</DialogContent>
      </Dialog>
    </div>
  );
};

export default RegisterConnectionModal;
