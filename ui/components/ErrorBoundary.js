import React from 'react';
import { Button, ErrorBoundary as SistentErrorBoundary } from '@layer5/sistent';

/**
 * ErrorBoundary is a React component that catches JavaScript errors in its child components and renders a fallback UI when an error occurs.
 * It should be used as a wrapper around components that might throw errors.
 */

function FallbackComponent({ error, resetErrorBoundary }) {
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

/**
 * Error handler function that logs the error.
 */
function handleError(error) {
  console.error('Error in Spaces Preferences Component', error);
}

function ErrorBoundary({ children }) {
  return (
    <SistentErrorBoundary FallbackComponent={FallbackComponent} onError={handleError}>
      {children}
    </SistentErrorBoundary>
  );
}

export default ErrorBoundary;