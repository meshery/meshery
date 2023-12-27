import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { FormGroup, FormControlLabel, Grid } from '@material-ui/core';
import NoSsr from '@material-ui/core/NoSsr';
import { setOrganization } from '../../lib/store';
import { EVENT_TYPES } from '../../lib/event-types';
import { useNotification } from '../../utils/hooks/useNotification';
import { useGetOrgsQuery } from '../../rtk-query/organization';
import OrgIcon from '../../assets/icons/OrgIcon';
import ErrorBoundary from '../ErrorBoundary';
import { Provider } from 'react-redux';
import { store } from '../../store';
import { withStyles } from '@material-ui/core';
import { Select, MenuItem, FormControl, FormLabel } from '@material-ui/core';
import styles from './style';
import theme from '../../themes/app';

const SpacesPreferences = (props) => {
  const {
    data: orgsResponse,
    isSuccess: isOrgsSuccess,
    isError: isOrgsError,
    error: orgsError,
  } = useGetOrgsQuery({});
  let orgs = orgsResponse?.organizations || [];
  const { organization, setOrganization, classes } = props;
  console.log('props', props);

  const { notify } = useNotification();

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
  };

  return (
    <NoSsr>
      <>
        {isOrgsSuccess && orgs && (
          <div className={classes.formContainerWrapper}>
            <FormControl component="fieldset" className={classes.formControlWrapper}>
              <FormLabel component="legend" className={classes.formLabelWrapper}>
                Spaces
              </FormLabel>
              <FormGroup>
                <FormControlLabel
                  key="SpacesPreferences"
                  control={
                    <Grid container spacing={1} alignItems="flex-end">
                      <Grid item xs={12} data-cy="mesh-adapter-url">
                        <Select
                          value={organization.id}
                          onChange={handleOrgSelect}
                          SelectDisplayProps={{ style: { display: 'flex', padding: '10px' } }}
                          className={classes.orgSelect}
                        >
                          {orgs?.map((org) => {
                            return (
                              <MenuItem key={org.id} value={org.id} className={classes.selectItem}>
                                <div className={classes.orgIconWrapper}>
                                  <OrgIcon
                                    width="24"
                                    height="24"
                                    secondaryFill={theme.palette.darkSlateGray}
                                  />
                                </div>
                                <span className={classes.org}>{org.name}</span>
                              </MenuItem>
                            );
                          })}
                        </Select>
                      </Grid>
                    </Grid>
                  }
                />
              </FormGroup>
            </FormControl>
          </div>
        )}
      </>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  setOrganization: bindActionCreators(setOrganization, dispatch),
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
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in Spaces Prefs Component', e)}
      >
        <Provider store={store}>
          <SpacesPreferences {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(SpacesPreferencesWithErrorBoundary)),
);
