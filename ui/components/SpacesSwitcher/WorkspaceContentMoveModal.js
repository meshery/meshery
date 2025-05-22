import PatternIcon from '@/assets/icons/Pattern';
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
  CircularProgress,
  DesignIcon,
  Modal,
  ModalBody,
  ModalFooter,
  PrimaryActionButtons,
  Skeleton,
  styled,
  Typography,
  ViewIcon,
} from '@layer5/sistent';
import { iconMedium } from 'css/icons.styles';
import React, { useContext, useState } from 'react';
import { useSelector } from 'react-redux';
import { getIconBasedOnMode } from './components';

const WorkspaceItem = styled('div')(({ theme, selected }) => ({
  padding: '12px 16px',
  borderRadius: '8px',
  cursor: 'pointer',
  marginBottom: '8px',
  transition: 'all 0.2s ease',
  backgroundColor: selected ? `${theme.palette.background.brand.default}20` : 'transparent',
  border: selected
    ? `1px solid ${theme.palette.background.brand.default}`
    : '1px solid transparent',
  '&:hover': {
    backgroundColor: `${theme.palette.background.brand.default}10`,
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
  const { setMultiSelectedContent, multiSelectedContent } = useContext(WorkspaceModalContext);
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
      skip: !orgId ? true : false,
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

  return (
    <Modal
      open={workspaceContentMoveModal}
      headerIcon={getIconBasedOnMode({ mode: type })}
      closeModal={() => {
        setWorkspaceContentMoveModal(false);
      }}
      title={getModalTitle(type, selectedContent, multiSelectedContent)}
    >
      <ModalBody>
        <CurrentWorkspaceSection>
          Current Workspace: {currentWorkspace.name}
        </CurrentWorkspaceSection>

        {isLoading ? (
          <CircularProgress size={24} />
        ) : (
          <>
            <Typography style={{ marginBottom: '0.5rem' }}>Select destination workspace</Typography>
            {filteredWorkspaces?.map((workspace) => (
              <WorkspaceItem
                key={workspace.id}
                selected={selectedWorkspaceForMove?.id === workspace.id}
                onClick={() => setSelectedWorkspaceForMove(workspace)}
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
            onClick: () => {
              setWorkspaceContentMoveModal(false);
              {
                if (RESOURCE_TYPE.DESIGN === type) {
                  multiSelectedContent.length > 0 &&
                    multiSelectedContent.map((id) => {
                      return assignDesignToWorkspace({
                        workspaceId: selectedWorkspaceForMove.id,
                        designId: id,
                      })
                        .unwrap()
                        .then(() => {
                          unAssignDesignFromWorkspace({
                            workspaceId: currentWorkspace.id,
                            designId: id,
                          });
                        })
                        .then(() => {
                          setMultiSelectedContent([]);
                        });
                    });
                  if (selectedContent) {
                    assignDesignToWorkspace({
                      workspaceId: selectedWorkspaceForMove.id,
                      designId: selectedContent.id,
                    })
                      .unwrap()
                      .then(() => {
                        unAssignDesignFromWorkspace({
                          workspaceId: currentWorkspace.id,
                          designId: selectedContent.id,
                        });
                      });
                  }
                } else {
                  multiSelectedContent.length > 0 &&
                    multiSelectedContent.map((id) => {
                      return assignViewToWorkspace({
                        workspaceId: selectedWorkspaceForMove.id,
                        viewId: id,
                      })
                        .unwrap()
                        .then(() => {
                          unAssignViewFromWorkspace({
                            workspaceId: currentWorkspace.id,
                            viewId: id,
                          });
                        })
                        .then(() => {
                          setMultiSelectedContent([]);
                        });
                    });
                  if (selectedContent) {
                    assignViewToWorkspace({
                      workspaceId: selectedWorkspaceForMove.id,
                      viewId: selectedContent.id,
                    })
                      .unwrap()
                      .then(() => {
                        unAssignViewFromWorkspace({
                          workspaceId: currentWorkspace.id,
                          viewId: selectedContent.id,
                        });
                      });
                  }
                }
              }
            },
            disabled: isLoading || !selectedWorkspaceForMove,
          }}
          secondaryButtonProps={{
            onClick: () => {
              setWorkspaceContentMoveModal(false);
            },
          }}
        />
      </ModalFooter>
    </Modal>
  );
};

export default WorkspaceContentMoveModal;
