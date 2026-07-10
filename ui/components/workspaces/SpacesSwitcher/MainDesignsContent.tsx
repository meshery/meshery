import { getDesign, useUpdatePatternFileMutation } from '@/rtk-query/design';
import { getUserAccessToken, getUserProfile, useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  DeleteIcon,
  Divider,
  ExportIcon,
  InfoIcon,
  ListItem,
  ListItemText,
  MergeOutlinedIcon,
  Modal,
  OutlinedPatternIcon,
  PromptComponent,
  ShareIcon,
  WorkspaceContentMoveModal,
  useModal,
  useRoomActivity,
  useTheme,
} from '@sistent/sistent';
import React, { useCallback, useContext, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, {
  handleUpdatePatternVisibility,
  useContentDelete,
  useContentDownload,
} from './hooks';
import { MenuComponent } from './MenuComponent';
import { DesignList, GhostContainer, GhostImage, GhostText, LoadingContainer } from './styles';
import ExportDesignModal from '../../designs/export/ExportDesignModal';
import { RESOURCE_TYPE } from '@/utils/Enum';
import ShareModal from '../ShareWorkspaceModal';
import InfoModal from '../../shared/Modal/Information/InfoModal';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import {
  isDesignOpenInExtension,
  mergeDesignWithCurrent,
  openDesignInExtension,
  useIsDesignerEnabled,
} from '@/utils/utils';
import Router, { useRouter } from 'next/router';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import { useSelector } from 'react-redux';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetWorkspacesQuery,
} from '@/rtk-query/workspace';
import { useNotification } from '@/utils/hooks/useNotification';

