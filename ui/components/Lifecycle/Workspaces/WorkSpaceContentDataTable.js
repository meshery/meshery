import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import {
  Box,
  DesignIcon,
  DesignTable,
  ErrorBoundary,
  Tab,
  Tabs,
  useTheme,
  ViewIcon,
  WorkspaceViewsTable,
} from '@layer5/sistent';
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
import downloadContent from '@/utils/fileDownloader';
import { iconMedium } from 'css/icons.styles';

const WorkSpaceContentDataTable = ({ workspaceId, workspaceName }) => {
  const [designSearch, setDesignSearch] = useState('');
  const { notify } = useNotification();

  const { data: designsOfWorkspace, refetch: refetchPatternData } = useGetDesignsOfWorkspaceQuery(
    {
      workspaceId: workspaceId,
      page: 0,
      pageSize: 10,
      expandUser: true,
      search: designSearch,
    },
    {
      skip: !workspaceId,
    },
  );

  const handleCopyUrl = (type, designName, designId) => {
    notify({
      message: `Link to "${designName}" copied to clipboard`,
      event_type: EVENT_TYPES.INFO,
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

  function CustomTabPanel(props) {
    const { children, value, index, ...other } = props;

    return (
      <div hidden={value !== index} {...other}>
        {value === index && <>{children}</>}
      </div>
    );
  }

  const [value, setValue] = useState(() => {
    if (isDesignsVisible) return 0;
    if (isViewVisible) return 1;
    return 0;
  });

  const handleChange = (event, newValue) => {
    event.stopPropagation();
    setValue(newValue);
  };

  const shouldRenderTabs = isDesignsVisible && isViewVisible;

  const handleOpenDesignInPlayground = (designId, designName) => {
    window.open(`/extension/meshmap?mode=design&type=design&id=${designId}&name=${designName}`);
  };
  const handleOpenViewInPlayground = (designId, designName) => {
    window.open(`/extension/meshmap?mode=operator&type=view&id=${designId}&name=${designName}`);
  };
  return (
    <ErrorBoundary>
      {shouldRenderTabs && (
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={value} onChange={handleChange}>
            {isDesignsVisible && (
              <Tab
                label="Design"
                style={{ width: '50%', maxWidth: '50%', gap: '0.5rem' }}
                icon={<DesignIcon />}
                iconPosition="start"
              />
            )}
            {isViewVisible && (
              <Tab
                label="Views"
                style={{ width: '50%', maxWidth: '50%', gap: '0.5rem' }}
                icon={<ViewIcon {...iconMedium} fill={theme.palette.icon.brand} />}
                iconPosition="start"
              />
            )}
          </Tabs>
        </Box>
      )}

      {isDesignsVisible && (
        <CustomTabPanel value={value} index={0}>
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
            handleShowDetails={handleOpenDesignInPlayground}
            handleDownload={handleDesignDownloadModal}
            handlePublish={handlePublish}
            setDesignSearch={setDesignSearch}
          />
        </CustomTabPanel>
      )}

      {isViewVisible && (
        <CustomTabPanel value={value} index={1}>
          <WorkspaceViewsTable
            isAssignAllowed={CAN(
              keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
              keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
            )}
            isRemoveAllowed={CAN(
              keys.REMOVE_VIEWS_FROM_WORKSPACE.action,
              keys.REMOVE_VIEWS_FROM_WORKSPACE.subject,
            )}
            handleShowDetails={handleOpenViewInPlayground}
            useAssignViewToWorkspaceMutation={useAssignViewToWorkspaceMutation}
            useGetViewsOfWorkspaceQuery={useGetViewsOfWorkspaceQuery}
            useUnassignViewFromWorkspaceMutation={useUnassignViewFromWorkspaceMutation}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
          />
        </CustomTabPanel>
      )}

      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDownload}
      />
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
