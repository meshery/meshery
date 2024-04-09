import Button from '@mui/material/Button';
import React, { useState } from 'react';

/**
 * ErrorBoundary is a React component that catches JavaScript errors in its child components and renders a fallback UI when an error occurs.
 * It should be used as a wrapper around components that might throw errors.
 */
const ErrorBoundary = ({ children }) => {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState(null);

  /** Update state so the next render will show the fallback UI. */
  const getDerivedStateFromError = (error) => {
    setHasError(true);
    setError(error);
  };

  const resetErrorBoundary = () => {
    setHasError(false);
    setError(null);
  };

  /** You can render any custom fallback UI */
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

  return children;
};

export default ErrorBoundary;
