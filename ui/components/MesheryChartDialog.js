import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@mui/styles';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

const styles = () => ({});

class MesheryChartDialog extends React.Component {

  render() {
    const {
      classes, open, title, handleClose, content,
    } = this.props;

    return (
      <React.Fragment>
        <Dialog
          fullWidth
          maxWidth="md"
          open={open}
          onClose={handleClose}
          aria-labelledby="chart-dialog-title"
        >
          <DialogTitle id="chart-dialog-title">{title && title.length
            ? title
            : 'Comparison'}</DialogTitle>
          <DialogContent>
            <DialogContentText className={classes.dialogContent}>
              {content}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">
              Close
            </Button>
          </DialogActions>
        </Dialog>
      </React.Fragment>
    );
  }
}

MesheryChartDialog.propTypes = {
  classes : PropTypes.object.isRequired,
  open : PropTypes.bool.isRequired,
  handleClose : PropTypes.func.isRequired,
  content : PropTypes.node.isRequired,
};

export default withStyles(styles)(MesheryChartDialog);
