import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import MesheryPerformanceComponent from '../MesheryPerformanceComponent';

export default function PerformanceDialog({open, handleClose, urlForModal}) {
  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="performance-dialog-title"
      aria-describedby="performance-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{"Performance Test"}</DialogTitle>
      <DialogContent>
        <DialogContentText id="performance-component">
          <MesheryPerformanceComponent
            customURL={urlForModal}
          />
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
}