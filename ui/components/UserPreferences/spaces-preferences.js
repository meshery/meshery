import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { FormGroup, FormControlLabel, Grid, FormLabel } from '@layer5/sistent';
import NoSsr from '@material-ui/core/NoSsr';
import { setOrganization, setKeys } from '../../lib/store';
import { EVENT_TYPES } from '../../lib/event-types';
import { useNotification } from '../../utils/hooks/useNotification';
import { useGetOrgsQuery } from '../../rtk-query/organization';
import OrgIcon from '../../assets/icons/OrgIcon';
import { ErrorBoundary as SistentErrorBoundary } from '@layer5/sistent';
import { Provider } from 'react-redux';
import { store } from '../../store';
import theme from '../../themes/app';
import { useGetCurrentAbilities } from '../../rtk-query/ability';
import {
  StyledFormContainer,
  StyledFormControl,
  StyledOrgSelect,
  StyledSelectItem,
  StyledOrgIconWrapper,
  StyledOrgText,
} from './style';
import { UsesSistent } from '@/components/SistentWrapper';

const SpacesPreferences = (props) => {
  const {
    data: orgsResponse,
    isSuccess: isOrgsSuccess,
    isError: isOrgsError,
    error: orgsError,
  } = useGetOrgsQuery({});
  let orgs = orgsResponse?.organizations || [];
  const { organization, setOrganization } = props;
  const [skip, setSkip] = React.useState(true);

  const { notify } = useNotification();

  useGetCurrentAbilities(organization, props.setKeys, skip);

  useEffect(() => {
    if (isOrgsError) {
      notify({
        message: `There was an error fetching available data ${orgsError?.data}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [orgsError]);

  const handleOrgSelect = (e) => {
    const id = e.target.value;
    const selected = orgs.find((org) => org.id === id);
    setOrganization({ organization: selected });
    setSkip(false);
  };

  return (
    <NoSsr>
      <>
        {isOrgsSuccess && orgs && (
          <UsesSistent>
            <StyledFormContainer>
              <StyledFormControl>
                <FormLabel component="legend" sx={{ fontSize: 20 }}>
                  Spaces
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="SpacesPreferences"
                    control={
                      <Grid container spacing={1} alignItems="flex-end">
                        <Grid item xs={12} data-cy="mesh-adapter-url">
                          <StyledOrgSelect
                            value={organization.id}
                            onChange={handleOrgSelect}
                            SelectDisplayProps={{ style: { display: 'flex', padding: '10px' } }}
                          >
                            {orgs?.map((org) => {
                              return (
                                <StyledSelectItem key={org.id} value={org.id}>
                                  <StyledOrgIconWrapper>
                                    <OrgIcon
                                      width="24"
                                      height="24"
                                      secondaryFill={theme.palette.darkSlateGray}
                                    />
                                  </StyledOrgIconWrapper>
                                  <StyledOrgText>{org.name}</StyledOrgText>
                                </StyledSelectItem>
                              );
                            })}
                          </StyledOrgSelect>
                        </Grid>
                      </Grid>
                    }
                  />
                </FormGroup>
              </StyledFormControl>
            </StyledFormContainer>
          </UsesSistent>
        )}
      </>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  setOrganization: bindActionCreators(setOrganization, dispatch),
  setKeys: bindActionCreators(setKeys, dispatch),
});

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  return {
    organization,
  };
};

const SpacesPreferencesWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <SistentErrorBoundary>
        <Provider store={store}>
          <SpacesPreferences {...props} />
        </Provider>
      </SistentErrorBoundary>
    </NoSsr>
  );
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(SpacesPreferencesWithErrorBoundary));
