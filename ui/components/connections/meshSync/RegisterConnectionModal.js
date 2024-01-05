import React from 'react';
import { DialogContent, Dialog } from '@material-ui/core';

import theme from '../../../themes/app.js';
import CustomizedSteppers from './Stepper/index.js';
import dataFetch from '../../../lib/data-fetch.js';

export const cancelConnectionRegister = (id) => {
  dataFetch(
    '/api/integrations/connections/register',
    {
      method: 'DELETE',
      credentials: 'include',
      body: JSON.stringify({
        id: id,
      }),
    },
    (result) => {
      console.log(result);
    },
  );
};

const RegisterConnectionModal = ({
  openRegistrationModal,
  connectionData,
  handleRegistrationModalClose,
}) => {
  const [sharedData, setSharedData] = React.useState(null);
  const handleClose = () => {
    handleRegistrationModalClose();
    cancelConnectionRegister(sharedData?.connection?.id);
  };

  const handleRegistrationComplete = () => {
    dataFetch(`/api/system/meshsync/resources/${connectionData?.resourceID}`, {
      method: 'DELETE',
      credentials: 'include',
    });
  };

  return (
    <div style={{ marginBottom: '1rem' }}>
      <Dialog
        open={openRegistrationModal}
        onClose={handleClose}
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
          <CustomizedSteppers
            formConnectionIdRef
            onClose={handleClose}
            connectionData={connectionData}
            sharedData={sharedData}
            setSharedData={setSharedData}
            handleRegistrationComplete={handleRegistrationComplete}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default RegisterConnectionModal;
