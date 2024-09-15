//import useState from "react"
import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withStyles } from '@material-ui/core/styles';
import CopyIcon from '../../assets/icons/CopyIcon';
import _ from 'lodash';
import {
  Typography,
  Grid,
  FormControl,
  FormLabel,
  FormGroup,
  FormControlLabel,
  Switch,
  IconButton,
  Card,
  CardContent,
  CardHeader,
  Box,
} from '@material-ui/core';
import { CustomTooltip } from '@layer5/sistent';
import NoSsr from '@material-ui/core/NoSsr';
import { updateUser, updateProgress, toggleCatalogContent } from '../../lib/store';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { Paper } from '@material-ui/core';
import SettingsRemoteIcon from '@material-ui/icons/SettingsRemote';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
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

const styles = (theme) => ({
  statsWrapper: {
    // padding : theme.spacing(2),
    maxWidth: '100%',
    height: 'auto',
    borderTopLeftRadius: 0,
    borderTopRightRadius: 0,
    borderBottomLeftRadius: 3,
    borderBottomRightRadius: 3,
  },
  paperRoot: {
    flexGrow: 1,
    maxWidth: '100%',
    marginLeft: 0,
    borderTopLeftRadius: 3,
    borderTopRightRadius: 3,
  },
  tabs: {
    width: '100%',
    marginLeft: 0,
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  tab: {
    width: '42%',
    // maxWidth: 'min(33%, 200px)',
    // minWidth: '50px',
    margin: 0,
    '&.Mui-selected': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  icon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(1.75),
    marginLeft: theme.spacing(0.5),
  },
  iconText: {
    display: 'inline',
    verticalAlign: 'middle',
  },
  backToPlay: { margin: theme.spacing(2) },
  link: { cursor: 'pointer' },
  formContainer: {
    display: 'flex',
    'flex-wrap': 'wrap',
    'justify-content': 'space-evenly',
    padding: 50,
  },
  formGrp: {
    padding: 20,
    border: '1.5px solid #969696',
    display: 'flex',
    width: '70%',
  },
  formLegend: { fontSize: 20 },
  formLegendSmall: { fontSize: 16 },
  switchBase: {
    color: '#647881',
    '&$checked': { color: '#00b39f' },
    '&$checked + $track': { backgroundColor: 'rgba(0,179,159,0.5)' },
  },
  track: { backgroundColor: 'rgba(100,120,129,0.5)' },
  checked: {},
  tabLabel: {
    [theme.breakpoints.up('sm')]: {
      fontSize: '1em',
    },
    [theme.breakpoints.between('xs', 'sm')]: {
      fontSize: '0.8em',
    },
  },
  hideScrollbar: {
    overflowX: 'auto',
    '&::-webkit-scrollbar': {
      display: 'none',
    },
    '&::-moz-scrollbar': {
      display: 'none',
    },
  },
  card: {
    border: '1px solid rgba(0,179,159,0.3)',
    margin: '20px 0px',
    backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#C9DBE3',
    // display: 'flex',
    // flexWrap: 'wrap',
  },
  box: {
    width: '100%',
    display: 'flex',
    justifyContent: 'space-between',
    whiteSpace: 'nowrap',
    paddingRight: '10px',
  },
  gridCapabilityHeader: {
    backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
  },
  gridExtensionHeader: {
    backgroundColor: theme.palette.type === 'dark' ? '#293B43 ' : '#C9DBE3',
  },
  gridExtensionItem: {
    backgroundColor: theme.palette.type === 'dark' ? '#3D4F57 ' : '#E7EFF3',
  },
  line: {
    border: '1px solid rgba(116,147,161, 0.3)',
    width: '100%',
    margin: '30px 0',
  },
  root: {
    width: '100%',
    paddingLeft: theme.spacing(15),
    paddingRight: theme.spacing(15),
    paddingBottom: theme.spacing(10),
    paddingTop: theme.spacing(5),
  },
});

const ThemeToggler = ({ classes }) => {
  const Component = ({ mode, toggleTheme }) => {
    return (
      <div>
        <Switch
          color="primary"
          classes={{
            switchBase: classes.switchBase,
            track: classes.track,
            checked: classes.checked,
            font: classes.checked,
          }}
          checked={mode === 'dark'}
          onChange={toggleTheme}
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
        <div className={props.classes.root}>
          <Typography variant="h5">Provider Information</Typography>
          <Grid container spacing={2}>
            {providerInfo &&
              Object.entries(providerInfo).map(
                ([providerName, provider], index) =>
                  (index < 2 || index === 3) && (
                    <Grid key={index} item md={4} xs={12}>
                      <Card className={props.classes.card}>
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
                          <Box className={props.classes.box}>
                            <Typography
                              variant="body1"
                              className={props.classes.hideScrollbar}
                              style={{ marginRight: '20px' }}
                            >
                              {provider}
                            </Typography>
                          </Box>
                        </CardContent>
                      </Card>
                    </Grid>
                  ),
              )}
          </Grid>
          {providerInfo &&
            Object.entries(providerInfo).map(
              ([providerName, provider], index) =>
                (index === 2 || index === 5) && (
                  <Card key={index} className={props.classes.card} sx={{ margin: '20px' }}>
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
                      <Box className={props.classes.box}>
                        <Typography
                          variant="body1"
                          className={props.classes.hideScrollbar}
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
                      </Box>
                    </CardContent>
                  </Card>
                ),
            )}

          <Card className={props.classes.card}>
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
          </Card>

          <hr className={props.classes.line} />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Capabilities
          </Typography>

          <Grid container spacing={2} style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Grid
              item
              xs={6}
              className={props.classes.gridCapabilityHeader}
              style={{ borderRadius: '10px 0 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Feature
              </Typography>
            </Grid>
            <Grid
              item
              xs={6}
              className={props.classes.gridCapabilityHeader}
              style={{ borderRadius: '0 10px 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Endpoint
              </Typography>
            </Grid>
            {providerInfo.capabilities &&
              providerInfo.capabilities.map((capability, index) => (
                <>
                  <Grid
                    item
                    key={`${index}-${capability.feature}`}
                    xs={6}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        props.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                          ? '#E7EFF3'
                          : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.feature}</Typography>
                  </Grid>
                  <Grid
                    item
                    key={`${index}-${capability.endpoint}`}
                    xs={6}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        props.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                          ? '#E7EFF3'
                          : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.endpoint}</Typography>
                  </Grid>
                </>
              ))}
          </Grid>
          <hr className={props.classes.line} />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Extensions
          </Typography>
          {providerInfo.extensions &&
            Object.entries(providerInfo.extensions).map(([extensionName, extension], index) => (
              <div key={index} margin="20px 0px">
                <Typography variant="h6"> {convertToTitleCase(extensionName)}</Typography>
                <Grid container spacing={2} style={{ margin: '10px 0 20px 0' }}>
                  <Grid
                    item
                    xs={6}
                    className={props.classes.gridExtensionHeader}
                    style={{
                      borderRadius: '10px 0 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Component
                    </Typography>
                  </Grid>
                  <Grid
                    item
                    xs={6}
                    className={props.classes.gridExtensionHeader}
                    style={{
                      borderRadius: '0 10px 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Type
                    </Typography>
                  </Grid>

                  <Grid
                    item
                    xs={6}
                    className={props.classes.gridExtensionItem}
                    style={{
                      borderRadius: '0 0 0 10px',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{extension[0].component}</Typography>
                  </Grid>
                  <Grid
                    item
                    xs={6}
                    className={props.classes.gridExtensionItem}
                    style={{
                      borderRadius: '0 0 10px 0',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{convertToTitleCase(extension[0].type)}</Typography>
                  </Grid>
                </Grid>
              </div>
            ))}
        </div>
      </NoSsr>
    );
  };

  const handleUpdateUserPref = (key, value) => {
    const updates = _.set(_.cloneDeep(userData), key, value);
    updateUserPrefWithContext(updates);
  };
  return (
    <NoSsr>
      <Paper square className={props.classes.paperRoot}>
        <Tabs
          value={tabVal}
          onChange={handleTabValChange}
          variant={width < 600 ? 'scrollable' : 'fullWidth'}
          scrollButtons="on"
          allowScrollButtonsMobile={true}
          indicatorColor="primary"
          textColor="primary"
          className={props.classes.tabs}
          centered
        >
          <CustomTooltip title="General preferences" placement="top">
            <Tab
              className={props.classes.tab}
              icon={<SettingsCellIcon style={iconMedium} />}
              label={<span className={props.classes.tabLabel}>General</span>}
            />
          </CustomTooltip>
          <CustomTooltip title="Choose Performance Test Defaults" placement="top">
            <Tab
              className={props.classes.tab}
              icon={<FontAwesomeIcon icon={faTachometerAlt} style={iconMedium} />}
              label={<span className={props.classes.tabLabel}>Performance</span>}
            />
          </CustomTooltip>
          {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
          {userPrefs && providerType != 'local' && (
            <CustomTooltip title="Remote Provider preferences" placement="top">
              <Tab
                className={props.classes.tab}
                icon={<SettingsRemoteIcon style={iconMedium} />}
                label={<span className={props.classes.tabLabel}>Remote Provider</span>}
              />
            </CustomTooltip>
          )}
        </Tabs>
      </Paper>
      <Paper className={props.classes.statsWrapper}>
        {tabVal === 0 && (
          <>
            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Extensions
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="CatalogContentPreference"
                    control={
                      <Switch
                        checked={catalogContent}
                        onChange={handleCatalogContentToggle}
                        color="primary"
                        classes={{
                          switchBase: props.classes.switchBase,
                          track: props.classes.track,
                          checked: props.classes.checked,
                        }}
                        data-cy="CatalogContentPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Meshery Catalog Content"
                  />
                </FormGroup>
              </FormControl>
            </div>
            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Analytics and Improvement Program
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="UsageStatsPreference"
                    control={
                      <Switch
                        checked={anonymousStats}
                        onChange={handleToggle('anonymousUsageStats')}
                        color="primary"
                        classes={{
                          switchBase: props.classes.switchBase,
                          track: props.classes.track,
                          checked: props.classes.checked,
                        }}
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
                        classes={{
                          switchBase: props.classes.switchBase,
                          track: props.classes.track,
                          checked: props.classes.checked,
                        }}
                        data-cy="PerfResultPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Send Anonymous Performance Results"
                  />
                </FormGroup>
              </FormControl>
            </div>

            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Theme
                </FormLabel>

                <FormGroup>
                  <FormControlLabel
                    key="ThemePreference"
                    control={
                      <ThemeToggler
                        handleUpdateUserPref={handleUpdateUserPref}
                        classes={props.classes}
                      />
                    }
                    labelPlacement="end"
                  />
                </FormGroup>
              </FormControl>
            </div>
          </>
        )}
        {tabVal === 1 && <MesherySettingsPerformanceComponent />}
        {tabVal === 2 && userPrefs && providerType !== 'local' && (
          <>
            <Tabs
              value={value}
              onChange={handleValChange}
              variant={width < 600 ? 'scrollable' : 'fullWidth'}
              scrollButtons="on"
              allowScrollButtonsMobile={true}
              indicatorColor="primary"
              textColor="primary"
              className={props.classes.tabs}
              centered
            >
              <CustomTooltip title="Details" placement="top">
                <Tab
                  className={props.classes.tab}
                  label={<span className={props.classes.tabLabel}>Details</span>}
                />
              </CustomTooltip>
              <CustomTooltip title="Plugins" placement="top">
                <Tab
                  className={props.classes.tab}
                  label={<span className={props.classes.tabLabel}>Plugins</span>}
                />
              </CustomTooltip>
            </Tabs>
            <Paper className={props.classes.statsWrapper}>
              {value === 0 && <RemoteProviderInfoTab />}

              {value === 1 && (
                <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
              )}
            </Paper>
          </>
        )}
      </Paper>
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

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(UserPreference)),
);
