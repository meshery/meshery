import React, { useEffect, useRef, useState } from 'react';
import IconButton from '@material-ui/core/IconButton';
import { Provider, useDispatch, useSelector } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import {
  Drawer,
  Divider,
  ClickAwayListener,
  Typography,
  alpha,
  Chip,
  Button,
  Badge,
  CircularProgress,
  Box,
  useTheme,
  Tooltip,
} from '@material-ui/core';
import Filter from './filter';
import BellIcon from '../../assets/icons/BellIcon.js';
import { iconMedium } from '../../css/icons.styles';
import {
  NOTIFICATION_CENTER_TOGGLE_CLASS,
  SEVERITY,
  SEVERITY_STYLE,
  STATUS,
  STATUS_STYLE,
} from './constants';
import classNames from 'classnames';
import Notification from './notification';
import { store } from '../../store';
import { useNavNotificationIconStyles, useStyles } from './notificationCenter.style';
import {
  closeNotificationCenter,
  loadEvents,
  loadNextPage,
  selectEvents,
  toggleNotificationCenter,
} from '../../store/slices/events';
import {
  useGetEventsSummaryQuery,
  useLazyGetEventsQuery,
} from '../../rtk-query/notificationCenter';
import _ from 'lodash';
import DoneIcon from '../../assets/icons/DoneIcon';
import { ErrorBoundary, withErrorBoundary } from '../General/ErrorBoundary';
import { hasClass } from '../../utils/Elements';

const getSeverityCount = (count_by_severity_level, severity) => {
  return count_by_severity_level.find((item) => item.severity === severity)?.count || 0;
};

const EmptyState = () => {
  const theme = useTheme().palette.secondary;
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '1rem',
        marginY: '5rem',
      }}
    >
      <DoneIcon height="10rem" width="8rem" fill={theme.icon2} />
      <Typography variant="h6" style={{ margin: 'auto', color: theme.icon2 }}>
        {' '}
        No notifications to show{' '}
      </Typography>
    </Box>
  );
};

const NavbarNotificationIcon = withErrorBoundary(() => {
  const { data } = useGetEventsSummaryQuery();
  const count_by_severity_level = data?.count_by_severity_level || [];

  const currentTopSeverity =
    getSeverityCount(count_by_severity_level, SEVERITY.ERROR) > 0
      ? SEVERITY.ERROR
      : getSeverityCount(count_by_severity_level, SEVERITY.WARNING) > 0
      ? SEVERITY.WARNING
      : null;
  const currentSeverityStyle = currentTopSeverity ? SEVERITY_STYLE[currentTopSeverity] : null;
  const topSeverityCount = getSeverityCount(count_by_severity_level, currentTopSeverity);
  const classes = useNavNotificationIconStyles({
    badgeColor: currentSeverityStyle?.color,
  });
  if (currentTopSeverity) {
    return (
      <Badge id="notification-badge" badgeContent={topSeverityCount} className={classes.root}>
        <currentSeverityStyle.icon {...iconMedium} fill="#fff" />
      </Badge>
    );
  }
  return <BellIcon className={iconMedium} fill="#fff" />;
});

const NotificationCountChip = withErrorBoundary(
  ({ classes, notificationStyle, count, type, handleClick }) => {
    const chipStyles = {
      fill: notificationStyle?.color,
      height: '20px',
      width: '20px',
    };
    count = Number(count).toLocaleString('en', { useGrouping: true });
    return (
      <Tooltip title={type} placement="bottom">
        <Button style={{ backgroundColor: alpha(chipStyles.fill, 0.2) }} onClick={handleClick}>
          <div className={classes.severityChip}>
            {<notificationStyle.icon {...chipStyles} />}
            <span>{count}</span>
          </div>
        </Button>
      </Tooltip>
    );
  },
);

const Header = withErrorBoundary(({ handleFilter, handleClose }) => {
  const { data } = useGetEventsSummaryQuery();
  const { count_by_severity_level, total_count } = data || {
    count_by_severity_level: [],
    total_count: 0,
  };
  const classes = useStyles();
  const onClickSeverity = (severity) => {
    handleFilter({
      severity: [severity],
    });
  };

  const onClickStatus = (status) => {
    handleFilter({
      status: status,
    });
  };

  const unreadCount =
    total_count - count_by_severity_level.reduce((acc, item) => acc + item.count, 0);
  return (
    <div className={classNames(classes.container, classes.header)}>
      <div className={classes.title}>
        <div className={classes.titleBellIcon} onClick={handleClose}>
          <BellIcon height="30" width="30" fill="#fff" />
        </div>
        <Typography variant="h6"> Notifications</Typography>
      </div>
      <div className={classes.severityChips}>
        {Object.values(SEVERITY).map((severity) => (
          <NotificationCountChip
            key={severity}
            classes={classes}
            handleClick={() => onClickSeverity(severity)}
            notificationStyle={SEVERITY_STYLE[severity]}
            type={`Unread ${severity}(s)`}
            count={getSeverityCount(count_by_severity_level, severity)}
          />
        ))}
        <NotificationCountChip
          classes={classes}
          notificationStyle={STATUS_STYLE[STATUS.READ]}
          handleClick={() => onClickStatus(STATUS.READ)}
          type={STATUS.READ}
          count={unreadCount}
        />
      </div>
    </div>
  );
});

const Loading = () => {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'center' }}>
      <CircularProgress />
    </Box>
  );
};

