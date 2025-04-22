import React, { useEffect, useState } from 'react';
import {
  FormControl,
  FormControlLabel,
  FormGroup,
  Grid,
  styled,
  MenuItem,
  Modal,
  WorkspaceIcon,
  ModalBody,
} from '@layer5/sistent';
import { WorkspacesComponent } from '../../components/Lifecycle';
import SettingsIcon from '@mui/icons-material/Settings';
import { NoSsr } from '@layer5/sistent';
import { useLegacySelector } from '../../lib/store';
import { StyledSelect } from './SpaceSwitcher';

import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { iconMedium } from 'css/icons.styles';
import { useRouter } from 'next/router';

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

function WorkspaceSwitcher({ open }) {
  const [defaultWorkspace, setDefaultWorkspace] = useState(null);
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const orgId = useLegacySelector((state) => state.get('organization'))?.id;
  const { data: workspacesData, isError: isWorkspacesError } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      search: '',
      order: '',
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );
  useEffect(() => {
    if (workspacesData && workspacesData.workspaces.length > 0) {
      const defaultWorkspace = workspacesData.workspaces[0];
      setDefaultWorkspace(defaultWorkspace);
    }
  }, [workspacesData]);

  const router = useRouter();

  const handleChangeWorkspace = (e) => {
    if (!workspacesData || !workspacesData.workspaces) return;

    const selectedWorkspace = workspacesData.workspaces.find(
      (workspace) => workspace.id === e.target.value,
    );
    setDefaultWorkspace(selectedWorkspace);
    setWorkspaceModal(true);
    router.push(
      {
        pathname: router.pathname,
        query: {
          ...router.query,
          view: 'table',
          id: selectedWorkspace.id,
          name: selectedWorkspace.name,
        },
      },
      undefined,
      { shallow: true },
    );
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
                        value={defaultWorkspace?.id || ''}
                        onChange={(e) => {
                          if (e.target.value !== defaultWorkspace?.id) {
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
                              if (works.id === defaultWorkspace?.id) {
                                handleChangeWorkspace({ target: { value: works.id } });
                              }
                            }}
                          >
                            <span>{works.name}</span>
                            <span className="setting-icon">
                              <SettingsIcon {...iconMedium} />
                            </span>
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
      <Modal
        closeModal={() => {
          setWorkspaceModal(false);
          if (router.query.id) {
            router.back();
          }
        }}
        open={workspaceModal}
        maxWidth="xl"
        headerIcon={<WorkspaceIcon {...iconMedium} />}
        title="Workspaces"
      >
        <ModalBody style={{ maxHeight: '80vh', overflowY: 'auto' }}>
          {workspaceModal && <WorkspacesComponent />}
        </ModalBody>
      </Modal>
    </NoSsr>
  );
}

export default WorkspaceSwitcher;
