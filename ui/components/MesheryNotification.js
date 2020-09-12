import IconButton from '@material-ui/core/IconButton';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import {
  Badge, Drawer, Tooltip, Divider, Typography, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button,
} from '@material-ui/core';
import NotificationsIcon from '@material-ui/icons/Notifications';
import ClearAllIcon from '@material-ui/icons/ClearAll';
import { withStyles } from '@material-ui/core/styles';
import classNames from 'classnames';
import amber from '@material-ui/core/colors/amber';
import { eventTypes } from '../lib/event-types';
import MesheryEventViewer from './MesheryEventViewer';
// import { bindActionCreators } from 'redux';
// import { updateSMIResults } from '../lib/store';
import dataFetch from '../lib/data-fetch';


const styles = (theme) => ({
  sidelist: {
    width: 350,
  },
  notificationButton: {
    height: '100%',
  },
  notificationDrawer: {
    backgroundColor: '#FFFFFF',
  },
  notificationTitle: {
    textAlign: 'center',
    paddingTop: theme.spacing(2),
    paddingLeft: theme.spacing(3),
    paddingBottom: theme.spacing(2),
  },
  icon: {
    fontSize: 20,
  },
  iconVariant: {
    opacity: 0.9,
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1) * 3 / 4,
  },
  error: {
    backgroundColor: theme.palette.error.dark,
  },
  info: {
    backgroundColor: theme.palette.primary.dark,
  },
  warning: {
    backgroundColor: amber[700],
  },
  message: {
    display: 'flex',
    // alignItems: 'center',
  },
  clearAllButton: {
    position: 'fixed',
    top: theme.spacing(1),
    right: theme.spacing(1),
  },
});

class MesheryNotification extends React.Component {
  state = {
    events: [],
    open: false,
    dialogShow: false,
    k8sConfig: {
      inClusterConfig: false,
      clusterConfigured: false,
      contextName: '',
    },
    meshAdapters: [],
    createStream: false,
  }

  handleToggle = () => {
    this.setState((state) => ({ open: !state.open }));
  };

  handleClose() {
    const self = this;
    return (event) => {
      if (self.anchorEl.contains(event.target)) {
        return;
      }
      self.setState({ open: false });
    };
  }

  static getDerivedStateFromProps(props, state) {
    if (JSON.stringify(props.k8sConfig) !== JSON.stringify(state.k8sConfig)
        || JSON.stringify(props.meshAdapters) !== JSON.stringify(state.meshAdapters)) {
      return {
        createStream: true,
        k8sConfig: props.k8sConfig,
        meshAdapters: props.meshAdapters,
      };
    }
    return null;
  }

  componentDidUpdate() {
    const { createStream, k8sConfig, meshAdapters } = this.state;
    if (!k8sConfig.clusterConfigured || meshAdapters.length === 0) {
      this.closeEventStream();
    }
    if (createStream && k8sConfig.clusterConfigured && typeof meshAdapters !== 'undefined' && meshAdapters.length > 0) {
      this.startEventStream();
    }
  }

  async startEventStream() {
    this.closeEventStream();
    this.eventStream = new EventSource('/api/events');
    this.eventStream.onmessage = this.handleEvents();
    this.eventStream.onerror = this.handleError();
    this.setState({ createStream: false });
  }

  handleEvents() {
    const self = this;
    return (e) => {
      const { events } = this.state;
      const data = JSON.parse(e.data);
      events.push(data);
      //Temperory Hack
      // if(data.summary==="Smi conformance test completed successfully"){
      //   self.props.updateSMIResults({smi_result: data,});
      //   console.log("HandleEvents",{smi_result: data,});
      // }
      self.setState({ events });
    };
  }

  handleError() {
    const self = this;
    return () => {
      self.closeEventStream();
      // check if server is available
      dataFetch('/api/user', { credentials: 'same-origin' }, () => {
        // attempting to reestablish connection
        // setTimeout(() => function() {
        self.startEventStream();
        // }, 2000);
      }, () => {
        // do nothing here
      });
    };
  }

  closeEventStream() {
    if (this.eventStream && this.eventStream.close) {
      this.eventStream.close();
    }
  }

  deleteEvent = (ind) => () => {
    const { events } = this.state;
    if (events[ind]) {
      events.splice(ind, 1);
    }
    this.setState({ events, dialogShow: false });
  }

