import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import {
  AssignmentModal,
  Box,
  DesignIcon,
  EnvironmentIcon,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  useDesignAssignment,
  useTheme,
  useViewAssignment,
} from '@layer5/sistent';
import React, { useCallback, useState } from 'react';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import MainViewsContent from './MainViewsContent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import {
  AssignDesignViewButton,
  SortBySelect,
  TableListHeader,
  VisibilitySelect,
} from './components';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import { getDefaultFilterType } from './hooks';

const WorkspaceContent = ({ workspace }) => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);

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
  const viewAssignment = useViewAssignment({
    workspaceId: workspace?.id,
    useGetViewsOfWorkspaceQuery,
    useUnassignViewFromWorkspaceMutation,
    useAssignViewToWorkspaceMutation,
    isViewsVisible: CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject),
  });

  const designAssignment = useDesignAssignment({
    workspaceId: workspace?.id,
    useAssignDesignToWorkspaceMutation,
    useUnassignDesignFromWorkspaceMutation,
    useGetDesignsOfWorkspaceQuery: useGetDesignsOfWorkspaceQuery,
    isDesignsVisible: CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject),
  });
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
              filters.type === RESOURCE_TYPE.VIEW ? (
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
          <AssignDesignViewButton
            type={filters.type}
            handleAssign={(e) => {
              e.stopPropagation();
              if (filters.type === RESOURCE_TYPE.DESIGN) {
                designAssignment.handleAssignModal();
              } else {
                viewAssignment.handleAssignModal();
              }
            }}
            disabled={
              filters.type === RESOURCE_TYPE.DESIGN
                ? !CAN(
                    keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
                    keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
                  )
                : !CAN(
                    keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
                    keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
                  )
            }
          />
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
              page={filters.designsPage}
              setPage={setDesignsPage}
              isLoading={isLoading}
              isFetching={isFetching}
              designs={designsData?.designs}
              hasMore={designsData?.total_count > designsData?.page_size * (designsData?.page + 1)}
              total_count={designsData?.total_count}
              workspaceId={workspace?.id}
              refetch={() => setDesignsPage(0)}
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
              workspaceId={workspace?.id}
              refetch={() => setViewsPage(0)}
            />
          )}
        </Box>
      </Box>
      <AssignmentModal
        open={viewAssignment.assignModal}
        onClose={viewAssignment.handleAssignModalClose}
        title={`Assign Views to ${workspace?.name}`}
        headerIcon={<EnvironmentIcon height="40" width="40" fill={'white'} />}
        name="Views"
        assignableData={viewAssignment.data}
        handleAssignedData={viewAssignment.handleAssignData}
        originalAssignedData={viewAssignment.workspaceData}
        emptyStateIcon={<EnvironmentIcon height="5rem" width="5rem" fill={'#808080'} />}
        handleAssignablePage={viewAssignment.handleAssignablePage}
        handleAssignedPage={viewAssignment.handleAssignedPage}
        originalLeftCount={viewAssignment.data?.length || 0}
        originalRightCount={viewsData?.total_count || 0}
        onAssign={viewAssignment.handleAssign}
        disableTransfer={viewAssignment.disableTransferButton}
        helpText={`Assign Views to ${workspace?.name}`}
        isAssignAllowed={CAN(
          keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
          keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
        )}
        isRemoveAllowed={CAN(
          keys.REMOVE_VIEWS_FROM_WORKSPACE.action,
          keys.REMOVE_VIEWS_FROM_WORKSPACE.subject,
        )}
      />
      <AssignmentModal
        open={designAssignment.assignModal}
        onClose={designAssignment.handleAssignModalClose}
        title={`Assign Designs to ${workspace?.name}`}
        headerIcon={<DesignIcon height="40" width="40" />}
        name="Designs"
        assignableData={designAssignment.data}
        handleAssignedData={designAssignment.handleAssignData}
        originalAssignedData={designAssignment.workspaceData}
        emptyStateIcon={<DesignIcon height="5rem" width="5rem" secondaryFill={'#808080'} />}
        handleAssignablePage={designAssignment.handleAssignablePage}
        handleAssignedPage={designAssignment.handleAssignedPage}
        originalLeftCount={designAssignment.data?.length || 0}
        originalRightCount={designAssignment.assignedItems?.length || 0}
        onAssign={designAssignment.handleAssign}
        disableTransfer={designAssignment.disableTransferButton}
        helpText={`Assign Designs to ${workspace?.name}`}
        isAssignAllowed={CAN(
          keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
          keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
        )}
        isRemoveAllowed={CAN(
          keys.REMOVE_DESIGNS_FROM_WORKSPACE.action,
          keys.REMOVE_DESIGNS_FROM_WORKSPACE.subject,
        )}
        showViews={false}
      />
    </>
  );
};

export default WorkspaceContent;
