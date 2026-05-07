import React from 'react';
import { Avatar, useTheme } from '@sistent/sistent';
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
import {
  CONNECTION_STATE_TO_TRANSITION_MAP,
  CONNECTION_STATES,
  CONTROLLER_STATES,
} from '../../utils/Enum';
import { normalizeStaticImagePath } from '@/utils/fallback';
import { CustomTooltip } from '@sistent/sistent';
import {
  ChipWrapper,
  ConnectedChip,
  DeletedChip,
  DisconnectedChip,
  DiscoveredChip,
  IgnoredChip,
  MaintainanceChip,
  NotFoundChip,
  RegisteredChip,
} from './styles';
import { iconMedium, iconSmall } from 'css/icons.styles';
import ConnectionIcon from '@/assets/icons/Connection';

export const ConnectionChip = ({ handlePing, onDelete, iconSrc, status, title, width }) => {
  const chipStyle = { width };
  const theme = useTheme();
  const normalizedIconSrc = normalizeStaticImagePath(iconSrc) || undefined;

  const STATUS_LEVEL_MAP = Object.fromEntries([
    ...[CONNECTION_STATES.CONNECTED, CONTROLLER_STATES.DEPLOYED].map((status) => [
      status.toLowerCase(),
      'healthy',
    ]),
    ...[
      CONTROLLER_STATES.ENABLED,
      CONTROLLER_STATES.RUNNING,
      CONTROLLER_STATES.DEPLOYING,
      CONNECTION_STATES.REGISTERED,
    ].map((status) => [status.toLowerCase(), 'partial']),
  ]);

  const getStatusLevel = (status) => {
    if (!status) return 'error';
    return STATUS_LEVEL_MAP[status.toLowerCase()] || 'error';
  };

  const getStatusColor = (statusLevel) => {
    switch (statusLevel) {
      case 'healthy':
        return theme.palette.background.brand.default;
      case 'partial':
        return theme.palette.background.warning.default;
      case 'error':
      default:
        return theme.palette.text.disabled;
    }
  };

  return (
    <ChipWrapper
      label={title}
      onClick={(e) => {
        e.stopPropagation(); // Prevent event propagation
        handlePing(); // Call your custom handler
      }}
      onDelete={onDelete}
      avatar={
        status ? (
          <BadgeAvatars color={getStatusColor(getStatusLevel(status))}>
            <Avatar
              src={normalizedIconSrc}
              style={{ ...(status ? {} : { opacity: 0.2 }), ...iconMedium }}
            >
              <ConnectionIcon {...iconSmall} />
            </Avatar>
          </BadgeAvatars>
        ) : (
          <Avatar src={normalizedIconSrc} sx={iconMedium}>
            <ConnectionIcon {...iconSmall} />
          </Avatar>
        )
      }
      // variant="filled"
      data-cy="chipContextName"
      style={chipStyle}
    />
  );
};

export const TooltipWrappedConnectionChip = (props) => {
  return (
    <CustomTooltip title={props.tooltip || props.title} placement="left">
      <div style={{ display: 'inline-block' }}>
        <ConnectionChip {...props} />
      </div>
    </CustomTooltip>
  );
};

const DiscoveredStateChip = ({ value }) => {
  return (
    <DiscoveredChip
      avatar={<ExploreIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const RegisteredStateChip = ({ value }) => {
  return (
    <RegisteredChip
      avatar={<AssignmentTurnedInIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const ConnectedStateChip = ({ value }) => {
  return (
    <ConnectedChip
      avatar={<CheckCircleIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const DisconnectedStateChip = ({ value }) => {
  return (
    <DisconnectedChip
      avatar={<DisconnectIcon fill={notificationColors.lightwarning} width={24} height={24} />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const IgnoredStateChip = ({ value }) => {
  // const classes = styles/();
  return (
    <IgnoredChip
      avatar={<RemoveCircleIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const DeletedStateChip = ({ value }) => {
  return (
    <DeletedChip
      avatar={<DeleteForeverIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const MaintainanceStateChip = ({ value }) => {
  return (
    <MaintainanceChip
      avatar={<HandymanIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
  );
};

const NotFoundStateChip = ({ value }) => {
  return (
    <NotFoundChip
      avatar={<NotInterestedRoundedIcon />}
      label={value}
      // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
    />
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
  return <DiscoveredChip avatar={<ExploreIcon />} label={value} />;
};

function getStatusChip(status, actionable) {
  const value = actionable ? CONNECTION_STATE_TO_TRANSITION_MAP[status] : status;
  switch (status) {
    case 'ignored':
      return <IgnoredStateChip value={value} />;
    case 'connected':
      return <ConnectedStateChip value={value} />;
    case 'registered':
      return <RegisteredStateChip value={value} />;
    case 'discovered':
      return <DiscoveredStateChip value={value} />;
    case 'deleted':
      return <DeletedStateChip value={value} />;
    case 'maintenance':
      return <MaintainanceStateChip value={value} />;
    case 'disconnected':
      return <DisconnectedStateChip value={value} />;
    case 'not found':
      return <NotFoundStateChip value={value} />;
    default:
      return <Default value={value} />;
  }
}

export const ConnectionStateChip = ({ status, actionable }) => {
  return <>{getStatusChip(status, actionable)}</>;
};
