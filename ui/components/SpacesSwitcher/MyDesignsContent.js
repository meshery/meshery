//@ts-check
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
    setSearchQuery(e.target.value);
    setPage(0);
  }, []);

  const [page, setPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: page,
      pagesize: 10,
      order: 'updated_at desc',
      user_id: currentUser?.id,
      metrics: true,
      visibility: visibility,
      search: searchQuery,
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
          value={searchQuery}
          onChange={onSearchChange}
          endAdornment={
            <p style={{ color: theme.palette.text.default }}>Total Designs: {total_count}</p>
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
        <ImportButton />
      </Box>
      <TableListHeader />
      <MainDesignsContent
        designs={designsData?.patterns}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        hasMore={hasMore}
        total_count={total_count}
      />
    </Box>
  );
};

export default MyDesignsContent;
