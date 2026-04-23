import React, { useContext } from 'react';
import { useGetOrgsQuery } from '@/rtk-query/organization';
import {
  useTheme,
  useMediaQuery,
  CircularProgress,
  FormControl,
  FormGroup,
  FormControlLabel,
  Grid2,
} from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import {
  useGetSelectedOrganization,
  useUpdateSelectedOrganizationMutation,
} from '@/rtk-query/user';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import OrgOutlinedIcon from '@/assets/icons/OrgOutlinedIcon';
import { StyledMenuItem, StyledSelect } from './styles';
import _ from 'lodash';

export function OrgMenu(props) {
  const { data: orgsResponse, isSuccess: isOrgsSuccess } = useGetOrgsQuery({});
  let orgs = orgsResponse?.organizations || [];
  let uniqueOrgs = _.uniqBy(orgs, 'id');

  const theme = useTheme();
  const isSmallScreen = useMediaQuery('(max-width:400px)');
  const { selectedOrganization: selectedOrgFromPref } = useGetSelectedOrganization();
  const { currentLoadedResource } = useContext(WorkspaceModalContext);

  const selectedOrganization = currentLoadedResource?.org?.id
    ? currentLoadedResource.org
    : selectedOrgFromPref;

  const [updateSelectedOrg, { isLoading: isUpdatingOrg }] = useUpdateSelectedOrganizationMutation();

  if (isUpdatingOrg) {
    return <CircularProgress size={24} style={{ color: theme.palette.background.brand.default }} />;
  }

  if (!selectedOrganization) return null;

  const organization = selectedOrganization;
  const { open, fromMobileView } = props;

  const handleOrgSelect = (e) => {
    const id = e.target.value;
    updateSelectedOrg({ orgId: id });
  };

  return (
    <NoSsr>
      {isOrgsSuccess && orgs && open && (
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
                key="OrgPreferences"
                control={
                  <Grid2 container spacing={1} alignItems="flex-end" size="grow">
                    <Grid2 data-cy="mesh-adapter-url" size={{ xs: 12 }}>
                      <StyledSelect
                        value={organization?.id ?? ''}
                        onChange={handleOrgSelect}
                        SelectDisplayProps={{
                          style: {
                            display: 'flex',
                            flexDirection: 'row',
                            fill: theme.palette.background.constant.white,
                            paddingBlock: '9px 8px',
                            paddingInline: '18px 34px',
                            color: fromMobileView
                              ? theme.palette.text.default
                              : theme.palette.background.constant.white,
                          },
                        }}
                        renderValue={() => {
                          return <span>{selectedOrganization?.name || 'Private Org'}</span>;
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
                          style: {
                            fill: theme.palette.text.secondary,
                          },
                        }}
                      >
                        {uniqueOrgs?.map((org) => (
                          <StyledMenuItem key={org.id} value={org.id}>
                            <OrgOutlinedIcon
                              width="24"
                              height="24"
                              className="OrgClass"
                              style={{ marginRight: '1rem', color: theme.palette.icon.default }}
                            />
                            <span>{org.name || 'Private Org'}</span>
                          </StyledMenuItem>
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
    </NoSsr>
  );
}
