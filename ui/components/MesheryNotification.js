import React from "react";
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
  Tabs,
  Tab,
  ClickAwayListener
} from '@material-ui/core';
import BellIcon from '@material-ui/icons/Notifications';
import InfoIcon from '@mui/icons-material/Info';
import ClearIcon from "../assets/icons/ClearIcon";
import ErrorIcon from '@material-ui/icons/Error';
import { withStyles } from '@material-ui/core/styles';
import amber from '@material-ui/core/colors/amber';
import { eventTypes } from '../lib/event-types';
import MesheryEventViewer from './MesheryEventViewer';
import dataFetch from '../lib/data-fetch';
import { withSnackbar } from 'notistack'
import { bindActionCreators } from "redux";
import { updateEvents } from "../lib/store";
import { iconMedium } from "../css/icons.styles";

const styles = (theme) => ({
  sidelist : { width : 450, },
  notificationButton : { height : '100%', },
  notificationDrawer : {
    backgroundColor : theme.palette.secondary.sideBar,
    display : 'flex',
    flexDirection : 'column',
    justifyContent : 'space-between'
  },
  listTop : {
    display : 'grid',
    alignItems : 'center',
    gridTemplateColumns : "2fr 6fr 2fr",
    paddingTop : theme.spacing(2),
    paddingLeft : theme.spacing(1),
    paddingRight : theme.spacing(1),
    paddingBottom : theme.spacing(2),
  },
  notificationTitle : { textAlign : 'left', },
  notifSelector : { display : 'flex', },
  icon : { fontSize : 20, },
  iconVariant : {
    opacity : 0.9,
    marginRight : theme.spacing(1),
    marginTop : theme.spacing(1) * 3 / 4,
  },
  error : { backgroundColor : theme.palette.error.dark, },
  info : { backgroundColor : theme.palette.primary.dark, },
  warning : { backgroundColor : amber[700], },
  message : {
    display : 'flex',
    // alignItems: 'center',
  },
  clearAllButton : {
    display : 'flex',
    justifyContent : 'flex-end'
  },
  drawerButton : {
    padding : '0.45rem',
    margin : '0.2rem',
    backgroundColor : theme.palette.secondary.dark,
    color : '#FFFFFF',
    "&:hover" : {
      backgroundColor : '#FFFFFF',
      color : theme.palette.secondary.dark
    }
  },
  fullView : {
    right : 0,
  },
  peekView : {
    right : "-26.1rem"
  }
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
  if (type === "success") return events.filter(ev => ev.event_type === 0)

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
 * @param {{ type: string,className: string }} props
 */
function NotificationIcon({ type, className }) {
  if (type === "error") return <ErrorIcon id="error-icon" className={className} />

  return <BellIcon className={className} style={iconMedium} />
}

class MesheryNotification extends React.Component {
  state = {
    open : false,
    dialogShow : false,
    displayEventType : "*",
    tabValue : 0,
    anchorEl : false,
    showFullNotificationCenter : false,
    eventIdToOpenInNotification : undefined,
  }

  handleToggle = () => {
    this.setState({ showFullNotificationCenter : !this.state.showFullNotificationCenter })
  };

  handleClose = () => {
    this.setState({ anchorEl : false });
    this.setState({ showFullNotificationCenter : false })
  }

  /**
   * notificationDispatcher dispatches the notifications
   * @param {number} type type of the event
   * @param {string} message message to be displayed
   */
  notificationDispatcher(type, message, operation_id) {
    const self = this;
    self.props.enqueueSnackbar(message, {
      variant : eventTypes[type]?.type,
      autoHideDuration : 5000,
      action : (key) => (
        <>
          <IconButton
            key="eye"
            aria-label="eye"
            color="inherit"
            onClick={() => self.openEventInNotificationCenter(operation_id)}
          >
            <InfoIcon style={iconMedium} />
          </IconButton>
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => self.props.closeSnackbar(key)}
          >
            <CloseIcon style={iconMedium} />
          </IconButton>
        </>
      ),
    });
  }

  componentDidMount() {
    this.startEventStream();
  }

  async startEventStream() {
    this.closeEventStream();
    this.eventStream = new EventSource('/api/events');
    this.eventStream.onmessage = this.handleEvents();
    this.eventStream.onerror = this.handleError();
  }

  handleEvents() {
    const self = this;
    return (e) => {
      const { events, updateEvents } = this.props;
      const data = JSON.parse(e.data);
      // set null event field as success
      data.event_type = data.event_type || 0

      // Dispatch the notification
      self.notificationDispatcher(data.event_type, data.summary, data.operation_id)
      //Temperory Hack
      // if(data.summary==="Smi conformance test completed successfully"){
      //   self.props.updateSMIResults({smi_result: data,});
      //   console.log("HandleEvents",{smi_result: data,});
      // }
      updateEvents({ events : [...events, data] })
    };
  }

  handleError() {
    const self = this;
    return () => {
      self.closeEventStream();
      // check if server is available
      dataFetch('/api/user', { credentials : 'same-origin' }, () => {
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
    const { events, updateEvents } = this.props;
    if (events[ind]) {
      events.splice(ind, 1);
    }
    updateEvents({ events : events })
    this.setState({ dialogShow : false });
  }

  handleDialogClose = () => {
    this.setState({ dialogShow : false });
  };

  handleClearAllNotifications() {
    const self = this;
    const { updateEvents } = this.props;
    return () => {
      updateEvents({ events : [] })
      self.setState({ open : false });
    };
  }

  handleNotifFiltering(type) {
    return () => {
      this.setState({ displayEventType : type })
    }
  }

  handleTabChange = (event, newTabValue) => {
    this.setState({ tabValue : newTabValue })
  }

  handleBellButtonClick = () => {
    this.setState({
      tabValue : 0,
      displayEventType : '*'
    })
  }

  openEventInNotificationCenter = (operation_id) => {
    this.setState({
      showFullNotificationCenter : true,
      eventIdToOpenInNotification : operation_id,
    })
  }

  render() {
    const { classes, events } = this.props;
    const { anchorEl, showFullNotificationCenter } = this.state;
    const self = this;
    let open = Boolean(anchorEl);
    if (showFullNotificationCenter) {
      open = showFullNotificationCenter;
    }

    let toolTipMsg;
    if (typeof events?.length !== 'undefined') {
      toolTipMsg = `There are ${events.length} events`;
      switch (events?.length) {
        case 0:
          toolTipMsg = 'There are no events';
          break;
        case 1:
          toolTipMsg = 'There is 1 event';
          break;
      }
    } else { // takes care of case when (typeof events.length === undefined)
      toolTipMsg = 'There are no events';
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
            <>
              <IconButton
                id="notification-button"
                className={classes.notificationButton}
                buttonRef={(node) => {
                  this.anchorEl = node;
                }}
                color="inherit"

                onClick={this.handleToggle}

                onMouseOver={(e) => {
                  e.preventDefault();
                  this.setState({ anchorEl : true })
                }}

                onMouseLeave={(e) => {
                  e.preventDefault();
                  this.setState({ anchorEl : false })
                }}
              >
                <Badge id="notification-badge" badgeContent={getNotificationCount(events)} color={badgeColorVariant}>
                  <NotificationIcon  style={iconMedium}  type={badgeColorVariant} />
                </Badge>
              </IconButton>
            </>
          </Tooltip>

          <ClickAwayListener onClickAway={(e) => {
            if (e.target.className.baseVal !== "" && e.target.className.baseVal !== "MuiSvgIcon-root" &&
              ((typeof e.target.className === "string")? !e.target.className?.includes("MesheryNotification"): null)) {
              this.handleClose();
            }
          }}>
            <Drawer
              anchor="right"
              variant="persistent"
              open={open}
              classes={{
                paper : classes.notificationDrawer,
                paperAnchorRight : showFullNotificationCenter? classes.fullView : classes.peekView,
              }}
            >
              <div>
                <div>
                  <div className={classes.sidelist}>
                    <div className={classes.listTop}>
                      <div className={classes.notifSelector}>
                        <Tooltip title="Show all notifications">
                          <IconButton
                            color="inherit"
                            className={classes.drawerButton}
                            onClick={this.handleBellButtonClick}
                          >
                            <BellIcon  style={iconMedium}  />
                          </IconButton>
                        </Tooltip>
                      </div>
                      <div className={classes.notificationTitle}>
                        <Typography variant="subtitle1" align="center">
                            Notifications
                        </Typography>
                      </div>
                      <div
                        className={classes.clearAllButton}>
                        <Tooltip title="Clear all notifications">
                          <IconButton
                            color="inherit"
                            className={classes.drawerButton}
                            onClick={this.handleClearAllNotifications()}
                          >
                            <ClearIcon width={'1em'} height={'1em'} fill={'white'}/>
                          </IconButton>
                        </Tooltip>
                      </div>
                    </div>
                    <Divider light />
                    <Tabs
                      value={this.state.tabValue}
                      onChange={this.handleTabChange}
                      indicatorColor="primary"
                      textColor="primary"
                      variant="fullWidth"
                    >
                      <Tab label="All" onClick={this.handleNotifFiltering('*')} style={{ minWidth : "15%" }} />
                      <Tab label="Error" onClick={this.handleNotifFiltering('error')} style={{ minWidth : "15%" }} />
                      <Tab label="Warning" onClick={this.handleNotifFiltering('warning')} style={{ minWidth : "15%" }} />
                      <Tab label="Success" onClick={this.handleNotifFiltering('success')} style={{ minWidth : "15%" }} />
                    </Tabs>
                    {getNotifications(this.props.events, this.state.displayEventType).map((event, ind) => (
                      <MesheryEventViewer
                        key={ind}
                        eventVariant={event.event_type}
                        eventSummary={event.summary}
                        deleteEvent={self.deleteEvent(ind)}
                        eventDetails={event.details || "Details Unavailable"}
                        eventCause={event.probable_cause}
                        eventRemediation={event.suggested_remediation}
                        eventErrorCode={event.error_code}
                        componentType={event.component}
                        componentName={event.component_name}
                        expand={(this.state.eventIdToOpenInNotification && this.state.eventIdToOpenInNotification === event.operation_id) ? true : false}
                      />
                    ))}
                  </div>
                </div>
              </div>
            </Drawer>
          </ClickAwayListener>
        </NoSsr>
      </div>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updateEvents : bindActionCreators(updateEvents, dispatch)
})

// const mapDispatchToProps = (dispatch) => ({
//   updateSMIResults: bindActionCreators(updateSMIResults, dispatch),
// });

const mapStateToProps = (state) => {
  const events = state.get("events");
  return { events };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withSnackbar(MesheryNotification)));
