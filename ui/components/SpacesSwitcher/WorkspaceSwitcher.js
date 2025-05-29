import React, { useContext } from 'react';
import {
  FormControl,
  FormControlLabel,
  FormGroup,
  Grid2,
  styled,
  MenuItem,
  CircularProgress,
  useTheme,
  OpenFileIcon,
} from '@layer5/sistent';
import { NoSsr } from '@layer5/sistent';
import { StyledSelect } from './SpaceSwitcher';
import WorkspaceModal from './WorkspaceModal';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import {
  useGetSelectedOrganization,
  useGetSelectedWorkspace,
  useUpdateSelectedWorkspaceMutation,
} from '@/rtk-query/user';

export const HoverMenuItem = styled(MenuItem)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: '1rem',
}));

function WorkspaceSwitcher({ open }) {
  const theme = useTheme();
  const { selectedOrganization } = useGetSelectedOrganization();
  const {
    selectedWorkspace,
    allWorkspaces,
    error: workspaceError,
    isLoading: isLoadingWorkspaces,
  } = useGetSelectedWorkspace();

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
    // openWorkspaceModal(true);
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
        <div
          style={{
            width: open ? 'auto' : 0,
            overflow: open ? '' : 'hidden',
            transition: 'all 1s',
          }}
        >
          <FormControl component="fieldset">
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
                        renderValue={() => {
                          return selectedWorkspace?.name || '';
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
                            selected={works.id === selectedWorkspace?.id}
                          >
                            <span>{works.name}</span>
                            <span
                              onClick={(e) => {
                                if (works.id === selectedWorkspace?.id) {
                                  e.stopPropagation();
                                }
                                openWorkspaceModal(true);
                              }}
                              style={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}
                            >
                              <OpenFileIcon
                                fill={
                                  works.id === selectedWorkspace?.id
                                    ? 'currentColor'
                                    : theme.palette.grey[600]
                                }
                                style={{
                                  transition: 'none',
                                  fill:
                                    works.id === selectedWorkspace?.id
                                      ? 'currentColor'
                                      : theme.palette.grey[600],
                                }}
                              />
                            </span>
                          </HoverMenuItem>
                        ))}
                      </StyledSelect>
                    </Grid2>
                  </Grid2>
                }
              />
            </FormGroup>
          </FormControl>
        </div>
      )}
      <WorkspaceModal workspaceModal={workspaceModal} closeWorkspaceModal={closeWorkspaceModal} />
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
