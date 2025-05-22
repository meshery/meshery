import {
  getDesign,
  useDeletePatternFileMutation,
  useUpdatePatternFileMutation,
} from '@/rtk-query/design';
import { getUserAccessToken, getUserProfile, useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ListItem,
  ListItemText,
  Divider,
  PROMPT_VARIANTS,
  PromptComponent,
  OutlinedPatternIcon,
  useModal,
  Modal,
  ShareIcon,
  useTheme,
  useRoomActivity,
} from '@layer5/sistent';
import React, { useCallback, useContext, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, { handleUpdatePatternVisibility } from './hooks';
import { MenuComponent } from './MenuComponent';
import { DesignList, GhostContainer, GhostImage, GhostText, LoadingContainer } from './styles';
import ExportModal from '../ExportModal';
import downloadContent from '@/utils/fileDownloader';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { RESOURCE_TYPE } from '@/utils/Enum';
import ShareModal from './ShareModal';
import InfoModal from '../Modals/Information/InfoModal';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import { openDesignInKanvas, useIsKanvasDesignerEnabled } from '@/utils/utils';
import Router from 'next/router';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DeleteIcon from '@mui/icons-material/Delete';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoIcon from '@mui/icons-material/Info';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import { updateProgress } from '@/store/slices/mesheryUi';
import { useSelector } from 'react-redux';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import WorkspaceContentMoveModal from './WorkspaceContentMoveModal';

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
}) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [selectedDesign, setSelectedDesign] = useState(null);
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setInfoModal] = useState({ open: false, userId: '' });
  const [moveModal, setMoveModal] = useState(false);
  const { notify } = useNotification();
  const modalRef = useRef(true);
  const [deletePatternFile] = useDeletePatternFileMutation();
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
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: design,
    }));
  };
  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };

  const handleDelete = async (design) => {
    const response = await modalRef.current.show({
      title: `Delete catalog item?`,
      subtitle: `Are you sure you want to delete ${design?.name}?`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
      showInfoIcon:
        "Unpublishing a catolog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want. Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the [documentation](https://docs.meshery.io/concepts/designs).",
    });
    if (response === 'DELETE') {
      const selectedPattern = design;
      const { name, id } = selectedPattern;
      deletePatternFile({
        id: id,
      })
        .unwrap()
        .then(() => {
          notify({ message: `"${name}" Design deleted`, event_type: EVENT_TYPES.SUCCESS });
        })
        .catch(() => {
          notify({ message: `Unable to delete "${name}" Design`, event_type: EVENT_TYPES.ERROR });
        });
    }
  };

  const handleRemove = async (design) => {
    setMoveModal(true);
    setSelectedDesign(design);
  };

  const handleShare = async (design) => {
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
    if (!isKanvasDesignerAvailable) {
      notify({
        message: 'Kanvas Designer is not available',
        event_type: EVENT_TYPES.ERROR,
      });
      return;
    }
    if (workspaceSwitcherContext?.closeModal) {
      workspaceSwitcherContext.closeModal();
    }

    openDesignInKanvas(designId, designName, Router);
  };
  const theme = useTheme();
  const DESIGN_ACTIONS = {
    EXPORT_DESIGN: {
      id: 'export_design',
      title: 'Export Design',
      icon: <GetAppIcon style={{ fill: theme.palette.icon.default }} />,
      enabled: () => CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
    },
    DELETE_DESIGN: {
      id: 'delete',
      title: 'Delete Design',
      icon: <DeleteIcon fill={theme.palette.icon.default} />,
      enabled: () => CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject),
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
  };

  const getMenuOptions = ({
    design,
    handleDesignDownloadModal,
    handleDelete,
    handleRemove,
    handleShare,
    handleInfoModal,
  }) => {
    const options = [
      {
        ...DESIGN_ACTIONS.EXPORT_DESIGN,
        handler: () => handleDesignDownloadModal(design),
      },
      {
        ...DESIGN_ACTIONS.DELETE_DESIGN,
        handler: () => handleDelete(design),
      },
      {
        ...DESIGN_ACTIONS.SHARE_DESIGN,
        handler: () => handleShare(design),
      },

      {
        ...DESIGN_ACTIONS.INFO_DESIGN,
        handler: () => handleInfoModal(design),
      },
    ];

    if (workspace) {
      options.unshift({
        ...DESIGN_ACTIONS.REMOVE_DESIGN,
        handler: () => handleRemove(design),
      });
    }
    return options.filter((option) => option.enabled({ design }));
  };
  const isInitialFetch = isFetching && page === 0;
  const isEmpty = total_count === 0;
  const shouldRenderDesigns = !isEmpty && !isInitialFetch;
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const providerUrl = capabilitiesRegistry?.provider_url;
  const activeUsers = useRoomActivity({
    provider_url: providerUrl,
    getUserAccessToken: getUserAccessToken,
    getUserProfile: getUserProfile,
  });
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
                  activeUsers={activeUsers?.[design?.id]}
                  type="design"
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
              .map((_, index) => <DesignViewListItemSkeleton key={index} />)
          ) : isFetching ? (
            <DesignViewListItemSkeleton />
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
        handleDesignDownload={handleDownload}
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
        />
      )}
    </>
  );
};

export default MainDesignsContent;
