import React from 'react';
import { ErrorBoundary as SistentErrorBoundary } from '@layer5/sistent';

function handleError(error) {
  console.error('Error in Spaces Preferences Component', error);
}

function ErrorBoundary({ children }) {
  return (
    <SistentErrorBoundary onError={handleError}>
      {children}
    </SistentErrorBoundary>
  );
}

export default ErrorBoundary;
