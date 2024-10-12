import { Button } from '@material-ui/core';
import React from 'react';

/**
 * ErrorBoundary is a React component that catches JavaScript errors in its child components and renders a fallback UI when an error occurs.
 * It should be used as a wrapper around components that might throw errors.
 * @deprecated use error boundary from sistent instead
 */
const ErrorBoundary = ({ children }) => {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState(null);

  const resetErrorBoundary = () => {
    setHasError(false);
    setError(null);
  };

  useEffect(() => {
    if (hasError && error) {
      // Log the error, send to analytics, etc.
      console.error(error);
    }
  }, [hasError, error]);

  // React Error boundary behavior
  const ErrorCatcher = ({ children }) => {
    try {
      return children;
    } catch (err) {
      setHasError(true);
      setError(err);
      return null;
    }
  };

  if (hasError && error) {
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

  return <ErrorCatcher>{children}</ErrorCatcher>;
};

export default ErrorBoundary;
