import React, { useContext } from 'react';
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
} from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import { StyledSelect } from './SpaceSwitcher';
import { iconMedium } from 'css/icons.styles';
import WorkspaceModal from './WorkspaceModal';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import {
  useGetSelectedOrganization,
  useGetSelectedWorkspace,
  useUpdateSelectedWorkspaceMutation,
} from '@/rtk-query/user';

export const HoverMenuItem = styled(MenuItem)(() => ({
  display: 'flex',
  alignItems: 'center',
  gap: '1rem',
  '& .workspace-icon': {
    display: 'flex',
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

function WorkspaceSwitcher({ open }) {
  const { selectedOrganization } = useGetSelectedOrganization();
  const {
    selectedWorkspace,
    allWorkspaces,
    error: workspaceError,
    isLoading: isLoadingWorkspaces,
  } = useGetSelectedWorkspace();
  const isSmallScreen = useMediaQuery('(max-width:400px)');

  const [updateSelectedWorkspace, { isLoading: isUpdatingSelectedWorkspace }] =
    useUpdateSelectedWorkspaceMutation();

  const {
    open: workspaceModal,
    setSelectedWorkspace,
    openModal: openWorkspaceModal,
    closeModal: closeWorkspaceModal,
  } = useContext(WorkspaceModalContext);

  // useEffect(() => {
  //   if (selectedWorkspace?.id) {
  //     setSelectedWorkspace(selectedWorkspace);
  //   }
  // }, [selectedWorkspace, setSelectedWorkspace]);

  const handleChangeWorkspace = (e) => {
    const newId = e.target.value;
    setSelectedWorkspace(allWorkspaces.find((w) => w.id === newId));
    updateSelectedWorkspace(selectedOrganization.id, newId);
    openWorkspaceModal(true);
  };

  if (workspaceError) {
    return <div>Error: {workspaceError.message}</div>;
  }

  if (isLoadingWorkspaces || isUpdatingSelectedWorkspace) {
    return <CircularProgress height="1.5rem" width="1.5rem" />;
  }

  return (
    <NoSsr>
      {!isLoadingWorkspaces && allWorkspaces?.length > 0 && (
        <Grid2
          sx={{
            width: isSmallScreen ? '80%' : open ? 'auto' : 0,
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
                  <Grid2 container spacing={1} alignItems="flex-end" size="grow">
                    <Grid2 size={{ xs: 12 }} data-cy="mesh-adapter-url">
                      <StyledSelect
                        size="small"
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
                        renderValue={() => {
                          return <span>{selectedWorkspace?.name || ''}</span>;
                        }}
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
                      </StyledSelect>
                    </Grid2>
                  </Grid2>
                }
              />
            </FormGroup>
          </FormControl>
        </Grid2>
      )}
      <WorkspaceModal workspaceModal={workspaceModal} closeWorkspaceModal={closeWorkspaceModal} />
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
