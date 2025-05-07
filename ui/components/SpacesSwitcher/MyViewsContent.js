import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { Box, useTheme } from '@layer5/sistent';
import React, { useCallback, useState } from 'react';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { VISIBILITY } from '@/utils/Enum';
import MainViewsContent from './MainViewsContent';
import { StyledSearchBar } from '@layer5/sistent';
import { SortBySelect, TableListHeader, VisibilitySelect } from './components';

const MyViewsContent = () => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

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
  return (
    <Box display={'flex'} flexDirection="column" gap="1rem">
      <Box display={'flex'} alignItems="center" marginBottom="1rem" gap={'1rem'}>
        <StyledSearchBar
          sx={{
            backgroundColor: 'transparent',
          }}
          width="auto"
          placeholder={'Search Views'}
          value={filters.searchQuery}
          onChange={onSearchChange}
          endAdornment={
            <p style={{ color: theme.palette.text.default }}>Total Views: {total_count}</p>
          }
        />
        <Box sx={{ minWidth: 200 }}>
          <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
        </Box>
        <Box sx={{ minWidth: 150 }}>
          <VisibilitySelect
            visibility={filters.visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Box>
      </Box>
      <TableListHeader />
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
      />
    </Box>
  );
};

export default MyViewsContent;
