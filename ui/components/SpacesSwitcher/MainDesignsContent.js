import { getDesign, useUpdatePatternFileMutation } from '@/rtk-query/design';
import { getUserAccessToken, getUserProfile, useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ListItem,
  ListItemText,
  Divider,
  PromptComponent,
  OutlinedPatternIcon,
  useModal,
  Modal,
  ShareIcon,
  useTheme,
  useRoomActivity,
  ExportIcon,
  DeleteIcon,
  InfoIcon,
  WorkspaceContentMoveModal,
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
import ExportModal from '../ExportModal';
import { RESOURCE_TYPE } from '@/utils/Enum';
import ShareModal from './ShareModal';
import InfoModal from '../General/Modals/Information/InfoModal';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import { openDesignInKanvas, useIsKanvasDesignerEnabled } from '@/utils/utils';
import Router, { useRouter } from 'next/router';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
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
  total_count,
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
      userId: selectedDesignWithPatternFile?.data?.user_id,
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
  const isKanvasDesignerAvailable = useIsKanvasDesignerEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const handleOpenDesignInDesigner = (designId, designName) => {
    if (workspaceSwitcherContext?.closeModal) {
      workspaceSwitcherContext.closeModal();
    }
    if (!isKanvasDesignerAvailable) {
      router.push(`/configuration/designs/configurator?design_id=${designId}`);
      return;
    }

    openDesignInKanvas(designId, designName, Router);
  };
  const theme = useTheme();
  const DESIGN_ACTIONS = {
    EXPORT_DESIGN: {
      id: 'export_design',
      title: 'Export Design',
      icon: <ExportIcon fill={theme.palette.icon.default} />,
      enabled: () => CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
    },

    REMOVE_DESIGN: {
      id: 'move',
      title: 'Move Design',
      icon: <MoveFileIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        CAN(keys.REMOVE_DESIGNS_FROM_WORKSPACE.action, keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject),
    },
    SHARE_DESIGN: {
      id: 'share',
      title: 'Share Design',
      icon: <ShareIcon fill={theme.palette.icon.default} />,
      enabled: ({ design }) =>
        design?.visibility !== 'published' &&
        CAN(keys.SHARE_DESIGN.action, keys.SHARE_DESIGN.subject),
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
      enabled: () => CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject),
    },
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
  const isEmpty = total_count === 0;
  const shouldRenderDesigns = !isEmpty && !isInitialFetch;
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { organization: currentOrganization } = useSelector((state) => state.ui);
  const providerUrl = capabilitiesRegistry?.provider_url;
  const [activeUsers] = useRoomActivity({
    provider_url: providerUrl,
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
            const isOwner = currentUser?.id === design?.user_id;
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
                    handleOpenDesignInDesigner(design?.id, design?.name);
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
        <GhostImage src="/static/img/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>
      <ExportModal
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
            keys.CREATE_WORKSPACE.action,
            keys.CREATE_WORKSPACE.subject,
          )}
          isMoveDesignAllowed={CAN(
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
          )}
          isMoveViewAllowed={CAN(
            keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
            keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
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
