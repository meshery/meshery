import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ListItem,
  ListItemText,
  Divider,
  ShareIcon,
  InfoIcon,
  DeleteIcon,
  PROMPT_VARIANTS,
  PromptComponent,
  useTheme,
} from '@layer5/sistent';
import React, { useCallback, useContext, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, { handleUpdateViewVisibility } from './hooks';
import { MenuComponent } from './MenuComponent';
import { RESOURCE_TYPE } from '@/utils/Enum';
import GetAppIcon from '@mui/icons-material/GetApp';
import { DesignList, LoadingContainer, GhostContainer, GhostImage, GhostText } from './styles';
import { downloadFileFromContent } from '@/utils/fileDownloader';
import { getView, useDeleteViewMutation, useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import ShareModal from './ShareModal';
import { ViewsInfoModal } from '../ViewInfoModal';
import { openViewInKanvas, useIsOperatorEnabled } from '@/utils/utils';
import { WorkspaceSwitcherContext } from './WorkspaceSwitcher';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { Router } from 'next/router';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useUnassignViewFromWorkspaceMutation } from '@/rtk-query/workspace';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';

const MainViewsContent = ({
  page,
  setPage,
  isLoading,
  isFetching,
  views,
  hasMore,
  total_count,
  workspaceId,
  refetch,
}) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setinfoModal] = useState(null);
  const [unassignViewFromWorkspace] = useUnassignViewFromWorkspaceMutation();

  const [selectedView, setSetselectedView] = useState(null);
  const [updateView] = useUpdateViewVisibilityMutation();
  const handleOpenShareModal = (view) => {
    setSetselectedView(view);
    setShareModal(true);
  };

  const handleCloseShareModal = () => {
    setSetselectedView(null);
    setShareModal(false);
  };

  const handleOpenInfoModal = (view) => {
    setSetselectedView(view);
    setinfoModal(true);
  };

  const handleCloseInfoModal = () => {
    setSetselectedView(null);
    setinfoModal(false);
  };

  const modalRef = useRef(true);

  const handleOpenDeleteModal = async (view) => {
    const response = await modalRef.current.show({
      title: `Delete catalog item?`,
      subtitle: `Are you sure you want to delete ${view?.name}?`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      const selectedView = view;
      const { name, id } = selectedView;
      deleteView({
        id: id,
      })
        .unwrap()
        .then(() => {
          notify({ message: `"${name}" View deleted`, event_type: EVENT_TYPES.SUCCESS });
        })
        .catch(() => {
          notify({ message: `Unable to delete "${name}" View`, event_type: EVENT_TYPES.ERROR });
        });
    }
  };
  const handleRemove = async (view, workspaceId) => {
    unassignViewFromWorkspace({
      workspaceId: workspaceId,
      viewId: view.id,
    })
      .unwrap()
      .then(() => {
        setPage(0);
        notify({
          message: 'View removed from workspace',
          event_type: EVENT_TYPES.SUCCESS,
        });
      });
  };

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

  const [deleteView] = useDeleteViewMutation();
  const theme = useTheme();

  const VIEW_ACTIONS = {
    EXPORT_VIEW: {
      id: 'EXPORT_VIEW',
      title: 'Export View',
      icon: <GetAppIcon style={{ fill: theme.palette.icon.default }} />,
      handler: async ({ view }) => {
        const res = await getView({ viewId: view.id });
        downloadFileFromContent(JSON.stringify(res.data), `${view.name}.json`, 'application/json');
      },
      enabled: () => true,
    },
    DELETE_VIEW: {
      id: workspaceId ? 'MOVE_VIEW' : 'DELETE_VIEW',
      title: workspaceId ? 'Move View' : 'Delete View',
      icon: workspaceId ? (
        <MoveFileIcon fill={theme.palette.icon.default} />
      ) : (
        <DeleteIcon fill={theme.palette.icon.default} />
      ),
      enabled: ({ view, userId }) =>
        workspaceId
          ? CAN(keys.REMOVE_VIEWS_FROM_WORKSPACE.action, keys.REMOVE_VIEWS_FROM_WORKSPACE.subject)
          : CAN(keys.DELETE_VIEW.action, keys.DELETE_VIEW.subject) && view.user_id === userId,
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
  };
  const getMenuOptions = ({
    view,
    user,
    handleRemove,
    handleOpenInfoModal,
    handleOpenShareModal,
    handleOpenDeleteModal,
  }) => {
    const options = [
      {
        ...VIEW_ACTIONS.EXPORT_VIEW,
        handler: () => VIEW_ACTIONS.EXPORT_VIEW.handler({ view }),
      },
      {
        ...VIEW_ACTIONS.DELETE_VIEW,
        handler: () =>
          workspaceId ? handleRemove(view, workspaceId) : handleOpenDeleteModal(view),
      },

      {
        ...VIEW_ACTIONS.SHARE_VIEW,
        handler: () => handleOpenShareModal(view, RESOURCE_TYPE.VIEW),
      },
      {
        ...VIEW_ACTIONS.VIEW_INFO,
        handler: () => handleOpenInfoModal(view, user),
      },
    ];

    return options.filter((option) => option.enabled({ view, userId: user?.id }));
  };
  const isKanvasDesignerAvailable = useIsOperatorEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
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
                        handleRemove,
                        handleOpenInfoModal,
                        handleOpenShareModal,
                        handleOpenDeleteModal,
                      })}
                    />
                  }
                />
                <Divider light />
              </React.Fragment>
            );
          })}
        <LoadingContainer ref={loadingRef}>
          {isLoading || isInitialFetch ? (
            Array(10)
              .fill()
              .map((_, index) => <DesignViewListItemSkeleton key={index} />)
          ) : isFetching ? (
            <DesignViewListItemSkeleton />
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
        <ViewsInfoModal
          open={infoModal}
          closeModal={handleCloseInfoModal}
          view_id={selectedView?.id}
          view_name={selectedView?.name}
          metadata={selectedView?.metadata}
        />
      )}
      <PromptComponent ref={modalRef} />
    </>
  );
};

export default MainViewsContent;
