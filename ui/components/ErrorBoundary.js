import React from 'react';
import { ErrorBoundary as SistentErrorBoundary } from '@layer5/sistent';

  useEffect(() => {
    if (hasError && error) {
      // Log the error, send to analytics, etc.
      console.error(error);
    }
  }, [hasError, error]);

function ErrorBoundary({ children }) {
  return <SistentErrorBoundary onError={handleError}>{children}</SistentErrorBoundary>;
}

export default ErrorBoundary;
