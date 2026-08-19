import {
  BottomSheet,
  ListItemIcon,
  MenuItem,
  MenuList,
  Popover,
  Typography,
  SettingsIcon,
  CopyLinkIcon,
  useMediaQuery,
  useTheme,
} from '@sistent/sistent';
import { ActionButton, ActionListItem } from './styles';
import { Keys } from '@meshery/schemas/permissions';
import { iconMedium } from '../../css/icons.styles';

type ConnectionActionMenuProps = {
  anchorEl: HTMLElement | null;
  open: boolean;
  onClose: () => void;
  onConfigure?: () => void;
  onConfigureControllers?: () => void;
  onCopyLink?: () => void;
};

export const ConnectionActionMenu = ({
  anchorEl,
  open,
  onClose,
  onConfigure,
  onConfigureControllers,
  onCopyLink,
}: ConnectionActionMenuProps) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const actions = [
    {
      show: !!onConfigure,
      label: 'Configure',
      icon: <SettingsIcon {...iconMedium} />,
      onClick: () => {
        onConfigure?.();
        onClose();
      },
      testId: 'btnConfigureConnection',
      permissionKey: Keys.LifecycleManagementEditConnection,
    },
    {
      show: !!onConfigureControllers,
      label: 'Configure Controllers',
      icon: <SettingsIcon {...iconMedium} />,
      onClick: () => {
        onConfigureControllers?.();
        onClose();
      },
      testId: 'btnConfigureConnectionControllers',
      permissionKey: Keys.LifecycleManagementEditConnection,
    },
    {
      show: !!onCopyLink,
      label: 'Copy link',
      icon: <CopyLinkIcon {...iconMedium} />,
      onClick: () => {
        onCopyLink?.();
        onClose();
      },
      testId: 'btnCopyConnectionLink',
    },
  ].filter((a) => a.show);

  if (isMobile) {
    return (
      <BottomSheet open={open} onClose={onClose} title="Connection Actions">
        <MenuList disablePadding>
          {actions.map(({ label, icon, onClick, testId, permissionKey }) => (
            <MenuItem key={label} onClick={onClick} data-cy={testId} permissionKey={permissionKey}>
              <ListItemIcon>{icon}</ListItemIcon>
              {label}
            </MenuItem>
          ))}
        </MenuList>
      </BottomSheet>
    );
  }

  return (
    <Popover
      open={open}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
    >
      {actions.map(({ label, icon, onClick, testId, permissionKey }) => (
        <ActionListItem key={label}>
          <ActionButton
            type="button"
            onClick={onClick}
            data-cy={testId}
            permissionKey={permissionKey}
          >
            {icon}
            <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
              {label}
            </Typography>
          </ActionButton>
        </ActionListItem>
      ))}
    </Popover>
  );
};
