import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import { DesignTable, ErrorBoundary, Grid, WorkspaceViewsTable } from '@layer5/sistent';
import { useState } from 'react';
import { useDeletePattern, usePublishPattern } from './hooks';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { WORKSPACE_ACTION_TYPES } from '.';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import Modal from '@/components/Modal';
import { useNotification } from '@/utils/hooks/useNotification';

const WorkSpaceContentDataTable = ({ workspaceId, workspaceName }) => {
  const [designSearch, setDesignSearch] = useState('');
  const { notify } = useNotification();

  const { data: designsOfWorkspace, refetch: refetchPatternData } = useGetDesignsOfWorkspaceQuery({
    workspaceId: workspaceId,
    page: 0,
    pageSize: 10,
    expandUser: true,
    search: designSearch,
  });

  const handleCopyUrl = (type, designName, designId) => {
    notify({
      message: `Link to "${designName}" copied to clipboard`,
      event_type: WORKSPACE_ACTION_TYPES.INFO,
    });
    navigator.clipboard.writeText(`${MESHERY_CLOUD_PROD}/catalog/content/${type}/${designId}`);
  };

  const { handleWorkspaceDesignDeleteModal, handleBulkWorkspaceDesignDeleteModal } =
    useDeletePattern();

  const { data: meshModelModelsData } = useGetMeshModelsQuery({
    page: 0,
    pagesize: 'all',
  });

  const { handlePublish, handlePublishModal: publishModalHandler } = usePublishPattern(
    meshModelModelsData,
    refetchPatternData,
  );
  return (
    <ErrorBoundary>
      <Grid
        container
        xs={12}
        spacing={1}
        sx={{
          margin: 'auto',
          backgroundColor: '#f3f1f1',
          paddingLeft: '0.5rem',
          borderRadius: '0.25rem',
          width: 'inherit',
          display: 'flex',
          flexDirection: 'column',
          gap: '1rem',
        }}
      >
        <DesignTable
          GenericRJSFModal={Modal}
          designsOfWorkspace={designsOfWorkspace}
          handleBulkWorkspaceDesignDeleteModal={handleBulkWorkspaceDesignDeleteModal}
          handleWorkspaceDesignDeleteModal={handleWorkspaceDesignDeleteModal}
          isAssignAllowed={CAN(
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
          )}
          isRemoveAllowed={CAN(
            keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject,
            keys.REMOVE_DESIGNS_FROM_WORKSPACE.action,
          )}
          isDownloadAllowed={CAN(keys.DOWNLOAD_DESIGN.subject, keys.DOWNLOAD_DESIGN.action)}
          isCopyLinkAllowed={CAN(keys.CLONE_DESIGN.subject, keys.CLONE_DESIGN.action)}
          isDeleteAllowed={CAN(keys.DELETE_A_DESIGN.subject, keys.DELETE_A_DESIGN.action)}
          isPublishAllowed={CAN(keys.PUBLISH_DESIGN.subject, keys.PUBLISH_DESIGN.action)}
          isUnpublishAllowed={CAN(keys.UNPUBLISH_DESIGN.subject, keys.UNPUBLISH_DESIGN.action)}
          publishModalHandler={publishModalHandler}
          useAssignDesignToWorkspaceMutation={useAssignDesignToWorkspaceMutation}
          useUnassignDesignFromWorkspaceMutation={useUnassignDesignFromWorkspaceMutation}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          useGetWorkspaceDesignsQuery={useGetDesignsOfWorkspaceQuery}
          meshModelModelsData={meshModelModelsData}
          handleCopyUrl={handleCopyUrl}
          handleShowDetails={() => {}}
          handlePublish={handlePublish}
          setDesignSearch={setDesignSearch}
        />
        <WorkspaceViewsTable
          isAssignAllowed={CAN(
            keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
            keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
          )}
          isRemoveAllowed={CAN(
            keys.REMOVE_VIEWS_FROM_WORKSPACE.subject,
            keys.REMOVE_VIEWS_FROM_WORKSPACE.action,
          )}
          handleShowDetails={() => {}}
          useAssignViewToWorkspaceMutation={useAssignViewToWorkspaceMutation}
          useGetViewsOfWorkspaceQuery={useGetViewsOfWorkspaceQuery}
          useUnassignViewFromWorkspaceMutation={useUnassignViewFromWorkspaceMutation}
          workspaceId={workspaceId}
          workspaceName={workspaceName}
        />
      </Grid>
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
