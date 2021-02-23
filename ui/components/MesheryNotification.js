import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import {
  Badge, 
  Drawer, 
  Tooltip, 
  Divider, 
  Typography, 
} from '@material-ui/core';
import BellIcon from '@material-ui/icons/Notifications';
import ErrorIcon from '@material-ui/icons/Error';
import ClearAllIcon from '@material-ui/icons/ClearAll';
import { withStyles } from '@material-ui/core/styles';
import amber from '@material-ui/core/colors/amber';
import { eventTypes } from '../lib/event-types';
import MesheryEventViewer from './MesheryEventViewer';
// import { bindActionCreators } from 'redux';
// import { updateSMIResults } from '../lib/store';
import dataFetch from '../lib/data-fetch';
import { withSnackbar } from 'notistack'

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
  listTop: {
    display: 'grid',
    alignItems: 'center',
    gridTemplateColumns: "2fr 6fr 2fr",
    paddingTop: theme.spacing(2),
    paddingLeft: theme.spacing(1),
    paddingRight: theme.spacing(1),
    paddingBottom: theme.spacing(2),
  },
  notificationTitle: {
    textAlign: 'center',
  },
  notifSelector: {
    display: 'flex',
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
    display: 'flex',
    justifyContent: 'flex-end'
  },
});

/**
 * getNotifications filters the notifications based on the 
 * given type and returns an array of filtered notifications
 * @param {{ 
 *  event_type: number, 
 *  summary: string, 
 *  detail: string,
 *  operation_id: string 
 * }[]} events 
 * 
 * @returns {{ 
 *  event_type: number, 
 *  summary: string, 
 *  details: string,
 *  operation_id: string 
 * }[]}
 */
function getNotifications(events, type) {
  if (!Array.isArray(events)) return [];

  if (type === "error") return events.filter(ev => ev.event_type === 2);
  if (type === "warning") return events.filter(ev => ev.event_type === 1)
  
  return events;
}

/**
 * 
 * @param {{ 
 *  event_type: number, 
 *  summary: string, 
 *  details: string,
 *  operation_id: string 
 * }[]} events 
 */
function getNotificationCount(events) {
  if (!Array.isArray(events)) return 0;

  const errorEventCount = events.filter(ev => ev.event_type === 2).length;
  const totalEventsCount = events.length;

  return errorEventCount || totalEventsCount;
}

/**
 * NotificationIcon is a wrapper react component for rendering
 * icons based on the "type" prop
 * @param {{ type: string }} props 
 */
function NotificationIcon ({ type }) {
  if (type === "error") return <ErrorIcon />

  return <BellIcon />
}

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
    displayEventType: "*",
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

  /**
   * notificationDispatcher dispatches the notifications
   * @param {number} type type of the event
   * @param {string} message message to be displayed
   */
  notificationDispatcher(type, message) {
    const self = this;
    self.props.enqueueSnackbar(message, {
      variant: eventTypes[type]?.type,
      autoHideDuration: 5000,
      action: (key) => (
        <IconButton 
          key="close" 
          aria-label="Close" 
          color="inherit" 
          onClick={() => self.props.closeSnackbar(key)}
        >
          <CloseIcon />
        </IconButton>
      ),
    });
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

      // Add the event to the state
      events.push(data);

      // Dispatch the notification
      self.notificationDispatcher(data.event_type || 0, data.summary)
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

  handleDialogClose = () => {
    this.setState({ dialogShow: false });
  };

  handleClearAllNotifications() {
    const self = this;
    return () => {
      self.setState({ events: [], open: false });
    };
  }

  handleNotifFiltering(type) {
    return () => {
      this.setState({ displayEventType: type })
    }
  }

  render() {
    const { classes } = this.props;
    const {
      open, events
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
              <Badge badgeContent={getNotificationCount(events)} color={badgeColorVariant}>
                <NotificationIcon type={badgeColorVariant} />
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
            <div>
              <div className={classes.sidelist}>
                <div className={classes.listTop}>
                  <div className={classes.notifSelector}>
                    <Tooltip title="Show all notifications">
                      <IconButton
                        color="inherit"
                        style={{ padding: '0 0.25rem 0 0' }}
                        onClick={this.handleNotifFiltering('*')}
                      >
                        <BellIcon />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Show error notifications">
                      <IconButton
                        color="inherit"
                        style={{ padding: '0 0 0 0.25rem' }}
                        onClick={this.handleNotifFiltering('error')}
                      >
                        <ErrorIcon />
                      </IconButton>
                    </Tooltip>
                  </div>
                  <div className={classes.notificationTitle}>
                    <Typography variant="subtitle2">
                      Notifications
                    </Typography>
                  </div>
                  <div className={classes.clearAllButton}>
                    <Tooltip title="Clear all notifications">
                      <IconButton
                        color="inherit"
                        style={{ padding: '0' }}
                        onClick={this.handleClearAllNotifications()}
                      >
                        <ClearAllIcon />
                      </IconButton>
                    </Tooltip>
                  </div>
                </div>
                <Divider light />
                {getNotifications(events, this.state.displayEventType).map((event, ind) => (
                  <MesheryEventViewer
                    eventVariant={event.event_type}
                    eventSummary={event.summary}
                    deleteEvent={self.deleteEvent(ind)}
                    eventDetails={event.details || "Details Unavailable"}
                  />
                ))}
              </div>
            </div>
          </Drawer>
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
)(withSnackbar(MesheryNotification)));