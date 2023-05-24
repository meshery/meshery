//import useState from 'react';
import React, { useState, useLayoutEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import { withStyles } from '@material-ui/core/styles';
import CloseIcon from '@material-ui/icons/Close';
import {
  IconButton, FormControl, FormLabel, FormGroup, FormControlLabel, Switch
} from '@material-ui/core';
import NoSsr from '@material-ui/core/NoSsr';
import dataFetch from '../lib/data-fetch';
import { updateUser, updateProgress, toggleCatalogContent } from '../lib/store';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { Paper, Tooltip } from '@material-ui/core';
import SettingsRemoteIcon from '@material-ui/icons/SettingsRemote';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import ExtensionSandbox from "./ExtensionSandbox";
import RemoteComponent from "./RemoteComponent";
import ExtensionPointSchemaValidator from "../utils/ExtensionPointSchemaValidator";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTachometerAlt } from '@fortawesome/free-solid-svg-icons';
import MesherySettingsPerformanceComponent from './MesherySettingsPerformanceComponent';
import { ctxUrl } from '../utils/multi-ctx';
import { iconMedium } from '../css/icons.styles';
import { getTheme,setTheme } from "../utils/theme";
import { isExtensionOpen } from "../pages/_app";


const styles = (theme) => ({
  statsWrapper : {
    // padding : theme.spacing(2),
    maxWidth : "100%",
    height : 'auto',
    borderTopLeftRadius : 0,
    borderTopRightRadius : 0,
    borderBottomLeftRadius : 3,
    borderBottomRightRadius : 3,
  },
  paperRoot : {
    flexGrow : 1,
    maxWidth : "100%",
    marginLeft : 0,
    borderTopLeftRadius : 3,
    borderTopRightRadius : 3,
  },
  tabs : { marginLeft : 0,
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tab : {
    maxWidth : 'min(33%, 200px)',
    minWidth : '50px',
    margin : 0,
    "&.Mui-selected" : {
      color : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    }
  },
  icon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(1.75),
    marginLeft : theme.spacing(0.5),
  },
  iconText : {
    display : 'inline',
    verticalAlign : 'middle',
  },
  backToPlay : { margin : theme.spacing(2), },
  link : { cursor : 'pointer', },
  formContainer : {
    display : 'flex',
    'flex-wrap' : 'wrap',
    'justify-content' : 'space-evenly',
    padding : 50
  },
  formGrp : {
    padding : 20,
    border : '1.5px solid #969696',
    display : 'flex',
    width : "70%",
  },
  formLegend : { fontSize : 20, },
  formLegendSmall : { fontSize : 16 },
  switchBase : {
    color : '#647881',
    "&$checked" : { color : '#00b39f' },
    "&$checked + $track" : { backgroundColor : 'rgba(0,179,159,0.5)' },
  },
  track : { backgroundColor : 'rgba(100,120,129,0.5)', },
  checked : {},
  tabLabel : {
    [theme.breakpoints.up("sm")] : {
      fontSize : '1em'
    },
    [theme.breakpoints.between("xs", 'sm')] : {
      fontSize : '0.8em'
    }
  }
});

function ThemeToggler({ theme, themeSetter, enqueueSnackbar, classes }) {
  const [themeToggle, setthemeToggle] = useState(false);
  const defaultTheme = "light";
  const handle = () => {
    if (isExtensionOpen()) {
      return;
    }

    theme === "dark" ? setthemeToggle(true)  : setthemeToggle(false);
    setTheme(theme)
  };

  useLayoutEffect(() => {
    if (isExtensionOpen()) {
      if (getTheme() && getTheme() !== defaultTheme) {
        themeSetter(defaultTheme);
      }
      return;
    }

    themeSetter(getTheme() || defaultTheme);

  }, []);

  useLayoutEffect(handle, [theme]);

  const themeToggler = () => {
    if (isExtensionOpen()) {
      enqueueSnackbar("Toggling between themes is not supported in MeshMap", {
        variant : "info",
        preventDuplicate : true,
      });
      return;
    }
    theme === "light" ? themeSetter("dark")  : themeSetter("light");
  };

  return (
    <div onClick={themeToggler}>
      <Switch
        color="primary"
        classes={{
          switchBase : classes.switchBase,
          track : classes.track,
          checked : classes.checked,
          font : classes.checked,
        }}
        checked={themeToggle}
        onChange={themeToggler}
      />{" "}
      {themeToggle ? "Light Mode" : "Dark Mode"}
    </div>
  );
}


