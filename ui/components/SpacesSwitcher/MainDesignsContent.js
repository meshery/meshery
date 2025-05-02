//@ts-check
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { styled, List, ListItem, ListItemText, Divider, CircularProgress } from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll from './hooks';
import MenuComponent from './MenuComponent';
import { MoreVert } from '@mui/icons-material';
import { DesignList, GhostContainer, GhostImage, GhostText, LoadingContainer } from './styles';
import ExportModal from '../ExportModal';
import { updateProgress } from 'lib/store';
import downloadContent from '@/utils/fileDownloader';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

const MainDesignsContent = ({ setPage, isLoading, isFetching, designs, hasMore, total_count }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const { notify } = useNotification();

  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage((prevPage) => prevPage + 1);
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
                          downloadHandler: () => handleDesignDownloadModal(design),
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
          {isLoading ? (
            Array(3)
              .fill()
              .map((_, index) => <DesignViewListItemSkeleton key={index} />)
          ) : isFetching ? (
            <DesignViewListItemSkeleton />
          ) : (
            <></>
          )}
          {!hasMore && designs?.length > 0 && total_count > 0 && (
            <ListItemText secondary={`No more designs to load`} />
          )}
        </LoadingContainer>
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
    </>
  );
};

export default MainDesignsContent;
