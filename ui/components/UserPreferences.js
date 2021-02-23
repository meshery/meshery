import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import { withStyles } from '@material-ui/core/styles';
import CloseIcon from '@material-ui/icons/Close';
import { IconButton, FormControl, FormLabel, FormGroup, FormControlLabel, Switch } from '@material-ui/core';
import NoSsr from '@material-ui/core/NoSsr';
import dataFetch from '../lib/data-fetch';
import { updateUser, updateProgress } from '../lib/store';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { Paper, Tooltip } from '@material-ui/core';
import SettingsRemoteIcon from '@material-ui/icons/SettingsRemote';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import ExtensionSandbox from "./ExtensionSandbox";
import RemoteUserPref from "./RemoteUserPref";
import ExtensionPointSchemaValidator from "../utils/ExtensionPointSchemaValidator";


const styles = (theme) => ({
  root: {
    maxWidth: "100%",
    height: 'auto',
    borderTopLeftRadius: 0,
    borderTopRightRadius: 0,
    borderBottomLeftRadius: 3,
    borderBottomRightRadius: 3,
  },
  paperRoot: {
    flexGrow: 1,
    maxWidth: "20%",
    marginLeft: 0,
    borderTopLeftRadius: 3,
    borderTopRightRadius: 3,
  },
  tabs: {
    marginLeft: 0
  },
  tab: {
    minWidth: "50%",
    margin: 0
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
  backToPlay: {
    margin: theme.spacing(2),
  },
  link: {
    cursor: 'pointer',
  },
  formContainer: {
    display: 'flex',
    'flex-wrap': 'wrap',
    'justify-content': 'space-evenly',
    padding: 50
  },
  formGrp: {
    padding: 20,
    border: '1.5px solid #969696',
  },
  formLegend: {
    fontSize: 20,
  },
  switchBase: {
    color: '#647881',
    "&$checked": {
      color: '#00b39f'
    },
    "&$checked + $track": {
      backgroundColor: 'rgba(0,179,159,0.5)'
    },
  },
  track: {
    backgroundColor: 'rgba(100,120,129,0.5)',
  },
  checked: {},
});

class UserPreference extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      anonymousStats: props.anonymousStats,
      perfResultStats: props.perfResultStats,
      startOnZoom: props.startOnZoom,
      tabVal: 0,
      userPrefs: ExtensionPointSchemaValidator("user_prefs")()
    };
  }

  handleToggle = (name) => () => {
    const self = this;
    if (name == 'anonymousUsageStats') {
      self.setState((state) => ({ anonymousStats: !state.anonymousStats }));
    } else if (name == 'anonymousPerfResults') {
      self.setState((state) => ({ perfResultStats: !state.perfResultStats }));
    } else {
      self.setState((state) => ({ startOnZoom: !state.startOnZoom }));
    }

    this.handleChange(name);
  };

  handleError = (msg) => () => {
    const self = this;
    this.props.updateProgress({ showProgress: false });
    this.props.enqueueSnackbar(msg, {
      variant: 'error',
      action: (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => self.props.closeSnackbar(key)}
        >
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration: 8000,
    });
  }

  handleChange = (name) => {
    const self = this;
    const { anonymousStats, perfResultStats, startOnZoom } = this.state;
    let val, msg;
    if (name == 'anonymousUsageStats') {
      val = anonymousStats;
      msg = !val ? "Sending anonymous usage statistics was enabled"
        : "Sending anonymous usage statistics was disabled";

    } else if (name == 'anonymousPerfResults') {
      val = perfResultStats;
      msg = !val ? "Sending anonymous performance results was enabled"
        : "Sending anonymous performance results was disabled";
    } else {
      val = startOnZoom;
      msg = !val ? "Start on Zoom was enabled"
        : "Start on Zoom was disabled";
    }

    const params = `${encodeURIComponent(name)}=${encodeURIComponent(!val)}`;
    this.props.updateProgress({ showProgress: true });
    dataFetch('/api/user/stats', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
      },
      body: params,
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        this.props.enqueueSnackbar(msg, {
          variant: !val ? 'success' : 'info',
          autoHideDuration: 4000,
          action: (key) => (
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
      }
    }, self.handleError('There was an error sending your preference'));
  }

  handleTabValChange = (event, newVal) => {
    this.setState({ tabVal: newVal });
  }

  componentDidMount = () => {
    dataFetch(
      "/api/provider/capabilities",
      {
        credentials: "same-origin",
        method: "GET",
        credentials: "include",
      },
      (result) => {
        if (result) {
          this.setState({
            userPrefs: ExtensionPointSchemaValidator("user-prefs")(result?.extensions?.user_prefs)
          })
        }
      },
      err => console.error(err)
    )
  }

  render() {
    const { anonymousStats, perfResultStats, tabVal, startOnZoom, userPrefs } = this.state;
    const { classes } = this.props;

    const handleToggle = this.handleToggle('startOnZoom');

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
                label="General"
              />
            </Tooltip>
            {userPrefs &&
              <Tooltip title="Remote Provider preferences" placement="top">
                <Tab
                  className={classes.tab}
                  icon={
                    <SettingsRemoteIcon />
                  }
                  label="Remote Provider"
                />
              </Tooltip>
            }
          </Tabs>
        </Paper>
        <Paper className={classes.root}>
          {tabVal == 0 &&
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
                        classes={{
                          switchBase: classes.switchBase,
                          track: classes.track,
                          checked: classes.checked,
                        }}
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
                        classes={{
                          switchBase: classes.switchBase,
                          track: classes.track,
                          checked: classes.checked,
                        }}
                        data-cy="PerfResultPreference"
                      />
                    )}
                    labelPlacement="end"
                    label="Send Anonymous Performance Results"
                  />
                </FormGroup>
              </FormControl>
            </div>
          }
          {tabVal == 1 && userPrefs &&
            <ExtensionSandbox type="user_prefs" Extension={() => RemoteUserPref({startOnZoom, handleToggle})}/>
          }
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updateUser: bindActionCreators(updateUser, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps,
)(withRouter(withSnackbar(UserPreference))));