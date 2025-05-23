import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { Box, Grid, PromptComponent, useTheme } from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import MainViewsContent from './MainViewsContent';
import { StyledSearchBar } from '@layer5/sistent';
import {
  MultiContentSelectToolbar,
  SortBySelect,
  TableListHeader,
  VisibilitySelect,
} from './components';
import { useContentDelete, useContentDownload } from './hooks';
import ShareModal from './ShareModal';

const MyViewsContent = () => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];
  const [shareModal, setShareModal] = useState({ open: false, content: null });

  const [filters, setFilters] = useState({
    visibility: visibilityItems,
    searchQuery: '',
    sortBy: 'updated_at desc',
    page: 0,
  });

  const handleSortByChange = useCallback((event) => {
    setFilters((prev) => ({
      ...prev,
      sortBy: event.target.value,
      page: 0,
    }));
  }, []);

  const handleVisibilityChange = useCallback((event) => {
    const value = event.target.value;
    setFilters((prev) => ({
      ...prev,
      visibility: typeof value === 'string' ? value.split(',') : value,
      page: 0,
    }));
  }, []);

  const onSearchChange = useCallback((e) => {
    setFilters((prev) => ({
      ...prev,
      searchQuery: e.target.value,
      page: 0,
    }));
  }, []);

  const setPage = useCallback((newPage) => {
    setFilters((prev) => ({
      ...prev,
      page: newPage,
    }));
  }, []);

  const {
    data: viewsData,
    isLoading,
    isFetching,
  } = useFetchViewsQuery(
    {
      page: filters.page,
      pagesize: 10,
      order: filters.sortBy,
      visibility: filters.visibility,
      search: filters.searchQuery,
      user_id: currentUser?.id,
    },
    {
      skip: !currentUser?.id,
    },
  );

  const views = viewsData?.views || [];
  const hasMore = viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1);
  const total_count = viewsData?.total_count || 0;
  const theme = useTheme();
  const modalRef = useRef(null);
  const { handleDelete } = useContentDelete(modalRef);
  const { handleViewDownload } = useContentDownload();
  return (
    <Box display={'flex'} flexDirection="column" gap="1rem">
      <Grid container spacing={2} alignItems="center" marginBottom="1rem">
        {/* Search Bar */}
        <Grid item xs={12} md={6}>
          <StyledSearchBar
            sx={{ backgroundColor: 'transparent' }}
            placeholder="Search Views"
            value={filters.searchQuery}
            onChange={onSearchChange}
            endAdornment={
              <p style={{ color: theme.palette.text.default }}>Total Views: {total_count}</p>
            }
          />
        </Grid>

        {/* Sort By Select */}
        <Grid item xs={6} md={3}>
          <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
        </Grid>

        {/* Visibility Select */}
        <Grid item xs={6} md={3}>
          <VisibilitySelect
            visibility={filters.visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Grid>
      </Grid>
      <MultiContentSelectToolbar
        type={RESOURCE_TYPE.VIEW}
        handleDelete={handleDelete}
        handleViewDownload={handleViewDownload}
        refetch={() => setPage(0)}
        handleShare={(multiSelectedContent) => {
          setShareModal({
            open: true,
            content: multiSelectedContent,
          });
        }}
      />
      <TableListHeader isMultiSelectMode content={views} />
      <MainViewsContent
        key={'my-views'}
        page={filters.page}
        hasMore={hasMore}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        views={views}
        total_count={total_count}
        refetch={() => setPage(0)}
        isMultiSelectMode={true}
      />
      <PromptComponent ref={modalRef} />
      {shareModal.open && (
        <ShareModal
          resource={shareModal.content}
          handleClose={() => setShareModal(false)}
          type={RESOURCE_TYPE.VIEW}
        />
      )}
    </Box>
  );
};

export default MyViewsContent;
