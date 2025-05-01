//@ts-check
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  styled,
  List,
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
  ShareIcon,
  InfoIcon,
  DeleteIcon,
} from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import DesignViewListItem from './DesignViewListItem';
import useInfiniteScroll from './hooks';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { GeorgeMenu } from './MenuComponent';
import { VISIBILITY } from '@/utils/Enum';
import GetAppIcon from '@mui/icons-material/GetApp';
import { MoreVert } from '@mui/icons-material';
import { DesignList, LoadingContainer, GhostContainer, GhostImage, GhostText } from './styles';

const MainViewsContent = ({ viewsData, setPage, isLoading, isFetching, views }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});

  const hasMore = viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1);

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

  const total_count = viewsData?.total_count || 0;

  const viewIsPublic = (v) => v?.visibility === VISIBILITY.PUBLIC;
  const viewIsPrivate = (v) => v?.visibility === VISIBILITY.PRIVATE;
  const viewIsOwnedByUser = (v, userId) => v?.user_id === userId;

  const VIEW_ACTIONS = {
    LOAD_VIEW: {
      id: 'LOAD_VIEW',
      title: 'Load View',
      icon: GetAppIcon,
      handler: ({ dispatchCmdToEditor, view }) => {},
      enabled: ({ view, userId }) => viewIsOwnedByUser(view, userId) || viewIsPublic(view),
    },

    EXPORT_VIEW: {
      id: 'EXPORT_VIEW',
      title: 'Export View',
      icon: GetAppIcon,
      handler: async ({ view }) => {},
      enabled: () => true,
    },

    DELETE_VIEW: {
      id: 'DELETE_VIEW',
      title: 'Delete View',
      icon: DeleteIcon,
      handler: async ({ dispatchCmdToEditor, view }) => {},
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
      icon: () => null,
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
      enabled: ({ view, userId }) => {
        return viewIsOwnedByUser(view, userId);
      },
    },
  };
  const getGeorgeOptions = () => {
    const options = [
      {
        ...VIEW_ACTIONS.DELETE_VIEW,
        handler: () => {},
      },
      {
        ...VIEW_ACTIONS.EXPORT_VIEW,
        handler: () => {},
      },
      {
        ...VIEW_ACTIONS.SHARE_VIEW,
        enabled: () => true,
        handler: () => {},
      },
      {
        id: 'VIEW_INFO',
        title: 'View Info',
        icon: InfoIcon,
        enabled: () => true,
        handler: () => {},
      },
    ];

    return options;
  };
  return (
    <>
      <DesignList data-testid="designs-list-item">
        {(!isFetching || !isLoading) && total_count === 0 && (
          <ListItem>
            <ListItemText primary={`No Designs found`} style={{ textAlign: 'center' }} />
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
                  onVisibilityChange={() => {}}
                  MenuComponent={<GeorgeMenu options={getGeorgeOptions()} triggerIcon={MoreVert} />}
                />
                <Divider light />
              </React.Fragment>
            );
          })}
        <LoadingContainer ref={loadingRef}>
          {(isFetching || isLoading) && <CircularProgress size={24} />}
          {!hasMore && views?.length > 0 && total_count > 0 && (
            <ListItemText secondary={`No more designs to load`} />
          )}
        </LoadingContainer>
      </DesignList>
      <GhostContainer ref={ghostRef}>
        <GhostImage src="/static/img/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>
    </>
  );
};

export default MainViewsContent;
