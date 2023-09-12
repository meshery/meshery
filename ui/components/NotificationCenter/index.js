import React, { useEffect, useRef, useState } from "react";
import IconButton from "@material-ui/core/IconButton";
import { Provider, connect, useDispatch, useSelector } from "react-redux";
import NoSsr from "@material-ui/core/NoSsr";
import { Drawer, Tooltip, Divider, ClickAwayListener, Typography, alpha } from "@material-ui/core";
import Filter from "./filter";
import BellIcon from "../../assets/icons/BellIcon.js"
import { loadEventsFromPersistence, toggleNotificationCenter, updateEvents } from "../../lib/store";
import { iconMedium } from "../../css/icons.styles";
import { bindActionCreators } from "redux";
import { SEVERITY, SEVERITY_STYLE, STATUS, STATUS_STYLE } from "./constants";
import classNames from "classnames";
import Notification from "./notification";
import { store } from "../../store";
import { useStyles } from "./notificationCenter.style";
import { loadEvents, loadNextPage, selectEvents } from "../../store/slices/events";
import { useGetEventsSummaryQuery, useLazyGetEventsQuery } from "../../rtk-query/notificationCenter";




const NotificationCountChip = ({ classes, notificationStyle, count }) => {
  const chipStyles = {
    fill: notificationStyle.color,
    height: "20px",
    width: "20px",
  }
  count = Number(count).toLocaleString('en', { useGrouping: true })
  return (
    <div className={classes.severityChip} style={{ backgroundColor: alpha(chipStyles.fill, 0.20) }} >
      {<notificationStyle.icon {...chipStyles} />}
      {count}
    </div>
  )
}

const Header = () => {

  const { data } = useGetEventsSummaryQuery();
  const { count_by_severity_level, total_count } = data || {
    count_by_severity_level: [],
    total_count: 0
  }
  const classes = useStyles()
  const getSeverityCount = (severity) => {
    return count_by_severity_level.find((item) => item.severity === severity)?.count || 0
  }

  const archivedCount = count_by_severity_level
    .reduce((acc, item) => acc + item.count, 0) - total_count
  return (
    <div className={classNames(classes.container, classes.header)}>
      <div className={classes.title}>
        <div className={classes.titleBellIcon}>
          <BellIcon height="30" width="30" fill="#fff" />
        </div>
        <Typography variant="h6"> Notifications</Typography>
      </div>
      <div className={classes.severityChips}>
        {Object.values(SEVERITY).map(severity => (
          <NotificationCountChip key={severity} classes={classes}
            notificationStyle={SEVERITY_STYLE[severity]}
            count={getSeverityCount(severity)} />)
        )}
        <NotificationCountChip classes={classes} notificationStyle={STATUS_STYLE[STATUS.UNREAD]} count={archivedCount} />

      </div>
    </div>
  )
}



const EventsView = ({ handleLoadNextPage, isLoading, hasMore }) => {
  const events = useSelector(selectEvents)
  // const page = useSelector((state) => state.events.current_view.page);

  const lastEventRef = useRef(null)
  const intersectionObserver = useRef(new IntersectionObserver((entries) => {
    const firstEntry = entries[0]
    if (firstEntry.isIntersecting) {
      handleLoadNextPage()
    }
  }, { threshold: 1 }))

  useEffect(() => {
    const currentObserver = intersectionObserver.current;
    if (lastEventRef.current) {
      currentObserver.observe(lastEventRef.current);
    }
    return () => {
      if (lastEventRef.current) {
        currentObserver.unobserve(lastEventRef.current);
      }
    };
  }, [lastEventRef.current]);

  return (
    <>
      {events.map((event, idx) => <div key={event.id+idx}  >
        <Notification event={event} />
      </div>)}
      {!isLoading && hasMore &&
        <div ref={lastEventRef} ></div>}
      {isLoading && <div>Loading...</div>}
    </>
  )
}


const MesheryNotification = (props) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const dispatch = useDispatch()

  const [fetchEvents, { isLoading }] = useLazyGetEventsQuery()
  const hasMore = useSelector((state) => state.events.current_view.has_more);

  useEffect(() => {
    dispatch(loadEvents(fetchEvents, 1, {
      status: STATUS.UNREAD,
    }))
  }, [])

  const loadMore = () => {
    dispatch(loadNextPage(fetchEvents))
  }

  const handleToggle = () => {
    props.toggleOpen();
  };

  const handleClose = () => {
    if (!props.showFullNotificationCenter) {
      return;
    }
    props.toggleOpen();
    setAnchorEl(null);
  };
  const classes = useStyles()
  const { showFullNotificationCenter } = props;
  const open = Boolean(anchorEl) || showFullNotificationCenter;


  const handleFilter = (filters) => {
    dispatch(loadEvents(fetchEvents, 1, filters))
  }

  return (
    <NoSsr>
      <div>
        <Tooltip title={"errors"}>
          <IconButton
            id="notification-button"
            className={classes.notificationButton}
            color="inherit"
            onClick={handleToggle}
            onMouseOver={(e) => {
              e.preventDefault();
              setAnchorEl(e.currentTarget);
            }}
            onMouseLeave={(e) => {
              e.preventDefault();
              setAnchorEl(null);
            }}
          >
            <BellIcon className={iconMedium} fill="#fff" />
            {/* <Badge id="notification-badge" badgeContent={getNotificationCount(events)} color={newNotificationsType}>
              </Badge> */}
          </IconButton>
        </Tooltip>
      </div>

      <ClickAwayListener
        onClickAway={(e) => {
          if (
            e.target.className.baseVal !== "" &&
            e.target.className.baseVal !== "MuiSvgIcon-root" &&
            (typeof e.target.className === "string" ? !e.target.className?.includes("MesheryNotification") : null)
          ) {
            handleClose();
          }
        }}
      >
        <Drawer
          anchor="right"
          variant="persistent"
          open={open}
          classes={{
            paper: classes.notificationDrawer,
            paperAnchorRight: showFullNotificationCenter ? classes.fullView : classes.peekView,
          }}
        >
          <div>
            <div>
              <div className={classes.sidelist}>
                <Header ></Header>
                <Divider light />
                <div className={classes.container}>
                  <Filter handleFilter={handleFilter}  ></Filter>
                  <EventsView handleLoadNextPage={loadMore} isLoading={isLoading} hasMore={hasMore} />
                </div>
              </div>
            </div>
          </div>
        </Drawer>
      </ClickAwayListener>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateEvents: bindActionCreators(updateEvents, dispatch),
  toggleOpen: bindActionCreators(toggleNotificationCenter, dispatch),
  loadEventsFromPersistence: bindActionCreators(loadEventsFromPersistence, dispatch),
});

const mapStateToProps = (state) => {
  const events = state.get("events");
  return {
    user: state.get("user"),
    events: events.toJS(),
    openEventId: state.get("notificationCenter").get("openEventId"),
    showFullNotificationCenter: state.get("notificationCenter").get("showFullNotificationCenter"),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)((props) => {

  return (
    <>
      <Provider store={store} >
        <MesheryNotification {...props} />
      </Provider >
    </>
  )

});