import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import {
  Tab,
  Tabs,
  Typography,
  Grid,
  FormGroup,
  FormControlLabel,
  Switch,
  IconButton,
  CardContent,
  CardHeader,
  CustomTooltip,
} from '@layer5/sistent';
import CopyIcon from '../../assets/icons/CopyIcon';
import _ from 'lodash';
import NoSsr from '@mui/material/NoSsr';
import {
  StatsWrapper,
  ProviderCard,
  RootContainer,
  BoxWrapper,
  Divider,
  GridCapabilityHeader,
  GridExtensionHeader,
  GridExtensionItem,
  StyledPaper,
  TabLabel,
  HideScrollbar,
  IconStyled,
  FormLegend,
  FormContainerWrapper,
  FormGroupWrapper,
} from './style';
import { updateUser, updateProgress, toggleCatalogContent } from '../../lib/store';
import SettingsRemoteIcon from '@mui/icons-material/SettingsRemote';
import SettingsCellIcon from '@mui/icons-material/SettingsCell';
import ExtensionSandbox from '../ExtensionSandbox';
import RemoteComponent from '../RemoteComponent';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTachometerAlt } from '@fortawesome/free-solid-svg-icons';
import MesherySettingsPerformanceComponent from '../MesherySettingsPerformanceComponent';
import { iconMedium } from '../../css/icons.styles';
import { EVENT_TYPES } from '../../lib/event-types';
import { useNotification } from '../../utils/hooks/useNotification';
import { useWindowDimensions } from '@/utils/dimension';
import {
  useGetProviderCapabilitiesQuery,
  useGetUserPrefQuery,
  useUpdateUserPrefMutation,
  useUpdateUserPrefWithContextMutation,
} from '@/rtk-query/user';
import { ThemeTogglerCore } from '@/themes/hooks';

import { SecondaryTab, SecondaryTabs } from '../DashboardComponent/style';

const ThemeToggler = ({ handleUpdateUserPref }) => {
  const Component = ({ mode, toggleTheme }) => {
    return (
      <div>
        <Switch
          color="primary"
          checked={mode === 'dark'}
          onChange={() => {
            toggleTheme();
            handleUpdateUserPref(mode === 'dark' ? 'light' : 'dark');
          }}
        />
        Dark Mode
      </div>
    );
  };

  return <ThemeTogglerCore Component={Component}></ThemeTogglerCore>;
};

