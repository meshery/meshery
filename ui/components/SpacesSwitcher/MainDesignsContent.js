//@ts-check
import { useDeletePatternFileMutation, useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  styled,
  List,
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
  PROMPT_VARIANTS,
  PromptComponent,
} from '@layer5/sistent';
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
import { RESOURCE_TYPE } from '@/utils/Enum';
import ShareModal from './ShareModal';

const MainDesignsContent = ({ setPage, isLoading, isFetching, designs, hasMore, total_count }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [selectedShareDesign, setSelectedShareDesign] = useState(null);
  const [shareModal, setShareModal] = useState(false);

  const { notify } = useNotification();
  const modalRef = useRef(true);
  const [deletePatternFile] = useDeletePatternFileMutation();
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

  const handleDelete = async (design) => {
    setPage(0);
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

  const handleShare = async (design) => {
    setShareModal(true);
    setSelectedShareDesign(design);
  };

  const handleShareClose = () => {
    setShareModal(false);
    setSelectedShareDesign(null);
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
                          deleteHandler: () => handleDelete(design),
                          downloadHandler: () => handleDesignDownloadModal(design),
                          cloneHandler: () => {},
                          shareHandler: () => handleShare(design),
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
      {shareModal && (
        <ShareModal
          resource={selectedShareDesign}
          handleClose={handleShareClose}
          type={RESOURCE_TYPE.DESIGN}
        />
      )}
      <PromptComponent ref={modalRef} />
    </>
  );
};

export default MainDesignsContent;
