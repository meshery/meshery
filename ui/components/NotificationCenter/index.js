import React, { useContext, useEffect, useRef, useState } from 'react';
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
  Checkbox,
  Collapse,
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
import { DarkBackdrop, useNavNotificationIconStyles, useStyles } from './notificationCenter.style';
import {
  closeNotificationCenter,
  loadEvents,
  loadNextPage,
  selectAreAllEventsChecked,
  selectCheckedEvents,
  selectEvents,
  selectSeverity,
  toggleNotificationCenter,
  updateCheckAllEvents,
} from '../../store/slices/events';
import {
  useDeleteEventsMutation,
  useGetEventsSummaryQuery,
  useLazyGetEventsQuery,
  useUpdateEventsMutation,
} from '../../rtk-query/notificationCenter';
import _ from 'lodash';
import DoneIcon from '../../assets/icons/DoneIcon';
import { ErrorBoundary, withErrorBoundary } from '../General/ErrorBoundary';
import { hasClass } from '../../utils/Elements';
import ReadIcon from '../../assets/icons/ReadIcon';
import UnreadIcon from '../../assets/icons/UnreadIcon';
import DeleteIcon from '../../assets/icons/DeleteIcon';
import { useNotification } from '../../utils/hooks/useNotification';
import { useActorRef } from '@xstate/react';
import { operationsCenterActor } from 'machines/operationsCenter';
import { useSelectorRtk } from '@/store/hooks';

export const NotificationCenterContext = React.createContext({
  drawerAnchorEl: null,
  setDrawerAnchor: () => {},
  toggleButtonRef: null,
  operationsCenterActorRef: null,
});

