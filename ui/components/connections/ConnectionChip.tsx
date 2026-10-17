import React, { memo } from 'react';
import {
  Avatar,
  AssignmentTurnedInIcon,
  CustomTooltip,
  notificationColors,
  useTheme,
} from '@sistent/sistent';
import {
  CheckCircle as CheckCircleIcon,
  Explore as ExploreIcon,
  RemoveCircle as RemoveCircleIcon,
  DeleteForever as DeleteForeverIcon,
  Handyman as HandymanIcon,
  // Warning — not Cancel/X. CancelIcon was re-exported as NotInterestedRounded
  // and read as a destructive delete action on "not found" chips.
  Warning as WarningIcon,
} from '@/assets/icons';
import BadgeAvatars from '../general/CustomAvatar';
import DisconnectIcon from '../../assets/icons/disconnect';
import {
  CONNECTION_STATE_TO_TRANSITION_MAP,
  CONNECTION_STATES,
  CONTROLLER_STATES,
} from '../../utils/Enum';
import { normalizeStaticImagePath } from '@/utils/fallback';
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

type ConnectionChipProps = {
  handlePing?: () => void;
  onDelete?: () => void;
  iconSrc?: string;
  status?: string;
  title?: React.ReactNode;
  width?: string | number;
  disabled?: boolean;
};

type TooltipWrappedConnectionChipProps = ConnectionChipProps & {
  tooltip?: React.ReactNode;
};

type ConnectionStateChipProps = {
  status?: string;
  actionable?: boolean;
};

// Status-dot color buckets. Keep in sync with ConnectionStateChip styling
// intent: green = healthy, yellow = transitional, orange = degraded, gray = inactive.
const HEALTHY_STATUSES = new Set(
  [CONNECTION_STATES.CONNECTED, CONTROLLER_STATES.DEPLOYED, CONTROLLER_STATES.CONNECTED].map(
    (status) => status.toLowerCase(),
  ),
);
// Yellow / amber — not fully healthy yet, but not failed.
const PARTIAL_STATUSES = new Set(
  [
    CONTROLLER_STATES.ENABLED,
    CONTROLLER_STATES.RUNNING,
    CONTROLLER_STATES.DEPLOYING,
    CONNECTION_STATES.REGISTERED,
    CONNECTION_STATES.DISCOVERED,
  ].map((status) => status.toLowerCase()),
);
// Orange — reachable-but-off / under maintenance (matches Disconnected/Maintenance chips).
const WARNING_STATUSES = new Set(
  [CONNECTION_STATES.DISCONNECTED, CONNECTION_STATES.MAINTENANCE].map((status) =>
    status.toLowerCase(),
  ),
);

const STATE_CHIP_CONFIG = {
  [CONNECTION_STATES.IGNORED]: {
    Component: IgnoredChip,
    avatar: <RemoveCircleIcon />,
  },
  [CONNECTION_STATES.CONNECTED]: {
    Component: ConnectedChip,
    avatar: <CheckCircleIcon />,
  },
  [CONNECTION_STATES.REGISTERED]: {
    Component: RegisteredChip,
    avatar: <AssignmentTurnedInIcon />,
  },
  [CONNECTION_STATES.DISCOVERED]: {
    Component: DiscoveredChip,
    avatar: <ExploreIcon />,
  },
  [CONNECTION_STATES.DELETED]: {
    Component: DeletedChip,
    avatar: <DeleteForeverIcon />,
  },
  [CONNECTION_STATES.MAINTENANCE]: {
    Component: MaintainanceChip,
    avatar: <HandymanIcon />,
  },
  [CONNECTION_STATES.DISCONNECTED]: {
    Component: DisconnectedChip,
    avatar: <DisconnectIcon fill={notificationColors.warning.light} width={24} height={24} />,
  },
  [CONNECTION_STATES.NOTFOUND]: {
    Component: NotFoundChip,
    // Explicit fill + size: WarningIcon has no path fill by default, so without
    // currentColor the glyph can disappear on dark/disabled chip backgrounds.
    avatar: (
      <WarningIcon
        data-testid="not-found-warning-icon"
        fill="currentColor"
        width={20}
        height={20}
      />
    ),
  },
} as const;

