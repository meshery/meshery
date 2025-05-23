import { getUserAccessToken, getUserProfile, useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ListItem,
  ListItemText,
  Divider,
  ShareIcon,
  InfoIcon,
  DeleteIcon,
  PromptComponent,
  useTheme,
  useRoomActivity,
} from '@layer5/sistent';
import React, { useCallback, useContext, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, {
  handleUpdateViewVisibility,
  useContentDelete,
  useContentDownload,
} from './hooks';
import { MenuComponent } from './MenuComponent';
import { RESOURCE_TYPE } from '@/utils/Enum';
import GetAppIcon from '@mui/icons-material/GetApp';
import { DesignList, LoadingContainer, GhostContainer, GhostImage, GhostText } from './styles';
import { useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import ShareModal from './ShareModal';
import { ViewInfoModal } from '../ViewInfoModal';
import { openViewInKanvas, useIsOperatorEnabled } from '@/utils/utils';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { Router } from 'next/router';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import { useSelector } from 'react-redux';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import WorkspaceContentMoveModal from './WorkspaceContentMoveModal';

const MainViewsContent = ({
  page,
  setPage,
  isLoading,
  isFetching,
  views,
  hasMore,
  total_count,
  workspace,
  refetch,
  isMultiSelectMode,
}) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setinfoModal] = useState(null);
  const [moveModal, setMoveModal] = useState(false);

  const [selectedView, setSelectedView] = useState(null);
  const [updateView] = useUpdateViewVisibilityMutation();
  const handleOpenShareModal = (view) => {
    setSelectedView(view);
    setShareModal(true);
  };

  const handleCloseShareModal = () => {
    setSelectedView(null);
    setShareModal(false);
  };

  const handleOpenInfoModal = (view) => {
    setSelectedView(view);
    setinfoModal(true);
  };

  const handleCloseInfoModal = () => {
    setSelectedView(null);
    setinfoModal(false);
  };
  const handleMoveModal = (view) => {
    setSelectedView(view);
    setMoveModal(true);
  };

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

  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);

  const theme = useTheme();

  const { handleViewDownload } = useContentDownload();
  const VIEW_ACTIONS = {
    EXPORT_VIEW: {
      id: 'EXPORT_VIEW',
      title: 'Export View',
      icon: <GetAppIcon style={{ fill: theme.palette.icon.default }} />,
      handler: ({ view }) => {
        handleViewDownload(view);
      },
      enabled: () => true,
    },

    MOVE_VIEW: {
      id: 'MOVE_VIEW',
      title: 'Move View',
      icon: <MoveFileIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        CAN(keys.REMOVE_VIEWS_FROM_WORKSPACE.action, keys.REMOVE_VIEWS_FROM_WORKSPACE.subject),
    },

    VIEW_INFO: {
      id: 'VIEW_INFO',
      title: 'View Info',
      icon: <InfoIcon fill={theme.palette.icon.default} />,
      enabled: () => true,
    },
    SHARE_VIEW: {
      id: 'SHARE_VIEW',
      title: 'Share View',
      icon: <ShareIcon fill={theme.palette.icon.default} />,
      enabled: () => true,
    },
    DELETE_VIEW: {
      id: 'DELETE_VIEW',
      title: 'Delete View',
      icon: <DeleteIcon fill={theme.palette.icon.default} />,
      enabled: ({ view, userId }) =>
        CAN(keys.DELETE_VIEW.action, keys.DELETE_VIEW.subject) && view.user_id === userId,
    },
  };
  const getMenuOptions = ({
    view,
    user,
    handleMoveModal,
    handleOpenInfoModal,
    handleOpenShareModal,
    handleDelete,
    refetch,
  }) => {
    const options = [
      {
        ...VIEW_ACTIONS.EXPORT_VIEW,
        handler: () => VIEW_ACTIONS.EXPORT_VIEW.handler({ view }),
      },

      {
        ...VIEW_ACTIONS.SHARE_VIEW,
        handler: () => handleOpenShareModal(view, RESOURCE_TYPE.VIEW),
      },
      {
        ...VIEW_ACTIONS.VIEW_INFO,
        handler: () => handleOpenInfoModal(view, user),
      },
      {
        ...VIEW_ACTIONS.DELETE_VIEW,
        handler: () => handleDelete([view], RESOURCE_TYPE.VIEW, refetch),
      },
    ];
    if (workspace) {
      options.unshift({
        ...VIEW_ACTIONS.MOVE_VIEW,
        handler: () => handleMoveModal(view),
      });
    }
    return options.filter((option) => option.enabled({ view, userId: user?.id }));
  };
  const isKanvasDesignerAvailable = useIsOperatorEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const { notify } = useNotification();
  const handleOpenViewInOperator = (viewId, viewName) => {
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

    openViewInKanvas(viewId, viewName, Router);
  };
  const isInitialFetch = isFetching && page === 0;
  const isEmpty = total_count === 0;
  const shouldRenderDesigns = !isEmpty && !isInitialFetch;
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const providerUrl = capabilitiesRegistry?.provider_url;
  const [activeUsers] = useRoomActivity({
    provider_url: providerUrl,
    getUserAccessToken: getUserAccessToken,
    getUserProfile: getUserProfile,
  });
  return (
    <>
      <DesignList data-testid="designs-list-item">
        {shouldRenderDesigns &&
          views?.map((view) => {
            const isPublished = view?.visibility === 'published';
            const isOwner = currentUser?.id === view?.user_id;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={view?.id}>
                <DesignViewListItem
                  type="view"
                  selectedItem={view}
                  handleItemClick={() => {
                    handleOpenViewInOperator(view?.id, view?.name);
                  }}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={async (value, selectedItem) => {
                    await handleUpdateViewVisibility({
                      value: value,
                      selectedResource: selectedItem,
                      updateView: updateView,
                    });
                    refetch();
                  }}
                  MenuComponent={
                    <MenuComponent
                      options={getMenuOptions({
                        view,
                        user: currentUser,
                        handleMoveModal,
                        handleOpenInfoModal,
                        handleOpenShareModal,
                        handleDelete,
                        refetch,
                      })}
                    />
                  }
                  activeUsers={activeUsers?.[view?.id]}
                  isMultiSelectMode={isMultiSelectMode}
                />
                <Divider light />
              </React.Fragment>
            );
          })}
        <LoadingContainer ref={loadingRef}>
          {isLoading || isInitialFetch ? (
            Array(10)
              .fill()
              .map((_, index) => (
                <DesignViewListItemSkeleton key={index} isMultiSelectMode={isMultiSelectMode} />
              ))
          ) : isFetching ? (
            <DesignViewListItemSkeleton isMultiSelectMode={isMultiSelectMode} />
          ) : (
            <></>
          )}
          {!hasMore && !isLoading && !isFetching && views?.length > 0 && !isEmpty && (
            <ListItemText secondary={`No more views to load`} style={{ padding: '1rem' }} />
          )}
        </LoadingContainer>

        {!isLoading && isEmpty && (
          <ListItem>
            <ListItemText primary={`No Views found`} style={{ textAlign: 'center' }} />
          </ListItem>
        )}
      </DesignList>
      <GhostContainer ref={ghostRef}>
        <GhostImage src="/static/img/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>
      {shareModal && (
        <ShareModal
          resource={selectedView}
          handleClose={handleCloseShareModal}
          type={RESOURCE_TYPE.VIEW}
        />
      )}
      {infoModal && (
        <ViewInfoModal
          open={infoModal}
          closeModal={handleCloseInfoModal}
          view_id={selectedView?.id}
          view_name={selectedView?.name}
          metadata={selectedView?.metadata}
        />
      )}
      {moveModal && (
        <WorkspaceContentMoveModal
          currentWorkspace={workspace}
          selectedContent={selectedView}
          setWorkspaceContentMoveModal={setMoveModal}
          type={RESOURCE_TYPE.VIEW}
          workspaceContentMoveModal={moveModal}
          refetch={refetch}
        />
      )}
      <PromptComponent ref={modalRef} />
    </>
  );
};

export default MainViewsContent;
