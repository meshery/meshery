import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import {
  AssignmentModal,
  Box,
  Button,
  DeleteIcon,
  DesignIcon,
  DownloadIcon,
  EnvironmentIcon,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  PromptComponent,
  Select,
  Typography,
  useDesignAssignment,
  useTheme,
  useViewAssignment,
} from '@layer5/sistent';
import React, { useCallback, useContext, useRef, useState } from 'react';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import MainViewsContent from './MainViewsContent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import { ImportButton, SortBySelect, TableListHeader, VisibilitySelect } from './components';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignViewFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import { getDefaultFilterType, useContentDelete, useContentDownload } from './hooks';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import WorkspaceContentMoveModal from './WorkspaceContentMoveModal';
import { iconMedium } from 'css/icons.styles';
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
  const { multiSelectedContent } = useContext(WorkspaceModalContext);
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

  return (
    <>
      <Box style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        <Grid container spacing={2} alignItems="center">
          {/* Search Bar */}
          <Grid item xs={12} sm={12} md={6} lg={5}>
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
          <Grid item xs={6} sm={6} md={2} lg={1.5}>
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
          <Grid item xs={6} sm={6} md={2} lg={1.5}>
            <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
          </Grid>

          {/* Visibility Select */}
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <VisibilitySelect
              visibility={filters.visibility}
              handleVisibilityChange={handleVisibilityChange}
              visibilityItems={visibilityItems}
            />
          </Grid>

          {/* Assign Button */}
          <Grid item xs={12} sm={6} md={3} lg={2}>
            {/* <AssignDesignViewButton
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
            /> */}
            {filters.type == RESOURCE_TYPE.DESIGN && (
              <ImportButton
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
          {multiSelectedContent.length > 0 && (
            <Box
              width={'100%'}
              sx={{ backgroundColor: theme.palette.background.default }}
              height={'4rem'}
              borderRadius={'0.5rem'}
              display={'flex'}
              justifyContent={'space-between'}
              alignItems={'center'}
              paddingInline={'1rem'}
            >
              <Typography>
                {multiSelectedContent.length} {filters.type} selected
              </Typography>
              <Box style={{ display: 'flex', gap: '0.5rem' }}>
                <Button
                  variant="contained"
                  startIcon={<MoveFileIcon style={iconMedium} />}
                  onClick={() => setWorkspaceContentMoveModal(true)}
                  disabled={!multiSelectedContent.length}
                >
                  Move
                </Button>
                <Button
                  variant="contained"
                  startIcon={<DownloadIcon style={iconMedium} fill={theme.palette.common.white} />}
                  onClick={() =>
                    filters.type === RESOURCE_TYPE.DESIGN
                      ? handleDownloadModalOpen(multiSelectedContent)
                      : handleViewDownload(multiSelectedContent)
                  }
                  disabled={!multiSelectedContent.length}
                >
                  Download
                </Button>
                <Button
                  color="error"
                  variant="contained"
                  onClick={() => {
                    handleDelete(
                      multiSelectedContent,
                      filters.type === RESOURCE_TYPE.DESIGN
                        ? RESOURCE_TYPE.DESIGN
                        : RESOURCE_TYPE.VIEW,
                    );
                  }}
                  sx={{
                    backgroundColor: `${theme.palette.error.dark} !important`,
                  }}
                  startIcon={<DeleteIcon style={iconMedium} fill={theme.palette.common.white} />}
                >
                  Delete
                </Button>
              </Box>
            </Box>
          )}
          <WorkspaceContentMoveModal
            workspaceContentMoveModal={workspaceContentMoveModal}
            setWorkspaceContentMoveModal={setWorkspaceContentMoveModal}
            currentWorkspace={workspace}
            type={filters.type}
          />

          <TableListHeader content={designsData?.designs} isMultiSelectMode={true} />

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
              refetch={() => setDesignsPage(0)}
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
              refetch={() => setViewsPage(0)}
              isMultiSelectMode={true}
            />
          )}
        </>
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
