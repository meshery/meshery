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
  const [visibility, setVisibility] = useState(visibilityItems);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState('updated_at desc');

  const handleSortByChange = useCallback((event) => {
    setSortBy(event.target.value);
    setPage(0);
  }, []);

  const handleVisibilityChange = useCallback((event) => {
    const value = event.target.value;
    setVisibility(typeof value === 'string' ? value.split(',') : value);
    setPage(0);
  }, []);

  const onSearchChange = useCallback((e) => {
    setPage(0);
    setSearchQuery(e.target.value);
  }, []);

  const [page, setPage] = useState(0);
  const {
    data: viewsData,
    isLoading,
    isFetching,
  } = useFetchViewsQuery(
    {
      page: page,
      pagesize: 10,
      order: sortBy,
      visibility: visibility,
      search: searchQuery,
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
          value={searchQuery}
          onChange={onSearchChange}
          endAdornment={
            <p style={{ color: theme.palette.text.default }}>Total Views: {total_count}</p>
          }
        />
        <Box sx={{ minWidth: 200 }}>
          <SortBySelect sortBy={sortBy} handleSortByChange={handleSortByChange} />
        </Box>
        <Box sx={{ minWidth: 150 }}>
          <VisibilitySelect
            visibility={visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Box>
      </Box>
      <TableListHeader />
      <MainViewsContent
        hasMore={hasMore}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        views={views}
        total_count={total_count}
      />
    </Box>
  );
};

export default MyViewsContent;
