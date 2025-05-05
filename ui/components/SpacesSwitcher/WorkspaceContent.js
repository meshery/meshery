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
import { VISIBILITY } from '@/utils/Enum';
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

const WorkspaceContent = ({ workspace }) => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);

  const [searchQuery, setSearchQuery] = useState('');
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];
  const [type, setType] = React.useState('design');
  const [sortBy, setSortBy] = useState('updated_at desc');
  const [visibility, setVisibility] = useState(visibilityItems);

  const handleTypeChange = useCallback((event) => {
    setType(event.target.value);
    setDesignsPage(0);
    setViewsPage(0);
  }, []);
  const handleSortByChange = useCallback((event) => {
    setDesignsPage(0);
    setViewsPage(0);
    setSortBy(event.target.value);
  }, []);

  const handleVisibilityChange = useCallback((event) => {
    const value = event.target.value;
    setVisibility(typeof value === 'string' ? value.split(',') : value);
    setDesignsPage(0);
    setViewsPage(0);
  }, []);

  const onSearchChange = useCallback((e) => {
    setDesignsPage(0);
    setViewsPage(0);
    setSearchQuery(e.target.value);
  }, []);

  const [designsPage, setDesignsPage] = useState(0);
  const [viewsPage, setViewsPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetDesignsOfWorkspaceQuery(
    {
      infiniteScroll: true,
      workspaceId: workspace?.id,
      search: searchQuery,
      page: designsPage,
      pagesize: 10,
      order: sortBy,
      visibility: visibility,
    },
    {
      skip: type !== 'design',
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
      search: searchQuery,
      page: viewsPage,
      pagesize: 10,
      visibility: visibility,
      order: 'updated_at desc',
    },
    {
      skip: type !== 'view' || !workspace?.id,
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
          />{' '}
          <AssignDesignViewButton
            type={type}
            handleAssign={(e) => {
              e.stopPropagation();
              if (type === 'design') {
                designAssignment.handleAssignModal();
              } else {
                viewAssignment.handleAssignModal();
              }
            }}
            disabled={
              type === 'design'
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
                value={type}
                label="Type"
                onChange={handleTypeChange}
                sx={{
                  '& .MuiSelect-select': {
                    paddingBlock: '0.85rem',
                  },
                }}
              >
                {isDesignsVisible && <MenuItem value={'design'}>Design</MenuItem>}
                {isViewVisible && <MenuItem value={'view'}>View</MenuItem>}
              </Select>
            </FormControl>
          </Box>

          <Box sx={{ minWidth: 120 }}>
            <SortBySelect sortBy={sortBy} handleSortByChange={handleSortByChange} />
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
              designs={designsData?.designs}
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
