import React, { useContext, useEffect, useState } from 'react';
import {
  FormControl,
  FormControlLabel,
  FormGroup,
  Grid2,
  styled,
  MenuItem,
  CircularProgress,
  WorkspaceIcon,
  useMediaQuery,
  useTheme,
  Divider,
  Button,
  Box,
} from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import { StyledSelect } from './SpaceSwitcher';
import { iconMedium } from 'css/icons.styles';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { Keys } from '@meshery/schemas/permissions';
import {
  useGetSelectedOrganization,
  useGetSelectedWorkspace,
  useUpdateSelectedWorkspaceMutation,
} from '@/rtk-query/user';
import { BottomSheetInlineSelect } from './BottomSheetInlineSelect';

export const HoverMenuItem = styled(MenuItem)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '1rem',
  '& .workspace-icon': {
    display: 'flex',
  },
  '&.Mui-selected': {
    backgroundColor: theme.palette.action.selected,
    '&:hover': {
      backgroundColor: theme.palette.action.selected + '!important',
    },
  },
}));

const WorkspaceIconWrapper = styled('div')(({ theme }) => ({
  '& svg': {
    color: theme.palette.icon.secondary,
  },
  '& svg:hover': {
    fill: theme.palette.icon.secondary + ' !important',
  },
}));

function WorkspaceSwitcher({ open, fromMobileView, expanded, onExpandedChange, onMobileSelect }) {
  const { selectedOrganization } = useGetSelectedOrganization();
  const {
    selectedWorkspace: selectedWorkspacePref,
    allWorkspaces,
    error: workspaceError,
    isLoading: isLoadingWorkspaces,
  } = useGetSelectedWorkspace();
  const {
    setSelectedWorkspace,
    openModal: openWorkspaceModal,
    setCreateNewWorkspaceModalOpen,
    currentLoadedResource,
  } = useContext(WorkspaceModalContext);
  const isSmallScreen = useMediaQuery('(max-width:400px)');
  const theme = useTheme();
  const selectedWorkspace = currentLoadedResource?.workspace?.id
    ? currentLoadedResource.workspace
    : selectedWorkspacePref;

  const [updateSelectedWorkspace, { isLoading: isUpdatingSelectedWorkspace }] =
    useUpdateSelectedWorkspaceMutation();
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    if (!open) {
      setMenuOpen(false);
    }
  }, [open]);

  const handleChangeWorkspace = (e) => {
    setMenuOpen(false);
    const newId = e.target.value;
    setSelectedWorkspace(allWorkspaces.find((w) => w.id === newId));
    updateSelectedWorkspace(selectedOrganization.id, newId);
    onMobileSelect?.();
    // Mobile UX: selecting a workspace should just switch context and close the sheet.
    // Desktop UX: selecting should open the workspace explorer modal.
    if (!fromMobileView) {
      openWorkspaceModal(true);
    }
  };

  if (workspaceError) {
    return <div>Error: {workspaceError.message}</div>;
  }

  if (isLoadingWorkspaces || isUpdatingSelectedWorkspace) {
    return <CircularProgress height="1.5rem" width="1.5rem" />;
  }

  if (fromMobileView) {
    const workspaceIcon = (
      <WorkspaceIcon
        {...iconMedium}
        fill={theme.palette.icon.default}
        secondaryFill={theme.palette.icon.default}
      />
    );

    return (
      <NoSsr>
        {!isLoadingWorkspaces && allWorkspaces?.length > 0 && open && (
          <BottomSheetInlineSelect
            data-cy="mesh-adapter-url"
            value={selectedWorkspace?.id || ''}
            options={allWorkspaces.map((works) => ({
              id: works.id,
              label: works.name,
            }))}
            onSelect={(id) => handleChangeWorkspace({ target: { value: id } })}
            expanded={expanded}
            onExpandedChange={onExpandedChange}
            renderOption={(option) => (
              <>
                {workspaceIcon}
                <Box
                  component="span"
                  sx={{
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  }}
                >
                  {option.label}
                </Box>
              </>
            )}
          />
        )}
      </NoSsr>
    );
  }

  return (
    <NoSsr>
      {!isLoadingWorkspaces && allWorkspaces?.length > 0 && (
        <Grid2
          sx={{
            width: isSmallScreen ? '100%' : open ? 'auto' : 0,
            minWidth: 0,
            overflow: open ? '' : 'hidden',
            transition: 'all 1s',
          }}
        >
          <FormControl
            sx={{
              width: isSmallScreen ? '100%' : 'auto',
            }}
            component="fieldset"
          >
            <FormGroup>
              <FormControlLabel
                key="SpacesPreferences"
                control={
                  <Grid2 container spacing={1} size="grow" sx={{ alignItems: 'flex-end' }}>
                    <Grid2 size={{ xs: 12 }} data-cy="mesh-adapter-url">
                      <StyledSelect
                        size="small"
                        open={menuOpen}
                        onOpen={() => setMenuOpen(true)}
                        onClose={() => setMenuOpen(false)}
                        value={selectedWorkspace?.id || ''}
                        onChange={(e) => {
                          if (e.target.value !== selectedWorkspace?.id) {
                            handleChangeWorkspace(e); // only call for new selection
                          }
                        }}
                        SelectDisplayProps={{
                          style: {
                            display: 'flex',
                            flexDirection: 'row',
                            fill: '#eee',
                            paddingBlock: '9px 8px',
                            paddingInline: '18px 34px',
                          },
                        }}
                        renderValue={() => (
                          <span>{selectedWorkspace?.name || 'Private Workspace'}</span>
                        )}
                        MenuProps={{
                          anchorOrigin: {
                            vertical: 'bottom',
                            horizontal: 'left',
                          },
                          transformOrigin: {
                            vertical: 'top',
                            horizontal: 'left',
                          },
                          getContentAnchorEl: null,
                          style: { zIndex: theme.zIndex.modal + 200 },
                        }}
                      >
                        {allWorkspaces?.map((works) => (
                          <HoverMenuItem
                            key={works.id}
                            value={works.id}
                            onClick={() => {
                              if (works.id === selectedWorkspace?.id) {
                                handleChangeWorkspace({ target: { value: works.id } });
                              }
                            }}
                          >
                            <WorkspaceIconWrapper className="workspace-icon">
                              <WorkspaceIcon {...iconMedium} />
                            </WorkspaceIconWrapper>
                            <span>{works.name}</span>
                          </HoverMenuItem>
                        ))}
                        <Divider />
                        <Box
                          sx={{
                            gap: 2,
                            px: 2,
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                          }}
                        >
                          <Button
                            variant="contained"
                            onClick={(e) => {
                              e.stopPropagation();
                              setMenuOpen(false);
                              openWorkspaceModal(true);
                            }}
                            permissionKey={Keys.WorkspaceManagementViewWorkspace}
                          >
                            Explore Workspaces
                          </Button>
                          <Button
                            variant="outlined"
                            onClick={(e) => {
                              e.stopPropagation();
                              setMenuOpen(false);
                              setSelectedWorkspace({
                                id: 'All Workspaces',
                                name: 'All Workspaces',
                              });
                              setCreateNewWorkspaceModalOpen(true);
                              openWorkspaceModal(true);
                            }}
                            permissionKey={Keys.WorkspaceManagementCreateWorkspace}
                          >
                            Create Workspace
                          </Button>
                        </Box>
                      </StyledSelect>
                    </Grid2>
                  </Grid2>
                }
              />
            </FormGroup>
          </FormControl>
        </Grid2>
      )}
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
