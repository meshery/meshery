import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Box, FormControl, InputLabel, MenuItem, Select, useTheme } from '@layer5/sistent';
import React, { useCallback, useState } from 'react';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import MainViewsContent from './MainViewsContent';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import {
  SortBySelect,
  TableListHeader,
  UserSearchAutoComplete,
  VisibilitySelect,
} from './components';
import { getDefaultFilterType } from './hooks';

const SharedContent = () => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);

  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

  const [filters, setFilters] = useState({
    type: getDefaultFilterType(),
    searchQuery: '',
    sortBy: 'updated_at desc',
    author: '',
    visibility: visibilityItems,
    designsPage: 0,
    viewsPage: 0,
  });

  const handleTypeChange = useCallback((event) => {
    setFilters((prev) => ({
      ...prev,
      type: event.target.value,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const handleSortByChange = useCallback((event) => {
    setFilters((prev) => ({
      ...prev,
      sortBy: event.target.value,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const handleVisibilityChange = useCallback((event) => {
    const value = event.target.value;
    setFilters((prev) => ({
      ...prev,
      visibility: typeof value === 'string' ? value.split(',') : value,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const handleAuthorChange = useCallback((user_id) => {
    setFilters((prev) => ({
      ...prev,
      author: user_id,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const onSearchChange = useCallback((e) => {
    setFilters((prev) => ({
      ...prev,
      searchQuery: e.target.value,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const setDesignsPage = useCallback((page) => {
    setFilters((prev) => ({
      ...prev,
      designsPage: page,
    }));
  }, []);

  const setViewsPage = useCallback((page) => {
    setFilters((prev) => ({
      ...prev,
      viewsPage: page,
    }));
  }, []);

  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: filters.designsPage,
      pagesize: 10,
      order: filters.sortBy,
      metrics: true,
      search: filters.searchQuery,
      visibility: filters.visibility,
      shared: true,
      user_id: filters.author,
    },
    {
      skip: filters.type !== RESOURCE_TYPE.DESIGN,
    },
  );

  const {
    data: viewsData,
    isLoading: isViewLoading,
    isFetching: isViewFetching,
  } = useFetchViewsQuery(
    {
      page: filters.viewsPage,
      pagesize: 10,
      order: filters.sortBy,
      visibility: filters.visibility,
      search: filters.searchQuery,
      user_id: filters.author,
      shared: true,
    },
    {
      skip: filters.type !== RESOURCE_TYPE.VIEW,
    },
  );

  const theme = useTheme();

  return (
    <>
      <Box style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        <Box display={'flex'} alignItems="center" gap={'1rem'}>
          <StyledSearchBar
            sx={{
              backgroundColor: 'transparent',
            }}
            width="auto"
            placeholder={filters.type === RESOURCE_TYPE.DESIGN ? 'Search Designs' : 'Search Views'}
            value={filters.searchQuery}
            onChange={onSearchChange}
            endAdornment={
              filters.type === RESOURCE_TYPE.DESIGN ? (
                <p style={{ color: theme.palette.text.default }}>
                  Total Designs: {designsData?.total_count ?? 0}
                </p>
              ) : (
                <p style={{ color: theme.palette.text.default }}>
                  Total Views: {viewsData?.total_count ?? 0}
                </p>
              )
            }
          />{' '}
        </Box>
        <Box display={'flex'} alignItems="center" marginBottom="1rem" gap={'1rem'}>
          <Box sx={{ minWidth: 120 }}>
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={filters.type}
                label="Type"
                onChange={handleTypeChange}
                sx={{
                  '& .MuiSelect-select': {
                    paddingBlock: '0.85rem',
                  },
                }}
              >
                {isDesignsVisible && <MenuItem value={RESOURCE_TYPE.DESIGN}>Design</MenuItem>}
                {isViewVisible && <MenuItem value={RESOURCE_TYPE.VIEW}>View</MenuItem>}
              </Select>
            </FormControl>
          </Box>

          <Box sx={{ minWidth: 120 }}>
            <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
          </Box>
          <Box sx={{ minWidth: 300 }}>
            <FormControl fullWidth>
              <UserSearchAutoComplete handleAuthorChange={handleAuthorChange} />
            </FormControl>
          </Box>
          <Box sx={{ minWidth: 120 }}>
            <VisibilitySelect
              visibility={filters.visibility}
              handleVisibilityChange={handleVisibilityChange}
              visibilityItems={visibilityItems}
            />
          </Box>
        </Box>
        <Box minWidth={'50rem'}>
          <TableListHeader />

          {filters.type == RESOURCE_TYPE.DESIGN && (
            <MainDesignsContent
              key={'shared-designs-content'}
              page={filters.designsPage}
              setPage={setDesignsPage}
              isLoading={isLoading}
              isFetching={isFetching}
              designs={designsData?.patterns}
              hasMore={
                designsData?.total_count > (filters.designsPage + 1) * designsData?.page_size
              }
              total_count={designsData?.total_count}
              refetch={() => setDesignsPage(0)}
            />
          )}
          {filters.type == RESOURCE_TYPE.VIEW && (
            <MainViewsContent
              key={'shared-views-content'}
              page={filters.viewsPage}
              setPage={setViewsPage}
              isLoading={isViewLoading}
              isFetching={isViewFetching}
              views={viewsData?.views}
              hasMore={viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1)}
              total_count={viewsData?.total_count}
              refetch={() => setViewsPage(0)}
            />
          )}
        </Box>
      </Box>
    </>
  );
};

export default SharedContent;
