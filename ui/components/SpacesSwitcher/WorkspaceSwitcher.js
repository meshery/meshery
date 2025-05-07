import React, { useEffect, useState } from 'react';
import { FormControl, FormControlLabel, FormGroup, Grid, styled, MenuItem } from '@layer5/sistent';
import SettingsIcon from '@mui/icons-material/Settings';
import { NoSsr } from '@layer5/sistent';
import { StyledSelect } from './SpaceSwitcher';
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { iconMedium } from 'css/icons.styles';
import WorkspaceModal from './WorkspaceModal';
import { useCurrentOrganization } from '@/utils/hooks/useCurrentOrganization';

export const HoverMenuItem = styled(MenuItem)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: '1rem',
  '& .setting-icon': {
    visibility: 'hidden',
    display: 'flex',
  },
  '&:hover .setting-icon': {
    visibility: 'visible',
  },
}));

const SettingsIconWrapper = styled('div')(({ theme }) => ({
  '& svg': {
    color: theme.palette.icon.secondary,
  },
  '& svg:hover': {
    fill: theme.palette.icon.secondary + ' !important',
  },
}));

export const WorkspaceSwitcherContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => {},
});

function WorkspaceSwitcher({ open }) {
  const [_defaultWorkspace, setDefaultWorkspace] = useState(null);
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  const orgId = useCurrentOrganization()?.id;
  const { data: workspacesData, isError: isWorkspacesError } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      order: 'updated_at desc',
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );

  useEffect(() => {
    if (workspacesData && workspacesData.workspaces?.length > 0) {
      const defaultWorkspace = workspacesData.workspaces[0];
      setDefaultWorkspace(defaultWorkspace);
    }
  }, [workspacesData]);

  const handleChangeWorkspace = (e) => {
    if (!workspacesData || !workspacesData.workspaces) return;

    const selectedWorkspace = workspacesData.workspaces.find(
      (workspace) => workspace.id === e.target.value,
    );
    setDefaultWorkspace(selectedWorkspace);
    setSelectedWorkspace({ id: selectedWorkspace.id, name: selectedWorkspace.name });
    setWorkspaceModal(true);
  };

  return (
    <NoSsr>
      <WorkspaceSwitcherContext.Provider
        value={{
          open: workspaceModal,
          openModal: () => setWorkspaceModal(true),
          closeModal: () => {
            setWorkspaceModal(false);
          },
          selectedWorkspace,
          setSelectedWorkspace,
        }}
      >
        {!isWorkspacesError && workspacesData && workspacesData.workspaces && (
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
                    <Grid container spacing={1} alignItems="flex-end">
                      <Grid item xs={12} data-cy="mesh-adapter-url">
                        <StyledSelect
                          size="small"
                          value={_defaultWorkspace?.id || ''}
                          onChange={(e) => {
                            if (e.target.value !== _defaultWorkspace?.id) {
                              handleChangeWorkspace(e); // only call for new selection
                            }
                          }}
                          renderValue={(selected) => {
                            const workspace = workspacesData?.workspaces?.find(
                              (w) => w.id === selected,
                            );
                            return workspace ? <span>{workspace.name}</span> : '';
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
                          {workspacesData?.workspaces?.map((works) => (
                            <HoverMenuItem
                              key={works.id}
                              value={works.id}
                              onClick={() => {
                                if (works.id === _defaultWorkspace?.id) {
                                  handleChangeWorkspace({ target: { value: works.id } });
                                }
                              }}
                            >
                              <span>{works.name}</span>
                              <SettingsIconWrapper className="setting-icon">
                                <SettingsIcon {...iconMedium} />
                              </SettingsIconWrapper>
                            </HoverMenuItem>
                          ))}
                        </StyledSelect>
                      </Grid>
                    </Grid>
                  }
                />
              </FormGroup>
            </FormControl>
          </div>
        )}
        <WorkspaceModal workspaceModal={workspaceModal} setWorkspaceModal={setWorkspaceModal} />
      </WorkspaceSwitcherContext.Provider>
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
