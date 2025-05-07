import React, { useState } from 'react';
import {
  Grid,
  L5DeleteIcon,
  Modal,
  Pagination,
  PaginationItem,
  ModalBody,
  ModalButtonPrimary,
  ModalFooter,
  ModalButtonSecondary,
  useTheme,
  ErrorBoundary,
} from '@layer5/sistent';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { useDeleteWorkspaceMutation } from '../../../rtk-query/workspace';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { useNotificationHandlers } from '@/utils/hooks/useNotification';
import { UserCommonBox } from './styles';
import MesheryWorkspaceCard from './MesheryWorkspaceCard';
import { debounce } from 'lodash';

const WorkspaceGridView = ({
  workspacesData,
  handleWorkspaceModalOpen,
  handleDeleteWorkspaceConfirm,
  totalPages,
  page,
  setPage,
}) => {
  const [deleteWorkspacesModal, setDeleteWorkspacesModal] = useState(false);
  const [selectedWorkspaces, setSelectedWorkspaces] = useState([]);
  const [deleteWorkspace] = useDeleteWorkspaceMutation();

  const handleDeleteWorkspacesModalClose = () => {
    setDeleteWorkspacesModal(false);
  };

  const handleDeleteWorkspacesModalOpen = () => {
    setDeleteWorkspacesModal(true);
  };

  const { handleSuccess, handleError } = useNotificationHandlers();
  const handleDeleteWorkspace = (id) => {
    deleteWorkspace({
      workspaceId: id,
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace deleted`))
      .catch((error) => handleError(`Workspace Delete Error: ${error?.data}`));
  };

  const handleBulkDeleteEnv = () => {
    selectedWorkspaces.map((workspaceId) => {
      handleDeleteWorkspace(workspaceId);
    });
    setSelectedWorkspaces([]);
    handleDeleteWorkspacesModalClose();
  };

  const handleBulkSelect = (e, id) => {
    const isChecked = e.target.checked;
    if (isChecked) {
      setSelectedWorkspaces([...selectedWorkspaces, id]);
    } else {
      const newSelectedEnv = selectedWorkspaces.filter((env) => env !== id);
      setSelectedWorkspaces(newSelectedEnv);
    }
  };

  const theme = useTheme();
  return (
    <ErrorBoundary>
      <Grid container spacing={2} sx={{ marginTop: '-16px' }}>
        {selectedWorkspaces?.length > 0 && (
          <UserCommonBox
            sx={{
              width: '100%',
              p: '0.8rem',
              justifyContent: 'space-between',
              marginLeft: '1rem',
              marginTop: '0.18rem',
            }}
          >
            <span>
              {selectedWorkspaces?.length > 1
                ? `${selectedWorkspaces?.length} workspaces selected`
                : `${selectedWorkspaces?.length} workspace selected`}
            </span>

            <L5DeleteIcon
              onClick={handleDeleteWorkspacesModalOpen}
              disabled={
                selectedWorkspaces.length > 0
                  ? !CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)
                  : true
              }
            />
          </UserCommonBox>
        )}
        {workspacesData?.map((workspace) => (
          <Grid item xs={12} md={6} key={workspace.id}>
            <MesheryWorkspaceCard
              key={workspace.id}
              workspaceDetails={workspace}
              handleWorkspaceModalOpen={handleWorkspaceModalOpen}
              handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
              handleBulkSelect={handleBulkSelect}
              selectedWorkspaces={selectedWorkspaces}
            />
          </Grid>
        ))}
      </Grid>
      <Grid
        container
        sx={{ padding: '2rem 0' }}
        flex
        justifyContent="center"
        spacing={2}
        position="relative"
      >
        <Pagination
          count={totalPages}
          page={page + 1}
          onChange={debounce((_, page) => setPage(page - 1), 150)}
          boundaryCount={3}
          renderItem={(item) => (
            <PaginationItem
              slots={{ previous: ChevronLeftIcon, next: ChevronRightIcon }}
              {...item}
            />
          )}
        />
      </Grid>
      {CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject) && (
        <Modal
          open={deleteWorkspacesModal}
          closeModal={handleDeleteWorkspacesModalClose}
          title={'Delete Workspace'}
        >
          <ModalBody
            style={{
              padding: '2rem 4rem',
            }}
          >
            <div>{`Do you want to delete ${selectedWorkspaces.length} workspace(s) ?`}</div>
          </ModalBody>
          <ModalFooter variant="filled">
            <ModalButtonSecondary
              style={{ marginRight: '16PX' }}
              onClick={handleDeleteWorkspacesModalClose}
            >
              Cancel
            </ModalButtonSecondary>
            <ModalButtonPrimary
              style={{
                background: theme.palette.background.error.default,
              }}
              onClick={handleBulkDeleteEnv}
            >
              Delete
            </ModalButtonPrimary>
          </ModalFooter>
        </Modal>
      )}
    </ErrorBoundary>
  );
};

export default WorkspaceGridView;
