import * as React from 'react';
import {
  Avatar,
  Box,
  Button,
  Collapse,
  Grid,
  Hidden,
  IconButton,
  Popover,
  Slide,
  Tooltip,
  Typography,
  alpha,
  useTheme,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core';
import { SEVERITY, SEVERITY_STYLE, STATUS } from './constants';
import { iconLarge, iconMedium } from '../../css/icons.styles';
import { MoreVert as MoreVertIcon } from '@material-ui/icons';
import FacebookIcon from '../../assets/icons/FacebookIcon';
import LinkedInIcon from '../../assets/icons/LinkedInIcon';
import TwitterIcon from '../../assets/icons/TwitterIcon';
import ShareIcon from '../../assets/icons/ShareIcon';
import DeleteIcon from '../../assets/icons/DeleteIcon';
import moment from 'moment';
import {
  useUpdateStatusMutation,
  useDeleteEventMutation,
} from '../../rtk-query/notificationCenter';
import { useDispatch, useSelector } from 'react-redux';
import {
  selectEventById,
  selectIsEventVisible,
  updateIsEventChecked,
} from '../../store/slices/events';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import { FacebookShareButton, LinkedinShareButton, TwitterShareButton } from 'react-share';
import ReadIcon from '../../assets/icons/ReadIcon';
import UnreadIcon from '../../assets/icons/UnreadIcon';
import {
  ErrorBoundary,
  withErrorBoundary,
  withSuppressedErrorBoundary,
} from '../General/ErrorBoundary';
import { FormattedMetadata } from './metadata';
import theme from '../../themes/app';
import { truncate } from 'lodash';
import { Checkbox } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';

const useStyles = makeStyles(() => ({
  root: (props) => ({
    width: '100%',
    borderRadius: '0.25rem',
    border: `0.1rem solid ${props.notificationColor}`,
    borderLeftWidth: props.status === STATUS.UNREAD ? '0.5rem' : '0.1rem',
    marginBlock: '0.5rem',
  }),

  summary: (props) => ({
    paddingBlock: '0.5rem',
    paddingInline: '0.25rem',
    cursor: 'pointer',
    backgroundColor: alpha(props.notificationColor, 0.2),
  }),

  gridItem: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },

  message: {
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    overflowWrap: 'break-word',
    // max of min of 20rem or 50vw
    maxWidth: 'min(25rem, 50vw)',
    width: '100%',
  },
  expanded: {
    paddingBlock: '0.75rem',
    paddingInline: '0.2rem',
    [theme.breakpoints.down('md')]: {
      padding: '0.5rem',
    },
  },
  actorAvatar: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'start',
    paddingTop: '1rem',
  },

  descriptionHeading: {
    fontWeight: 'bolder !important',
    textTransform: 'uppercase',
    fontSize: '0.9rem',
  },
}));

export const eventPreventDefault = (e) => {
  e.preventDefault();
};

export const eventstopPropagation = (e) => {
  e.stopPropagation();
};

export const MAX_NOTIFICATION_DESCRIPTION_LENGTH = 45;

export const canTruncateDescription = (description) => {
  return description.length > MAX_NOTIFICATION_DESCRIPTION_LENGTH;
};

const AvatarStack = ({ avatars, direction }) => {
  const theme = useTheme();
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: direction,
        alignItems: 'center',
        justifyContent: 'center',
        gap: '0',

        '& .MuiAvatar-root': {
          width: '2rem',
          height: '2rem',
          border: '0.05rem solid ' + theme.palette.secondary.menuActionText,
        },
      }}
    >
      {avatars.map((avatar, index) => (
        <Tooltip title={avatar.name} placement="top" key={index}>
          <div style={{ zIndex: avatars.length - index, marginTop: '-0.4rem' }}>
            <Avatar alt={avatar.name} src={avatar.avatar_url} />
          </div>
        </Tooltip>
      ))}
    </Box>
  );
};

