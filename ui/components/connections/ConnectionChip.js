import { Avatar, useTheme } from '@layer5/sistent';
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
import { CustomTooltip } from '@layer5/sistent';
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
  ConnectionStyledMenuItem,
} from './styles';
import { iconMedium } from 'css/icons.styles';
import { UsesSistent } from '../SistentWrapper';

export const _ConnectionChip = ({ handlePing, onDelete, iconSrc, status, title, width }) => {
  const chipStyle = { width };
  const theme = useTheme();
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
          <BadgeAvatars
            color={
              status === CONNECTION_STATES.CONNECTED || status === CONTROLLER_STATES.DEPLOYED
                ? theme.palette.background.brand.default
                : theme.palette.text.disabled
            }
          >
            <Avatar src={iconSrc} style={(status ? {} : { opacity: 0.2 }, iconMedium)}>
              <img style={iconMedium} src="/static/img/kubernetes.svg" />
            </Avatar>
          </BadgeAvatars>
        ) : (
          <Avatar src={iconSrc} sx={iconMedium}>
            <img style={iconMedium} src="/static/img/kubernetes.svg" />
          </Avatar>
        )
      }
      // variant="filled"
      data-cy="chipContextName"
      style={chipStyle}
    />
  );
};

export const TootltipWrappedConnectionChip = (props) => {
  return (
    <UsesSistent>
      <CustomTooltip title={props.tooltip || props.title} placement="left">
        <div>
          <_ConnectionChip {...props} />
        </div>
      </CustomTooltip>
    </UsesSistent>
  );
};

const DiscoveredStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value} sx={{ padding: 0 }}>
      <DiscoveredChip
        avatar={<ExploreIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="1-discovered" /> */}
    </ConnectionStyledMenuItem>
  );
};

const RegisteredStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <RegisteredChip
        avatar={<AssignmentTurnedInIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="2-registered" /> */}
    </ConnectionStyledMenuItem>
  );
};

const ConnectedStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <ConnectedChip
        avatar={<CheckCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="3-connected" /> */}
    </ConnectionStyledMenuItem>
  );
};

const DisconnectedStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <DisconnectedChip
        avatar={<DisconnectIcon fill={notificationColors.lightwarning} width={24} height={24} />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
    </ConnectionStyledMenuItem>
  );
};

const IgnoredStateChip = ({ value }) => {
  // const classes = styles/();
  return (
    <ConnectionStyledMenuItem value={value}>
      <IgnoredChip
        avatar={<RemoveCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="4-ignored" /> */}
    </ConnectionStyledMenuItem>
  );
};

const DeletedStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <DeletedChip
        avatar={<DeleteForeverIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </ConnectionStyledMenuItem>
  );
};

const MaintainanceStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <MaintainanceChip
        avatar={<HandymanIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </ConnectionStyledMenuItem>
  );
};

const NotFoundStateChip = ({ value }) => {
  return (
    <ConnectionStyledMenuItem value={value}>
      <NotFoundChip
        avatar={<NotInterestedRoundedIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </ConnectionStyledMenuItem>
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
  return (
    <ConnectionStyledMenuItem value={value}>
      <DiscoveredChip value={value} avatar={<ExploreIcon />} label={value} />
    </ConnectionStyledMenuItem>
  );
};

function getStatusChip(status) {
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
  return <UsesSistent>{getStatusChip(status)}</UsesSistent>;
};
