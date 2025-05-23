import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetWorkspacesQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { RESOURCE_TYPE } from '@/utils/Enum';
import {
  AddCircleIcon,
  Button,
  CircularProgress,
  Modal,
  ModalBody,
  ModalFooter,
  PrimaryActionButtons,
  styled,
  Typography,
  useTheme,
} from '@layer5/sistent';
import React, { useContext, useState } from 'react';
import { useSelector } from 'react-redux';
import { useGetIconBasedOnMode } from './hooks';
import { useRouter } from 'next/router';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

const WorkspaceItem = styled('div')(({ theme, selected, disabled }) => ({
  padding: '12px 16px',
  borderRadius: '8px',
  marginBottom: '8px',
  transition: 'all 0.2s ease',
  backgroundColor: selected ? `${theme.palette.background.brand.default}20` : 'transparent',
  border: selected
    ? `1px solid ${theme.palette.background.brand.default}`
    : '1px solid transparent',
  opacity: disabled ? 0.5 : 1,
  '&:hover': {
    backgroundColor: disabled ? 'transparent' : `${theme.palette.background.brand.default}10`,
  },
}));

const CurrentWorkspaceSection = styled('div')(({ theme }) => ({
  padding: '1rem',
  borderRadius: '0.5rem',
  marginBottom: '1rem',
  display: 'flex',
  border: `1px solid ${theme.palette.border.strong}`,
  fontWeight: 'bolder',
}));

const NoWorkspacesContainer = styled('div')(({ theme }) => ({
  padding: '1rem',
  textAlign: 'center',
  color: theme.palette.text.secondary,
  border: `1px dashed ${theme.palette.border.strong}`,
  borderRadius: '8px',
  margin: '1rem 0',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: '1rem',
}));

const getModalTitle = (type, selectedContent, multiSelectedContent) => {
  const itemCount = multiSelectedContent?.length || 1;
  const resourceName = selectedContent?.name || itemCount;
  const resourceType = type === RESOURCE_TYPE.DESIGN ? 'Design' : 'View';

  return `Move ${resourceName} ${resourceType}${itemCount > 1 ? 's' : ''}`;
};

