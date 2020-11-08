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


const styles = () => ({
  formContainer: {
    margin: 50,
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
    };
  }

  handleToggle = (name) => () => {
    const self = this;
    if (name == 'anonymousUsageStats') {
      self.setState((state) => ({ anonymousStats: !state.anonymousStats }));
    } else {
      self.setState((state) => ({ perfResultStats: !state.perfResultStats }));
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
    const { anonymousStats, perfResultStats } = this.state;
    let val, msg;
    if (name == 'anonymousUsageStats') {
      val = anonymousStats;
      msg = !val ? "Sending anonymous usage statistics was enabled"
        : "Sending anonymous usage statistics was disabled";

    } else {
      val = perfResultStats;
      msg = !val ? "Sending anonymous performance results was enabled"
        : "Sending anonymous performance results was disabled";
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
          variant: !val ? 'success': 'info',
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

  render() {
    const { anonymousStats, perfResultStats } = this.state;
    const { classes } = this.props;

    return (
      <NoSsr>
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
