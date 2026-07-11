import { Button, Popover, Typography, SettingsIcon, CopyLinkIcon } from '@sistent/sistent';
import { ActionListItem } from './styles';
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
      {onConfigure && (
        <ActionListItem>
          <Button type="button" onClick={onConfigure} data-cy="btnConfigureConnection">
            <SettingsIcon {...iconMedium} />
            <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
              Configure
            </Typography>
          </Button>
        </ActionListItem>
      )}
      {onConfigureControllers && (
        <ActionListItem>
          <Button
            type="button"
            onClick={onConfigureControllers}
            data-cy="btnConfigureConnectionControllers"
          >
            <SettingsIcon {...iconMedium} />
            <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
              Configure Controllers
            </Typography>
          </Button>
        </ActionListItem>
      )}
      {onCopyLink && (
        <ActionListItem>
          <Button
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
          </Button>
        </ActionListItem>
      )}
    </Popover>
  );
};
