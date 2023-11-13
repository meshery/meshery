import { Chip, MenuItem, Tooltip, makeStyles } from '@material-ui/core';
import classNames from 'classnames';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import HelpIcon from '@material-ui/icons/Help';
import { iconSmall } from '../../css/icons.styles';

export const ConnectionChip = ({ handlePing, title, icon }) => (
  <Chip label={title} onClick={() => handlePing()} icon={icon} variant="outlined" />
);

const styles = makeStyles((theme) => ({
  statusCip: {
    minWidth: '120px !important',
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
      color: theme.palette.secondary.warning,
    },
    background: `${theme.palette.secondary.warning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.warning} !important`,
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
}));

const DiscoveredStateChip = ({ value }) => {
  const classes = styles();
  return (
    <MenuItem value={value}>
      <Chip
        className={classNames(classes.statusCip, classes.discovered)}
        avatar={<ExploreIcon />}
        label={value}
        helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      <HelpToolTip classes={classes} value="1-discovered" />
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
        helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      <HelpToolTip classes={classes} value="2-registered" />
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
        helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      <HelpToolTip classes={classes} value="3-connected" />
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
        helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
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
        helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
    </MenuItem>
  );
};

const HelpToolTip = ({ classes, value }) => {
  const url = `https://docs.meshery.io/concepts/connections#${value}`;
  const onClick = () => (e) => {
    e.preventDefault();
    window.open(url, '_blank');
  };
  return (
    <Tooltip onClick={onClick()} title={url}>
      <HelpIcon className={classes.helpIcon} style={{ fontSize: '1.45rem', ...iconSmall }} />
    </Tooltip>
  );
};

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
    default:
      return <Default value={status} />;
  }
}

export const ConnectionStateChip = ({ status }) => {
  return getStatusChip(status);
};
