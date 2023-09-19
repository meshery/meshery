import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

const styles = () => ({});

class MesheryChartDialog extends React.Component {
  render() {
    const { classes, open, title, handleClose, content } = this.props;

    return (
      <React.Fragment>
        <Dialog
          fullWidth
          maxWidth="md"
          open={open}
          onClose={handleClose}
          aria-labelledby="chart-dialog-title"
        >
          <DialogTitle id="chart-dialog-title">
            {title && title.length ? title : 'Comparison'}
          </DialogTitle>
          <DialogContent>
            <DialogContentText className={classes.dialogContent}>{content}</DialogContentText>
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
  classes: PropTypes.object.isRequired,
  open: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  content: PropTypes.node.isRequired,
};

export default withStyles(styles)(MesheryChartDialog);
