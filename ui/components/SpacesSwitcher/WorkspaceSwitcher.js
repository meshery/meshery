import React, { useContext, useEffect, useState } from 'react';
import { FormControl, FormControlLabel, FormGroup, Grid, styled, MenuItem } from '@layer5/sistent';
import { NoSsr } from '@layer5/sistent';
import { StyledSelect } from './SpaceSwitcher';
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import WorkspaceModal from './WorkspaceModal';
import { useSelector } from 'react-redux';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';

export const HoverMenuItem = styled(MenuItem)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  gap: '1rem',
}));

function WorkspaceSwitcher({ open }) {
  const [_defaultWorkspace, setDefaultWorkspace] = useState(null);
  const { organization: currentOrg } = useSelector((state) => state.ui);
  const { id: orgId } = currentOrg;
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
  const {
    setSelectedWorkspace,
    open: workspaceModal,
    openModal: openWorkspaceModal,
    closeModal: closeWorkspaceModal,
  } = useContext(WorkspaceModalContext);

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
    openWorkspaceModal(true);
  };

  return (
    <NoSsr>
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
                            selected={works.id === _defaultWorkspace?.id}
                            onClick={() => {
                              if (works.id === _defaultWorkspace?.id) {
                                handleChangeWorkspace({ target: { value: works.id } });
                              }
                            }}
                          >
                            <span>{works.name}</span>
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
      <WorkspaceModal workspaceModal={workspaceModal} closeWorkspaceModal={closeWorkspaceModal} />
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