const useMenuStyles = makeStyles((theme) => {
  return {
    paper: {
      color: theme.palette.secondary.iconMain,
      boxShadow: theme.shadows[4],
      borderRadius: '0.25',
      paddingInline: '0.5rem',
      paddingBlock: '0.25rem',
      width: '12.5rem',
    },

    list: {
      display: 'flex',
      flexDirection: 'column',
      gridGap: '0.5rem',
      marginBlock: '0.5rem',
      borderRadius: '0.25rem',
      backgroundColor: theme.palette.secondary.honeyComb,
      '&:hover': {
        backgroundColor: alpha(theme.palette.secondary.link2, 0.25),
      },
    },

    listItem: {
      display: 'flex',
      gridGap: '0.5rem',
      alignItems: 'center',
      justifyContent: 'space-around',
    },
    socialListItem: {
      display: 'flex',
      backgroundColor: alpha(theme.palette.secondary.honeyComb, 0.25),
      alignItems: 'center',
      justifyContent: 'space-around',
      padding: '.65rem',
    },

    button: {
      height: '100%',
      width: '100%',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'start',
    },
  };
});

const formatTimestamp = (utcTimestamp) => {
  const currentUtcTimestamp = moment.utc().valueOf();

  const timediff = currentUtcTimestamp - moment(utcTimestamp).valueOf();

  if (timediff >= 24 * 60 * 60 * 1000) {
    return moment(utcTimestamp).local().format('MMM DD, YYYY');
  }
  return moment(utcTimestamp).fromNow();
};

const BasicMenu = withSuppressedErrorBoundary(({ event }) => {
  const classes = useMenuStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);

  const handleClick = (event) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };
  const handleClose = (e) => {
    e.stopPropagation();
    setAnchorEl(null);
  };

  const [isSocialShareOpen, setIsSocialShareOpen] = React.useState(false);
  const toggleSocialShare = (e) => {
    e.stopPropagation();
    setIsSocialShareOpen((prev) => !prev);
  };

  const theme = useTheme();
  return (
    <div className="mui-fixed" onClick={(e) => e.stopPropagation()}>
      <IconButton
        id="basic-button"
        aria-controls={open ? 'basic-menu' : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
        onClick={handleClick}
      >
        <MoreVertIcon />
      </IconButton>
      <Popover
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
      >
        <Box className={classes.paper}>
          <div className={classes.list}>
            <Box className={classes.listItem} sx={{ width: '100%' }}>
              <Button onClick={toggleSocialShare} className={classes.button}>
                <ShareIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
                  Share
                </Typography>
              </Button>
            </Box>
            <Collapse in={isSocialShareOpen}>
              <Box className={classes.socialListItem}>
                <FacebookShareButton url={'https://meshery.io'} quote={event.description || ''}>
                  <FacebookIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </FacebookShareButton>
                <LinkedinShareButton url={'https://meshery.io'} summary={event.description || ''}>
                  <LinkedInIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </LinkedinShareButton>
                <TwitterShareButton url={'https://meshery.io'} title={event.description || ''}>
                  <TwitterIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </TwitterShareButton>
              </Box>
            </Collapse>
          </div>

          <DeleteEvent event={event} />
          <ChangeStatus event={event} />
        </Box>
      </Popover>
    </div>
  );
});

export const DeleteEvent = ({ event }) => {
  const classes = useMenuStyles();
  const [deleteEventMutation] = useDeleteEventMutation();
  const theme = useTheme();
  const handleDelete = (e) => {
    e.stopPropagation();
    deleteEventMutation({ id: event.id });
  };
  return (
    <div className={classes.list}>
      <Button className={classes.button} onClick={handleDelete}>
        <DeleteIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
        <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
          {' '}
          Delete{' '}
        </Typography>
      </Button>
    </div>
  );
};

