import React from "react";
import IconButton from '@material-ui/core/IconButton';
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
import ClearIcon from "../assets/icons/ClearIcon";
import ErrorIcon from '@material-ui/icons/Error';
import { withStyles } from '@material-ui/core/styles';
import amber from '@material-ui/core/colors/amber';
import {  EVENT_TYPES, eventTypes } from '../lib/event-types';
import MesheryEventViewer from './MesheryEventViewer';
import dataFetch from '../lib/data-fetch';
import { bindActionCreators } from "redux";
import { toggleNotificationCenter, updateEvents } from "../lib/store";
import { iconMedium } from "../css/icons.styles";
import { cursorNotAllowed } from "../css/disableComponent.styles";
import { v4 } from "uuid";
import moment from "moment";
import { withNotify } from "../utils/hooks/useNotification";
import { parseJson } from "./ConnectionWizard/helpers/jsonParser";

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
    transition : '0.3s ease-in-out !important'
  },
  peekView : {
    right : "-26.1rem",
    transition : '0.3s ease-in-out !important'
  },
  tabs : {
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tab : {
    "&.Mui-selected" : {
      color : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
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

  if (type === "error") return events.filter(ev => getEventType(ev).type == EVENT_TYPES.ERROR.type);
  if (type === "warning") return events.filter(ev => getEventType(ev).type == EVENT_TYPES.WARNING.type)
  if (type === "success") return events.filter(ev =>
    ((getEventType(ev).type == EVENT_TYPES.SUCCESS.type)
                                      || getEventType(ev).type == EVENT_TYPES.DEFAULT.type))

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

  const errorEventCount = events.filter(ev => getEventType(ev) == EVENT_TYPES.ERROR ).length;
  const totalEventsCount = events.length;
  return errorEventCount || totalEventsCount;
}


const  getEventType = (event) => {
  // checks if an event_type is as cardinal (0 , 1 ,2 ) or as a event_type object
  // return the event_type object
  let eventVariant = event.event_type
  eventVariant = typeof eventVariant  == "number" ? eventTypes[eventVariant] : eventVariant
  return eventVariant  ? eventVariant :  EVENT_TYPES.DEFAULT
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

//TODO: Convert To functional Compoent
class MesheryNotification extends React.Component {
  state = {
    open : false,
    dialogShow : false,
    displayEventType : "*",
    tabValue : 0,
    anchorEl : false,
    // showFullNotificationCenter : false,
  }

  handleToggle = () => {
    // this.setState({ showFullNotificationCenter : !this.state.showFullNotificationCenter })
    this.props.toggleOpen()
  };

  handleClose = () => {
    if (! this.props.showFullNotificationCenter) {
      return
    }
    this.setState({ anchorEl : false });
    this.props.toggleOpen()

  }

  /**
   * notificationDispatcher dispatches the notifications
   * @param {number} type type of the event
   * @param {string} message message to be displayed
   */

  // const {notify} = useNotification()

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
    const self = this
    return (e) => {
      const data = parseJson(e.data);
      const event = {
        ...data,
        event_type : getEventType(data),
        timestamp : data.timestamp || moment.utc().valueOf() ,
        id : data.id || v4() ,
      }
      self.props.notify({
        message : event.summary ,
        event_type : event.event_type,
        details : event.details,
        customEvent : event
      })
    }
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
      self.handleClose();
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

  render() {
    const { classes, events ,showFullNotificationCenter } = this.props;
    const { anchorEl, show } = this.state;
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
      if ( getEventType(eev).type == EVENT_TYPES.ERROR.type ) {
        badgeColorVariant = 'error';
      }
    });

    return (
      <NoSsr>
        <div style={ show ? cursorNotAllowed : {}}>
          <Tooltip title={toolTipMsg}>
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
          </Tooltip>
        </div>

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
                    className={classes.tabs}
                    textColor="primary"
                    variant="fullWidth"
                  >
                    <Tab label="All" className={classes.tab} onClick={this.handleNotifFiltering('*')} style={{ minWidth : "15%" }} />
                    <Tab label="Error" className={classes.tab} onClick={this.handleNotifFiltering('error')} style={{ minWidth : "15%" }} />
                    <Tab label="Warning" className={classes.tab} onClick={this.handleNotifFiltering('warning')} style={{ minWidth : "15%" }} />
                    <Tab label="Success" className={classes.tab} onClick={this.handleNotifFiltering('success')} style={{ minWidth : "15%" }} />
                  </Tabs>
                  {getNotifications(this.props.events, this.state.displayEventType).map((event, ind) => (
                    <MesheryEventViewer
                      key={ind}
                      eventVariant={getEventType(event)}
                      eventSummary={event.summary}
                      deleteEvent={self.deleteEvent(ind)}
                      eventDetails={event.details || "Details Unavailable"}
                      eventCause={event.probable_cause}
                      eventRemediation={event.suggested_remediation}
                      eventErrorCode={event.error_code}
                      componentType={event.component}
                      componentName={event.component_name}
                      eventTimestamp={event.timestamp}
                      expand={(this.props.openEventId && this.props.openEventId === ( event.id) ) ? true : false}
                    />
                  ))}
                </div>
              </div>
            </div>
          </Drawer>
        </ClickAwayListener>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updateEvents : bindActionCreators(updateEvents, dispatch),
  toggleOpen : bindActionCreators(toggleNotificationCenter,dispatch)
})

// const mapDispatchToProps = (dispatch) => ({
//   updateSMIResults: bindActionCreators(updateSMIResults, dispatch),
// });

const mapStateToProps = (state) => {
  //TODO: Sort While Storing
  const events = state.get('events')
  return {
    events : events.toJS(),
    openEventId : state.get("notificationCenter").get("openEventId"),
    showFullNotificationCenter : state.get("notificationCenter").get("showFullNotificationCenter")
  };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withNotify(MesheryNotification)));