import React from 'react';
import {
  Button,
  Popover,
  Typography,
  SyncAltIcon,
  SettingsIcon,
  CopyLinkIcon,
} from '@sistent/sistent';
import { ActionListItem } from './styles';
import { iconMedium } from '../../css/icons.styles';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { MESHSYNC_DEPLOYMENT_TYPE } from '../../utils/Enum';

type ConnectionActionMenuProps = {
  anchorEl: HTMLElement | null;
  open: boolean;
  onClose: () => void;
  onFlushMeshSync: () => void;
  onDeploymentModeAnchor: (event: React.MouseEvent<HTMLElement>) => void;
  onConfigure?: () => void;
  onConfigureControllers?: () => void;
  onCopyLink?: () => void;
};

export const ConnectionActionMenu = ({
  anchorEl,
  open,
  onClose,
  onFlushMeshSync,
  onDeploymentModeAnchor,
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
      <ActionListItem>
        <Button
          type="submit"
          onClick={onFlushMeshSync}
          data-cy="btnResetDatabase"
          disabled={!CAN(keys.FLUSH_MESHSYNC_DATA.action, keys.FLUSH_MESHSYNC_DATA.subject)}
        >
          <SyncAltIcon {...iconMedium} />
          <Typography variant="body1" style={{ marginLeft: '0.5rem' }}>
            Flush MeshSync
          </Typography>
        </Button>
      </ActionListItem>
      <ActionListItem>
        <Button
          type="submit"
          onClick={(e) => {
            e.stopPropagation();
            onDeploymentModeAnchor(e);
          }}
          data-cy="btnChangeDeploymentMode"
        >
          <Typography variant="body1">Modify Deployment Mode</Typography>
        </Button>
      </ActionListItem>
    </Popover>
  );
};

type ConnectionDeploymentModeMenuProps = {
  anchorEl: HTMLElement | null;
  open: boolean;
  onClose: () => void;
  onSelectMode: (mode: string) => void;
};

export const ConnectionDeploymentModeMenu = ({
  anchorEl,
  open,
  onClose,
  onSelectMode,
}: ConnectionDeploymentModeMenuProps) => {
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
      <ActionListItem>
        <Button onClick={() => onSelectMode(MESHSYNC_DEPLOYMENT_TYPE.OPERATOR)}>
          <Typography variant="body1">Operator</Typography>
        </Button>
      </ActionListItem>
      <ActionListItem>
        <Button onClick={() => onSelectMode(MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED)}>
          <Typography variant="body1">Embedded</Typography>
        </Button>
      </ActionListItem>
    </Popover>
  );
};
