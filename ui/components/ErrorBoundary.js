import React from 'react';
import { Button, ErrorBoundary as SistentErrorBoundary, useTheme } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

/**
 * ErrorBoundary is a React component that catches JavaScript errors in its child components and renders a fallback UI when an error occurs.
 * It should be used as a wrapper around components that might throw errors.
 */
const ErrorBoundary = ({ children }) => {
  const theme = useTheme();

  return (
    <UsesSistent>
      <SistentErrorBoundary
        fallback={({ error, resetErrorBoundary }) => (
          <div
            className="alert alert-danger"
            style={{
              background: theme.palette.error.light,
              color: theme.palette.error.contrastText,
            }}
          >
            <p>Couldn&apos;t open form. Encountered the following error:</p>
            <pre>{error.message}</pre>
            <Button color="primary" variant="contained" onClick={resetErrorBoundary}>
              Refresh Form
            </Button>
          </div>
        )}
      >
        {children}
      </SistentErrorBoundary>
    </UsesSistent>
  );
};

export default ErrorBoundary;
