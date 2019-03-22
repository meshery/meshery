import NoSsr from '@material-ui/core/NoSsr';
// import { Snackbar, Icon, Button, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@material-ui/core';
import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { eventTypes } from '../lib/event-types';

const styles = theme => ({
    event: {
        margin: theme.spacing(1),
    }
})

class MesheryEventViewer extends React.Component {

  state = {
    dialogShow: false,
  }

//   handleToggle = () => {
//     this.setState(state => ({ dialogShow: !state.dialogShow }));
//   };

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    this.props.deleteEvent();
    // this.setState({ eventShow: false });
  };

  handleSnackbarClick = () => {
    // this.setState({dialogShow: true});
    this.props.onClick();
  }

  render() {
    const {classes, eventVariant, eventSummary, onClick, eventDetail} = this.props;
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
                    key={`event_-_${eventVariant}`}
                    variant={eventTypes[eventVariant]?eventTypes[eventVariant].type:eventTypes[0].type}
                    message={eventSummary}
                    onClose={this.handleSnackbarClose}
                    // onClick={this.handleSnackbarClick}
                    onClick={onClick}
                    className={classes.event}
                    />
                {/* </Snackbar> */}

                
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
    onClick: PropTypes.func.isRequired,
};

export default withStyles(styles)(MesheryEventViewer);