  clickEvent = (event, ind) => () => {
    const { events } = this.state;
    let fInd = -1;
    events.forEach((ev, i) => {
      if (ev.event_type === event.event_type && ev.summary === event.summary && ev.details === event.details) {
        fInd = i;
      }
    });
    if (fInd === ind) {
      this.setState({
        open: true, dialogShow: true, ev: event, ind,
      });
    }
  }

  handleDialogClose = () => {
    this.setState({ dialogShow: false });
  };

  handleClearAllNotifications() {
    const self = this;
    return () => {
      self.setState({ events: [], open: false });
    };
  }

  viewEventDetails = () => {
    const { classes } = this.props;
    const { ev, ind, dialogShow } = this.state;
    if (ev && typeof ind !== 'undefined') {
      // console.log(`decided icon class: ${JSON.stringify(eventTypes[ev.event_type]?eventTypes[ev.event_type].icon:eventTypes[0].icon)}`);
      const Icon = eventTypes[ev.event_type] ? eventTypes[ev.event_type].icon : eventTypes[0].icon;
      return (
        <Dialog
          fullWidth
          maxWidth="md"
          open={dialogShow}
          onClose={this.handleDialogClose}
          aria-labelledby="event-dialog-title"
        >
          <DialogTitle id="event-dialog-title">
            <span id="client-snackbar" className={classes.message}>
              <Icon className={classNames(classes.icon, classes.iconVariant)} fontSize="large" />
              {ev.summary}
            </span>
          </DialogTitle>
          <Divider light variant="fullWidth" />
          <DialogContent>
            <DialogContentText>
              {ev.details && ev.details.split('\n').map((det) => (
                <div>{det}</div>
              ))}
            </DialogContentText>
          </DialogContent>
          <Divider light variant="fullWidth" />
          <DialogActions>
            <Button onClick={this.deleteEvent(ind)} color="secondary" variant="outlined">
              Dismiss
            </Button>
            <Button onClick={this.handleDialogClose} color="primary" variant="outlined">
              Close
            </Button>
          </DialogActions>
        </Dialog>
      );
    }
    return null;
  }

  render() {
    const { classes } = this.props;
    const {
      open, events, ev, ind,
    } = this.state;
    const self = this;

    let toolTipMsg = `There are ${events.length} events`;
    switch (events.length) {
      case 0:
        toolTipMsg = 'There are no events';
        break;
      case 1:
        toolTipMsg = 'There is 1 event';
        break;
    }
    let badgeColorVariant = 'default';
    events.forEach((eev) => {
      if (eventTypes[eev.event_type] && eventTypes[eev.event_type].type === 'error') {
        badgeColorVariant = 'error';
      }
    });

    return (
      <div>
        <NoSsr>
          <Tooltip title={toolTipMsg}>
            <IconButton
              className={classes.notificationButton}
              buttonRef={(node) => {
                this.anchorEl = node;
              }}
              color="inherit"
              onClick={this.handleToggle}
            >
              <Badge badgeContent={events.length} color={badgeColorVariant}>
                <NotificationsIcon />
              </Badge>
            </IconButton>
          </Tooltip>

          <Drawer
            anchor="right"
            open={open}
            onClose={this.handleClose()}
            classes={{
              paper: classes.notificationDrawer,
            }}
          >
            <div
              tabIndex={0}
              role="button"
            // onClick={this.handleClose()}
            // onKeyDown={this.handleClose()}
            >
              <div className={classes.sidelist}>
                <div className={classes.notificationTitle}>
                  <Typography variant="subtitle2">
                    Notifications
                  </Typography>
                  <Tooltip title="Clear all notifications">
                    <IconButton
                      className={classes.clearAllButton}
                      color="inherit"
                      onClick={this.handleClearAllNotifications()}
                    >
                      <ClearAllIcon />
                    </IconButton>
                  </Tooltip>
                </div>
                <Divider light />
                {events && events.map((event, ind) => (
                  <MesheryEventViewer
                    eventVariant={event.event_type}
                    eventSummary={event.summary}
                    deleteEvent={self.deleteEvent(ind)}
                    onClick={self.clickEvent(event, ind)}
                  />
                ))}
              </div>
            </div>
          </Drawer>
          {this.viewEventDetails(ev, ind)}
        </NoSsr>
      </div>
    );
  }
}

// const mapDispatchToProps = (dispatch) => ({
//   updateSMIResults: bindActionCreators(updateSMIResults, dispatch),
// });

const mapStateToProps = (state) => {
  const k8sConfig = state.get('k8sConfig').toJS();
  const meshAdapters = state.get('meshAdapters').toJS();
  return { k8sConfig, meshAdapters };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  null,
  // mapDispatchToProps,
)(MesheryNotification));