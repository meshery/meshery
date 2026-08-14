/**
 * Workspace switcher modal.
 *
 * Renders the modal chrome around the drawer + content layout used to switch
 * between Recents, Designs, Views, Shared, and individual workspaces. The
 * navigation surface itself lives in `WorkspaceFormModalNav.tsx`.
 */
import React, { useState, FC } from 'react';
import { Box, ModalBody, WorkspaceIcon, useMediaQuery } from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { iconMedium } from 'css/icons.styles';
import { Navigation, HeaderInfo } from './WorkspaceFormModalNav';

const BodyShell = styled(Box)({
  flex: 1,
  minHeight: 0,
  padding: 0,
  display: 'flex',
  flexDirection: 'column',
});

// Match RegistryModal: fixed viewport footprint on desktop, full-screen sheet on mobile.
const StyledWorkspaceModal = styled(Modal)(({ theme }) => ({
  zIndex: theme.zIndex.modal,
  '& .MuiDialog-paperFullScreen': {
    margin: 0,
  },
  '& .MuiDialog-paperFullWidth': {
    width: '90%',
    height: '80%',
    maxHeight: '80vh',
  },
  '& .MuiDialog-paper': {
    display: 'flex',
    flexDirection: 'column',
    maxWidth: '100%',
    [theme.breakpoints.down('md')]: {
      margin: 0,
      width: '100%',
      maxWidth: '100%',
      height: '100%',
      maxHeight: '100%',
      borderRadius: 0,
    },
  },
}));

const StyledModalBody = styled(ModalBody)(() => ({
  flex: 1,
  minHeight: 0,
  padding: 0,
  overflow: 'hidden',
  display: 'flex',
  flexDirection: 'column',
}));

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
    <StyledWorkspaceModal
      isOpen={workspaceModal}
      onClose={closeWorkspaceModal}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
      size="xl"
      isFullScreenModeAllowed={!isSmallScreen}
      disableBodyWrap
    >
      <StyledModalBody>
        <BodyShell>{workspaceModal && <Navigation setHeaderInfo={setHeaderInfo} />}</BodyShell>
      </StyledModalBody>
    </StyledWorkspaceModal>
  );
};

export default WorkspaceFormModal;
