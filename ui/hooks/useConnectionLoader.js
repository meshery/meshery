import { useState, useEffect } from 'react';
import { useSnackbar } from 'notistack';

const useConnectionLoader = () => {
  const [isVerifying, setIsVerifying] = useState(false);
  const [verificationMessage, setVerificationMessage] = useState('');
  const { enqueueSnackbar } = useSnackbar();

  useEffect(() => {
    // Listen for connection verification events
    const handleConnectionEvent = (event) => {
      if (event.category === 'connection' && event.action === 'verifying') {
        setIsVerifying(true);
        setVerificationMessage(event.description || 'Verifying connection...');
      } else if (event.category === 'connection' && event.action === 'create') {
        setIsVerifying(false);
        setVerificationMessage('');
        if (event.severity === 'informational') {
          enqueueSnackbar('Connection created successfully!', { variant: 'success' });
        } else if (event.severity === 'error') {
          enqueueSnackbar('Failed to create connection', { variant: 'error' });
        }
      } else if (event.category === 'connection' && event.action === 'update') {
        setIsVerifying(false);
        setVerificationMessage('');
        if (event.severity === 'informational') {
          enqueueSnackbar('Connection updated successfully!', { variant: 'success' });
        } else if (event.severity === 'error') {
          enqueueSnackbar('Failed to update connection', { variant: 'error' });
        }
      }
    };

    // Subscribe to events (assuming you have an event system)
    // This would need to be integrated with your existing event system
    if (window.mesheryEventBus) {
      window.mesheryEventBus.subscribe('connection', handleConnectionEvent);
    }

    return () => {
      if (window.mesheryEventBus) {
        window.mesheryEventBus.unsubscribe('connection', handleConnectionEvent);
      }
    };
  }, [enqueueSnackbar]);

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