import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import React, { useCallback, useRef, useState } from 'react';
import MainDesignsContent from './MainDesignsContent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import { Box, Grid, PromptComponent, useTheme } from '@layer5/sistent';
import { StyledSearchBar } from '@layer5/sistent';
import {
  ImportButton,
  MultiContentSelectToolbar,
  SortBySelect,
  TableListHeader,
  VisibilitySelect,
} from './components';
import { useContentDelete, useContentDownload } from './hooks';
import ExportModal from '../ExportModal';
import ShareModal from './ShareModal';

const MyDesignsContent = () => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE, VISIBILITY.PUBLISHED];
  // const { organization: currentOrganization } = useSelector((state) => state.ui);
  const [shareModal, setShareModal] = useState({ open: false, content: null });
  const [filters, setFilters] = useState({
    visibility: visibilityItems,
    searchQuery: '',
    sortBy: 'updated_at desc',
    page: 0,
  });

  const handleSortByChange = useCallback((event) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      sortBy: event.target.value,
      page: 0,
    }));
  }, []);

  const handleVisibilityChange = useCallback((event) => {
    const value = event.target.value;
    setFilters((prevFilters) => ({
      ...prevFilters,
      visibility: typeof value === 'string' ? value.split(',') : value,
      page: 0,
    }));
  }, []);

  const onSearchChange = useCallback((e) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      searchQuery: e.target.value,
      page: 0,
    }));
  }, []);

  const setPage = useCallback((newPage) => {
    setFilters((prevFilters) => ({
      ...prevFilters,
      page: newPage,
    }));
  }, []);

  const {
    data: designsData,
    isLoading,
    isFetching,
    refetch: refetchDesigns,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: filters.page,
      pagesize: 10,
      order: filters.sortBy,
      user_id: currentUser?.id,
      metrics: true,
      visibility: filters.visibility,
      search: filters.searchQuery,
      // orgId: currentOrganization?.id,
    },
    {
      skip: !currentUser?.id,
    },
  );
  const hasMore = designsData?.total_count > designsData?.page_size * (designsData?.page + 1);
  const total_count = designsData?.total_count || 0;
  const theme = useTheme();
  const modalRef = useRef(null);
  const { handleDelete } = useContentDelete(modalRef);
  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });
  const handleDownloadModalOpen = (content) => {
    setDownloadModal({
      open: true,
      content: content,
    });
  };
  const handleDownloadModalClose = () => {
    setDownloadModal({
      open: false,
      content: null,
    });
  };
  const { handleDesignDownload } = useContentDownload();
  const refetch = useCallback(() => {
    if (filters.page > 0) setPage(0);
    else refetchDesigns();
  }, [filters.page, refetchDesigns, setPage]);
  return (
    <Box display={'flex'} flexDirection="column" gap="1rem">
      <Grid container spacing={2} alignItems="center" marginBottom="1rem">
        {/* Search Bar */}
        <Grid item xs={12} md={7}>
          <StyledSearchBar
            sx={{ backgroundColor: 'transparent' }}
            placeholder="Search Designs"
            value={filters.searchQuery}
            onChange={onSearchChange}
            endAdornment={
              <p style={{ color: theme.palette.text.default }}>Total Designs: {total_count}</p>
            }
          />
        </Grid>

        {/* Sort By Select */}
        <Grid item xs={4} md={2}>
          <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
        </Grid>

        {/* Visibility Select */}
        <Grid item xs={4} md={2}>
          <VisibilitySelect
            visibility={filters.visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Grid>

        {/* Import Button */}
        <Grid item xs={4} md={1}>
          <ImportButton refetch={refetch} />
        </Grid>
      </Grid>
      <MultiContentSelectToolbar
        type={RESOURCE_TYPE.DESIGN}
        handleDelete={handleDelete}
        handleDownload={handleDownloadModalOpen}
        refetch={refetch}
        handleShare={(multiSelectedContent) => {
          setShareModal({
            open: true,
            content: multiSelectedContent,
          });
        }}
      />

      <TableListHeader isMultiSelectMode={true} content={designsData?.patterns} />
      <MainDesignsContent
        key={'my-designs-content'}
        page={filters.page}
        designs={designsData?.patterns}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        hasMore={hasMore}
        refetch={refetch}
        isMultiSelectMode={true}
        total_count={total_count}
      />
      <PromptComponent ref={modalRef} />
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadModalClose}
        handleDesignDownload={handleDesignDownload}
      />
      {shareModal.open && (
        <ShareModal
          resource={shareModal.content}
          handleClose={() => setShareModal(false)}
          type={RESOURCE_TYPE.DESIGN}
        />
      )}
    </Box>
  );
};

export default MyDesignsContent;
