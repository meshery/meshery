//import useState from "react"
import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import CopyIcon from '../../assets/icons/CopyIcon';
import _ from 'lodash';
import {FormLabel} from '@material-ui/core';
import {
  CustomTooltip,
  Typography,
  Grid,
  FormControl,
  FormGroup,
  FormControlLabel,
  Switch,
  IconButton,
  Card,
  CardContent,
  CardHeader,
  Box,
  Tabs,
  Tab,
  Paper,
  styled,
  useTheme,
} from '@layer5/sistent';
import NoSsr from '@material-ui/core/NoSsr';
import { updateUser, updateProgress, toggleCatalogContent } from '../../lib/store';
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
import { UsesSistent } from '@/components/SistentWrapper';

export const StyledPaper = styled(Paper)(() => {
  const theme = useTheme();
  return {
    flexGrow: 1,
    maxWidth: '100%',
    marginLeft: 0,
    backgroundColor:
    theme.palette.mode === 'dark'
      ? theme.palette?.constant?.disabled 
      : theme.palette?.constant?.disabled,
    borderTopLeftRadius: 3,
    borderTopRightRadius: 3,
  };
});

export const StyledBox = styled(Box)(() => ({
  width: '100%',
  display: 'flex',
  justifyContent: 'space-between',
  whiteSpace: 'nowrap',
  paddingRight: '10px',
}));

export const StyledTabs = styled(Tabs)(({ theme }) => ({
  width: '100%',
  marginLeft: 0,
  '& .MuiTabs-indicator': {
    backgroundColor: theme.palette.mode === 'dark' 
      ? theme.palette.background.brand.default  
      : theme.palette.primary.main,
  },
  '& .MuiTab-root': {
    backgroundColor: 'transparent',
    color: theme.palette.mode === 'dark' 
      ? '#ffffff'  
      : theme.palette.text.primary,
    '&.Mui-selected': {
      backgroundColor: 'transparent',
      color: theme.palette.mode === 'dark' 
        ? theme.palette.background.brand.default  
        : theme.palette.primary.main,
    },
  }
}));

export const StyledTab = styled(Tab)(({ theme }) => ({
  width: '42%',
  margin: 0,
  '&.Mui-selected': {
    color: theme.palette.mode === 'dark' ? theme.palette.background.brand.default : theme.palette.background.brand.default,
  },
  [theme.breakpoints.up('sm')]: {
    '& .MuiTab-labelIcon': {
      fontSize: '1em',
    },
  },
  [theme.breakpoints.between('xs', 'sm')]: {
    '& .MuiTab-labelIcon': {
      fontSize: '0.8em',
    },
  },
}));

export const StyledCard = styled(Card)(({ theme }) => ({
  border: '1px solid rgba(0,179,159,0.3)',
  margin: '20px 0px',
  backgroundColor: theme.palette.mode === 'dark' ? theme.palette.background.secondary : '#C9DBE3',
}));

export const StyledFormContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: theme.spacing(6),
}));

export const StyledFormGroup = styled(FormControl)(({ theme }) => ({
  padding: theme.spacing(2.5),
  border: '1.5px solid #969696',
  display: 'flex',
  width: '70%',
}));

export const StyledSwitch = styled(Switch)(() => ({
  '& .MuiSwitch-switchBase': {
    color: '#647881',
    '&.Mui-checked': {
      color: '#00b39f',
      '& + .MuiSwitch-track': {
        backgroundColor: 'rgba(0,179,159,0.5)',
      },
    },
  },
  '& .MuiSwitch-track': {
    backgroundColor: 'rgba(100,120,129,0.5)',
  },
}));

