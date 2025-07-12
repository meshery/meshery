import { useState, useEffect, useCallback } from 'react';
import { useSnackbar } from 'notistack';

const useConnectionLoader = () => {
  const [isVerifying, setIsVerifying] = useState(false);
  const [verificationMessage, setVerificationMessage] = useState('');
  const { enqueueSnackbar } = useSnackbar();

  const handleConnectionEvent = useCallback((event) => {
    if (event.category !== 'connection') {
      return;
    }
    switch (event.action) {
      case 'verifying':
        setIsVerifying(true);
        setVerificationMessage(event.description || 'Verifying connection...');
        break;
      case 'create':
      case 'update': {
        setIsVerifying(false);
        setVerificationMessage('');
        const successMessage = `Connection ${event.action}d successfully!`;
        const errorMessage = `Failed to ${event.action} connection.`;
        if (event.severity === 'informational') {
          enqueueSnackbar(successMessage, { variant: 'success' });
        } else if (event.severity === 'error') {
          enqueueSnackbar(errorMessage, { variant: 'error' });
        }
        break;
      }
      default:
        // Do nothing for other actions
        break;
    }
  }, [enqueueSnackbar]);

  useEffect(() => {
    if (window.mesheryEventBus) {
      window.mesheryEventBus.subscribe('connection', handleConnectionEvent);
    }
    return () => {
      if (window.mesheryEventBus) {
        window.mesheryEventBus.unsubscribe('connection', handleConnectionEvent);
      }
    };
  }, [handleConnectionEvent]);

  const startVerification = (message = 'Verifying connection...') => {
    setIsVerifying(true);
    setVerificationMessage(message);
  };

  const stopVerification = () => {
    setIsVerifying(false);
    setVerificationMessage('');
  };

  return {
    isVerifying,
    verificationMessage,
    startVerification,
    stopVerification,
  };
};

export default useConnectionLoader; 