const DEFAULT_STATE_CHIP = {
  Component: DiscoveredChip,
  avatar: <ExploreIcon />,
};

const getStatusLevel = (status?: string): 'healthy' | 'partial' | 'warning' | 'error' => {
  if (!status) {
    return 'error';
  }

  const normalizedStatus = status.toLowerCase();

  if (HEALTHY_STATUSES.has(normalizedStatus)) {
    return 'healthy';
  }

  if (PARTIAL_STATUSES.has(normalizedStatus)) {
    return 'partial';
  }

  if (WARNING_STATUSES.has(normalizedStatus)) {
    return 'warning';
  }

  // ignored / deleted / not found / unknown → gray
  return 'error';
};

const getStatusColor = (theme: ReturnType<typeof useTheme>, status?: string) => {
  switch (getStatusLevel(status)) {
    case 'healthy':
      // Green
      return theme.palette.background.brand.default;
    case 'partial':
      // Yellow / amber
      return theme.palette.background.warning.default;
    case 'warning':
      // Orange (same token used by the disconnected state chip icon)
      return notificationColors.warning.light;
    default:
      // Gray
      return theme.palette.text.disabled;
  }
};

const ConnectionChip_ = ({
  handlePing,
  onDelete,
  iconSrc,
  status,
  title,
  width,
  disabled = false,
}: ConnectionChipProps) => {
  const theme = useTheme();
  const normalizedIconSrc = normalizeStaticImagePath(iconSrc) || undefined;
  const isPingEnabled = Boolean(handlePing) && !disabled;

  const avatar = status ? (
    <BadgeAvatars color={getStatusColor(theme, status)}>
      <Avatar src={normalizedIconSrc} sx={iconMedium}>
        <ConnectionIcon {...iconSmall} />
      </Avatar>
    </BadgeAvatars>
  ) : (
    <Avatar src={normalizedIconSrc} sx={iconMedium}>
      <ConnectionIcon {...iconSmall} />
    </Avatar>
  );

  return (
    <ChipWrapper
      label={title}
      onClick={
        isPingEnabled
          ? (event) => {
              event.stopPropagation();
              handlePing?.();
            }
          : undefined
      }
      onDelete={disabled ? undefined : onDelete}
      disabled={disabled}
      avatar={avatar}
      data-cy="chipContextName"
      style={width ? { width } : undefined}
    />
  );
};

export const ConnectionChip = memo(ConnectionChip_);

const TooltipWrappedConnectionChip_ = ({
  tooltip,
  title,
  ...props
}: TooltipWrappedConnectionChipProps) => {
  const chip = (
    <span style={{ display: 'inline-block' }}>
      <ConnectionChip title={title} {...props} />
    </span>
  );

  if (!tooltip && !title) {
    return chip;
  }

  return (
    <CustomTooltip title={tooltip || title} placement="left">
      {chip}
    </CustomTooltip>
  );
};

export const TooltipWrappedConnectionChip = memo(TooltipWrappedConnectionChip_);

const ConnectionStateChip_ = ({ status, actionable = false }: ConnectionStateChipProps) => {
  const normalizedStatus = status?.toLowerCase() || '';
  const value =
    actionable && normalizedStatus
      ? CONNECTION_STATE_TO_TRANSITION_MAP[normalizedStatus] || status
      : status;
  const { Component, avatar } = STATE_CHIP_CONFIG[normalizedStatus] || DEFAULT_STATE_CHIP;

  return <Component avatar={avatar} label={value || ''} />;
};

export const ConnectionStateChip = memo(ConnectionStateChip_);