export const ChangeStatus = ({ event }) => {
  const classes = useMenuStyles();
  const newStatus = event.status === STATUS.READ ? STATUS.UNREAD : STATUS.READ;
  const [updateStatusMutation] = useUpdateStatusMutation();
  const theme = useTheme();

  const updateStatus = (e) => {
    e.stopPropagation();
    updateStatusMutation({ id: event.id, status: newStatus });
  };
  return (
    <div className={classes.list}>
      <Button className={classes.button} onClick={updateStatus}>
        {newStatus === STATUS.READ ? (
          <ReadIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
        ) : (
          <UnreadIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
        )}
        <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
          {' '}
          Mark as {newStatus}{' '}
        </Typography>
      </Button>
    </div>
  );
};

export const Notification = withErrorBoundary(({ event_id }) => {
  const event = useSelector((state) => selectEventById(state, event_id));
  const isVisible = useSelector((state) => selectIsEventVisible(state, event.id));
  const severityStyles = SEVERITY_STYLE[event.severity] || SEVERITY_STYLE[SEVERITY.INFO];
  const classes = useStyles({
    notificationColor: severityStyles?.color,
    status: event?.status,
  });
  const theme = useTheme();
  const dispatch = useDispatch();
  const [expanded, setExpanded] = React.useState(false);
  const handleExpandClick = (e) => {
    e.stopPropagation();
    setExpanded(!expanded);
  };

  const { data: user } = useGetUserByIdQuery(event.user_id || '');

  const userName = `${user?.first_name || ''} ${user?.last_name || ''}`;
  const userAvatarUrl = user?.avatar_url || '';

  const handleSelectEvent = (e, value) => {
    e.stopPropagation();
    dispatch(
      updateIsEventChecked({
        id: event.id,
        value,
      }),
    );
  };

  const eventActors = [
    ...(event.user_id && user
      ? [{ name: userName, avatar_url: userAvatarUrl, tooltip: userName }]
      : []),
    ...(event.system_id
      ? [
          {
            name: 'Meshery',
            avatar_url: '/static/img/meshery-logo.png',
            tooltip: `System ID: ${event.system_id}`,
          },
        ]
      : []),
  ];

  return (
    <Slide
      in={isVisible}
      timeout={250}
      direction="left"
      appear={false}
      enter={false}
      mountOnEnter
      unmountOnExit
    >
      <div className={classes.root}>
        <Grid container className={classes.summary} onClick={handleExpandClick}>
          <Grid
            item
            xs="auto"
            sm={2}
            className={classes.gridItem}
            style={{
              justifyContent: 'start',
              alignItems: 'center',
              gap: '0.25rem',
            }}
          >
            <UsesSistent>
              <Checkbox
                onClick={eventstopPropagation}
                checked={Boolean(event.checked)}
                onChange={handleSelectEvent}
                style={{ margin: '0rem', padding: '0rem' }}
              />
            </UsesSistent>
            <severityStyles.icon {...iconLarge} fill={severityStyles?.color} />
          </Grid>
          <Grid item xs={8} sm={6} className={classes.gridItem}>
            <Typography variant="body1" className={classes.message}>
              {truncate(event.description, {
                length: MAX_NOTIFICATION_DESCRIPTION_LENGTH,
              })}
            </Typography>
          </Grid>
          <Grid
            item
            xs={1}
            sm={4}
            className={classes.gridItem}
            style={{ display: 'flex', justifyContent: 'space-around', alignItems: 'center' }}
          >
            <Hidden smDown>
              <Typography variant="body1"> {formatTimestamp(event.created_at)} </Typography>
            </Hidden>
            <BasicMenu event={event} />
          </Grid>
        </Grid>
        <Collapse in={expanded}>
          <ErrorBoundary>
            <Grid container className={classes.expanded}>
              <Grid item sm={1} className={classes.actorAvatar}>
                <AvatarStack
                  avatars={eventActors}
                  direction={{
                    xs: 'row',
                    md: 'column',
                  }}
                />
              </Grid>
              <Grid
                item
                sm={10}
                style={{
                  color: theme.palette.secondary.textMain,
                }}
              >
                <FormattedMetadata event={event} classes={classes} />
              </Grid>
            </Grid>
          </ErrorBoundary>
        </Collapse>
      </div>
    </Slide>
  );
});

export default Notification;
