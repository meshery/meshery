import React, { useState, useCallback } from 'react';
import { Button } from '@material-ui/core';
import { ErrorBoundary as SistentErrorBoundary } from 'sistent';

/**
 * Custom ErrorBoundary component to catch errors and render fallback UI.
 * @deprecated use error boundary from sistent instead
 */
const ErrorBoundary = ({ children }) => {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState(null);

  // Handle error and update state
  const handleError = useCallback((error) => {
    setHasError(true);
    setError(error);
  }, []);

  // Reset error boundary state
  const resetErrorBoundary = () => {
    setHasError(false);
    setError(null);
  };

  // Render fallback UI if an error is present
  if (hasError) {
    return (
      <div className="alert alert-danger">
        <p>Couldn&apos;t open form. Encountered the following error:</p>
        <pre>{error.message}</pre>
        <Button color="primary" variant="contained" onClick={resetErrorBoundary}>
          Refresh Form
        </Button>
      </div>
    );
  }

  // Use Sistent's ErrorBoundary to handle errors in children
  return (
    <SistentErrorBoundary onError={handleError}>
      {children}
    </SistentErrorBoundary>
  );
};

export default ErrorBoundary;
