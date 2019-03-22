import classNames from 'classnames';
import NoSsr from '@material-ui/core/NoSsr';
import { Snackbar, Icon, Button, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@material-ui/core';
import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import WarningIcon from '@material-ui/icons/Warning';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';
import amber from '@material-ui/core/colors/amber';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

const eventTypes = {
    0: {
        icon: InfoIcon,
        type: 'info',
    },
    1: {
        icon: WarningIcon,
        type: 'warn',
    },
    2: {
        icon: ErrorIcon,
        type: 'error',
    },
}

const styles = theme => ({
    error: {
      backgroundColor: theme.palette.error.dark,
    },
    info: {
      backgroundColor: theme.palette.primary.dark,
    },
    warning: {
      backgroundColor: amber[700],
    },
    icon: {
      fontSize: 20,
    },
    iconVariant: {
        opacity: 0.9,
        marginRight: theme.spacing(1),
    }
})

class MesheryEventViewer extends React.Component {

  state = {
    eventShow: true,
    dialogShow: false,
  }

//   handleToggle = () => {
//     this.setState(state => ({ dialogShow: !state.dialogShow }));
//   };

  handleDialogClose = event => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ dialogShow: false });
  };

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    this.props.deleteEvent(this.props.eventIndex);
    this.setState({ eventShow: false });
  };

  handleSnackbarClick = () => {
    this.setState({dialogShow: true});
  }

  render() {
    const {classes, eventVariant, eventSummary, eventDetail} = this.props;
    const { eventShow, dialogShow } = this.state;

    return (
        <NoSsr>
            <React.Fragment>
                {/* <Snackbar
                    open={eventShow}
                    onClose={this.handleSnackbarClose}
                    onClick={this.handleSnackbarClick}
                    > */}
                    <MesherySnackbarWrapper 
                    variant={eventTypes[eventVariant]?eventTypes[eventVariant].type:eventTypes[0].type}
                    message={eventSummary}
                    onClose={this.handleSnackbarClose}
                    onClick={this.handleSnackbarClick}
                    />
                {/* </Snackbar> */}

                <Dialog
                    fullWidth
                    maxWidth='xs'
                    open={dialogShow}
                    onClose={this.handleDialogClose}
                    aria-labelledby="event-dialog-title"
                    >
                    <DialogTitle id="event-dialog-title">
                        <Icon className={classNames(classes.icon, eventTypes[eventVariant]?eventTypes[eventVariant].icon:eventTypes[0].icon)} /> 
                        {eventSummary}
                    </DialogTitle>
                    <DialogContent>
                        <DialogContentText>
                        {eventDetail}
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={this.handleSnackbarClose} color="secondary">
                        Dismiss
                        </Button>
                        <Button onClick={this.handleDialogClose} color="primary">
                        Close
                        </Button>
                    </DialogActions>
                </Dialog>
            </React.Fragment>
        </NoSsr>
    )
  }
}

MesheryEventViewer.propTypes = {
    classes: PropTypes.object.isRequired,
    eventVariant: PropTypes.oneOf([0,1,2]).isRequired,
    eventSummary: PropTypes.string.isRequired,
    eventDetail: PropTypes.string.isRequired,
    eventIndex: PropTypes.number.isRequired,
    deleteEvent: PropTypes.func.isRequired,
};

export default withStyles(styles)(MesheryEventViewer);