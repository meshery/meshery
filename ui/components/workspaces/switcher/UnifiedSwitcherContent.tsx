import React, { useCallback, useRef, useState, useMemo } from 'react';
import { useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import {
  Box,
  Grid2,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  useTheme,
  StyledSearchBar,
  PromptComponent,
  WorkspaceContentMoveModal,
} from '@sistent/sistent';
import { RESOURCE_TYPE, VISIBILITY } from '@/utils/Enum';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useFetchViewsQuery } from '@/rtk-query/view';
import {
  useGetDesignsOfWorkspaceQuery,
  useGetViewsOfWorkspaceQuery,
  useGetWorkspacesQuery,
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
} from '@/rtk-query/workspace';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  ImportButton,
  MultiContentSelectToolbar,
  SortBySelect,
  TableListHeader,
  VisibilitySelect,
  UserSearchAutoComplete,
} from './components';
import {
  getDefaultFilterType,
  useContentDelete,
  useContentDownload,
} from '@/components/workspaces/switcher/hooks';
import UnifiedResourceList from './UnifiedResourceList';
import ExportDesignModal from '@/components/designs/export/ExportDesignModal';
import ShareModal from '@/components/workspaces/ShareWorkspaceModal';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { useNotification } from '@/utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

export type SwitcherVariant = 'recent' | 'workspace' | 'my' | 'shared';

interface UnifiedSwitcherContentProps {
  variant: SwitcherVariant;
  workspace?: any;
  initialType?: RESOURCE_TYPE;
}

/**
 * Unified Switcher Content component that handles different views (recent, workspace, my, shared).
 * Standardizes filtering, querying, and resource listing.
 */