const WorkspaceContentMoveModal = ({
  workspaceContentMoveModal,
  setWorkspaceContentMoveModal,
  currentWorkspace,
  type,
  selectedContent,
}) => {
  const router = useRouter();
  const { notify } = useNotification();
  const { setMultiSelectedContent, multiSelectedContent, closeModal } =
    useContext(WorkspaceModalContext);
  const { organization: currentOrg } = useSelector((state) => state.ui);
  const { id: orgId } = currentOrg;
  const { data: workspaceData, isLoading } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      order: 'updated_at desc',
      orgId: orgId,
    },
    {
      skip: !orgId,
    },
  );

  const filteredWorkspaces = workspaceData?.workspaces?.filter(
    (workspace) => workspace.id !== currentWorkspace.id,
  );
  const [selectedWorkspaceForMove, setSelectedWorkspaceForMove] = useState(null);
  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();
  const [unAssignDesignFromWorkspace] = useUnassignDesignFromWorkspaceMutation();
  const [assignViewToWorkspace] = useAssignViewToWorkspaceMutation();
  const [unAssignViewFromWorkspace] = useUnassignViewFromWorkspaceMutation();

  const handleMove = async () => {
    setWorkspaceContentMoveModal(false);

    try {
      const moveDesign = async (designId) => {
        await assignDesignToWorkspace({
          workspaceId: selectedWorkspaceForMove.id,
          designId,
        }).unwrap();
        await unAssignDesignFromWorkspace({
          workspaceId: currentWorkspace.id,
          designId,
        });
      };

      const moveView = async (viewId) => {
        await assignViewToWorkspace({
          workspaceId: selectedWorkspaceForMove.id,
          viewId,
        }).unwrap();
        await unAssignViewFromWorkspace({
          workspaceId: currentWorkspace.id,
          viewId,
        });
      };

      if (RESOURCE_TYPE.DESIGN === type) {
        if (multiSelectedContent.length > 0) {
          await Promise.all(multiSelectedContent.map((design) => moveDesign(design.id)));
          setMultiSelectedContent([]);
        }
        if (selectedContent) {
          await moveDesign(selectedContent.id);
        }
      } else {
        if (multiSelectedContent.length > 0) {
          await Promise.all(multiSelectedContent.map((view) => moveView(view.id)));
          setMultiSelectedContent([]);
        }
        if (selectedContent) {
          await moveView(selectedContent.id);
        }
      }

      notify({
        message: `Successfully moved ${type === RESOURCE_TYPE.DESIGN ? 'design' : 'view'}${multiSelectedContent.length > 1 ? 's' : ''} to ${selectedWorkspaceForMove.name}`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (error) {
      notify({
        message: `Failed to move ${type === RESOURCE_TYPE.DESIGN ? 'design' : 'view'}. Please try again.`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  };

  const handleCreateWorkspace = () => {
    closeModal();
    setWorkspaceContentMoveModal(false);
    router.push('/management/workspaces');
  };

  const isMoveDesignAllowed =
    CAN(keys.ASSIGN_DESIGNS_TO_WORKSPACE.action, keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject) &&
    CAN(keys.REMOVE_DESIGNS_FROM_WORKSPACE.action, keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject);

  const isMoveViewAllowed =
    CAN(keys.ASSIGN_VIEWS_TO_WORKSPACE.action, keys.ASSIGN_VIEWS_TO_WORKSPACE.subject) &&
    CAN(keys.REMOVE_VIEWS_FROM_WORKSPACE.action, keys.REMOVE_VIEWS_FROM_WORKSPACE.subject);

  const isMoveAllowed = type === RESOURCE_TYPE.DESIGN ? isMoveDesignAllowed : isMoveViewAllowed;

  return (
    <Modal
      open={workspaceContentMoveModal}
      headerIcon={useGetIconBasedOnMode({ mode: type })}
      closeModal={() => setWorkspaceContentMoveModal(false)}
      title={getModalTitle(type, selectedContent, multiSelectedContent)}
    >
      <ModalBody>
        <CurrentWorkspaceSection>
          Current Workspace: {currentWorkspace.name}
        </CurrentWorkspaceSection>

        {isLoading ? (
          <CircularProgress size={24} />
        ) : !filteredWorkspaces?.length ? (
          <NoWorkspacesContainer>
            <Typography>No other workspaces available to move content to.</Typography>
            <Button
              variant="contained"
              color="primary"
              onClick={handleCreateWorkspace}
              disabled={!CAN(keys.CREATE_WORKSPACE.action, keys.CREATE_WORKSPACE.subject)}
            >
              Create Workspace
            </Button>
          </NoWorkspacesContainer>
        ) : (
          <>
            <Typography style={{ marginBottom: '0.5rem' }}>Select destination workspace</Typography>
            {filteredWorkspaces.map((workspace) => (
              <WorkspaceItem
                key={workspace.id}
                selected={selectedWorkspaceForMove?.id === workspace.id}
                onClick={() => isMoveAllowed && setSelectedWorkspaceForMove(workspace)}
                disabled={!isMoveAllowed}
              >
                {workspace.name}
              </WorkspaceItem>
            ))}
          </>
        )}
      </ModalBody>
      <ModalFooter variant="filled">
        <PrimaryActionButtons
          primaryText={'Move'}
          secondaryText="Cancel"
          primaryButtonProps={{
            onClick: handleMove,
            disabled:
              isLoading ||
              !selectedWorkspaceForMove ||
              !filteredWorkspaces?.length ||
              !isMoveAllowed,
          }}
          secondaryButtonProps={{
            onClick: () => setWorkspaceContentMoveModal(false),
          }}
        />
      </ModalFooter>
    </Modal>
  );
};

export default WorkspaceContentMoveModal;