export const NotificationCenterProvider = ({ children }) => {
  const [drawerAnchorEl, setDrawerAnchor] = useState(null);
  const toggleButtonRef = useRef(null);
  const { notify } = useNotification();
  const operationsCenterActorRef = useActorRef(operationsCenterActor, {
    input: {
      notify,
    },
  });
  return (
    <NotificationCenterContext.Provider
      value={{
        drawerAnchorEl,
        setDrawerAnchor,
        toggleButtonRef,
        operationsCenterActorRef,
      }}
    >
      {children}
      <NotificationCenter />
    </NotificationCenterContext.Provider>
  );
};

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
  ({ classes, notificationStyle, count, type, handleClick, severity }) => {
    const theme = useTheme();
    const selectedSeverity = useSelector(selectSeverity);
    const darkColor = notificationStyle?.darkColor || notificationStyle?.color;
    const chipStyles = {
      fill: theme.palette.type === 'dark' ? darkColor : notificationStyle?.color,
      height: '20px',
      width: '20px',
    };
    count = Number(count).toLocaleString('en', { useGrouping: true });
    return (
      <Tooltip title={type} placement="bottom">
        <Button
          style={{
            backgroundColor: alpha(chipStyles.fill, 0.2),
            border:
              selectedSeverity === severity
                ? `solid 2px ${chipStyles.fill}`
                : 'solid 2px transparent',
          }}
          onClick={handleClick}
        >
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
  const { data } = useGetEventsSummaryQuery({
    status: STATUS.UNREAD,
  });
  const { count_by_severity_level } = data || {
    count_by_severity_level: [],
    total_count: 0,
  };
  const {
    data: { total_count: read_count } = {
      total_count: 0,
    },
  } = useGetEventsSummaryQuery({
    status: STATUS.READ,
  });
  const classes = useStyles();
  const onClickSeverity = (severity) => {
    handleFilter({
      severity: [severity],
      status: STATUS.UNREAD,
    });
  };

  const onClickStatus = (status) => {
    handleFilter({
      status: status,
    });
  };

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
            severity={severity}
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
          severity={STATUS.READ}
          count={read_count}
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

const BulkActions = () => {
  const checkedEvents = useSelector(selectCheckedEvents);
  const noEventsPresent = useSelector((state) => selectEvents(state).length === 0);
  const [deleteEvents, { isLoading: isDeleting }] = useDeleteEventsMutation();
  const [updateEvents, { isLoading: isUpdatingStatus }] = useUpdateEventsMutation();

  // stores which update is currently going on , usefull to know which action is going
  // if multiple updates can be triggered from same mutator , only single bulk action is allowed at a time
  const [curentOngoingUpdate, setCurrentOngoingUpdate] = useState(null);
  const isActionInProgress = isDeleting || isUpdatingStatus;

  const dispatch = useDispatch();
  const areAllEventsChecked = useSelector(selectAreAllEventsChecked);
  const handleCheckboxChange = (_e, v) => {
    dispatch(updateCheckAllEvents(v));
  };
  const resetSelection = () => {
    dispatch(updateCheckAllEvents(false));
  };

  const handleDelete = () => {
    deleteEvents({
      ids: checkedEvents.map((e) => e.id),
    }).then(resetSelection);
  };

  const handleChangeStatus = (status) => {
    setCurrentOngoingUpdate(status);
    updateEvents({
      ids: checkedEvents.map((e) => e.id),
      updatedFields: {
        status,
      },
    }).then(resetSelection);
  };

  const BulkActionButton = ({ isLoading, isDisabled, tooltip, Icon, onClick }) => {
    const disabled = isDisabled || isActionInProgress;
    if (isLoading) {
      return (
        <div style={iconMedium}>
          <CircularProgress size={iconMedium.height} />
        </div>
      );
    }
    return (
      <Tooltip title={tooltip} placement="top">
        <IconButton onClick={onClick} disabled={disabled}>
          <Icon
            {...iconMedium}
            style={{
              opacity: disabled ? 0.5 : 1,
            }}
            fill="currentColor"
          />
        </IconButton>
      </Tooltip>
    );
  };

  if (noEventsPresent) {
    return null;
  }

  return (
    <Box
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '0.15rem',
      }}
    >
      <Box>
        <Checkbox checked={areAllEventsChecked} color="primary" onChange={handleCheckboxChange} />
      </Box>
      <Collapse in={checkedEvents.length > 0}>
        <Box style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <BulkActionButton
            tooltip="Delete selected notifications"
            Icon={DeleteIcon}
            isLoading={isDeleting}
            onClick={handleDelete}
          />
          <BulkActionButton
            tooltip="Mark selected notifications as read"
            Icon={ReadIcon}
            isLoading={isUpdatingStatus && curentOngoingUpdate == STATUS.READ}
            onClick={() => handleChangeStatus(STATUS.READ)}
          />
          <BulkActionButton
            tooltip="Mark selected notifications as unread"
            Icon={UnreadIcon}
            isLoading={isUpdatingStatus && curentOngoingUpdate == STATUS.UNREAD}
            onClick={() => handleChangeStatus(STATUS.UNREAD)}
          />
        </Box>
      </Collapse>
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

const NotificationCenterDrawer = () => {
  const dispatch = useDispatch();
  const {
    toggleButtonRef,
    drawerAnchorEl: anchorEl,
    setDrawerAnchor: setAnchorEl,
  } = useContext(NotificationCenterContext);
  const isNotificationCenterOpen = useSelector((state) => state.events.isNotificationCenterOpen);
  const [fetchEvents, { isFetching }] = useLazyGetEventsQuery();
  const hasMore = useSelector((state) => state.events.current_view.has_more);

  const [isLoadingFilters, setIsLoadingFilters] = useState(false); // whether we are loading filters and basically should show loading spinner as we are loading the whole page

  useEffect(() => {
    dispatch(
      loadEvents(fetchEvents, 0, {
        status: STATUS.UNREAD,
      }),
    );
  }, []);

  const loadMore = () => {
    dispatch(loadNextPage(fetchEvents));
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
  const handleFilter = async (filters) => {
    setIsLoadingFilters(true);
    await dispatch(loadEvents(fetchEvents, 0, filters));
    setIsLoadingFilters(false);
  };
  const drawerRef = useRef();
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
      <DarkBackdrop open={isNotificationCenterOpen} />
      <ClickAwayListener onClickAway={clickwayHandler}>
        <Drawer
          anchor="right"
          variant="persistent"
          open={open}
          ref={drawerRef}
          BackdropComponent={<DarkBackdrop open={isNotificationCenterOpen} />}
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
                  <BulkActions />

                  {isLoadingFilters ? (
                    <Loading />
                  ) : (
                    <EventsView
                      handleLoadNextPage={loadMore}
                      isFetching={isFetching}
                      hasMore={hasMore}
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
        </Drawer>
      </ClickAwayListener>
    </>
  );
};

const NotificationDrawerButton_ = () => {
  const classes = useStyles();
  const { setDrawerAnchor, toggleButtonRef } = useContext(NotificationCenterContext);
  const dispatch = useDispatch();
  const handleToggle = () => {
    dispatch(toggleNotificationCenter());
  };
  return (
    <div ref={toggleButtonRef}>
      <IconButton
        id="notification-button"
        className={classes.notificationButton}
        color="inherit"
        onClick={handleToggle}
        onMouseOver={(e) => {
          e.preventDefault();
          setDrawerAnchor(e.currentTarget);
        }}
        onMouseLeave={(e) => {
          e.preventDefault();
          setDrawerAnchor(null);
        }}
      >
        <NavbarNotificationIcon />
      </IconButton>
    </div>
  );
};

export const NotificationDrawerButton = () => {
  return (
    <Provider store={store}>
      <NotificationDrawerButton_ />
    </Provider>
  );
};

const NotificationCenter = (props) => {
  const isOpen = useSelectorRtk((state) => state.events.isNotificationCenterOpen);

  if (!isOpen) {
    return null;
  }

  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in NotificationCenter', e)}
      >
        <Provider store={store}>
          <NotificationCenterDrawer {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

export default NotificationCenter;