const UserPreference = (props) => {
  const [anonymousStats, setAnonymousStats] = useState(props.anonymousStats);
  const [perfResultStats, setPerfResultStats] = useState(props.perfResultStats);
  const [tabVal, setTabVal] = useState(0);
  const [userPrefs, setUserPrefs] = useState(ExtensionPointSchemaValidator('user_prefs')());
  const [providerType, setProviderType] = useState('');
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [capabilitiesLoaded, setCapabilitiesLoaded] = useState(false);
  const { width } = useWindowDimensions();
  const [value, setValue] = useState(0);
  const [providerInfo, setProviderInfo] = useState({});

  const {
    data: userData,
    isSuccess: isUserDataFetched,
    isError: isUserDataError,
    error: userDataError,
  } = useGetUserPrefQuery();

  const { data: capabilitiesData, isSuccess: isCapabilitiesDataFetched } =
    useGetProviderCapabilitiesQuery();

  const [updateUserPref] = useUpdateUserPrefMutation();
  const [updateUserPrefWithContext] = useUpdateUserPrefWithContextMutation();

  const { notify } = useNotification();

  const handleValChange = (event, newVal) => {
    setValue(newVal);
  };

  const handleCatalogContentToggle = () => {
    props.toggleCatalogContent({ catalogVisibility: !catalogContent });

    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  };

  const handleCatalogPreference = (catalogContent) => {
    let body = Object.assign({}, extensionPreferences);
    body['catalogContent'] = catalogContent;
    updateUserPref({ usersExtensionPreferences: body })
      .unwrap()
      .then(() => {
        notify({
          message: `Catalog Content was ${catalogContent ? 'enab' : 'disab'}led`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch(() => {
        handleError('There was an error sending your preference');
      });
  };

  const handleToggle = (name) => () => {
    if (name === 'anonymousUsageStats') {
      setAnonymousStats(!anonymousStats);
      handleChange(name, !anonymousStats);
    } else {
      setPerfResultStats(!perfResultStats);
      handleChange(name, !perfResultStats);
    }
  };

  const handleError = (name) => () => {
    props.updateProgress({ showProgress: false });

    notify({ message: name, event_type: EVENT_TYPES.ERROR });
  };

  const handleChange = (name, resultState) => {
    let val = resultState,
      msg;
    if (name === 'anonymousUsageStats') {
      msg = val
        ? 'Sending anonymous usage statistics was enabled'
        : 'Sending anonymous usage statistics was disabled';
    } else {
      msg = val
        ? 'Sending anonymous performance results was enabled'
        : 'Sending anonymous performance results was disabled';
    }

    const requestBody = JSON.stringify({
      anonymousUsageStats: name === 'anonymousUsageStats' ? val : anonymousStats,
      anonymousPerfResults: name === 'anonymousPerfResults' ? val : perfResultStats,
    });

    props.updateProgress({ showProgress: true });
    updateUserPrefWithContext({ body: requestBody })
      .unwrap()
      .then((result) => {
        props.updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: msg, event_type: val ? EVENT_TYPES.SUCCESS : EVENT_TYPES.INFO });
        }
      })
      .catch(() => {
        handleError('There was an error sending your preference');
      });
  };

  const handleTabValChange = (event, newVal) => {
    setTabVal(newVal);
  };

  useEffect(() => {
    if (props.capabilitiesRegistry && !capabilitiesLoaded) {
      setCapabilitiesLoaded(true); // to prevent re-compute
      setUserPrefs(
        ExtensionPointSchemaValidator('user_prefs')(
          props.capabilitiesRegistry?.extensions?.user_prefs,
        ),
      );
      setProviderType(props.capabilitiesRegistry?.provider_type);
    }
  }, [props.capabilitiesRegistry]);

  useEffect(() => {
    if (isUserDataFetched && userData) {
      setExtensionPreferences(userData?.usersExtensionPreferences);
      setCatalogContent(userData?.usersExtensionPreferences?.catalogContent);
    } else if (isUserDataError) {
      console.error(userDataError);
    }
  }, [isUserDataFetched, userData]);

  useEffect(() => {
    if (isCapabilitiesDataFetched && capabilitiesData) {
      setProviderInfo(capabilitiesData);
    }
  }, [isCapabilitiesDataFetched, capabilitiesData]);

  function convertToTitleCase(str) {
    const words = str.split('_');
    for (let i = 0; i < words.length; i++) {
      words[i] = words[i].charAt(0).toUpperCase() + words[i].slice(1);
    }
    return words.join(' ');
  }

  const RemoteProviderInfoTab = () => {
    const [copied, setCopied] = useState(false);
    const copyToClipboard = (text) => {
      navigator.clipboard
        .writeText(text)
        .then(() => {
          setCopied(true);

          setTimeout(() => {
            setCopied(false);
          }, 2000);
        })
        .catch((error) => {
          console.error('error copying to clipboard:', error);
        });
    };

    return (
      <NoSsr>
        <RootContainer>
          <Typography variant="h5">Provider Information</Typography>
          <Grid container spacing={2}>
            {providerInfo &&
              Object.entries(providerInfo).map(
                ([providerName, provider], index) =>
                  (index < 2 || index === 3) && (
                    <Grid key={index} item md={4} xs={12}>
                      <ProviderCard>
                        <CardHeader
                          title={
                            <Typography
                              variant="h6"
                              style={{
                                textDecoration: 'underline',
                                textDecorationColor: 'rgba(116,147,161,0.5)',
                                textUnderlineOffset: 10,
                              }}
                            >
                              {convertToTitleCase(providerName)}
                            </Typography>
                          }
                        />
                        <CardContent>
                          {' '}
                          <BoxWrapper>
                            <Typography
                              variant="body1"
                              component={HideScrollbar}
                              style={{ marginRight: '20px' }}
                            >
                              {provider}
                            </Typography>
                          </BoxWrapper>
                        </CardContent>
                      </ProviderCard>
                    </Grid>
                  ),
              )}
          </Grid>
          {providerInfo &&
            Object.entries(providerInfo).map(
              ([providerName, provider], index) =>
                (index === 2 || index === 5) && (
                  <ProviderCard key={index} sx={{ margin: '20px' }}>
                    <CardHeader
                      title={
                        <Typography
                          variant="h6"
                          style={{
                            textDecoration: 'underline',
                            textDecorationColor: 'rgba(116,147,161,0.5)',
                            textUnderlineOffset: 10,
                          }}
                        >
                          {convertToTitleCase(providerName)}
                        </Typography>
                      }
                    />
                    <CardContent>
                      {' '}
                      <BoxWrapper>
                        <Typography
                          variant="body1"
                          component={HideScrollbar}
                          style={{ marginRight: '20px' }}
                        >
                          {provider}
                        </Typography>

                        <CustomTooltip title={copied ? 'Copied!' : 'Copy'} placement="top">
                          <IconButton
                            onClick={() => copyToClipboard(provider)}
                            style={{ padding: '0.25rem', float: 'right' }}
                          >
                            <CopyIcon />
                          </IconButton>
                        </CustomTooltip>
                      </BoxWrapper>
                    </CardContent>
                  </ProviderCard>
                ),
            )}

          <ProviderCard>
            <CardHeader
              title={
                <Typography
                  variant="h6"
                  style={{
                    textDecoration: 'underline',
                    textDecorationColor: 'rgba(116,147,161,0.5)',
                    textUnderlineOffset: 10,
                  }}
                >
                  Description
                </Typography>
              }
            />
            <CardContent>
              <Typography>
                <ul>
                  {providerInfo.provider_description &&
                    providerInfo.provider_description.map((desc, index) => (
                      <li key={index}>
                        <Typography>{desc}</Typography>
                      </li>
                    ))}
                </ul>
              </Typography>
            </CardContent>
          </ProviderCard>

          <Divider />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Capabilities
          </Typography>

          <Grid container spacing={2} style={{ display: 'flex', justifyContent: 'space-between' }}>
            <GridCapabilityHeader
              item
              xs={6}
              style={{ borderRadius: '10px 0 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Feature
              </Typography>
            </GridCapabilityHeader>
            <GridCapabilityHeader
              item
              xs={6}
              style={{ borderRadius: '0 10px 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Endpoint
              </Typography>
            </GridCapabilityHeader>
            {providerInfo.capabilities &&
              providerInfo.capabilities.map((capability, index) => (
                <>
                  <GridCapabilityHeader
                    item
                    key={`${index}-${capability.feature}`}
                    xs={6}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        userData?.remoteProviderPreferences?.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                            ? '#E7EFF3'
                            : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.feature}</Typography>
                  </GridCapabilityHeader>
                  <GridCapabilityHeader
                    item
                    key={`${index}-${capability.endpoint}`}
                    xs={6}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        userData?.remoteProviderPreferences?.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                            ? '#E7EFF3'
                            : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.endpoint}</Typography>
                  </GridCapabilityHeader>
                </>
              ))}
          </Grid>
          <Divider />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Extensions
          </Typography>
          {providerInfo.extensions &&
            Object.entries(providerInfo.extensions).map(([extensionName, extension], index) => (
              <div key={index} margin="20px 0px">
                <Typography variant="h6"> {convertToTitleCase(extensionName)}</Typography>
                <Grid container spacing={2} style={{ margin: '10px 0 20px 0' }}>
                  <GridExtensionHeader
                    item
                    xs={6}
                    style={{
                      borderRadius: '10px 0 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Component
                    </Typography>
                  </GridExtensionHeader>
                  <GridExtensionHeader
                    item
                    xs={6}
                    style={{
                      borderRadius: '0 10px 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Type
                    </Typography>
                  </GridExtensionHeader>

                  <GridExtensionItem
                    item
                    xs={6}
                    style={{
                      borderRadius: '0 0 0 10px',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{extension[0].component}</Typography>
                  </GridExtensionItem>
                  <GridExtensionItem
                    item
                    xs={6}
                    style={{
                      borderRadius: '0 0 10px 0',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{convertToTitleCase(extension[0].type)}</Typography>
                  </GridExtensionItem>
                </Grid>
              </div>
            ))}
        </RootContainer>
      </NoSsr>
    );
  };

  const handleUpdateUserPref = (key, value) => {
    const updates = _.set(_.cloneDeep(userData), key, value);
    updateUserPrefWithContext(updates);
  };
  return (
    <NoSsr>
      <StyledPaper>
        <Tabs
          value={tabVal}
          onChange={handleTabValChange}
          variant={width < 600 ? 'scrollable' : 'fullWidth'}
          scrollButtons="on"
          allowScrollButtonsMobile={true}
          indicatorColor="primary"
          textColor="primary"
          centered
        >
          <CustomTooltip title="General preferences" placement="top">
            <Tab icon={<IconStyled as={SettingsCellIcon} />} label={<TabLabel>General</TabLabel>} />
          </CustomTooltip>
          <CustomTooltip title="Choose Performance Test Defaults" placement="top">
            <Tab
              icon={<FontAwesomeIcon icon={faTachometerAlt} style={iconMedium} />}
              label={<TabLabel>Performance</TabLabel>}
            />
          </CustomTooltip>
          {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
          {userPrefs && providerType != 'local' && (
            <CustomTooltip title="Remote Provider preferences" placement="top">
              <Tab
                icon={<SettingsRemoteIcon style={iconMedium} />}
                label={<TabLabel>Remote Provider</TabLabel>}
              />
            </CustomTooltip>
          )}
        </Tabs>
      </StyledPaper>
      <StatsWrapper>
        {tabVal === 0 && (
          <>
            <FormContainerWrapper>
              <FormGroupWrapper component="fieldset">
                <FormLegend component="legend">Extensions</FormLegend>
                <FormGroup>
                  <FormControlLabel
                    key="CatalogContentPreference"
                    control={
                      <Switch
                        checked={catalogContent}
                        onChange={handleCatalogContentToggle}
                        color="primary"
                        data-cy="CatalogContentPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Meshery Catalog Content"
                  />
                </FormGroup>
              </FormGroupWrapper>
            </FormContainerWrapper>
            <FormContainerWrapper>
              <FormGroupWrapper component="fieldset">
                <FormLegend component="legend">Analytics and Improvement Program</FormLegend>
                <FormGroup>
                  <FormControlLabel
                    key="UsageStatsPreference"
                    control={
                      <Switch
                        checked={anonymousStats}
                        onChange={handleToggle('anonymousUsageStats')}
                        color="primary"
                        data-cy="UsageStatsPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Send Anonymous Usage Statistics"
                  />
                  <FormControlLabel
                    key="PerfResultPreference"
                    control={
                      <Switch
                        checked={perfResultStats}
                        onChange={handleToggle('anonymousPerfResults')}
                        color="primary"
                        data-cy="PerfResultPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Send Anonymous Performance Results"
                  />
                </FormGroup>
              </FormGroupWrapper>
            </FormContainerWrapper>
            <FormContainerWrapper>
              <FormGroupWrapper component="fieldset">
                <FormLegend component="legend">Theme</FormLegend>

                <FormGroup>
                  <ThemeToggler
                    handleUpdateUserPref={handleUpdateUserPref}
                    classes={props.classes}
                  />
                </FormGroup>
              </FormGroupWrapper>
            </FormContainerWrapper>
          </>
        )}
        {tabVal === 1 && <MesherySettingsPerformanceComponent />}
        {tabVal === 2 && userPrefs && providerType !== 'local' && (
          <>
            <SecondaryTabs
              value={value}
              onChange={handleValChange}
              variant={width < 600 ? 'scrollable' : 'fullWidth'}
              scrollButtons="on"
              allowScrollButtonsMobile={true}
              indicatorColor="primary"
              textColor="primary"
              centered
            >
              <CustomTooltip title="Details" placement="top">
                <SecondaryTab label={<TabLabel>Details</TabLabel>} />
              </CustomTooltip>
              <CustomTooltip title="Plugins" placement="top">
                <SecondaryTab label={<TabLabel>Plugins</TabLabel>} />
              </CustomTooltip>
            </SecondaryTabs>
            <StatsWrapper>
              {value === 0 && <RemoteProviderInfoTab />}

              {value === 1 && (
                <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
              )}
            </StatsWrapper>
          </>
        )}
      </StatsWrapper>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateUser: bindActionCreators(updateUser, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
});

const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const catalogVisibility = state.get('catalogVisibility');
  const capabilitiesRegistry = state.get('capabilitiesRegistry');
  return {
    selectedK8sContexts,
    catalogVisibility,
    capabilitiesRegistry,
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(UserPreference));
