import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import React, { useCallback, useState } from 'react';
import MainDesignsContent from './MainDesignsContent';
import { VISIBILITY } from '@/utils/Enum';
import { Box, useTheme } from '@layer5/sistent';
import { StyledSearchBar } from '@layer5/sistent';
import { ImportButton, SortBySelect, TableListHeader, VisibilitySelect } from './components';

const MyDesignsContent = () => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE, VISIBILITY.PUBLISHED];

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
    },
    {
      skip: !currentUser?.id,
    },
  );
  const hasMore = designsData?.total_count > designsData?.page_size * (designsData?.page + 1);
  const total_count = designsData?.total_count || 0;
  const theme = useTheme();
  return (
    <Box display={'flex'} flexDirection="column" gap="1rem">
      <Box display={'flex'} alignItems="center" marginBottom="1rem" gap={'1rem'}>
        <StyledSearchBar
          sx={{
            backgroundColor: 'transparent',
          }}
          width="auto"
          placeholder={'Search Designs'}
          value={filters.searchQuery}
          onChange={onSearchChange}
          endAdornment={
            <p style={{ color: theme.palette.text.default }}>Total Designs: {total_count}</p>
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
        <ImportButton />
      </Box>
      <TableListHeader />
      <MainDesignsContent
        key={'my-designs-content'}
        page={filters.page}
        designs={designsData?.patterns}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        hasMore={hasMore}
        refetch={() => setPage(0)}
        total_count={total_count}
      />
    </Box>
  );
};

export default MyDesignsContent;
