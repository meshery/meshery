import React, { useEffect, useState } from "react";
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
import axios from "axios";
import classNames from "classnames";
import Notification from "./notification";
import { store } from "../../store";
import { useStyles } from "./notificationCenter.style";
import { clearEvents, setEvents, setEventsSummary } from "../../store/slices/events";




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
  useLoadEventsSummary()
  const { count_by_severity_level,  total_count } = useSelector((state) => state.events.summary);
  const classes = useStyles()
  const getSeverityCount = (severity) => {
    return count_by_severity_level.find((item) => item.severity === severity)?.count || 0
  }

  const archivedCount = count_by_severity_level
    .reduce((acc, item) => acc + item.count, 0) - total_count
  console.log("count_by_severity_level", count_by_severity_level,total_count)
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

const useLoadEvents = (filters, page) => {

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [hasMore, setHasMore] = useState(true)
  const events = useSelector((state) => state.events.events);
  const dispatch = useDispatch()
  useEffect(() => {
    setLoading(true)
    const parsedFilters = {}
    Object.keys(filters).forEach((key) => {
      if (filters[key]) {
        parsedFilters[key] = JSON.stringify(filters[key], (_key, value) => (value instanceof Set ? [...value] : value))
      }
    })
    axios.get(`/api/v2/events`,
      {
        params: {
          ...parsedFilters,
          page: page,
          page_size: 15
        }
      }
    ).then(({ data }) => {
      if (data.events.length === 0) {
        setHasMore(false)
        return
      }
      dispatch(setEvents([...events, ...data.events]))
    }).catch((err) => {
      setError(err)
    }).finally(() => {
      setLoading(false)
    })

  }, [page, filters])

  const reset = () => {
    dispatch(clearEvents())
    setHasMore(true)
  }

  return {
    loading,
    error,
    hasMore,
    reset,
  }
}

const useLoadEventsSummary = () => {


  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const dispatch = useDispatch()
  useEffect(() => {
    setLoading(true)
    axios.get(`/api/v2/events?page=$1&page_size=15`).then(({ data }) => {
      dispatch(setEventsSummary({
        count_by_severity_level: data.count_by_severity_level,
        total_count: data.total_count
      }))
    }).catch((err) => {
      setError(err)
    }).finally(() => {
      setLoading(false)
    })

  }, [])

  return {
    loading,
    error,
  }
}

const MesheryNotification = (props) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [filters, setFilters] = useState({})
  const [page, setPage] = useState(1)

  const events = useSelector((state) => state.events.events);
  const loadMore = () => {
    setPage(page => page + 1)
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

  const { loading, hasMore, reset } = useLoadEvents(filters, page)

  const loader = React.useRef(null);
  const handleObserver = React.useCallback((entries) => {
    const target = entries[0];
    if (target.isIntersecting && !loading && hasMore) {
      loadMore();
    }
  }, [loading]);


  useEffect(() => {
    const option = {
      root: null,
      rootMargin: "20px",
      threshold: 0
    };
    const observer = new IntersectionObserver(handleObserver, option);
    if (loader.current) observer.observe(loader.current);
  }, [handleObserver]);

  const handleFilter = (filters) => {
    reset()
    setFilters(filters)
    setPage(1)
  }

  const value = useSelector((state) => state.events.value);
  console.log("value", value)

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
                  {events.map((event) => <Notification key={event.id} event={event} />)}
                  {loading && <div>Loading...</div>}
                  {!loading && <button onClick={loadMore}>LoadMore</button>}
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