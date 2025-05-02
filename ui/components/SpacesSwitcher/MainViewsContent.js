//@ts-check
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
  ShareIcon,
  InfoIcon,
  DeleteIcon,
} from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, { handleUpdateViewVisibility } from './hooks';
import { GeorgeMenu } from './MenuComponent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import GetAppIcon from '@mui/icons-material/GetApp';
import { MoreVert } from '@mui/icons-material';
import { DesignList, LoadingContainer, GhostContainer, GhostImage, GhostText } from './styles';
import { downloadFileFromContent } from '@/utils/fileDownloader';
import { api } from '@/rtk-query/index';
import { getView, useDeleteViewMutation, useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import ShareModal from './ShareModal';
import { ViewsInfoModal } from '../ViewInfoModal';

const MainViewsContent = ({ setPage, isLoading, isFetching, views, hasMore, total_count }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setinfoModal] = useState(null);
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

  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage((prevPage) => prevPage + 1);
  }, [isLoading, isFetching]);

  const { loadingRef } = useInfiniteScroll({
    isLoading: isLoading || isFetching,
    hasMore,
    onLoadMore: loadNextPage,
  });

  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);

  const viewIsPublic = (v) => v?.visibility === VISIBILITY.PUBLIC;
  const viewIsPrivate = (v) => v?.visibility === VISIBILITY.PRIVATE;
  const viewIsOwnedByUser = (v, userId) => v?.user_id === userId;
  const [deleteView] = useDeleteViewMutation();

  const VIEW_ACTIONS = {
    EXPORT_VIEW: {
      id: 'EXPORT_VIEW',
      title: 'Export View',
      icon: GetAppIcon,
      handler: async ({ view }) => {
        const res = await getView({ viewId: view.id });
        console.log('amit res', res);
        downloadFileFromContent(JSON.stringify(res.data), `${view.name}.json`, 'application/json');
      },
      enabled: () => true,
    },
    DELETE_VIEW: {
      id: 'DELETE_VIEW',
      title: 'Delete View',
      icon: DeleteIcon,
      handler: async ({ view, setPage }) => {
        deleteView({ id: view.id })
          .unwrap()
          .then(() => {
            setPage(0);
          });
      },
      enabled: ({ view, userId }) => viewIsOwnedByUser(view, userId),
    },

    UPDATE_VISIBILITY: {
      id: 'UPDATE_VISIBILITY',
      title: 'Update Visibility',
      icon: () => null,
      handler: ({ view, visibility }) => {},
      enabled: ({ view, userId }) => {
        return viewIsOwnedByUser(view, userId);
      },
    },
    UPDATE_INFO: {
      id: 'UPDATE_INFO',
      title: 'Update Info',
      icon: InfoIcon,
      handler: ({ view, metadata }) => {},
      enabled: ({ view, userId }) => {
        return viewIsOwnedByUser(view, userId);
      },
    },
    SHARE_VIEW: {
      id: 'SHARE_VIEW',
      title: 'Share View',
      icon: ShareIcon,
      handler: ({ view, emails }) => {},
      enabled: () => true,
    },
  };
  const getGeorgeOptions = ({ view, user, handleOpenInfoModal, handleOpenShareModal, setPage }) => {
    const options = [
      {
        ...VIEW_ACTIONS.EXPORT_VIEW,
        handler: () => VIEW_ACTIONS.EXPORT_VIEW.handler({ view }),
      },
      {
        ...VIEW_ACTIONS.DELETE_VIEW,
        handler: () => VIEW_ACTIONS.DELETE_VIEW.handler({ view, setPage }),
      },

      {
        ...VIEW_ACTIONS.SHARE_VIEW,
        handler: () => handleOpenShareModal(view, RESOURCE_TYPE.VIEW),
      },
      {
        id: 'VIEW_INFO',
        title: 'View Info',
        icon: InfoIcon,
        enabled: () => true,
        handler: () => handleOpenInfoModal(view, user),
      },
    ];

    return options.filter((option) => option.enabled({ view, userId: user?.id }));
  };

  return (
    <>
      <DesignList data-testid="designs-list-item">
        {(!isFetching || !isLoading) && total_count === 0 && (
          <ListItem>
            <ListItemText primary={`No Views found`} style={{ textAlign: 'center' }} />
          </ListItem>
        )}
        {total_count !== 0 &&
          views?.map((view) => {
            const isPublished = view?.visibility === 'published';
            const isOwner = currentUser?.id === view?.user_id;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={view?.id}>
                <DesignViewListItem
                  selectedItem={view}
                  handleItemClick={(e) => {}}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={(value, selectedItem) => {
                    handleUpdateViewVisibility({
                      value: value,
                      selectedResource: selectedItem,
                      updateView: updateView,
                    });
                    setPage(0);
                  }}
                  MenuComponent={
                    <GeorgeMenu
                      options={getGeorgeOptions({
                        view,
                        user: currentUser,
                        handleOpenInfoModal: handleOpenInfoModal,
                        handleOpenShareModal: handleOpenShareModal,
                        setPage,
                      })}
                    />
                  }
                />
                <Divider light />
              </React.Fragment>
            );
          })}
        <LoadingContainer ref={loadingRef}>
          {isLoading ? (
            Array(3)
              .fill()
              .map((_, index) => <DesignViewListItemSkeleton key={index} />)
          ) : isFetching ? (
            <DesignViewListItemSkeleton />
          ) : (
            <></>
          )}{' '}
          {!hasMore && views?.length > 0 && total_count > 0 && (
            <ListItemText secondary={`No more views to load`} />
          )}
        </LoadingContainer>
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
    </>
  );
};

export default MainViewsContent;