const ThemeToggler = () => {
  const Component = ({ mode, toggleTheme }) => {
    return (
      <UsesSistent>
      <div>
        <StyledSwitch checked={mode === 'dark'} onChange={toggleTheme} />
        Dark Mode
      </div>
      </UsesSistent>
    );
  };

  return <ThemeTogglerCore Component={Component} />;
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

  const theme = useTheme();

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
           <UsesSistent>
        <StyledPaper
          sx={{
            width: '100%',
            paddingLeft: theme.spacing(15),
            paddingRight: theme.spacing(15),
            paddingBottom: theme.spacing(10),
            paddingTop: theme.spacing(5),
          }}
        >
          <Typography variant="h5">Provider Information</Typography>
          <Grid container spacing={2}>
            {providerInfo &&
              Object.entries(providerInfo).map(
                ([providerName, provider], index) =>
                  (index < 2 || index === 3) && (
                    <Grid key={index} item md={4} xs={12}>
                      <StyledCard>
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
                          <StyledBox>
                            <Typography
                              variant="body1"
                              sx={{
                                marginRight: '20px',
                                overflowX: 'auto',
                                '&::-webkit-scrollbar': {
                                  display: 'none',
                                },
                                '&::-moz-scrollbar': {
                                  display: 'none',
                                },
                              }}
                            >
                              {provider}
                            </Typography>
                          </StyledBox>
                        </CardContent>
                      </StyledCard>
                    </Grid>
                  ),
              )}
          </Grid>
          {providerInfo &&
            Object.entries(providerInfo).map(
              ([providerName, provider], index) =>
                (index === 2 || index === 5) && (
                  <StyledCard key={index} sx={{ margin: '20px' }}>
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
                      <StyledBox>
                        <Typography
                          variant="body1"
                          sx={{
                            marginRight: '20px',
                            overflowX: 'auto',
                            '&::-webkit-scrollbar': {
                              display: 'none',
                            },
                            '&::-moz-scrollbar': {
                              display: 'none',
                            },
                          }}
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
                      </StyledBox>
                    </CardContent>
                  </StyledCard>
                ),
            )}

          <StyledCard>
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
          </StyledCard>

          <hr
            sx={{ border: '1px solid rgba(116,147,161, 0.3)', width: '100%', margin: '30px 0' }}
          />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Capabilities
          </Typography>

          <Grid container spacing={2} style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Grid
              item
              xs={6}
              sx={{
                backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
                borderRadius: '10px 0 0 0',
                padding: '10px 20px',
              }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Feature
              </Typography>
            </Grid>
            <Grid
              item
              xs={6}
              sx={{
                backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
                borderRadius: '10px 0 0 0',
                padding: '10px 20px',
              }}
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
                  </Grid>
                  <Grid
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
                  </Grid>
                </>
              ))}
          </Grid>
          <Box
            sx={{ border: '1px solid rgba(116,147,161, 0.3)', width: '100%', margin: '30px 0' }}
          />
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
                    sx={{
                      backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
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
                    sx={{
                      backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
                      borderRadius: '10px 0 0 0',
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
                    sx={{
                      backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
                      borderRadius: '0 0 0 10px',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{extension[0].component}</Typography>
                  </Grid>
                  <Grid
                    item
                    xs={6}
                    sx={{
                      backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#7493A1',
                      borderRadius: '0 0 10px 0',
                      padding: '20px 20px',
                    }}
                  >
                    <Typography variant="body1">{convertToTitleCase(extension[0].type)}</Typography>
                  </Grid>
                </Grid>
              </div>
            ))}
        </StyledPaper>
        </UsesSistent>
      </NoSsr>
    );
  };

  const handleUpdateUserPref = (key, value) => {
    const updates = _.set(_.cloneDeep(userData), key, value);
    updateUserPrefWithContext(updates);
  };
  return (
    <NoSsr>
      <UsesSistent>
      <StyledPaper square>
        <StyledTabs
          value={tabVal}
          onChange={handleTabValChange}
          variant={width < 600 ? 'scrollable' : 'fullWidth'}
          scrollButtons={true}
          allowScrollButtonsMobile={true}
          indicatorColor="primary"
          textColor="primary"
          centered
          sx={{backgroundColor: theme.palette?.info?.default}}
        >
          <CustomTooltip title="General preferences" placement="top">
            <StyledTab
              icon={<SettingsCellIcon style={iconMedium} />}
              label={
                <span
                  sx={{
                    [theme.breakpoints.up('sm')]: {
                      fontSize: '1em',
                    },
                    [theme.breakpoints.between('xs', 'sm')]: {
                      fontSize: '0.8em',
                    },
                  }}
                >
                  General
                </span>
              }
            />
          </CustomTooltip>
          <CustomTooltip title="Choose Performance Test Defaults" placement="top">
            <StyledTab
              icon={<FontAwesomeIcon icon={faTachometerAlt} style={iconMedium} />}
              label={
                <span
                  sx={{
                    [theme.breakpoints.up('sm')]: {
                      fontSize: '1em',
                    },
                    [theme.breakpoints.between('xs', 'sm')]: {
                      fontSize: '0.8em',
                    },
                  }}
                >
                  Performance
                </span>
              }
            />
          </CustomTooltip>
          {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
          {userPrefs && providerType != 'local' && (
            <CustomTooltip title="Remote Provider preferences" placement="top">
              <StyledTab
                icon={<SettingsRemoteIcon style={iconMedium} />}
                label={
                  <span
                    sx={{
                      [theme.breakpoints.up('sm')]: {
                        fontSize: '1em',
                      },
                      [theme.breakpoints.between('xs', 'sm')]: {
                        fontSize: '0.8em',
                      },
                    }}
                  >
                    Remote Provider
                  </span>
                }
              />
            </CustomTooltip>
          )}
        </StyledTabs>
      </StyledPaper>
      <StyledPaper
        sx={{
          height: 'auto',
          borderTopLeftRadius: 0,
          borderTopRightRadius: 0,
          borderBottomLeftRadius: 3,
          borderBottomRightRadius: 3,
        }}
      >
        {tabVal === 0 && (
          <>
            <StyledFormContainer>
              <StyledFormGroup component="fieldset">
                <FormLabel component="legend" sx={{ fontSize: 20 }}>
                  Extensions
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="CatalogContentPreference"
                    control={
                      <StyledSwitch
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
              </StyledFormGroup>
            </StyledFormContainer>
            <StyledFormContainer>
              <StyledFormGroup component="fieldset">
                <FormLabel component="legend" sx={{ fontSize: 16 }}>
                  Analytics and Improvement Program
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="UsageStatsPreference"
                    control={
                      <StyledSwitch
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
                      <StyledSwitch
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
              </StyledFormGroup>
            </StyledFormContainer>

            <StyledFormContainer>
              <StyledFormGroup component="fieldset">
                <FormLabel component="legend" sx={{ fontSize: 16 }}>
                  Theme
                </FormLabel>

                <FormGroup>
                  <FormControlLabel
                    key="ThemePreference"
                    control={
                      <ThemeToggler
                        handleUpdateUserPref={handleUpdateUserPref}
                        // classes={props.classes}
                      />
                    }
                    labelPlacement="end"
                  />
                </FormGroup>
              </StyledFormGroup>
            </StyledFormContainer>
          </>
        )}
        {tabVal === 1 && <MesherySettingsPerformanceComponent />}
        {tabVal === 2 && userPrefs && providerType !== 'local' && (
          <>
            <StyledTabs
              value={value}
              onChange={handleValChange}
              variant={width < 600 ? 'scrollable' : 'fullWidth'}
              scrollButtons={true}
              allowScrollButtonsMobile={true}
              indicatorColor="primary"
              textColor="primary"
              centered
            >
              <CustomTooltip title="Details" placement="top">
                <StyledTab
                  label={
                    <span
                      sx={{
                        [theme.breakpoints.up('sm')]: {
                          fontSize: '1em',
                        },
                        [theme.breakpoints.between('xs', 'sm')]: {
                          fontSize: '0.8em',
                        },
                      }}
                    >
                      Details
                    </span>
                  }
                />
              </CustomTooltip>
              <CustomTooltip title="Plugins" placement="top">
                <StyledTab
                  label={
                    <span
                      sx={{
                        [theme.breakpoints.up('sm')]: {
                          fontSize: '1em',
                        },
                        [theme.breakpoints.between('xs', 'sm')]: {
                          fontSize: '0.8em',
                        },
                      }}
                    >
                      Plugins
                    </span>
                  }
                />
              </CustomTooltip>
            </StyledTabs>
            <StyledPaper
              sx={{
                height: 'auto',
                borderTopLeftRadius: 0,
                borderTopRightRadius: 0,
                borderBottomLeftRadius: 3,
                borderBottomRightRadius: 3,
              }}
            >
              {value === 0 && <RemoteProviderInfoTab />}

              {value === 1 && (
                <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
              )}
            </StyledPaper>
          </>
        )}
      </StyledPaper>
      </UsesSistent>
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


