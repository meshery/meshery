import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import {
  Box,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  PromptComponent,
  Select,
  useTheme,
} from '@layer5/sistent';
import React, { useCallback, useRef, useState } from 'react';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import MainViewsContent from './MainViewsContent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import {
  ImportButton,
  MultiContentSelectToolbar,
  SortBySelect,
  TableListHeader,
  VisibilitySelect,
} from './components';
import { useGetDesignsOfWorkspaceQuery, useGetViewsOfWorkspaceQuery } from '@/rtk-query/workspace';
import { getDefaultFilterType, useContentDelete, useContentDownload } from './hooks';
import WorkspaceContentMoveModal from './WorkspaceContentMoveModal';
import ExportModal from '../ExportModal';

const WorkspaceContent = ({ workspace }) => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

  const [filters, setFilters] = useState({
    type: getDefaultFilterType(),
    searchQuery: '',
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
  } = useGetDesignsOfWorkspaceQuery(
    {
      infiniteScroll: true,
      workspaceId: workspace?.id,
      search: filters.searchQuery,
      page: filters.designsPage,
      pagesize: 10,
      order: filters.sortBy,
      visibility: filters.visibility,
    },
    {
      skip: filters.type !== RESOURCE_TYPE.DESIGN || !workspace?.id,
    },
  );

  const {
    data: viewsData,
    isLoading: isViewLoading,
    isFetching: isViewFetching,
    refetch: refetchViews,
  } = useGetViewsOfWorkspaceQuery(
    {
      infiniteScroll: true,
      workspaceId: workspace?.id,
      search: filters.searchQuery,
      page: filters.viewsPage,
      pagesize: 10,
      visibility: filters.visibility,
      order: 'updated_at desc',
    },
    {
      skip: filters.type !== RESOURCE_TYPE.VIEW || !workspace?.id,
    },
  );

  const theme = useTheme();

  const [workspaceContentMoveModal, setWorkspaceContentMoveModal] = useState(false);
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
  const { handleDesignDownload, handleViewDownload } = useContentDownload();

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
        <Grid container spacing={2} alignItems="center">
          {/* Search Bar */}
          <Grid item xs={12} md={5}>
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
          </Grid>

          {/* Type Select */}
          <Grid item xs={3} md={2}>
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
          </Grid>

          {/* Sort By */}
          <Grid item xs={3} md={2}>
            <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
          </Grid>

          {/* Visibility Select */}
          <Grid item xs={3} md={2}>
            <VisibilitySelect
              visibility={filters.visibility}
              handleVisibilityChange={handleVisibilityChange}
              visibilityItems={visibilityItems}
            />
          </Grid>

          <Grid item xs={3} md={1}>
            {filters.type == RESOURCE_TYPE.DESIGN && (
              <ImportButton
                refetch={refetch}
                workspaceId={workspace?.id}
                disabled={
                  !CAN(
                    keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
                    keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
                  )
                }
              />
            )}
          </Grid>
        </Grid>

        <>
          <MultiContentSelectToolbar
            type={filters.type}
            handleDelete={handleDelete}
            handleDownload={handleDownloadModalOpen}
            handleViewDownload={handleViewDownload}
            handleContentMove={setWorkspaceContentMoveModal}
            refetch={refetch}
          />
          <WorkspaceContentMoveModal
            workspaceContentMoveModal={workspaceContentMoveModal}
            setWorkspaceContentMoveModal={setWorkspaceContentMoveModal}
            currentWorkspace={workspace}
            type={filters.type}
            refetch={refetch}
          />

          <TableListHeader
            content={
              RESOURCE_TYPE.DESIGN === filters.type ? designsData?.designs : viewsData?.views
            }
            isMultiSelectMode={true}
          />

          {filters.type == RESOURCE_TYPE.DESIGN && (
            <MainDesignsContent
              page={filters.designsPage}
              setPage={setDesignsPage}
              isLoading={isLoading}
              isFetching={isFetching}
              designs={designsData?.designs}
              hasMore={designsData?.total_count > designsData?.page_size * (designsData?.page + 1)}
              total_count={designsData?.total_count}
              workspace={workspace}
              refetch={refetch}
              isMultiSelectMode={true}
            />
          )}
          {filters.type == RESOURCE_TYPE.VIEW && (
            <MainViewsContent
              page={filters.viewsPage}
              setPage={setViewsPage}
              isLoading={isViewLoading}
              isFetching={isViewFetching}
              views={viewsData?.views}
              hasMore={viewsData?.total_count > viewsData?.page_size * (viewsData?.page + 1)}
              total_count={viewsData?.total_count}
              workspace={workspace}
              refetch={refetch}
              isMultiSelectMode={true}
            />
          )}
        </>
      </Box>
      <PromptComponent ref={modalRef} />
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadModalClose}
        handleDesignDownload={handleDesignDownload}
      />
    </>
  );
};

export default WorkspaceContent;
