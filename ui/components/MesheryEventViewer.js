import NoSsr from '@material-ui/core/NoSsr';
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

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    this.props.deleteEvent();
  };

  handleSnackbarClick = () => {
    this.props.onClick();
  }

  render() {
    const {classes, eventVariant, eventSummary, onClick} = this.props;

    return (
      <NoSsr>
        <React.Fragment>
          <MesherySnackbarWrapper 
            key={`event_-_${eventVariant}`}
            variant={eventTypes[eventVariant]?eventTypes[eventVariant].type:eventTypes[0].type}
            message={eventSummary}
            onClose={this.handleSnackbarClose}
            onClick={onClick}
            className={classes.event}
          />
        </React.Fragment>
      </NoSsr>
    )
  }
}

MesheryEventViewer.propTypes = {
  classes: PropTypes.object.isRequired,
  eventVariant: PropTypes.oneOf([0,1,2]).isRequired,
  eventSummary: PropTypes.string.isRequired,
  deleteEvent: PropTypes.func.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default withStyles(styles)(MesheryEventViewer);