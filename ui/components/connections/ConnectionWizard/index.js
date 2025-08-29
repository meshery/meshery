import React, { useState } from 'react';
import { Box, Modal, ModalBody } from '@sistent/sistent';
import { ConnectionWizardStepper } from './ConnectionWizardStepper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';


const ConnectionWizard = ({
  open,
  onClose,
  connectionType = 'kubernetes',
  onConnectionComplete,
  initialData = null,
}) => {
  const [wizardData, setWizardData] = useState({
    connectionType,
    uploadedFile: null,
    discoveredConnections: [],
    selectedConnections: {},
    preflightResults: {},
    executionResults: {},
    ...initialData,
  });

  const { notify } = useNotification();

  const handleWizardComplete = (results) => {
    notify({
      message: 'Connection wizard completed successfully!',
      event_type: EVENT_TYPES.SUCCESS,
    });

    if (onConnectionComplete) {
      onConnectionComplete(results);
    }

    onClose();
  };

  const handleWizardCancel = () => {
    // Clean up any pending operations
    setWizardData({
      connectionType,
      uploadedFile: null,
      discoveredConnections: [],
      selectedConnections: {},
      preflightResults: {},
      executionResults: {},
    });

    onClose();
  };

  if (!open) return null;

  return (
    <Modal
      open={open}
      closeModal={handleWizardCancel}
      maxWidth="lg"
      aria-labelledby="connection-wizard-title"
    >
      <ModalBody>
        <Box sx={{ minHeight: '70vh', width: '100%' }}>
          <ConnectionWizardStepper
            connectionType={connectionType}
            wizardData={wizardData}
            setWizardData={setWizardData}
            onComplete={handleWizardComplete}
            onCancel={handleWizardCancel}
          />
        </Box>
      </ModalBody>
    </Modal>
  );
};

export default ConnectionWizard;
export { ConnectionWizardStepper } from './ConnectionWizardStepper';