class UserPreference extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      anonymousStats : props.anonymousStats,
      perfResultStats : props.perfResultStats,
      tabVal : 0,
      userPrefs : ExtensionPointSchemaValidator("user_prefs")(),
      providerType : '',
      catalogContent : true,
      extensionPreferences : {},
      capabilitiesLoaded : false
    };
  }

  handleCatalogContentToggle = () => {
    this.props.toggleCatalogContent({ catalogVisibility : !this.state.catalogContent });
    this.setState((state) => ({ catalogContent : !state.catalogContent }), () => this.handleCatalogPreference());
  }

  handleCatalogPreference = () => {
    let body = Object.assign({}, this.state.extensionPreferences)
    body["catalogContent"] = this.state.catalogContent
    dataFetch("/api/user/prefs",
      { credentials : "include",
        method : "POST",
        body : JSON.stringify({ usersExtensionPreferences : body })
      },
      () => {
        this.props.enqueueSnackbar(`Catalog Content was ${this.state.catalogContent ? "enab" : "disab"}led`,
          { variant : 'success',
            autoHideDuration : 4000,
            action : (key) => (
              <IconButton
                key="close"
                aria-label="Close"
                color="inherit"
                onClick={() => this.props.closeSnackbar(key)}
              >
                <CloseIcon style={iconMedium} />
              </IconButton>
            ),
          });
      },
      this.handleError("There was an error sending your preference")
    )
  }

  handleToggle = (name) => () => {
    const self = this;
    if (name === 'anonymousUsageStats') {
      self.setState((state) => ({ anonymousStats : !state.anonymousStats }), () => this.handleChange(name));
    } else {
      self.setState((state) => ({ perfResultStats : !state.perfResultStats }), () => this.handleChange(name));
    }
  }

  handleError = (msg) => () => {
    const self = this;
    this.props.updateProgress({ showProgress : false });
    this.props.enqueueSnackbar(msg, { variant : 'error',
      action : (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => self.props.closeSnackbar(key)}
        >
          <CloseIcon style={iconMedium} />
        </IconButton>
      ),
      autoHideDuration : 8000, });
  }

  handleChange = (name) => {
    const self = this;
    const { anonymousStats, perfResultStats } = this.state;
    let val, msg;
    if (name === 'anonymousUsageStats') {
      val = anonymousStats;
      msg = val
        ? "Sending anonymous usage statistics was enabled"
        : "Sending anonymous usage statistics was disabled";
    } else {
      val = perfResultStats;
      msg = val
        ? "Sending anonymous performance results was enabled"
        : "Sending anonymous performance results was disabled";
    }

    const requestBody = JSON.stringify({
      "anonymousUsageStats" : anonymousStats,
      "anonymousPerfResults" : perfResultStats,
    });

    // console.log(requestBody,anonymousStats,perfResultStats);

    this.props.updateProgress({ showProgress : true });
    dataFetch(
      ctxUrl('/api/user/prefs', this.props.selectedK8sContexts), {
        credentials : 'include',
        method : 'POST',
        headers : { 'Content-Type' : 'application/json;charset=UTF-8', },
        body : requestBody,
      }, (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== 'undefined') {
          this.props.enqueueSnackbar(msg, { variant : val
            ? 'success'
            : 'info',
          autoHideDuration : 4000,
          action : (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key)}
            >
              <CloseIcon style={iconMedium} />
            </IconButton>
          ), });
        }
      }, self.handleError('There was an error sending your preference'));
  }

  handleTabValChange = (event, newVal) => {
    this.setState({ tabVal : newVal });
  }

  componentDidMount = () => {
    dataFetch(
      "/api/user/prefs",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (result) {
          this.setState({
            extensionPreferences : result?.usersExtensionPreferences,
            catalogContent : result?.usersExtensionPreferences?.catalogContent
          })
        }
      },
      err => console.error(err)
    )

  }

  componentDidUpdate() {
    const { capabilitiesRegistry } = this.props;
    if (capabilitiesRegistry && !this.state.capabilitiesLoaded) {
      this.setState({
        capabilitiesLoaded : true, // to prevent re-compute
        userPrefs : ExtensionPointSchemaValidator("user_prefs")(capabilitiesRegistry?.extensions?.user_prefs),
        providerType : capabilitiesRegistry?.provider_type,
      })
    }
  }

  render() {
    const {
      anonymousStats, perfResultStats, tabVal, userPrefs, providerType, catalogContent
    } = this.state;
    const { classes } = this.props;

    // const mainIconScale = 'grow-10';

    return (
      <NoSsr>
        <Paper square className={classes.paperRoot}>
          <Tabs
            value={tabVal}
            onChange={this.handleTabValChange}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
            className={classes.tabs}
          >
            <Tooltip title="General preferences" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <SettingsCellIcon style={iconMedium} />
                }
                label={<span className={classes.tabLabel}>General</span>}
              />
            </Tooltip>
            <Tooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faTachometerAlt} style={iconMedium}/>
                }
                label={<span className={classes.tabLabel}>Performance</span>}
              />
            </Tooltip>
            {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
            {userPrefs && providerType != 'local' &&
              <Tooltip title="Remote Provider preferences" placement="top">
                <Tab
                  className={classes.tab}
                  icon={
                    <SettingsRemoteIcon style={iconMedium} />
                  }
                  label={<span className={classes.tabLabel}>Remote Provider</span>}
                />
              </Tooltip>
            }
          </Tabs>
        </Paper>
        <Paper className={classes.statsWrapper}>
          {tabVal === 0 &&
          <>
            <div className={classes.formContainer}>
              <FormControl component="fieldset" className={classes.formGrp}>
                <FormLabel component="legend" className={classes.formLegend}>Extensions</FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="CatalogContentPreference"
                    control={(
                      <Switch
                        checked={catalogContent}
                        onChange={this.handleCatalogContentToggle}
                        color="primary"
                        classes={{ switchBase : classes.switchBase,
                          track : classes.track,
                          checked : classes.checked, }}
                        data-cy="CatalogContentPreference"
                      />
                    )}
                    labelPlacement="end"
                    label="Meshery Catalog Content"
                  />
                </FormGroup>
              </FormControl>
            </div>
            <div className={classes.formContainer}>
              <FormControl component="fieldset" className={classes.formGrp}>
                <FormLabel component="legend" className={classes.formLegend}>Analytics and Improvement Program</FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="UsageStatsPreference"
                    control={(
                      <Switch
                        checked={anonymousStats}
                        onChange={this.handleToggle('anonymousUsageStats')}
                        color="primary"
                        classes={{ switchBase : classes.switchBase,
                          track : classes.track,
                          checked : classes.checked, }}
                        data-cy="UsageStatsPreference"
                      />
                    )}
                    labelPlacement="end"
                    label="Send Anonymous Usage Statistics"
                  />
                  <FormControlLabel
                    key="PerfResultPreference"
                    control={(
                      <Switch
                        checked={perfResultStats}
                        onChange={this.handleToggle('anonymousPerfResults')}
                        color="primary"
                        classes={{ switchBase : classes.switchBase,
                          track : classes.track,
                          checked : classes.checked, }}
                        data-cy="PerfResultPreference"
                      />
                    )}
                    labelPlacement="end"
                    label="Send Anonymous Performance Results"
                  />
                </FormGroup>
              </FormControl>
            </div>

            <div className={classes.formContainer} >
              <FormControl component="fieldset" className={classes.formGrp}>
                <FormLabel component="legend" className={classes.formLegend}>
                    Theme
                </FormLabel>

                <FormGroup>
                  <FormControlLabel
                    key="ThemePreference"
                    control={(
                      <ThemeToggler
                        classes={classes}
                        theme={this.props.theme}
                        themeSetter={this.props.themeSetter}
                        enqueueSnackbar={this.props.enqueueSnackbar}

                      />
                    )}
                    labelPlacement="end"
                    // label="Theme"
                  />
                </FormGroup>
              </FormControl>
            </div>
          </>
          }
          {tabVal === 1 &&
            <MesherySettingsPerformanceComponent />
          }
          {tabVal === 2 && userPrefs && providerType !== 'local' &&
            <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
          }

        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({ updateUser : bindActionCreators(updateUser, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch), toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch) });

const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const catalogVisibility = state.get('catalogVisibility');
  const capabilitiesRegistry = state.get("capabilitiesRegistry")
  return {
    selectedK8sContexts,
    catalogVisibility,
    capabilitiesRegistry
  };
};


export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withSnackbar(UserPreference))));