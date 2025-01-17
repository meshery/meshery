import { MenuItem, Avatar } from '@layer5/sistent';
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
} from './styles';
import { iconMedium } from 'css/icons.styles';
import { UsesSistent } from '../SistentWrapper';

export const _ConnectionChip = ({ handlePing, onDelete, iconSrc, status, title, width }) => {
  const chipStyle = { width };
  return (
    // <Tooltip title={tooltip || title} placement="bottom">
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
                ? theme.palette.secondary.success
                : theme.palette.secondary.penColorSecondary
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
    // </Tooltip>
  );
};

export const TootltipWrappedConnectionChip = (props) => {
  return (
    <CustomTooltip title={props.tooltip || props.title} placement="left">
      <div>
        <_ConnectionChip {...props} />
      </div>
    </CustomTooltip>
  );
};

const DiscoveredStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <DiscoveredChip
        avatar={<ExploreIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="1-discovered" /> */}
    </MenuItem>
  );
};

const RegisteredStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <RegisteredChip
        avatar={<AssignmentTurnedInIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="2-registered" /> */}
    </MenuItem>
  );
};

const ConnectedStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <ConnectedChip
        avatar={<CheckCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="3-connected" /> */}
    </MenuItem>
  );
};

const DisconnectedStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <DisconnectedChip
        avatar={<DisconnectIcon fill={notificationColors.lightwarning} width={24} height={24} />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
    </MenuItem>
  );
};

const IgnoredStateChip = ({ value }) => {
  // const classes = styles/();
  return (
    <MenuItem value={value}>
      <IgnoredChip
        avatar={<RemoveCircleIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="4-ignored" /> */}
    </MenuItem>
  );
};

const DeletedStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <DeletedChip
        avatar={<DeleteForeverIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </MenuItem>
  );
};

const MaintainanceStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <MaintainanceChip
        avatar={<HandymanIcon />}
        label={value}
        // helpIcon={<HelpToolTip classes={classes} value="7-deleted" />}
      />
      {/* <HelpToolTip classes={classes} value="7-deleted" /> */}
    </MenuItem>
  );
};

const NotFoundStateChip = ({ value }) => {
  return (
    <MenuItem value={value}>
      <NotFoundChip
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
  return (
    <MenuItem value={value}>
      <DiscoveredChip value={value} avatar={<ExploreIcon />} label={value} />
    </MenuItem>
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
