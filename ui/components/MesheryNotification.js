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
  ClickAwayListener,
} from '@material-ui/core';
import BellIcon from '@material-ui/icons/Notifications';
import ClearIcon from "../assets/icons/ClearIcon";
import ErrorIcon from '@material-ui/icons/Error';
import { withStyles } from '@material-ui/core/styles';
import amber from '@material-ui/core/colors/amber';
import {  EVENT_TYPES, NOTIFICATION_STATUS, SERVER_EVENT_TYPES } from '../lib/event-types';
import Notification from "./NotificationCenter/Notification"
import dataFetch from '../lib/data-fetch';
import { bindActionCreators } from "redux";
import { loadEventsFromPersistence, toggleNotificationCenter, updateEvents } from "../lib/store";
import { iconMedium } from "../css/icons.styles";
import { cursorNotAllowed } from "../css/disableComponent.styles";
import { v4 } from "uuid";
import moment from "moment";
import { withNotify } from "../utils/hooks/useNotification";

const styles = (theme) => ({
  sidelist : { width : "35rem", },
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
    right : "-32.1rem",
    transition : '0.3s ease-in-out !important'
  },
  tabs : {

    "&.MuiTabs-flexContainer" : {
      display : "flex",
      justifyContent : "space-around",
    },
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tab : {
    padding : "0px",
    margin : "0px",
    "&.Mui-selected" : {
      color : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },

  notification : {
    margin : theme.spacing(0.5, 1)
  }
});




export const NOTIFICATION_FILTERS = {
  ALL : "all",
  ERROR : EVENT_TYPES.ERROR.type ,
  SUCCESS : EVENT_TYPES.SUCCESS.type ,
  WARNING : EVENT_TYPES.WARNING.type,
  HISTORY : "history" ,
}

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
function getNotifications(events, filter) {
  if (!Array.isArray(events)) return [];
  if (filter == NOTIFICATION_FILTERS.HISTORY) return events.filter(ev => ev.status === NOTIFICATION_STATUS.VIEWED)

  if (filter == NOTIFICATION_FILTERS.ALL) return events.filter(ev => ev.status !== NOTIFICATION_STATUS.VIEWED)

  if (filter === NOTIFICATION_FILTERS.SUCCESS) {
    return events.filter(ev => {
      const ev_type = getEventType(ev).type
      return (ev_type == (EVENT_TYPES.SUCCESS.type || ev_type == EVENT_TYPES.INFO.type) &&  ev.status !== NOTIFICATION_STATUS.VIEWED )
    })
  }
  return events.filter(ev => getEventType(ev).type == filter && ev.status !== NOTIFICATION_STATUS.VIEWED )

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

  const errorEventCount = getNotifications(events,NOTIFICATION_FILTERS.ERROR).length;
  const totalEventsCount = getNotifications(events,NOTIFICATION_FILTERS.ALL).length;
  return errorEventCount || totalEventsCount;
}


const  getEventType = (event) => {
  // checks if an event_type is as cardinal (0 , 1 ,2 ) or as a event_type object
  // return the event_type object
  let eventVariant = event.event_type
  eventVariant = typeof eventVariant  == "number" ? SERVER_EVENT_TYPES[eventVariant] : eventVariant
  return eventVariant  ? eventVariant :  EVENT_TYPES.INFO
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

const notificationBadgeTooltipMessage = (events) => {
  const total_unread = getNotifications(events,NOTIFICATION_FILTERS.ALL).length
  if (total_unread) {
    return `${total_unread} new notifications`
  }
  return `No new notifications`
}

const TabLabel = ({ filterType,events }) => {

  const notifCount = getNotifications(events,filterType).length
  return (
    <div style={{ margin : 0 , padding : 0 ,textTransform : "capitalize" , display : "flex" , gap : "0.3rem" }}>
      {notifCount ? `${filterType} (${notifCount})` : filterType }
    </div>
  )
}

//TODO: Convert To functional Compoent
class MesheryNotification extends React.Component {
  state = {
    open : false,
    dialogShow : false,
    displayEventType : NOTIFICATION_FILTERS.ALL,
    tabValue : 0,
    anchorEl : false,
  }

  handleToggle = () => {
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

  componentDidUpdate() {
    if (this.props.user.get("user_id") && this.props.events.length == 0) {
      this.props.loadEventsFromPersistence()
    }
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
      const data = JSON.parse(e.data);
      const event = {
        ...data,
        status : NOTIFICATION_STATUS.NEW,
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

  deleteEvent = (id) => {
    const { events, updateEvents } = this.props;
    const newEvents = events.filter(ev => ev.id !== id)
    updateEvents({ events : newEvents })
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

  handleTabChange = (_event, newTabValue) => {
    this.setState({ tabValue : newTabValue })
  }

  handleBellButtonClick = () => {
    this.setState({
      tabValue : 0,
      displayEventType : '*'
    })
  }

  markAsRead = (event) => {
    const events = this.props.events.filter(ev => ev.id !== event.id)
    events.push({
      ...event,
      status : NOTIFICATION_STATUS.VIEWED
    })
    this.props.updateEvents({
      events
    })
  }


  render() {
    const { classes, events ,showFullNotificationCenter  } = this.props;
    const { anchorEl, show } = this.state;
    let open = Boolean(anchorEl);
    if (showFullNotificationCenter) {
      open = showFullNotificationCenter;
    }

    const newErrors =  getNotifications(events,NOTIFICATION_FILTERS.ERROR)
    const newNotificationsType = newErrors.length > 0 ? "error" : "default"

    return (
      <NoSsr>
        <div style={ show ? cursorNotAllowed : {}}>
          <Tooltip title={notificationBadgeTooltipMessage(events)}>
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
              <Badge id="notification-badge" badgeContent={getNotificationCount(events)} color={newNotificationsType}>
                <NotificationIcon  style={iconMedium}  type={newNotificationsType}   />
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

                    {Object.values(NOTIFICATION_FILTERS).map(filter  =>
                      <Tab label={<TabLabel filterType={filter} events={events}/>}
                        className={classes.tab}
                        onClick={this.handleNotifFiltering(filter)}
                        style={{ minWidth : "8%" }} />
                    )}
                  </Tabs>
                  {getNotifications(this.props.events, this.state.displayEventType).map((event) => (
                    <Notification
                      className = {classes.notification}
                      key={event.id}
                      event={event}
                      onDeleteEvent = {() => this.deleteEvent(event.id)}
                      onMarkAsRead = {() => this.markAsRead(event) }
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
  toggleOpen : bindActionCreators(toggleNotificationCenter,dispatch),
  loadEventsFromPersistence : bindActionCreators(loadEventsFromPersistence,dispatch)
})

const mapStateToProps = (state) => {
  const events = state.get('events')
  return {
    user : state.get("user"),
    events : events.toJS(),
    openEventId : state.get("notificationCenter").get("openEventId"),
    showFullNotificationCenter : state.get("notificationCenter").get("showFullNotificationCenter")
  };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withNotify(MesheryNotification)));