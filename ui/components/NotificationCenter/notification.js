import * as React from 'react';
import {
  Avatar,
  Box,
  Button,
  Collapse,
  Grid,
  IconButton,
  Popover,
  Slide,
  Tooltip,
  Typography,
  useTheme,
  Checkbox,
  styled
} from '@layer5/sistent';
import { alpha } from '@mui/system';
import { SEVERITY, SEVERITY_STYLE, STATUS } from './constants';
import { MoreVert as MoreVertIcon } from '@mui/icons-material';
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
import theme from '../../themes/app';
import { truncate } from 'lodash';
import { UsesSistent } from '../SistentWrapper';

const StyledRoot = styled('div')(({ notificationColor, status }) => ({
  width: '100%',
  borderRadius: '0.25rem',
  border: `0.1rem solid ${notificationColor}`,
  borderLeftWidth: status === STATUS.UNREAD ? '0.5rem' : '0.1rem',
  marginBlock: '0.5rem',
}));

const StyledSummary = styled('div')(({ theme, notificationColor }) => ({
  paddingBlock: '0.5rem',
  paddingInline: '0.25rem',
  cursor: 'pointer',
  backgroundColor: alpha(notificationColor, 0.2),
}));

const StyledExpanded = styled(Grid)(({ theme }) => ({
  paddingBlock: '0.75rem',
  paddingInline: '0.2rem',
  [theme.breakpoints.down('md')]: {
    padding: '0.5rem',
  },
}));

const AvatarStack = ({ avatars, direction }) => {
  const theme = useTheme();
  return (
    <UsesSistent>
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
    </UsesSistent>
  );
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
        <Box sx={{ padding: '0.5rem', borderRadius: '0.25rem', width: '12.5rem', boxShadow: 4 }}>
          <Button onClick={toggleSocialShare} sx={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <ShareIcon />
            <Typography variant="body1">Share</Typography>
          </Button>
          <Collapse in={isSocialShareOpen}>
            <Box sx={{ display: 'flex', gap: '0.5rem', padding: '0.65rem' }}>
              <FacebookShareButton url={'https://meshery.io'} quote={event.description || ''}>
                <FacebookIcon />
              </FacebookShareButton>
              <LinkedinShareButton url={'https://meshery.io'} summary={event.description || ''}>
                <LinkedInIcon />
              </LinkedinShareButton>
              <TwitterShareButton url={'https://meshery.io'} title={event.description || ''}>
                <TwitterIcon />
              </TwitterShareButton>
            </Box>
          </Collapse>
        </Box>
      </Popover>
    </div>
    </UsesSistent>
  );
};

const Notification = ({ event_id }) => {
  const event = useSelector((state) => selectEventById(state, event_id));
  const isVisible = useSelector((state) => selectIsEventVisible(state, event.id));
  const severityStyles = SEVERITY_STYLE[event.severity] || SEVERITY_STYLE[SEVERITY.INFO];
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
    <UsesSistent>
    <StyledExpanded container>
      <Grid item sm={1}>
        <AvatarStack avatars={eventActors} direction={{ xs: 'row', md: 'column' }} />
      </Grid>
      <Grid item sm={10}>
        <FormattedMetadata event={event} />
      </Grid>
    </StyledExpanded>
    </UsesSistent>
  );

  return (
    <UsesSistent>
    <Slide in={isVisible} timeout={250} direction="left" appear={false} mountOnEnter unmountOnExit>
      <StyledRoot notificationColor={severityStyles.color} status={event.status}>
        <Grid container component={StyledSummary} notificationColor={severityStyles.color} onClick={handleExpandClick}>
          <Grid item xs={2} sm={2} sx={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
            <Checkbox
              onClick={(e) => e.stopPropagation()}
              checked={Boolean(event.checked)}
              onChange={handleSelectEvent}
              sx={{ margin: '0rem', padding: '0rem' }}
            />
            <severityStyles.icon />
          </Grid>
          <Grid item xs={8} sm={6}>
            <Typography variant="body1" sx={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', maxWidth: 'min(25rem, 50vw)' }}>
              {truncate(event.description, { length: MAX_NOTIFICATION_DESCRIPTION_LENGTH })}
            </Typography>
          </Grid>
          <Grid item xs={1} sm={4} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-around' }}>
            <Typography variant="body1">{moment(event.created_at).fromNow()}</Typography>
            <BasicMenu event={event} />
          </Grid>
        </Grid>
        <Collapse in={expanded}>{expanded && <Detail />}</Collapse>
      </StyledRoot>
    </Slide>
    </UsesSistent>
  );
};

export default Notification;
