import { Button } from '@material-ui/core';
import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
import SupportForm from './support-form';
import React, { useState } from 'react';
import useStyles from '../Modals/Information/styles';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Typography,
} from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

function Fallback({ error, resetErrorBoundary }) {
  // Call resetErrorBoundary() to reset the error boundary and retry the render.

  return (
    <div role="alert">
      <h2>Uh-oh!😔 Please pardon the mesh.</h2>
      <div
        style={{
          backgroundColor: '#1E2117',
          color: '#FFFFFF',
          padding: '.85rem',
          borderRadius: '.2rem',
          marginBlock: '.5rem',
        }}
      >
        <code>{error.message}</code>
      </div>
      <Button color="primary" variant="contained" onClick={resetErrorBoundary}>
        Try again
      </Button>
      <GetHelpButton />
    </div>
  );
}

const reportError = (error, info) => {
  // This is where you'd send the error to Sentry,etc
  console.log('Error Caught Inside Boundary --reportError', error, 'Info', info);
};

export const ErrorBoundary = ({ children, ...props }) => {
  return (
    <ReactErrorBoundary FallbackComponent={Fallback} onError={reportError} {...props}>
      {children}
    </ReactErrorBoundary>
  );
};

export const withErrorBoundary = (Component, errorHandlingProps = {}) => {
  const WrappedWithErrorBoundary = (props) => (
    <ErrorBoundary {...errorHandlingProps}>
      <Component {...props} />
    </ErrorBoundary>
  );

  return WrappedWithErrorBoundary;
};

export const withSuppressedErrorBoundary = (Component) => {
  const WrappedWithErrorBoundary = (props) => (
    <ErrorBoundary FallbackComponent={() => null}>
      <Component {...props} />
    </ErrorBoundary>
  );

  return WrappedWithErrorBoundary;
};

export const GetHelpButton = () => {
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleGetHelpClick = () => {
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
  };

  const classes = useStyles();

  return (
    <>
      <Button variant="contained" color="primary" onClick={handleGetHelpClick}>
        Get Help
      </Button>

      <Dialog open={isDialogOpen} onClose={handleCloseDialog}>
        <DialogTitle className={classes.dialogTitle}>
          <Typography className={classes.textHeader} variant="h6">
            Get Help
          </Typography>
          <IconButton
            aria-label="close"
            onClick={handleCloseDialog}
            component="button"
            style={{
              color: '#FFFFFF',
            }}
          >
            <CloseIcon className={classes.closing} />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <SupportForm />
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
};
