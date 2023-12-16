import { Chip, MenuItem, Tooltip, makeStyles, Avatar } from '@material-ui/core';
import classNames from 'classnames';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import HandymanIcon from '@mui/icons-material/Handyman';
import BadgeAvatars from '../CustomAvatar';
import { notificationColors } from '../../themes';
import DisconnectIcon from '../../assets/icons/disconnect';
import NotInterestedRoundedIcon from '@mui/icons-material/NotInterestedRounded';
import { CONNECTION_STATES, CONTROLLER_STATES } from '../../utils/Enum';
import theme from '../../themes/app';

const useChipStyles = makeStyles(() => ({
  Chip: {
    width: '13rem',
    maxWidth: '13rem',
    minWidth: '9rem',
    textAlign: 'left',
    cursor: 'pointer',
    '& .MuiChip-label': {
      flexGrow: 1,
    },
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
  },
  icon: {
    width: '1.5rem',
    height: '1.5rem',
  },
}));

export const _ConnectionChip = ({ handlePing, onDelete, iconSrc, status, title }) => {
  const classes = useChipStyles();
  return (
    // <Tooltip title={tooltip || title} placement="bottom">
    <Chip
      label={title}
      onClick={handlePing}
      onDelete={onDelete}
      avatar={
        status ? (
          <BadgeAvatars
            color={
              status === CONNECTION_STATES.CONNECTED || status === CONTROLLER_STATES.DEPLOYED
                ? theme.palette.secondary.success
                : theme.palette.secondary.penColorSecondary
            }
          >
            <Avatar src={iconSrc} className={classes.icon} style={status ? {} : { opacity: 0.2 }} />
          </BadgeAvatars>
        ) : (
          <Avatar src={iconSrc} className={classes.icon} />
        )
      }
      variant="filled"
      className={classes.Chip}
      data-cy="chipContextName"
    />
    // </Tooltip>
  );
};

export const TootltipWrappedConnectionChip = (props) => {
  return (
    <Tooltip title={props.tooltip || props.title} placement="bottom">
      <>
        <_ConnectionChip {...props} />
      </>
    </Tooltip>
  );
};
const styles = makeStyles((theme) => ({
  statusCip: {
    minWidth: '142px !important',
    maxWidth: 'max-content !important',
    display: 'flex !important',
    justifyContent: 'flex-start !important',
    textTransform: 'capitalize',
    borderRadius: '3px !important',
    padding: '6px 8px',
    '& .MuiChip-label': {
      paddingTop: '3px',
      fontWeight: '400',
    },
    '& .MuiSvgIcon-root': {
      marginLeft: '0px !important',
    },
    '&:hover': {
      boxShadow: '0px 1px 2px 0px rgba(0, 0, 0, 0.25)',
    },
  },
  helpIcon: {
    color: '#fff',
    opacity: '0.7',
    transition: 'opacity 200ms linear',
    '&:hover': {
      opacity: 1,
      background: 'transparent',
    },
    '&:focus': {
      opacity: 1,
      background: 'transparent',
    },
  },
  ignored: {
    '& .MuiChip-label': {
      color: `${theme.palette.secondary.default}`,
    },
    background: `${theme.palette.secondary.default}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.default} !important`,
    },
  },
  connected: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.success,
    },
    background: `${theme.palette.secondary.success}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.success} !important`,
    },
  },
  registered: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.primary,
    },
    background: `${theme.palette.secondary.primary}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.primary} !important`,
    },
  },
  discovered: {
    '& .MuiChip-label': {
      color: notificationColors.info,
    },
    background: `${notificationColors.info}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${notificationColors.info} !important`,
    },
  },
  deleted: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.error,
    },
    background: `${theme.palette.secondary.lightError}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.error} !important`,
    },
  },
  maintenance: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.warning,
    },
    background: `${theme.palette.secondary.warning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.warning} !important`,
    },
  },
  disconnected: {
    '& .MuiChip-label': {
      color: notificationColors.lightwarning,
    },
    background: `${notificationColors.lightwarning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${notificationColors.lightwarning} !important`,
    },
  },
  notfound: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.text,
    },
    background: `${theme.palette.secondary.disableButtonBg}60 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.iconMain} !important`,
    },
  },
}));

const DiscoveredStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.discovered)}
        avatar={<ExploreIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="1-discovered" /> */}
    </MenuItem>
  );
};

const RegisteredStateChip = ({ value }) => {
  const classes = styles();
  console.log('tetsppppp', value);
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.registered)}
        avatar={<AssignmentTurnedInIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="2-registered" /> */}
    </MenuItem>
  );
};

const ConnectedStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.connected)}
        avatar={<CheckCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="3-connected" /> */}
    </MenuItem>
  );
};

const DisconnectedStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.disconnected)}
        avatar={<DisconnectIcon fill={notificationColors.lightwarning} width={24} height={24} />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
    </MenuItem>
  );
};

const IgnoredStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.ignored)}
        avatar={<RemoveCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="4-ignored" /> */}
    </MenuItem>
  );
};

const DeletedStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.deleted)}
        avatar={<DeleteForeverIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </MenuItem>
  );
};

const MaintainanceStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.maintenance)}
        avatar={<HandymanIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </MenuItem>
  );
};

const NotFoundStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.notfound)}
        avatar={<NotInterestedRoundedIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </MenuItem>
  );
};

// const HelpToolTip = ({ classes, value }) => {
//   const url = `https://docs.meshery.io/concepts/connections#${value}`;
//   const onClick = () => (e) => {
//     e.preventDefault();
//     e.stopPropagation();
//     window.open(url, '_blank');
//   };
//   return (
//     <Tooltip title={url}>
//       <IconButton onClick={onClick()} aria-label="help">
//         <HelpIcon className={classes.helpIcon} style={{ fontSize: '1.45rem', ...iconSmall }} />
//       </IconButton>
//     </Tooltip>
//   );
// };

const Default = ({ value }) => {
  const classes = styles();

  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusChip, classes.discovered)}
        value={value}
        avatar={<ExploreIcon />}
        label={value}
      />
    </MenuItem>
  );
};

function getStatusChip(status) {
  console.log('inside func: ', status);

  switch (status) {
    case 'ignored':
      return <IgnoredStateChip value={status} />;
    case 'connected':
      return <ConnectedStateChip value={status} />;
    case 'registered':
      return <RegisteredStateChip value={status} />;
    case 'discovered':
      return <DiscoveredStateChip value={status} />;
    case 'deleted':
      return <DeletedStateChip value={status} />;
    case 'maintenance':
      return <MaintainanceStateChip value={status} />;
    case 'disconnected':
      return <DisconnectedStateChip value={status} />;
    case 'not found':
      return <NotFoundStateChip value={status} />;
    default:
      return <Default value={status} />;
  }
}

export const ConnectionStateChip = ({ status }) => {
  return getStatusChip(status);
};