const UnifiedSwitcherContent: React.FC<UnifiedSwitcherContentProps> = ({
  variant,
  workspace,
  initialType,
}) => {
  const theme = useTheme();
  const router = useRouter();
  const { notify } = useNotification();
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const { organization: currentOrganization } = useSelector((state: any) => state.ui);

  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isWorkspace = variant === 'workspace';
  const isMy = variant === 'my';
  const isShared = variant === 'shared';
  const isRecent = variant === 'recent';

  const visibilityItems = useMemo(() => {
    const base = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];
    if (isMy) base.push(VISIBILITY.PUBLISHED);
    return base;
  }, [isMy]);

  const [filters, setFilters] = useState({
    type: initialType || getDefaultFilterType(),
    searchQuery: '',
    sortBy: 'updated_at desc',
    author: '',
    visibility: visibilityItems,
    designsPage: 0,
    viewsPage: 0,
  });

  // State for modals
  const [shareModal, setShareModal] = useState({ open: false, content: null });
  const [downloadModal, setDownloadModal] = useState({ open: false, content: null });
  const [moveModalOpen, setMoveModalOpen] = useState(false);
  const modalRef = useRef<any>(null);

  const { handleDelete } = useContentDelete(modalRef);
  const { handleDesignDownload, handleViewDownload } = useContentDownload();

  // Handlers
  const handleTypeChange = useCallback((e: any) => {
    setFilters((prev) => ({ ...prev, type: e.target.value, designsPage: 0, viewsPage: 0 }));
  }, []);

  const handleSortByChange = useCallback((e: any) => {
    setFilters((prev) => ({ ...prev, sortBy: e.target.value, designsPage: 0, viewsPage: 0 }));
  }, []);

  const handleVisibilityChange = useCallback((e: any) => {
    const value = e.target.value;
    setFilters((prev) => ({
      ...prev,
      visibility: typeof value === 'string' ? value.split(',') : value,
      designsPage: 0,
      viewsPage: 0,
    }));
  }, []);

  const handleAuthorChange = useCallback((user_id: string) => {
    setFilters((prev) => ({ ...prev, author: user_id, designsPage: 0, viewsPage: 0 }));
  }, []);

  const onSearchChange = useCallback((e: any) => {
    setFilters((prev) => ({ ...prev, searchQuery: e.target.value, designsPage: 0, viewsPage: 0 }));
  }, []);

  // Queries
  const designQueryArgs = {
    page: filters.designsPage,
    pagesize: 10,
    order: filters.sortBy,
    search: filters.searchQuery,
    visibility: filters.visibility,
    user_id: isMy ? currentUser?.id : filters.author || undefined,
    orgID: currentOrganization?.id,
    expandUser: true,
    metrics: true,
    shared: isShared ? true : undefined,
    infiniteScroll: isWorkspace,
    workspaceId: isWorkspace ? workspace?.id : undefined,
  };

  const viewQueryArgs = {
    page: filters.viewsPage,
    pagesize: 10,
    order: filters.sortBy,
    search: filters.searchQuery,
    visibility: filters.visibility,
    user_id: isMy ? currentUser?.id : filters.author || undefined,
    shared: isShared ? true : undefined,
    infiniteScroll: isWorkspace,
    workspaceId: isWorkspace ? workspace?.id : undefined,
  };

  // Select query hooks
  const {
    data: designsData,
    isLoading: designsLoading,
    isFetching: designsFetching,
    refetch: refetchDesigns,
  } = isWorkspace
    ? useGetDesignsOfWorkspaceQuery(designQueryArgs, {
        skip: filters.type !== RESOURCE_TYPE.DESIGN || !workspace?.id,
      })
    : useGetUserDesignsQuery(designQueryArgs, {
        skip: filters.type !== RESOURCE_TYPE.DESIGN || (isMy && !currentUser?.id),
      });

  const {
    data: viewsData,
    isLoading: viewsLoading,
    isFetching: viewsFetching,
    refetch: refetchViews,
  } = isWorkspace
    ? useGetViewsOfWorkspaceQuery(viewQueryArgs, {
        skip: filters.type !== RESOURCE_TYPE.VIEW || !workspace?.id,
      })
    : useFetchViewsQuery(viewQueryArgs, {
        skip: filters.type !== RESOURCE_TYPE.VIEW || (isMy && !currentUser?.id),
      });

  const refetch = useCallback(() => {
    if (filters.type === RESOURCE_TYPE.DESIGN) {
      filters.designsPage > 0 ? setFilters((p) => ({ ...p, designsPage: 0 })) : refetchDesigns();
    } else {
      filters.viewsPage > 0 ? setFilters((p) => ({ ...p, viewsPage: 0 })) : refetchViews();
    }
  }, [filters.type, filters.designsPage, filters.viewsPage, refetchDesigns, refetchViews]);

  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();
  const [assignViewToWorkspace] = useAssignViewToWorkspaceMutation();

  const isDesign = filters.type === RESOURCE_TYPE.DESIGN;
  const currentItems = isDesign ? designsData?.patterns || designsData?.designs : viewsData?.views;
  const totalCount = (isDesign ? designsData?.total_count : viewsData?.total_count) || 0;
  const hasMore = isDesign
    ? designsData?.total_count > (filters.designsPage + 1) * (designsData?.page_size || 10)
    : viewsData?.total_count > (viewsData?.page_size || 10) * (viewsData?.page + 1);

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <Grid2 container spacing={2} alignItems="center">
        {/* Search Bar */}
        <Grid2 size={{ xs: 12, md: isRecent || isShared ? 5 : 6.5 }}>
          <StyledSearchBar
            sx={{ backgroundColor: 'transparent' }}
            placeholder={`Search ${isDesign ? 'Designs' : 'Views'}`}
            value={filters.searchQuery}
            onChange={onSearchChange}
            endAdornment={<p style={{ color: theme.palette.text.default }}>Total: {totalCount}</p>}
          />
        </Grid2>

        {/* Filters - Type Select (only for recent/shared/workspace) */}
        {!isMy && (
          <Grid2 size={{ xs: 6, md: 1.5 }}>
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={filters.type}
                label="Type"
                onChange={handleTypeChange}
                sx={{ '& .MuiSelect-select': { paddingBlock: '0.85rem' } }}
              >
                <MenuItem value={RESOURCE_TYPE.DESIGN}>Design</MenuItem>
                {isViewVisible && <MenuItem value={RESOURCE_TYPE.VIEW}>View</MenuItem>}
              </Select>
            </FormControl>
          </Grid2>
        )}

        {/* Sort By */}
        <Grid2 size={{ xs: 6, md: isMy ? 2 : 1.5 }}>
          <SortBySelect sortBy={filters.sortBy} handleSortByChange={handleSortByChange} />
        </Grid2>

        {/* Author Search (only for recent/shared) */}
        {(isRecent || isShared) && (
          <Grid2 size={{ xs: 12, md: 2.5 }}>
            <UserSearchAutoComplete handleAuthorChange={handleAuthorChange} />
          </Grid2>
        )}

        {/* Visibility */}
        <Grid2 size={{ xs: 6, md: isMy ? 2 : 1.5 }}>
          <VisibilitySelect
            visibility={filters.visibility}
            handleVisibilityChange={handleVisibilityChange}
            visibilityItems={visibilityItems}
          />
        </Grid2>

        {/* Import */}
        {isDesign && (
          <Grid2 size={{ xs: 6, md: 1 }}>
            <ImportButton refetch={refetch} workspaceId={workspace?.id} />
          </Grid2>
        )}
      </Grid2>

      <Box>
        <MultiContentSelectToolbar
          type={filters.type}
          handleDelete={handleDelete}
          handleDownload={(c) => setDownloadModal({ open: true, content: c })}
          handleViewDownload={handleViewDownload}
          handleContentMove={setMoveModalOpen}
          refetch={refetch}
          handleShare={(c) => setShareModal({ open: true, content: c })}
        />

        <TableListHeader
          isMultiSelectMode={true}
          content={currentItems}
          showOrganizationName={!isWorkspace}
          showWorkspaceName={!isWorkspace}
        />

        <UnifiedResourceList
          type={filters.type as any}
          items={currentItems}
          page={isDesign ? filters.designsPage : filters.viewsPage}
          setPage={(p) =>
            setFilters((prev) => ({ ...prev, [isDesign ? 'designsPage' : 'viewsPage']: p }))
          }
          isLoading={isDesign ? designsLoading : viewsLoading}
          isFetching={isDesign ? designsFetching : viewsFetching}
          hasMore={hasMore}
          total_count={totalCount}
          workspace={workspace}
          refetch={refetch}
          isMultiSelectMode={true}
          showWorkspaceName={!isWorkspace}
          showOrganizationName={!isWorkspace}
        />
      </Box>

      <PromptComponent ref={modalRef} />

      <ExportDesignModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={() => setDownloadModal({ open: false, content: null })}
        handleDesignDownload={handleDesignDownload}
      />

      {shareModal.open && (
        <ShareModal
          resource={shareModal.content}
          handleClose={() => setShareModal({ open: false, content: null })}
          type={filters.type}
        />
      )}

      {moveModalOpen && (
        <WorkspaceContentMoveModal
          workspaceContentMoveModal={moveModalOpen}
          setWorkspaceContentMoveModal={setMoveModalOpen}
          currentWorkspace={workspace}
          type={filters.type}
          refetch={refetch}
          useGetWorkspacesQuery={useGetWorkspacesQuery}
          WorkspaceModalContext={WorkspaceModalContext}
          assignDesignToWorkspace={assignDesignToWorkspace}
          assignViewToWorkspace={assignViewToWorkspace}
          isCreateWorkspaceAllowed={CAN(
            keys.CREATE_WORKSPACE.action,
            keys.CREATE_WORKSPACE.subject,
          )}
          isMoveDesignAllowed={CAN(
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
          )}
          isMoveViewAllowed={CAN(
            keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
            keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
          )}
          currentOrgId={currentOrganization?.id}
          notify={notify}
          router={router}
        />
      )}
    </Box>
  );
};

export default UnifiedSwitcherContent;
