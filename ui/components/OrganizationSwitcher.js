import { useGetOrgsQuery } from '@/rtk-query/organization';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import React, { useEffect, useState } from 'react';
import {
  Button,
  FormControl,
  FormControlLabel,
  FormGroup,
  Grid,
  MenuItem,
  NoSsr,
  Select,
  withStyles,
} from '@material-ui/core';
import theme from '@/themes/app';
import { setKeys, setOrganization, setWorkspace } from '../lib/store';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from 'redux';
import styles from './UserPreferences/style';
import { store } from '../store';
import { withRouter } from 'next/router';
import OrgOutlinedIcon from '@/assets/icons/OrgOutlinedIcon';
import { iconXLarge } from 'css/icons.styles';
import WorkspaceOutlinedIcon from '@/assets/icons/WorkspaceOutlinedIcon';
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';

function OrgMenu(props) {
  const {
    data: orgsResponse,
    isSuccess: isOrgsSuccess,
    isError: isOrgsError,
    error: orgsError,
  } = useGetOrgsQuery({});
  let orgs = orgsResponse?.organizations || [];
  const { organization, setOrganization, open } = props;

  const { notify } = useNotification();
  useEffect(() => {
    if (isOrgsError) {
      notify({
        message: `There was an error fetching available data ${orgsError?.data}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [isOrgsError, notify, orgsError]);

  const handleOrgSelect = (e) => {
    const id = e.target.value;
    const selected = orgs.find((org) => org.id === id);
    setOrganization({ organization: selected });
  };

  return (
    <NoSsr>
      {isOrgsSuccess && orgs && (
        <div
          style={{ width: open ? 'auto' : 0, overflow: 'hidden', transition: 'width 0.7s ease' }}
        >
          <FormControl component="fieldset">
            <FormGroup>
              <FormControlLabel
                key="SpacesPreferences"
                control={
                  <Grid container spacing={1} alignItems="flex-end">
                    <Grid item xs={12} data-cy="mesh-adapter-url">
                      <Select
                        value={organization.id}
                        onChange={handleOrgSelect}
                        SelectDisplayProps={{ style: { display: 'flex', flexDirection: 'row' } }}
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
                        {orgs?.map((org) => (
                          <MenuItem key={org.id} value={org.id}>
                            <div>
                              <OrgOutlinedIcon
                                width="24"
                                height="24"
                                secondaryFill={theme.palette.darkSlateGray}
                              />
                            </div>
                            <span>{org.name}</span>
                          </MenuItem>
                        ))}
                      </Select>
                    </Grid>
                  </Grid>
                }
              />
            </FormGroup>
          </FormControl>
        </div>
      )}
    </NoSsr>
  );
}

function WorkspaceSwitcher({ organization, open, workspace, setWorkspace }) {
  const [orgId, setOrgId] = useState('');
  const { data: workspacesData, isError: isWorkspacesError } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 10,
      search: '',
      order: '',
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );

  const handleWorkspaceSelect = (e) => {
    const id = e.target.value;
    const selected = workspacesData.workspaces.find((org) => org.id === id);
    setWorkspace({ workspace: selected });
  };

  useEffect(() => {
    setOrgId(organization?.id);
  }, [organization]);

  if (!organization || !workspace) {
    return null;
  }

  return (
    <NoSsr>
      {!isWorkspacesError && workspace && (
        <div style={{ display: open ? 'block' : 'none', transition: 'all 1s' }}>
          <FormControl component="fieldset">
            <FormGroup>
              <FormControlLabel
                key="SpacesPreferences"
                control={
                  <Grid container spacing={1} alignItems="flex-end">
                    <Grid item xs={12} data-cy="mesh-adapter-url">
                      <Select
                        value={workspace.id}
                        onChange={handleWorkspaceSelect}
                        SelectDisplayProps={{ style: { display: 'flex', flexDirection: 'row' } }}
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
                        {workspacesData?.workspaces?.map((org) => (
                          <MenuItem key={org.id} value={org.id}>
                            <div>
                              <WorkspaceOutlinedIcon
                                width="24"
                                height="24"
                                secondaryFill={theme.palette.darkSlateGray}
                              />
                            </div>
                            <span>{org.name}</span>
                          </MenuItem>
                        ))}
                      </Select>
                    </Grid>
                  </Grid>
                }
              />
            </FormGroup>
          </FormControl>
        </div>
      )}
    </NoSsr>
  );
}

function OrganizationSwitcher(props) {
  const [orgOpen, setOrgOpen] = useState(false);
  const [workspaceOpen, setWorkspaceOpen] = useState(false);

  return (
    <NoSsr>
      <Provider store={store}>
        <Button
          onClick={() => setOrgOpen(!orgOpen)}
          style={{ marginRight: orgOpen ? '1rem' : '0' }}
        >
          <OrgOutlinedIcon {...iconXLarge} />
        </Button>
        <OrgMenu {...props} open={orgOpen} />
        <Button
          onClick={() => setWorkspaceOpen(!workspaceOpen)}
          style={{ marginRight: workspaceOpen ? '1rem' : '0' }}
        >
          <WorkspaceOutlinedIcon {...iconXLarge} />
        </Button>
        <WorkspaceSwitcher {...props} open={workspaceOpen} />
      </Provider>
    </NoSsr>
  );
}

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  const workspace = state.get('workspace');
  return {
    organization,
    workspace,
  };
};

const mapDispatchToProps = (dispatch) => ({
  setOrganization: bindActionCreators(setOrganization, dispatch),
  setWorkspace: bindActionCreators(setWorkspace, dispatch),
  setKeys: bindActionCreators(setKeys, dispatch),
});

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(OrganizationSwitcher)),
);
