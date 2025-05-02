//@ts-check
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Box, FormControl, InputLabel, MenuItem, Select, useTheme } from '@layer5/sistent';
import React, { useState } from 'react';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import MainViewsContent from './MainViewsContent';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { VISIBILITY } from '@/utils/Enum';
import {
  SortBySelect,
  TableListHeader,
  UserSearchAutoComplete,
  VisibilitySelect,
} from './components';

const RecentContent = () => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);

  const [searchQuery, setSearchQuery] = useState('');
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

  const [type, setType] = React.useState('design');
  const [author, setAuthor] = React.useState('');
  const [sortBy, setSortBy] = useState('updated_at desc');
  const [visibility, setVisibility] = useState(visibilityItems);

  const handleTypeChange = (event) => {
    setType(event.target.value);
    setDesignsPage(0);
    setViewsPage(0);
  };
  const handleAuthorChange = (user_id) => {
    setDesignsPage(0);
    setViewsPage(0);
    setAuthor(user_id);
  };
  const handleSortByChange = (event) => {
    setDesignsPage(0);
    setViewsPage(0);
    setSortBy(event.target.value);
  };
  const handleVisibilityChange = (event) => {
    const value = event.target.value;
    setVisibility(typeof value === 'string' ? value.split(',') : value);
    setDesignsPage(0);
    setViewsPage(0);
  };
  const onSearchChange = (e) => {
    setDesignsPage(0);
    setViewsPage(0);
    setSearchQuery(e.target.value);
  };

  const [designsPage, setDesignsPage] = useState(0);
  const [viewsPage, setViewsPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: designsPage,
      pagesize: 10,
      order: sortBy,
      metrics: true,
      search: searchQuery,
      visibility: visibility,
      user_id: author,
    },
    {
      skip: type !== 'design',
    },
  );

  const {
    data: viewsData,
    isLoading: isViewLoading,
    isFetching: isViewFetching,
  } = useFetchViewsQuery(
    {
      page: viewsPage,
      pagesize: 10,
      order: sortBy,
      user_id: author,
      visibility: visibility,
      search: searchQuery,
    },
    {
      skip: type !== 'view',
    },
  );

  const theme = useTheme();
  return (
    <Box style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <StyledSearchBar
        sx={{
          backgroundColor: 'transparent',
        }}
        width="auto"
        placeholder={type === 'design' ? 'Search Designs' : 'Search Views'}
        value={searchQuery}
        onChange={onSearchChange}
        endAdornment={
          type === 'design' ? (
            <p style={{ color: theme.palette.text.default }}>
              Total Designs: {designsData?.total_count ?? 0}
            </p>
          ) : (
            <p style={{ color: theme.palette.text.default }}>
              Total Views: {viewsData?.total_count ?? 0}
            </p>
          )
        }
      />
      <Box display={'flex'} alignItems="center" marginBottom="1rem" gap={'1rem'}>
        <Box sx={{ minWidth: 120 }}>
          <FormControl fullWidth>
            <InputLabel>Type</InputLabel>
            <Select
              value={type}
              label="Type"
              onChange={handleTypeChange}
              sx={{
                '& .MuiSelect-select': {
                  paddingBlock: '0.85rem',
                },
              }}
            >
              <MenuItem value={'design'}>Design</MenuItem>
              <MenuItem value={'view'}>View</MenuItem>
            </Select>
          </FormControl>
        </Box>

        <Box sx={{ minWidth: 120 }}>
          <SortBySelect sortBy={sortBy} handleSortByChange={handleSortByChange} />
        </Box>
        <Box sx={{ minWidth: 300 }}>
          <FormControl fullWidth>
            <UserSearchAutoComplete handleAuthorChange={handleAuthorChange} />
          </FormControl>
        </Box>
        <Box sx={{ minWidth: 120 }}>
          <VisibilitySelect
            visibility={visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Box>
      </Box>
      <Box minWidth={'50rem'}>
        <TableListHeader />

        {type == 'design' && (
          <MainDesignsContent
            setPage={setDesignsPage}
            isLoading={isLoading}
            isFetching={isFetching}
            designs={designsData?.patterns}
            hasMore={designsData?.total_count > designsData?.page_size * (designsData?.page + 1)}
            total_count={designsData?.total_count}
          />
        )}
        {type == 'view' && (
          <MainViewsContent
            setPage={setViewsPage}
            isLoading={isViewLoading}
            isFetching={isViewFetching}
            views={viewsData?.views}
            hasMore={viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1)}
            total_count={viewsData?.total_count}
          />
        )}
      </Box>
    </Box>
  );
};

export default RecentContent;
