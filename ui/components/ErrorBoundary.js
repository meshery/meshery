import { Button } from '@material-ui/core';
import React from 'react';

/**
 * ErrorBoundary is a React component that catches JavaScript errors in its child components and renders a fallback UI when an error occurs.
 * It should be used as a wrapper around components that might throw errors.
 * @deprecated use error boundary from sistent instead
 */
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  /** Update state so the next render will show the fallback UI. */
  static getDerivedStateFromError(error) {
    return { hasError: true, error: error };
  }

  resetErrorBoundary = () => {
    this.setState({ hasError: false, error: null });
  };

  /** You can render any custom fallback UI */
  render() {
    if (this.state.hasError) {
      return (
        <div className="alert alert-danger">
          <p>Couldn&apos;t open form. Encountered the following error:</p>
          <pre>{this.state.error.message}</pre>
          <Button color="primary" variant="contained" onClick={this.resetErrorBoundary}>
            Refresh Form
          </Button>
        </div>
      );
    }
    return this.props.children;
  }
}

export default ErrorBoundary;
