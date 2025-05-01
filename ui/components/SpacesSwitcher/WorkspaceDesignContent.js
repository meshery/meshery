//@ts-check
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { styled, List, ListItem, ListItemText, Divider, CircularProgress } from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import DesignViewListItem from './DesignViewListItem';
import useInfiniteScroll from './hooks';
import MenuComponent from './MenuComponent';
import { MoreVert } from '@mui/icons-material';
import { DesignList, GhostContainer, GhostImage, GhostText, LoadingContainer } from './styles';
import { useGetDesignsOfWorkspaceQuery } from '@/rtk-query/workspace';

const WorkspaceDesignContent = ({ workspaceId }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [page, setPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetDesignsOfWorkspaceQuery(
    {
      workspaceId: workspaceId,
      page: page,
      pagesize: 10,
      order: 'updated_at desc',
    },
    {
      skip: !workspaceId,
    },
  );
  const designs = designsData?.designs || [];
  const hasMore = designsData?.total_count > designsData?.page_size * (designsData?.page + 1);
  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage((prevPage) => prevPage + 1);
  }, [isLoading, isFetching]);

  const { loadingRef } = useInfiniteScroll({
    isLoading: isLoading || isFetching,
    hasMore,
    onLoadMore: loadNextPage,
  });
  const total_count = designsData?.total_count || 0;

  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);

  return (
    <>
      <DesignList data-testid="designs-list-item">
        {(!isFetching || !isLoading) && total_count === 0 && (
          <ListItem>
            <ListItemText primary={`No Designs found`} style={{ textAlign: 'center' }} />
          </ListItem>
        )}

        {total_count !== 0 &&
          designs?.map((design) => {
            const isPublished = design?.visibility === 'published';
            const isOwner = currentUser?.id === design?.user_id;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={design?.id}>
                <DesignViewListItem
                  selectedItem={design}
                  handleItemClick={(e) => {}}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={() => {}}
                  MenuComponent={
                    <MenuComponent
                      // dataName={PATTERN_PLURAL}
                      iconType={MoreVert}
                      rowData={design}
                      // designOwnerId={design?.user_id}
                      visibility={design?.visibility}
                      items={[
                        {
                          deleteHandler: () => {},
                          downloadDesignHandler: () => {},
                          cloneHandler: () => {},
                          shareHandler: () => {},
                          infoHandler: () => {},
                          publishHandler: () => {},
                          unPublishHandler: () => {},
                        },
                      ]}
                    />
                  }
                />
                <Divider light />
              </React.Fragment>
            );
          })}

        <LoadingContainer ref={loadingRef}>
          {(isFetching || isLoading) && <CircularProgress size={24} />}
          {!hasMore && designs?.length > 0 && total_count > 0 && (
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

export default WorkspaceDesignContent;