const EventsView = withErrorBoundary(({ handleLoadNextPage, isFetching, hasMore }) => {
  const events = useSelector(selectEvents);
  // const page = useSelector((state) => state.events.current_view.page);

  const lastEventRef = useRef(null);
  const intersectionObserver = useRef(
    new IntersectionObserver(
      (entries) => {
        if (isFetching && !hasMore) {
          return;
        }
        const firstEntry = entries[0];
        if (firstEntry.isIntersecting) {
          handleLoadNextPage();
        }
      },
      { threshold: 1 },
    ),
  );

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
      {events.map((event, idx) => (
        <div key={event.id + idx}>
          <Notification event_id={event.id} />
        </div>
      ))}

      {events.length === 0 && <EmptyState />}

      <div ref={lastEventRef}></div>
      {isFetching && hasMore && <Loading />}
    </>
  );
});

const CurrentFilterView = withErrorBoundary(({ handleFilter }) => {
  const currentFilters = useSelector((state) => state.events.current_view.filters);
  const onDelete = (key, value) => {
    const newFilters = {
      ...currentFilters,
      [key]:
        typeof currentFilters[key] === 'string'
          ? null
          : currentFilters[key].filter((item) => item !== value),
    };
    handleFilter(newFilters);
  };

  const Chips = ({ type, value }) => {
    if (typeof value === 'string') {
      return (
        <Chip label={value} style={{ paddingTop: '0rem' }} onDelete={() => onDelete(type, value)} />
      );
    }

    if (_.isArray(value) && value.length > 0) {
      return (
        <div style={{ display: 'flex', gap: '0.2rem' }}>
          {value.map((item) => (
            <Chip key={item} label={item} onDelete={() => onDelete(type, item)} />
          ))}
        </div>
      );
    }

    return null;
  };

  return (
    <div>
      {Object.entries(currentFilters).map(([key, value]) => {
        if (value && value?.length > 0) {
          return (
            <div
              key={key}
              style={{
                display: 'flex',
                gap: '0.3rem',
                alignItems: 'center',
                marginLeft: '1rem',
                paddingTop: '.35rem',
              }}
            >
              <Typography variant="subtitle2" style={{ textTransform: 'capitalize' }}>
                {' '}
                {key}:
              </Typography>
              <Chips value={value} type={key} />
            </div>
          );
        }
      })}
    </div>
  );
});

const MesheryNotification = () => {
  const [anchorEl, setAnchorEl] = useState(null);
  const dispatch = useDispatch();
  const isNotificationCenterOpen = useSelector((state) => state.events.isNotificationCenterOpen);
  const [fetchEvents, { isFetching }] = useLazyGetEventsQuery();
  const hasMore = useSelector((state) => state.events.current_view.has_more);

  useEffect(() => {
    dispatch(
      loadEvents(fetchEvents, 1, {
        status: STATUS.UNREAD,
      }),
    );
  }, []);

  const loadMore = () => {
    dispatch(loadNextPage(fetchEvents));
  };

  const handleToggle = () => {
    dispatch(toggleNotificationCenter());
  };

  const handleClose = () => {
    if (!isNotificationCenterOpen) {
      return;
    }
    dispatch(closeNotificationCenter());
    setAnchorEl(null);
  };
  const classes = useStyles();
  // const { showFullNotificationCenter } = props;
  const open = Boolean(anchorEl) || isNotificationCenterOpen;

  const handleFilter = (filters) => {
    dispatch(loadEvents(fetchEvents, 1, filters));
  };
  const drawerRef = useRef();
  const toggleButtonRef = useRef();

  const clickwayHandler = (e) => {
    // checks if event has occured/bubbled up from clicking inside notificationcenter or on the bell icon
    if (drawerRef.current.contains(e.target) || toggleButtonRef.current.contains(e.target)) {
      return;
    }
    // check for element with toggle class
    if (hasClass(e.target, NOTIFICATION_CENTER_TOGGLE_CLASS)) {
      return;
    }
    // check for svg icon (special case) , not checking the toggle class as it is not added to svg
    if (e.target?.className?.baseVal?.includes('MuiSvgIcon')) {
      return;
    }
    handleClose();
  };

  return (
    <>
      <div ref={toggleButtonRef}>
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
          <NavbarNotificationIcon />
        </IconButton>
      </div>

      <ClickAwayListener onClickAway={clickwayHandler}>
        <Drawer
          anchor="right"
          variant="persistent"
          open={open}
          ref={drawerRef}
          classes={{
            paper: classes.notificationDrawer,
            paperAnchorRight: isNotificationCenterOpen ? classes.fullView : classes.peekView,
          }}
        >
          <div>
            <div>
              <div className={classes.sidelist}>
                <Header handleFilter={handleFilter} handleClose={handleClose}></Header>
                <Divider light />
                <div className={classes.container}>
                  <Filter handleFilter={handleFilter}></Filter>
                  <CurrentFilterView handleFilter={handleFilter} />
                  <EventsView
                    handleLoadNextPage={loadMore}
                    isFetching={isFetching}
                    hasMore={hasMore}
                  />
                </div>
              </div>
            </div>
          </div>
        </Drawer>
      </ClickAwayListener>
    </>
  );
};

// const mapDispatchToProps = (dispatch) => ({
//   updateEvents: bindActionCreators(updateEvents, dispatch),
//   toggleOpen: bindActionCreators(toggleNotificationCenter, dispatch),
//   loadEventsFromPersistence: bindActionCreators(loadEventsFromPersistence, dispatch),
// });
//
// const mapStateToProps = (state) => {
//   const events = state.get("events");
//   return {
//     user: state.get("user"),
//     events: events.toJS(),
//     openEventId: state.get("notificationCenter").get("openEventId"),
//     showFullNotificationCenter: state.get("notificationCenter").get("showFullNotificationCenter"),
//   };
// };
const NotificationCenter = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in NotificationCenter', e)}
      >
        <Provider store={store}>
          <MesheryNotification {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

export default NotificationCenter;
