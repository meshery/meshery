/**
 * Workspace switcher modal.
 *
 * Renders the modal chrome around the drawer + content layout used to switch
 * between Recents, Designs, Views, Shared, and individual workspaces. The
 * navigation surface itself lives in `WorkspaceFormModalNav.tsx`.
 */
import React, { useState, FC } from 'react';
import { Box, WorkspaceIcon, useMediaQuery } from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { iconMedium } from 'css/icons.styles';
import { Navigation, HeaderInfo } from './WorkspaceFormModalNav';

const BodyShell = styled(Box)({
  height: '100%',
  padding: 0,
});

export interface WorkspaceFormModalProps {
  workspaceModal: boolean;
  closeWorkspaceModal: () => void;
}

const WorkspaceFormModal: FC<WorkspaceFormModalProps> = ({
  workspaceModal,
  closeWorkspaceModal,
}) => {
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));
  const [headerInfo, setHeaderInfo] = useState<HeaderInfo>({
    title: 'All Workspaces',
    icon: <WorkspaceIcon {...iconMedium} secondaryFill={theme.palette.icon.neutral.default} />,
  });

  return (
    <Modal
      isOpen={workspaceModal}
      onClose={closeWorkspaceModal}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
      isFullScreenModeAllowed={!isSmallScreen}
      size="xl"
    >
      <BodyShell>{workspaceModal && <Navigation setHeaderInfo={setHeaderInfo} />}</BodyShell>
    </Modal>
  );
};

export default WorkspaceFormModal;
