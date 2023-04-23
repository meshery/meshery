import NoSsr from '@material-ui/core/NoSsr';
import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import { eventTypes } from '../lib/event-types';
const styles = (theme) => ({ event : { margin : theme.spacing(0.5, 1), }, });

class MesheryEventViewer extends React.Component {
  state = { dialogShow : false, }

  handleSnackbarClose = (_, reason) => {
    if (reason === 'clickaway') return

    this.props.deleteEvent();
  }

  render() {
    const {
      classes, eventVariant, eventSummary, eventDetails, eventCause, eventRemediation, eventErrorCode, componentType, componentName, expand,
    } = this.props;

    return (
      <NoSsr>
        <React.Fragment>
          <MesherySnackbarWrapper
            key={`event_-_${eventVariant}`}
            variant={eventTypes[eventVariant]
              ? eventTypes[eventVariant].type
              : eventTypes[0].type}
            message={eventSummary}
            details={eventDetails}
            onClose={this.handleSnackbarClose}
            className={classes.event}
            cause={eventCause}
            remedy={eventRemediation}
            errorCode={eventErrorCode}
            componentType={componentType}
            componentName={componentName}
            expand={expand}
          />
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryEventViewer.propTypes = {
  classes : PropTypes.object.isRequired,
  eventVariant : PropTypes.oneOf([0, 1, 2]).isRequired,
  eventSummary : PropTypes.string.isRequired,
  eventDetails : PropTypes.string.isRequired,
  deleteEvent : PropTypes.func.isRequired,
  expand : PropTypes.string.isRequired,
};

export default withStyles(styles)(MesheryEventViewer);
