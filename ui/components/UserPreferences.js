import React from 'react';
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


const styles = (theme) => ({
  statsWrapper : {
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
  tabs : { marginLeft : 0 },
  tab : {
    maxWidth : 'min(33%, 200px)',
    minWidth : '50px',
    margin : 0
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
  },
  formLegend : { fontSize : 20, },
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
      extensionPreferences : {}
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
                onClick={() => self.props.closeSnackbar(key)}
              >
                <CloseIcon />
              </IconButton>
            ),
          });
      },
      this.handleError("There was an error sending your preference")
    )
  }

  handleToggle = (name) => () => {
    const self = this;
    if (name == 'anonymousUsageStats') {
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
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 8000, });
  }

  handleChange = (name) => {
    const self = this;
    const { anonymousStats, perfResultStats } = this.state;
    let val, msg;
    if (name == 'anonymousUsageStats') {
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

    console.log(requestBody,anonymousStats,perfResultStats);

    this.props.updateProgress({ showProgress : true });
    dataFetch(
      ctxUrl('/api/user/prefs', this.props.selectedK8sContexts), {
        credentials : 'same-origin',
        method : 'POST',
        credentials : 'include',
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
              <CloseIcon />
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
      "/api/provider/capabilities",
      { credentials : "same-origin",
        method : "GET",
        credentials : "include", },
      (result) => {
        if (result) {
          this.setState({
            userPrefs : ExtensionPointSchemaValidator("user_prefs")(result?.extensions?.user_prefs),
            providerType : result?.provider_type
          })
        }
      },
      err => console.error(err)
    )

    dataFetch(
      "/api/user/prefs",
      { credentials : "same-origin",
        method : "GET",
        credentials : "include", },
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

  render() {
    const {
      anonymousStats, perfResultStats, tabVal, userPrefs, providerType, catalogContent
    } = this.state;
    const { classes } = this.props;

    const mainIconScale = 'grow-10';

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
                  <SettingsCellIcon />
                }
                label={<span className={classes.tabLabel}>General</span>}
              />
            </Tooltip>
            <Tooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faTachometerAlt} transform={mainIconScale} fixedWidth />
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
                    <SettingsRemoteIcon />
                  }
                  label={<span className={classes.tabLabel}>Remote Provider</span>}
                />
              </Tooltip>
            }
          </Tabs>
        </Paper>
        <Paper className={classes.statsWrapper}>
          {tabVal == 0 &&
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
          </>
          }
          {tabVal === 1 &&
            <MesherySettingsPerformanceComponent />
          }
          {tabVal == 2 && userPrefs && providerType != 'local' &&
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
  const catalogVisibility = state.get('catalogVisibility')
  return {
    selectedK8sContexts,
    catalogVisibility
  };
};


export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withSnackbar(UserPreference))));