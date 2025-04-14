import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import { DesignTable, ErrorBoundary, Grid, useTheme, WorkspaceViewsTable } from '@layer5/sistent';
import { useState } from 'react';
import { useDeletePattern, usePublishPattern } from './hooks';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import Modal from '@/components/Modal';
import { useNotification } from '@/utils/hooks/useNotification';
import ExportModal from '@/components/ExportModal';
import { EVENT_TYPES } from '@/utils/Enum';

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
      event_type: EVENT_TYPEs.INFO,
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

  const handleDesignDownloadModal = (pattern) => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: pattern,
    }));
  };

  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });

  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };
  const theme = useTheme();
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);
  return (
    <ErrorBoundary>
      <Grid
        container
        xs={12}
        spacing={1}
        sx={{
          margin: 'auto',
          width: 'inherit',
          display: 'flex',
          backgroundColor: theme.palette.background.elevatedComponents,
          flexDirection: 'column',
          gap: '1rem',
        }}
      >
        {isDesignsVisible && (
          <DesignTable
            GenericRJSFModal={Modal}
            designsOfWorkspace={designsOfWorkspace}
            handleBulkWorkspaceDesignDeleteModal={handleBulkWorkspaceDesignDeleteModal}
            handleWorkspaceDesignDeleteModal={handleWorkspaceDesignDeleteModal}
            isAssignAllowed={CAN(
              keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
              keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
            )}
            isRemoveAllowed={CAN(
              keys.REMOVE_DESIGNS_FROM_WORKSPACE.action,
              keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject,
            )}
            isDownloadAllowed={CAN(keys.DOWNLOAD_DESIGN.action, keys.DOWNLOAD_DESIGN.subject)}
            isCopyLinkAllowed={CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject)}
            isDeleteAllowed={CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
            isPublishAllowed={CAN(keys.PUBLISH_DESIGN.action, keys.PUBLISH_DESIGN.subject)}
            isUnpublishAllowed={CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject)}
            publishModalHandler={publishModalHandler}
            useAssignDesignToWorkspaceMutation={useAssignDesignToWorkspaceMutation}
            useUnassignDesignFromWorkspaceMutation={useUnassignDesignFromWorkspaceMutation}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
            useGetWorkspaceDesignsQuery={useGetDesignsOfWorkspaceQuery}
            meshModelModelsData={meshModelModelsData}
            handleCopyUrl={handleCopyUrl}
            handleShowDetails={() => {}}
            handleDownload={handleDesignDownloadModal}
            handlePublish={handlePublish}
            setDesignSearch={setDesignSearch}
          />
        )}
        {isViewVisible && (
          <WorkspaceViewsTable
            isAssignAllowed={CAN(
              keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
              keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
            )}
            isRemoveAllowed={CAN(
              keys.REMOVE_VIEWS_FROM_WORKSPACE.action,
              keys.REMOVE_VIEWS_FROM_WORKSPACE.subject,
            )}
            handleShowDetails={() => {}}
            useAssignViewToWorkspaceMutation={useAssignViewToWorkspaceMutation}
            useGetViewsOfWorkspaceQuery={useGetViewsOfWorkspaceQuery}
            useUnassignViewFromWorkspaceMutation={useUnassignViewFromWorkspaceMutation}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
          />
        )}
      </Grid>
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDownload}
      />
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
