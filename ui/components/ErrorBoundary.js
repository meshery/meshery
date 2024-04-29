import { Button } from '@material-ui/core';
import React, { useState, useEffect } from 'react';

const ErrorBoundary = ({ children }) => {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const handleErrors = (error) => {
      setHasError(true);
      setError(error);
    };
  }, []);

  const resetErrorBoundary = () => {
    setHasError(false);
    setError(null);
  };

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

