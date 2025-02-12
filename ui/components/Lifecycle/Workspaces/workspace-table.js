import {
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useAssignTeamToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetTeamsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUnassignTeamFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import {
  DesignTable,
  EnvironmentTable,
  ErrorBoundary,
  WorkspaceTeamsTable,
  WorkspaceViewsTable,
  Grid,
} from '@layer5/sistent';
import { TableCell } from '@mui/material';
import { useState } from 'react';
import { useDeletePattern, usePublishPattern } from './hooks';
import { useLegacySelector } from 'lib/store';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import React from 'react';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import Modal from '@/components/Modal';

import { useGetUsersForOrgQuery, useRemoveUserFromTeamMutation } from '@/rtk-query/user';
import { useNotification, useNotificationHandlers } from '@/utils/hooks/useNotification';
import downloadContent from '@/utils/fileDownloader';
import ExportModal from '@/components/ExportModal';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { EVENT_TYPES } from 'lib/event-types';

const WorkSpaceDataTable = ({ rowData }) => {
  const [designSearch, setDesignSearch] = useState('');
  const org_id = useLegacySelector((state) => state.get('organization')).id;
  const { notify } = useNotification();

  const { data: designsOfWorkspace, refetch: refetchPatternData } = useGetDesignsOfWorkspaceQuery({
    workspaceId: rowData[0],
    page: 0,
    pagesize: 10,
    expandUser: true,
    search: designSearch,
  });
  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });
  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDesignDownloadModal = (pattern) => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: pattern,
    }));
  };

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

  const handleCopyUrl = (type, designName, designId) => {
    notify({
      message: `Link to "${designName}" copied to clipboard`,
      event_type: EVENT_TYPES.INFO,
    });
    navigator.clipboard.writeText(`${MESHERY_CLOUD_PROD}/catalog/content/${type}/${designId}`);
  };

  const { handlePublish, handlePublishModal: publishModalHandler } = usePublishPattern(
    meshModelModelsData,
    refetchPatternData,
  );

  const { handleWorkspaceDesignDeleteModal, handleBulkWorkspaceDesignDeleteModal } =
    useDeletePattern();

  const { data: meshModelModelsData } = useGetMeshModelsQuery({
    page: 0,
    pagesize: 'all',
  });
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isTeamsVisible = CAN(keys.VIEW_TEAMS.action, keys.VIEW_TEAMS.subject);
  const isEnvironmentsVisible = CAN(keys.VIEW_ENVIRONMENTS.action, keys.VIEW_ENVIRONMENTS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);
  return (
    <>
      <ErrorBoundary>
        <TableCell
          colSpan={rowData.length + 1}
          sx={{
            padding: '0.5rem',
            backgroundColor: 'rgba(0, 0, 0, 0.05)',
          }}
        >
          <Grid
            container
            xs={12}
            spacing={1}
            sx={{
              margin: 'auto',
              borderRadius: '0.25rem',
              width: 'inherit',
              display: 'flex',
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
                isUnpublishAllowed={CAN(
                  keys.UNPUBLISH_DESIGN.action,
                  keys.UNPUBLISH_DESIGN.subject,
                )}
                publishModalHandler={publishModalHandler}
                useAssignDesignToWorkspaceMutation={useAssignDesignToWorkspaceMutation}
                useUnassignDesignFromWorkspaceMutation={useUnassignDesignFromWorkspaceMutation}
                workspaceId={rowData[0]}
                workspaceName={rowData[1]}
                useGetWorkspaceDesignsQuery={useGetDesignsOfWorkspaceQuery}
                meshModelModelsData={meshModelModelsData}
                handleCopyUrl={handleCopyUrl}
                handleShowDetails={() => {}}
                handlePublish={handlePublish}
                handleDownload={handleDesignDownloadModal}
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
                workspaceId={rowData[0]}
                workspaceName={rowData[1]}
              />
            )}
            {isTeamsVisible && (
              <WorkspaceTeamsTable
                workspaceId={rowData[0]}
                isAssignTeamAllowed={CAN(
                  keys.ASSIGN_TEAM_TO_WORKSPACE.action,
                  keys.ASSIGN_TEAM_TO_WORKSPACE.subject,
                )}
                isDeleteTeamAllowed={CAN(keys.DELETE_TEAM.action, keys.DELETE_TEAM.subject)}
                isEditTeamAllowed={CAN(keys.EDIT_TEAM.action, keys.EDIT_TEAM.subject)}
                isLeaveTeamAllowed={CAN(keys.LEAVE_TEAM.action, keys.LEAVE_TEAM.subject)}
                isRemoveTeamFromWorkspaceAllowed={CAN(
                  keys.REMOVE_TEAM_FROM_WORKSPACE.action,
                  keys.REMOVE_TEAM_FROM_WORKSPACE.subject,
                )}
                useAssignTeamToWorkspaceMutation={useAssignTeamToWorkspaceMutation}
                useGetTeamsOfWorkspaceQuery={useGetTeamsOfWorkspaceQuery}
                useUnassignTeamFromWorkspaceMutation={useUnassignTeamFromWorkspaceMutation}
                workspaceName={rowData[1]}
                fetchTeamUsers={() => {}}
                org_id={org_id}
                useGetUsersForOrgQuery={useGetUsersForOrgQuery}
                useNotificationHandlers={useNotificationHandlers}
                useRemoveUserFromTeamMutation={useRemoveUserFromTeamMutation}
              />
            )}
            {isEnvironmentsVisible && (
              <EnvironmentTable
                workspaceId={rowData[0]}
                workspaceName={rowData[1]}
                isAssignAllowed={CAN(
                  keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.action,
                  keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.subject,
                )}
                isRemoveAllowed={CAN(
                  keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.action,
                  keys.REMOVE_ENVIRONMENT_FROM_WORKSPACE.subject,
                )}
                useGetEnvironmentsOfWorkspaceQuery={useGetEnvironmentsOfWorkspaceQuery}
                useUnassignEnvironmentFromWorkspaceMutation={
                  useUnassignEnvironmentFromWorkspaceMutation
                }
                useAssignEnvironmentToWorkspaceMutation={useAssignEnvironmentToWorkspaceMutation}
              />
            )}
          </Grid>
        </TableCell>
      </ErrorBoundary>
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDownload}
      />
    </>
  );
};

export default WorkSpaceDataTable;