const MainDesignsContent = ({
  page,
  setPage,
  isLoading,
  isFetching,
  designs,
  hasMore,
  totalCount,
  workspace,
  refetch,
  isMultiSelectMode,
  showWorkspaceName = true,
  showOrganizationName = true,
}) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [selectedDesign, setSelectedDesign] = useState(null);
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setInfoModal] = useState({ open: false, userId: '' });
  const [moveModal, setMoveModal] = useState(false);
  const modalRef = useRef(true);
  const { handleDelete } = useContentDelete(modalRef);

  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage(page + 1);
  }, [isLoading, isFetching]);

  const { loadingRef } = useInfiniteScroll({
    isLoading: isLoading || isFetching,
    hasMore,
    onLoadMore: loadNextPage,
  });
  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });

  const handleDesignDownloadModal = (design) => {
    setDownloadModal({
      open: true,
      content: design,
    });
  };
  const handleDownloadDialogClose = () => {
    setDownloadModal({
      open: false,
      content: null,
    });
  };
  const { handleDesignDownload } = useContentDownload();

  const handleRemove = (design) => {
    setMoveModal(true);
    setSelectedDesign(design);
  };

  const handleShare = (design) => {
    setShareModal(true);
    setSelectedDesign(design);
  };

  const handleShareClose = () => {
    setShareModal(false);
    setSelectedDesign(null);
  };
  const sistentInfoModal = useModal({
    headerIcon: OutlinedPatternIcon,
  });

  const handleInfoModal = async (design) => {
    const selectedDesignWithPatternFile = await getDesign({
      design_id: design?.id,
    });

    setSelectedDesign(selectedDesignWithPatternFile?.data);

    sistentInfoModal.openModal({
      title: selectedDesign?.name,
    });
    setInfoModal({
      open: true,
      userId: selectedDesignWithPatternFile?.data?.userId,
    });
  };

  const { data: meshModelsData } = useGetMeshModelsQuery(
    {
      params: {
        pagesize: 'all',
        trim: true,
      },
    },
    {
      skip: !infoModal.open,
    },
  );

  const handleInfoModalClose = () => {
    setSelectedDesign(null);
    setInfoModal({ open: false, userId: '' });
  };

  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);
  const [updatePatterns] = useUpdatePatternFileMutation();
  const isDesignerAvailable = useIsDesignerEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const handleOpenDesignInExtension = (designId, designName) => {
    if (workspaceSwitcherContext?.closeModal) {
      workspaceSwitcherContext.closeModal();
    }
    if (!isDesignerAvailable) {
      router.push(`/configuration/designs/configurator?design_id=${designId}`);
      return;
    }

    openDesignInExtension(designId, designName, Router);
  };
  const theme = useTheme();
  const DESIGN_ACTIONS = {
    MERGE_DESIGN: {
      id: 'merge_design',
      title: 'Merge Into Current Design',
      icon: <MergeOutlinedIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        isDesignOpenInExtension() &&
        CAN(Keys.CatalogManagementEditDesign.id, Keys.CatalogManagementEditDesign.function),
    },
    EXPORT_DESIGN: {
      id: 'export_design',
      title: 'Export Design',
      icon: <ExportIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        CAN(
          Keys.CatalogManagementDownloadADesign.id,
          Keys.CatalogManagementDownloadADesign.function,
        ),
    },

    REMOVE_DESIGN: {
      id: 'move',
      title: 'Move Design',
      icon: <MoveFileIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        CAN(
          Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.id,
          Keys.WorkspaceManagementRemoveDesignsFromWorkspaces.function,
        ),
    },
    SHARE_DESIGN: {
      id: 'share',
      title: 'Share Design',
      icon: <ShareIcon fill={theme.palette.icon.default} />,
      enabled: ({ design }) =>
        design?.visibility !== 'published' &&
        CAN(Keys.CatalogManagementShareDesign.id, Keys.CatalogManagementShareDesign.function),
    },
    INFO_DESIGN: {
      id: 'info',
      title: 'Info',
      icon: <InfoIcon fill={theme.palette.icon.default} />,
      enabled: () => true,
    },
    DELETE_DESIGN: {
      id: 'delete',
      title: 'Delete Design',
      icon: <DeleteIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        CAN(Keys.CatalogManagementDeleteADesign.id, Keys.CatalogManagementDeleteADesign.function),
    },
  };

  const handleMerge = (design) => {
    mergeDesignWithCurrent(design.id, design.name);
    if (workspaceSwitcherContext?.closeModal) {
      workspaceSwitcherContext.closeModal();
    }
  };

  const getMenuOptions = ({
    design,
    handleDesignDownloadModal,
    handleDelete,
    handleRemove,
    handleShare,
    handleInfoModal,
    refetch,
  }) => {
    const options = [
      {
        ...DESIGN_ACTIONS.MERGE_DESIGN,
        handler: () => handleMerge(design),
      },
      {
        ...DESIGN_ACTIONS.EXPORT_DESIGN,
        handler: () => handleDesignDownloadModal(design),
      },

      {
        ...DESIGN_ACTIONS.SHARE_DESIGN,
        handler: () => handleShare(design),
      },

      {
        ...DESIGN_ACTIONS.INFO_DESIGN,
        handler: () => handleInfoModal(design),
      },
      {
        ...DESIGN_ACTIONS.DELETE_DESIGN,
        handler: () => handleDelete([design], RESOURCE_TYPE.DESIGN, refetch),
      },
    ];

    if (workspace) {
      options.unshift({
        ...DESIGN_ACTIONS.REMOVE_DESIGN,
        handler: () => handleRemove(design, refetch),
      });
    }
    return options.filter((option) => option.enabled({ design }));
  };
  const isInitialFetch = isFetching && page === 0;
  const isEmpty = totalCount === 0;
  const shouldRenderDesigns = !isEmpty && !isInitialFetch;
  const { providerCapabilities } = useSelector((state) => state.ui);
  const { organization: currentOrganization } = useSelector((state) => state.ui);
  const providerUrl = providerCapabilities?.providerUrl;
  const [activeUsers] = useRoomActivity({
    providerUrl,
    getUserAccessToken: getUserAccessToken,
    getUserProfile: getUserProfile,
  });
  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();
  const [assignViewToWorkspace] = useAssignViewToWorkspaceMutation();
  const { notify } = useNotification();
  const router = useRouter();
  return (
    <>
      <DesignList data-testid="designs-list-item">
        {shouldRenderDesigns &&
          designs?.map((design) => {
            const isPublished = design?.visibility === 'published';
            const isOwner = currentUser?.id === design?.userId;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={`${design?.id}-${design?.name}`}>
                <DesignViewListItem
                  showWorkspaceName={showWorkspaceName}
                  showOrganizationName={showOrganizationName}
                  activeUsers={activeUsers?.[design?.id]}
                  type={RESOURCE_TYPE.DESIGN}
                  selectedItem={design}
                  handleItemClick={() => {
                    handleOpenDesignInExtension(design?.id, design?.name);
                  }}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={async (value, selectedItem) => {
                    await handleUpdatePatternVisibility({
                      value,
                      selectedResource: selectedItem,
                      updatePatterns,
                    });
                    refetch();
                  }}
                  MenuComponent={
                    <MenuComponent
                      options={getMenuOptions({
                        design,
                        handleRemove,
                        handleDelete,
                        handleDesignDownloadModal,
                        handleShare,
                        handleInfoModal,
                        refetch,
                      })}
                    />
                  }
                  isMultiSelectMode={isMultiSelectMode}
                />
                <Divider light />
              </React.Fragment>
            );
          })}

        <LoadingContainer ref={loadingRef}>
          {isLoading || isInitialFetch ? (
            Array(10)
              .fill(null)
              .map((_, index) => (
                <DesignViewListItemSkeleton key={index} isMultiSelectMode={isMultiSelectMode} />
              ))
          ) : isFetching ? (
            <DesignViewListItemSkeleton isMultiSelectMode={isMultiSelectMode} />
          ) : null}

          {!hasMore && !isLoading && !isFetching && designs?.length > 0 && !isEmpty && (
            <ListItemText secondary="No more designs to load" sx={{ padding: '1rem' }} />
          )}
        </LoadingContainer>

        {!isLoading && isEmpty && (
          <ListItem>
            <ListItemText primary="No Designs found" style={{ textAlign: 'center' }} />
          </ListItem>
        )}
      </DesignList>
      <GhostContainer ref={ghostRef}>
        <GhostImage src="/static/img/designs/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>
      <ExportDesignModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDesignDownload}
      />
      {shareModal && (
        <ShareModal
          resource={selectedDesign}
          handleClose={handleShareClose}
          type={RESOURCE_TYPE.DESIGN}
        />
      )}
      {infoModal.open && (
        <Modal {...sistentInfoModal}>
          <InfoModal
            infoModalOpen={infoModal.open}
            handleInfoModalClose={handleInfoModalClose}
            selectedResource={selectedDesign}
            resourceOwnerID={infoModal.userId}
            currentUser={currentUser}
            meshModels={meshModelsData?.models}
            patternFetcher={refetch}
          />
        </Modal>
      )}
      <PromptComponent ref={modalRef} />
      {moveModal && (
        <WorkspaceContentMoveModal
          currentWorkspace={workspace}
          setWorkspaceContentMoveModal={setMoveModal}
          type={RESOURCE_TYPE.DESIGN}
          workspaceContentMoveModal={moveModal}
          selectedContent={selectedDesign}
          refetch={refetch}
          useGetWorkspacesQuery={useGetWorkspacesQuery}
          WorkspaceModalContext={WorkspaceModalContext}
          assignDesignToWorkspace={assignDesignToWorkspace}
          assignViewToWorkspace={assignViewToWorkspace}
          isCreateWorkspaceAllowed={CAN(
            Keys.WorkspaceManagementCreateWorkspace.id,
            Keys.WorkspaceManagementCreateWorkspace.function,
          )}
          isMoveDesignAllowed={CAN(
            Keys.WorkspaceManagementAssignDesignsToWorkspaces.id,
            Keys.WorkspaceManagementAssignDesignsToWorkspaces.function,
          )}
          isMoveViewAllowed={CAN(
            Keys.KanvasAssignViewsToWorkspace.id,
            Keys.KanvasAssignViewsToWorkspace.function,
          )}
          currentOrgId={currentOrganization?.id}
          notify={notify}
          router={router}
        />
      )}
    </>
  );
};

export default MainDesignsContent;
