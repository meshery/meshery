import * as React from 'react';
import {
  Avatar,
  Box,
  Collapse,
  Grid,
  IconButton,
  Slide,
  Tooltip,
  Typography,
  useTheme,
  Checkbox,
} from '@layer5/sistent';
import {
  OptionList,
  OptionListItem,
  MenuPaper,
  SocialListItem,
  ListButton,
  ActorAvatar,
  Expanded,
  GridItem,
  Message,
  StyledAvatarStack,
  Root,
  Summary,
} from './notificationCenter.style';
import { UsesSistent } from '../SistentWrapper';
import { Popover } from '@mui/material';
import { alpha } from '@mui/material';
import { SEVERITY, SEVERITY_STYLE, STATUS } from './constants';
import { iconLarge, iconMedium } from '../../css/icons.styles';
import MoreVertIcon from '@mui/icons-material/MoreVert';
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
import { FormattedMetadata } from './metadata';

import { truncate } from 'lodash';

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
  return (
    <UsesSistent>
      <StyledAvatarStack
        sx={{
          flexDirection: direction,
        }}
      >
        {avatars.map((avatar, index) => (
          <Tooltip title={avatar.name} placement="top" key={index}>
            <Box
              sx={{
                zIndex: avatars.length - index,
                mt: '-0.4rem',
              }}
            >
              <Avatar alt={avatar.name} src={avatar.avatar_url} />
            </Box>
          </Tooltip>
        ))}
      </StyledAvatarStack>
    </UsesSistent>
  );
};

const formatTimestamp = (utcTimestamp) => {
  const currentUtcTimestamp = moment.utc().valueOf();

  const timediff = currentUtcTimestamp - moment(utcTimestamp).valueOf();

  if (timediff >= 24 * 60 * 60 * 1000) {
    return moment(utcTimestamp).local().format('MMM DD, YYYY');
  }
  return moment(utcTimestamp).fromNow();
};

const BasicMenu = ({ event }) => {
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
    <UsesSistent>
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
          <MenuPaper>
            <OptionList>
              <OptionListItem sx={{ width: '100%' }}>
                <ListButton onClick={toggleSocialShare}>
                  <ShareIcon {...iconMedium} fill={theme.palette.icon.secondary} />
                  <Typography variant="body1" sx={{ marginLeft: '0.5rem' }}>
                    Share
                  </Typography>
                </ListButton>
              </OptionListItem>
              <Collapse in={isSocialShareOpen}>
                <SocialListItem>
                  <FacebookShareButton url={'https://meshery.io'} quote={event.description || ''}>
                    <FacebookIcon {...iconMedium} fill={theme.palette.icon.secondary} />
                  </FacebookShareButton>
                  <LinkedinShareButton url={'https://meshery.io'} summary={event.description || ''}>
                    <LinkedInIcon {...iconMedium} fill={theme.palette.icon.secondary} />
                  </LinkedinShareButton>
                  <TwitterShareButton url={'https://meshery.io'} title={event.description || ''}>
                    <TwitterIcon {...iconMedium} fill={theme.palette.icon.secondary} />
                  </TwitterShareButton>
                </SocialListItem>
              </Collapse>
            </OptionList>

            <DeleteEvent event={event} />
            <ChangeStatus event={event} />
          </MenuPaper>
        </Popover>
      </div>
    </UsesSistent>
  );
};

export const DeleteEvent = ({ event }) => {
  const theme = useTheme();
  const [deleteEventMutation] = useDeleteEventMutation();

  const handleDelete = (e) => {
    e.stopPropagation();
    deleteEventMutation({ id: event.id });
  };
  return (
    <UsesSistent>
      <OptionList>
        <ListButton onClick={handleDelete}>
          <DeleteIcon {...iconMedium} fill={theme.palette.icon.secondary} />
          <Typography variant="body1" sx={{ marginLeft: '0.5rem' }}>
            {' '}
            Delete{' '}
          </Typography>
        </ListButton>
      </OptionList>
    </UsesSistent>
  );
};

export const ChangeStatus = ({ event }) => {
  const newStatus = event.status === STATUS.READ ? STATUS.UNREAD : STATUS.READ;
  const [updateStatusMutation] = useUpdateStatusMutation();
  const theme = useTheme();
  const updateStatus = (e) => {
    e.stopPropagation();
    updateStatusMutation({ id: event.id, status: newStatus });
  };
  return (
    <UsesSistent>
      <OptionList>
        <ListButton onClick={updateStatus}>
          {newStatus === STATUS.READ ? (
            <ReadIcon {...iconMedium} fill={theme.palette.icon.secondary} />
          ) : (
            <UnreadIcon {...iconMedium} fill={theme.palette.icon.secondary} />
          )}
          <Typography variant="body1" sx={{ marginLeft: '0.5rem' }}>
            {' '}
            Mark as {newStatus}{' '}
          </Typography>
        </ListButton>
      </OptionList>
    </UsesSistent>
  );
};

export const Notification = ({ event_id }) => {
  const event = useSelector((state) => selectEventById(state, event_id));
  const isVisible = useSelector((state) => selectIsEventVisible(state, event.id));
  const severityStyles = SEVERITY_STYLE[event.severity] || SEVERITY_STYLE[SEVERITY.INFO];
  const eventStyle = SEVERITY_STYLE[event?.severity] || {};
  const notificationColor = severityStyles?.color;
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

  const Detail = () => (
    <Expanded container>
      <Grid
        item
        xs={12}
        sx={{
          color: theme.palette.text.default,
          borderTop: `1px solid ${notificationColor}`,
          padding: '1rem',
          backgroundColor: alpha(eventStyle.color, 0.1),
        }}
      >
        <ActorAvatar item sm={1}>
          <AvatarStack
            avatars={eventActors}
            direction={{
              xs: 'row',
              md: 'row',
            }}
          />
        </ActorAvatar>
        <FormattedMetadata event={event} />
      </Grid>
    </Expanded>
  );
  return (
    <UsesSistent>
      <Slide
        in={isVisible}
        timeout={250}
        direction="left"
        appear={false}
        enter={false}
        mountOnEnter
        unmountOnExit
      >
        <Root notificationcolor={notificationColor} status={event?.status}>
          <Summary container notificationcolor={notificationColor} onClick={handleExpandClick}>
            <GridItem item xs="auto" sm={2}>
              <Checkbox
                onClick={eventstopPropagation}
                checked={Boolean(event.checked)}
                onChange={handleSelectEvent}
                sx={{ margin: '0rem', padding: '0rem' }}
              />

              <severityStyles.icon {...iconLarge} fill={severityStyles?.color} />
            </GridItem>
            <GridItem item xs={8} sm={6}>
              <Message variant="body1">
                {truncate(event.description, {
                  length: MAX_NOTIFICATION_DESCRIPTION_LENGTH,
                })}
              </Message>
            </GridItem>
            <GridItem item xs={1} sm={4} sx={{ justifyContent: 'space-around' }}>
              <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
                <Typography variant="body1">{formatTimestamp(event.created_at)}</Typography>
              </Box>
              <BasicMenu event={event} />
            </GridItem>
          </Summary>
          <Collapse in={expanded}>{expanded && <Detail />}</Collapse>
        </Root>
      </Slide>
    </UsesSistent>
  );
};

export default Notification;
