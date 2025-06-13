import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Box, FormControl, Grid2, InputLabel, MenuItem, Select, useTheme } from '@sistent/sistent';
import React, { useCallback, useState } from 'react';
import { StyledSearchBar } from '@sistent/sistent';
import MainDesignsContent from './MainDesignsContent';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import MainViewsContent from './MainViewsContent';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import {
  ImportButton,
  SortBySelect,
  TableListHeader,
  UserSearchAutoComplete,
  VisibilitySelect,
} from './components';
import { getDefaultFilterType } from './hooks';

const RecentContent = () => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

  const [filters, setFilters] = useState({
    type: getDefaultFilterType(),
    searchQuery: '',
    author: '',
    sortBy: 'updated_at desc',
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

  const handleAuthorChange = useCallback((user_id) => {
    setFilters((prev) => ({
      ...prev,
      author: user_id,
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
    refetch: refetchDesigns,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: filters.designsPage,
      pagesize: 10,
      order: filters.sortBy,
      metrics: true,
      search: filters.searchQuery,
      visibility: filters.visibility,
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
    refetch: refetchViews,
  } = useFetchViewsQuery(
    {
      page: filters.viewsPage,
      pagesize: 10,
      order: filters.sortBy,
      user_id: filters.author,
      visibility: filters.visibility,
      search: filters.searchQuery,
    },
    {
      skip: filters.type !== RESOURCE_TYPE.VIEW,
    },
  );

  const theme = useTheme();
  const refetch = useCallback(() => {
    if (filters.type === RESOURCE_TYPE.DESIGN) {
      if (filters.designsPage > 0) setDesignsPage(0);
      else refetchDesigns();
    } else {
      if (filters.viewsPage > 0) setViewsPage(0);
      else refetchViews();
    }
  }, [filters.type, filters.designsPage, filters.viewsPage, refetchDesigns, refetchViews]);

  return (
    <>
      <Box style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        <Grid2 container spacing={2} alignItems="center" size="grow">
          {/* Search Bar */}
          <Grid2 size={{ xs: 12, sm: 12, md: 6, lg: 3.5 }}>
            <StyledSearchBar
              sx={{ backgroundColor: 'transparent' }}
              width="auto"
              placeholder={
                filters.type === RESOURCE_TYPE.DESIGN ? 'Search Designs' : 'Search Views'
              }
              value={filters.searchQuery}
              onChange={onSearchChange}
              endAdornment={
                filters.type === RESOURCE_TYPE.DESIGN ? (
                  <p style={{ color: theme.palette.text.default, paddingLeft: "0.25rem" }}>
                    Total: {designsData?.total_count ?? 0}
                  </p>
                ) : (
                  <p style={{ color: theme.palette.text.default, paddingLeft: "0.25rem" }}>
                    Total: {viewsData?.total_count ?? 0}
                  </p>
                )
              }
            />
          </Grid2>

          {/* Type Select */}
          <Grid2 size={{ xs: 6, sm: 6, md: 3, lg: 1.5 }}>
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
                <MenuItem value={RESOURCE_TYPE.DESIGN}>Design</MenuItem>
                {isViewVisible && <MenuItem value={RESOURCE_TYPE.VIEW}>View</MenuItem>}
              </Select>
            </FormControl>
          </Grid2>

          {/* Sort By */}
          <Grid2 size={{ xs: 6, sm: 6, md: 3, lg: 1.5 }}>
            <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
          </Grid2>

          {/* Author Search */}
          <Grid2 size={{ xs: 12, sm: 6, md: 6, lg: 2.5 }}>
            <FormControl fullWidth>
              <UserSearchAutoComplete handleAuthorChange={handleAuthorChange} />
            </FormControl>
          </Grid2>

          {/* Visibility */}
          <Grid2 size={{ xs: 6, sm: 3, md: 3, lg: 1.5 }}>
            <VisibilitySelect
              visibility={filters.visibility}
              handleVisibilityChange={handleVisibilityChange}
              visibilityItems={visibilityItems}
            />
          </Grid2>

          {/* Import Button */}
          {filters.type === RESOURCE_TYPE.DESIGN && (
            <Grid2 size={{ xs: 6, sm: 3, md: 3, lg: 1 }}>
              <ImportButton refetch={refetch} />
            </Grid2>
          )}
        </Grid2>

        <>
          <TableListHeader />

          {filters.type == RESOURCE_TYPE.DESIGN && (
            <MainDesignsContent
              key={RESOURCE_TYPE.DESIGN}
              page={filters.designsPage}
              setPage={setDesignsPage}
              isLoading={isLoading}
              isFetching={isFetching}
              designs={designsData?.patterns}
              hasMore={
                designsData?.total_count > (filters.designsPage + 1) * designsData?.page_size
              }
              total_count={designsData?.total_count}
              refetch={refetch}
            />
          )}
          {filters.type == RESOURCE_TYPE.VIEW && (
            <MainViewsContent
              key={RESOURCE_TYPE.VIEW}
              page={filters.viewsPage}
              setPage={setViewsPage}
              isLoading={isViewLoading}
              isFetching={isViewFetching}
              views={viewsData?.views}
              hasMore={viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1)}
              total_count={viewsData?.total_count}
              refetch={refetch}
            />
          )}
        </>
      </Box>
    </>
  );
};

export default RecentContent;
