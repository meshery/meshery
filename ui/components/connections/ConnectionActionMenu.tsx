import {
  BottomSheet,
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

const renderMenuItems = ({
  onConfigure,
  onConfigureControllers,
  onCopyLink,
  onClose,
}: Omit<ConnectionActionMenuProps, 'anchorEl' | 'open'>) => (
  <>
    {onConfigure && (
      <ActionListItem>
        <ActionButton
          type="button"
          onClick={onConfigure}
          data-cy="btnConfigureConnection"
          permissionKey={Keys.LifecycleManagementEditConnection}
        >
          <SettingsIcon {...iconMedium} />
          <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
            Configure
          </Typography>
        </ActionButton>
      </ActionListItem>
    )}
    {onConfigureControllers && (
      <ActionListItem>
        <ActionButton
          type="button"
          onClick={onConfigureControllers}
          data-cy="btnConfigureConnectionControllers"
          permissionKey={Keys.LifecycleManagementEditConnection}
        >
          <SettingsIcon {...iconMedium} />
          <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
            Configure Controllers
          </Typography>
        </ActionButton>
      </ActionListItem>
    )}
    {onCopyLink && (
      <ActionListItem>
        <ActionButton
          type="button"
          onClick={() => {
            onCopyLink();
            onClose();
          }}
          data-cy="btnCopyConnectionLink"
        >
          <CopyLinkIcon {...iconMedium} />
          <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
            Copy link
          </Typography>
        </ActionButton>
      </ActionListItem>
    )}
  </>
);

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
  const menuItems = renderMenuItems({
    onConfigure,
    onConfigureControllers,
    onCopyLink,
    onClose,
  });

  if (isMobile) {
    return (
      <BottomSheet open={open} onClose={onClose} title="Connection Actions">
        {menuItems}
      </BottomSheet>
    );
  }

  return (
    <Popover
      open={open}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'left',
      }}
    >
      {menuItems}
    </Popover>
  );
};
